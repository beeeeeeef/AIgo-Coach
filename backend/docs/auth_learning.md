# 登录注册后端 — 学习笔记

> 记录开发顺序、每层职责、架构决策及其原因。

---

## 第 1 步：数据库建表（migration）

### 写了什么

- `backend/migration/001_create_user.sql` — users 表
- `backend/migration/002_create_verify_codes.sql` — verify_codes 表

### 这一层的作用

定义数据库的物理结构。表结构是所有代码的底层基础——Model 照搬表字段，Repository 操作表数据。

### 为什么第一步做

1. 表结构决定了 Model 长什么样
2. Model 决定了 Repository 接口的入参和返回值
3. 自顶向下写代码容易"凭空想象"，先有表结构再写代码每一步都有依据

### 踩坑

`001_create_user.sql` 最后一列 `created_at` 后面多了一个逗号，SQL 会报错。后来修正了。

---

## 第 2 步：定义 Model

### 写了什么

- `backend/model/user.go` — User 实体
- `backend/model/verify_code.go` — VerifyCode 实体

### 这一层的作用

Model 是 Go 侧的数据实体，一行代码对应数据库里的一行数据。各层之间用 Model 传递数据，它是"通用语言"。

```go
type User struct {
    ID           int64     `db:"id"`
    Email        string    `db:"email"`
    PasswordHash string    `db:"password_hash"`
    CreatedAt    time.Time `db:"created_at"`
}
```

### 为什么 Model 不放业务逻辑

Model 只负责"数据结构"，不负责"该怎么用"。业务规则放在 Service 层，这样同样的 Model 可以被不同的 Service 复用。

### 踩坑

1. VerifyCode 最初没有 `ID` 字段，但 Repository 接口里有 `IncrementApplyTimes(id int64)` 和 `MarkUsed(id int64)`，需要 ID 来定位记录。后来补上了。
2. 字段名 `ExpiredAt` vs `expires_at`：Go 用 `ExpiredAt`（Expired + At），数据库用 `expired_at`，统一确认了这个命名。

---

## 第 3 步：定义 DTO

### 写了什么

- `backend/handler/auth.dto.go` — 请求 DTO + 响应 DTO

### 这一层的作用

DTO（Data Transfer Object）只负责前后端传输结构，不等于数据库实体。

| DTO | 方向 | 作用 |
|-----|------|------|
| `SendRegisterCodeRequest` | 前端 → 后端 | 发送注册验证码 |
| `RegisterRequest` | 前端 → 后端 | 注册 |
| `LoginByPasswordRequest` | 前端 → 后端 | 密码登录 |
| `LoginByCodeRequest` | 前端 → 后端 | 验证码登录 |
| `LoginResponse` | 后端 → 前端 | 登录成功返回 token + 用户信息 |
| `MeResponse` | 后端 → 前端 | 当前用户信息 |

### 为什么 DTO 不等于 Model

前端需要的字段 ≠ 数据库字段。例如：

- 注册时前端传 `password`（明文），数据库存 `password_hash`（bcrypt 加密后）。DTO 里不出现 `password_hash`。
- 登录成功返回 `token`，但 `token` 不存在于 User 表里，它是 JWT 生成的。

DTO 隔离了这种差异——前端和后端各自用自己舒服的结构，互不污染。

---

## 第 4 步：定义 Repository 接口

### 写了什么

- `backend/repository/user_repository.go` — UserRepository + VerifyCodeRepository 接口

### 这一层的作用

接口定义了"数据访问能做什么"，不做业务判断。它是 Service 层和数据库之间的**契约**。

```go
type UserRepository interface {
    GetByEmail(email string) (*model.User, error)
    GetByID(id int64) (*model.User, error)
    Create(user *model.User) error
}

type VerifyCodeRepository interface {
    Create(code *model.VerifyCode) error
    GetLatestByEmailAndType(email string, codeType string) (*model.VerifyCode, error)
    RefreshCode(id int64, code string, expiredAt time.Time) error
    IncrementApplyTimes(id int64) error
    MarkUsed(id int64) error
}
```

### 为什么先写接口，后写实现

1. **解耦**：Service 只依赖接口，不关心底层是 MySQL 还是 PostgreSQL 还是 Mock
2. **可测试**：测试 Service 时可以注入一个假的 Repository，不需要真实数据库
3. **接口即文档**：看接口就知道数据访问层提供了哪些能力

### 设计原则

Repository 只负责 **"查/增/改"**，不负责 **"该不该查"**。

例如：
- `GetLatestByEmailAndType` 只负责查最新一条记录 —— ✅
- 判断验证码是否过期 —— ❌ 这是 Service 的事
- 判断能不能发送 —— ❌ 这也是 Service 的事

---

## 第 5 步：实现 Repository（SQL）

### 写了什么

- `backend/repository/user_SQL.go` — UserRepository 的 MySQL 实现
- `backend/repository/verify_code_SQL.go` — VerifyCodeRepository 的 MySQL 实现

### 这一层的作用

把接口翻译成具体的 SQL 语句，真正操作数据库。

每个方法做的事：
- `GetByEmail` → `SELECT ... WHERE email = ?`
- `Create` → `INSERT INTO ...`
- `RefreshCode` → `UPDATE ... SET code = ?, expired_at = ?, apply_times = apply_times + 1, used = false WHERE id = ?`
- `MarkUsed` → `UPDATE ... SET used = true WHERE id = ?`

