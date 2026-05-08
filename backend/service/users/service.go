package users
import(
    "aigo-coach/backend/repository"
)

// UserService 处理用户注册和登录的所有业务逻辑
type UserService struct {
	UserRepo repository.UserRepository
	VerifyCodeRepo repository.VerifyCodeRepository
}

// NewUserService 创建一个新的 UserService 实例， 需要注入 UserRepository 和 VerifyCodeRepository
func NewUserService(userRepo repository.UserRepository, verifyCodeRepo repository.VerifyCodeRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
		VerifyCodeRepo: verifyCodeRepo,
	}
}

