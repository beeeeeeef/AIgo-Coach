package users

import (		
	"aigo-coach/backend/model"
	"aigo-coach/backend/utils"
	"errors"
	"time"
)

var (
	ErrUserAlreadyExists = errors.New("邮箱已注册")
	ErrCodeNotFound = errors.New("验证码不存在")
	ErrCodeExpired = errors.New("验证码已过期")
	ErrCodeUsed = errors.New("验证码已被使用")
	ErrInvalidCode = errors.New("验证码错误")
	ErrCodeSendTooFrequent = errors.New("验证码发送过于频繁")
    ErrCodeSendLimitExceeded = errors.New("验证码发送次数过多")
)

// SendRegisterCode 发送注册验证码
// 1.判断邮箱是否已经注册
// 2.判断 60 秒内是否发过验证码
// 3.检查是否到达发送上限
// 4.生成验证码，10 分钟后过期，保存到数据库
func (s *UserService) SendRegisterCode(email string) error {
	// 1.判断邮箱是否已经注册
	user,err := s.UserRepo.GetByEmail(email)
	if err != nil {
		return err
	}
	if user != nil {
		return ErrUserAlreadyExists
	}
	// 2.判断 60 秒内是否发过验证码
	latest,err := s.VerifyCodeRepo.GetLatestByEmailAndType(email, "register")
	if err != nil {
		return err
	}
	if latest != nil && time.Since(latest.CreatedAt) < 60*time.Second {
		return ErrCodeSendTooFrequent
	}
	// 3.检查是否到达发送上限
	if latest != nil && latest.ApplyTimes >= 5 {
		return ErrCodeSendLimitExceeded
	}
	// 4.生成验证码，10 分钟后过期，保存或刷新到数据库
	code := utils.GenerateCode()
	expiredAt := time.Now().Add(10 * time.Minute)

	if latest == nil {
		// 首次发送：创建新记录
		verifyCode := &model.VerifyCode{
			Email:      email,
			Code:       code,
			Type:       "register",
			ExpiredAt:  expiredAt,
			ApplyTimes: 1,
		}
		return s.VerifyCodeRepo.Create(verifyCode)
	}

	// 非首次发送：刷新已有记录（更新 code、过期时间，apply_times+1，重置 used）
	return s.VerifyCodeRepo.RefreshCode(latest.Id, code, expiredAt,)
}

// Register 处理用户注册
// 1.验证验证码：存在、未过期、未使用、匹配
// 2.再次检查邮箱是否已经注册（防止并发注册）
// 3.bcrypt 加密密码
// 4.保存用户到数据库
// 5.标记验证码已使用
func (s *UserService) Register(email,code,password string) (*model.User,error){
	// 校验验证码
	latest,err :=s.VerifyCodeRepo.GetLatestByEmailAndType(email,"register")
	if err != nil {
		return nil,err
	}
	if latest == nil {
		return nil, ErrCodeNotFound
	}
	if time.Now().After(latest.ExpiredAt) {
		return nil, ErrCodeExpired
	}
	if latest.Used {
		return nil, ErrCodeUsed
	}
	if latest.Code != code {
		_ = s.VerifyCodeRepo.IncrementApplyTimes(latest.Id)
		return nil, ErrInvalidCode
	}
	// 再次检查邮箱是否已经注册（防止并发注册）
	user,err := s.UserRepo.GetByEmail(email)
	if err != nil {
		return nil,err
	}
	if user != nil {
		return nil, ErrUserAlreadyExists
	}

	// bcrypt 加密密码
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	// 保存用户到数据库
	newUser := &model.User{
		Email: email,
		PasswordHash: hashedPassword,
	}
	err = s.UserRepo.Create(newUser)
	if err != nil {
		return nil, err
	}
	// 标记验证码已使用（失败不阻断注册）
	_ = s.VerifyCodeRepo.MarkUsed(latest.Id)

	return newUser, nil
}