### 为什么 `GetByEmail` 返回 `(*User, nil)` 而不是报错

没找到用户不是"系统错误"，是"正常情况"。返回 `nil, nil` 让 Service 层自己判断用户是否存在，而不是被迫处理一个 error。

### 为什么有 `RefreshCode` 这个方法

用户 60 秒后重新点"发送验证码"，不应该 INSERT 一条新记录（会产生大量垃圾数据），而是 UPDATE 已有记录——换一个新的 code、重置过期时间、`apply_times + 1`、重置 `used`。

---

## 第 6 步：写工具函数（Utils）

### 写了什么

- `backend/utils/password.go` — bcrypt 密码哈希 + 验证
- `backend/utils/code.go` — 6 位数字验证码生成

### 这一层的作用

纯函数，不依赖任何其他层，可被任何层调用。

| 函数 | 作用 | 被谁调用 |
|------|------|---------|
| `HashPassword(password)` | 明文 → bcrypt 哈希 | Service（注册时） |
| `CheckPasswordHash(password, hash)` | 验证密码是否匹配 | Service（密码登录时） |
| `GenerateCode()` | 生成 6 位数字验证码 | Service（发送验证码时） |

### 为什么抽出来

1. 密码哈希和验证码生成不是**业务逻辑**，是**工具**
2. 如果以后换加密算法（bcrypt → argon2），只改这一个文件
3. Service 层专注业务流程，不需要关心 bcrypt 的具体 API

---

## 第 7 步：写 Service 层

### 写了什么

- `backend/service/users/service.go` — UserService 结构体 + 依赖注入
- `backend/service/users/register.go` — 注册相关业务
- `backend/service/users/login.go` — 登录相关业务

### 这一层的作用

Service 是**整个系统的核心**——写业务流程和规则，判断"该不该这么做"。

### 每个方法的职责

| 方法 | 做的事情 | 不做什么 |
|------|---------|---------|
| `SendRegisterCode(email)` | 检查邮箱未注册 → 检查冷却期 → 检查发送上限 → 生成/刷新验证码 | 不发邮件（那是另一个服务的事） |
| `Register(email, code, password)` | 校验验证码 → 检查邮箱未注册 → 加密密码 → 创建用户 → 标记验证码已使用 | 不返回 token（token 是 Handler 的事） |
| `LoginByPassword(email, password)` | 查用户 → 验证密码 | 不返回 token |
| `SendLoginCode(email)` | 检查邮箱已注册 → 检查冷却期 → 生成/刷新验证码 | 不发邮件 |
| `LoginByCode(email, code)` | 校验验证码 → 返回用户 | 不返回 token |

### 为什么 Service 不返回 token

签发 token 不是业务规则，是**基础设施**。Service 只回答"这个人能不能登录"，至于登录成功之后怎么发通行证，由 Handler 决定（调用 JWT 工具生成 token）。

### 架构决策

#### 决策 1：验证码维护两条记录（register + login）

同一个邮箱可以同时有两条验证码记录，一条 `type="register"`，一条 `type="login"`。

**原因**：
- 注册和登录是两种不同的业务场景，验证规则不同（注册要求邮箱不存在，登录要求邮箱存在）
- 共享一条记录会导致互相覆盖：用户刚收到注册码，切到登录页点了"发送验证码"，注册码就被登录码覆盖了
- 安全边界更清晰：即使攻击者拿到验证码，也只能用于指定场景

#### 决策 2：非首次发送用 RefreshCode，而不是新建记录

用户 60 秒后重新点"发送验证码"时，UPDATE 已有记录而不是 INSERT 新记录。

**原因**：
- 减少垃圾数据
- `apply_times` 能正确累计，实现发送次数上限（5 次）
- `RefreshCode` 不需要更新 `type` 字段，因为查询时已经用 `type` 过滤，记录的 type 本来就是目标值

#### 决策 3：验证码错误时调用 IncrementApplyTimes

输入错误验证码时，`apply_times + 1`。

**原因**：防止暴力破解。如果有人不断尝试不同验证码，apply_times 会快速增长，达到上限后可以触发额外限制。

#### 决策 4：MarkUsed 失败不阻断注册

注册成功后标记验证码已使用，但如果标记失败，不 return error，用户仍然注册成功。

**原因**：用户已创建是**核心结果**，标记验证码是**辅助操作**。辅助操作失败不应该回滚核心结果。

#### 决策 5：Service 注入 Repository 接口，不注入具体实现

```go
type UserService struct {
    UserRepo       repository.UserRepository       // 接口，不是 *UserRepositorySQL
    VerifyCodeRepo repository.VerifyCodeRepository // 接口，不是 *VerifyCodeRepositorySQL
}
```

**原因**：
- Service 不需要知道底层是什么数据库
- 测试时可以注入假的 Repository
- 符合"依赖倒置原则"——高层模块不依赖低层模块，两者都依赖抽象

---

## 待续

以下章节将在后续开发中补充：

- 第 8 步：写 JWT 工具（`utils/jwt.go`）
- 第 9 步：写认证中间件（`middleware/auth.go`）
- 第 10 步：写 Handler（`handler/auth_handler.go`）
- 第 11 步：组装 main.go 并启动服务
