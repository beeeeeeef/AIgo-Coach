package users

import (
	"aigo-coach/backend/model"
	"aigo-coach/backend/utils"
	"errors"
	"time"
)

var(
	ErrUserNotFound = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("密码错误")
)

// LoginByPassword 用户登录
// 1.根据邮箱获取用户记录，判断用户是否存在
// 2.比较密码哈希值，判断密码是否正确

func (s *UserService) LoginByPassword(email, password string) (*model.User, error) {
	// 1.根据邮箱获取用户记录，判断用户是否存在
	user, err := s.UserRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	// 2.比较密码哈希值，判断密码是否正确
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, ErrInvalidPassword
	}
	return user, nil
}

// SendLoginCode 发送登录验证码
// 1.根据邮箱获取用户记录，判断用户是否存在
// 2.判断 60 秒冷却期
// 3.生成验证码，10 分钟后过期，保存或刷新到数据库
func (s *UserService) SendLoginCode(email string) error {
	// 1.判断用户是否存在
	user, err := s.UserRepo.GetByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 2.判断 60 秒冷却期
	latest, err := s.VerifyCodeRepo.GetLatestByEmailAndType(email, "login")
	if err != nil {
		return err
	}
	if latest != nil && time.Since(latest.CreatedAt) < 60*time.Second {
		return ErrCodeSendTooFrequent
	}

	// 3.生成验证码，10 分钟后过期，保存或刷新到数据库
	code := utils.GenerateCode()
	expiredAt := time.Now().Add(10 * time.Minute)

	if latest == nil {
		// 首次发送：创建新记录
		verifyCode := &model.VerifyCode{
			Email:      email,
			Code:       code,
			Type:       "login",
			ExpiredAt:  expiredAt,
			ApplyTimes: 1,
		}
		return s.VerifyCodeRepo.Create(verifyCode)
	}

	// 非首次发送：刷新已有记录（更新 code、过期时间，apply_times+1，重置 used）
	return s.VerifyCodeRepo.RefreshCode(latest.Id, code, expiredAt)
}

// LoginByCode 验证码登录
// 1.根据邮箱获取用户记录，判断用户是否存在
// 2.校验验证码：存在、未过期、未使用、匹配
// 3.标记验证码已使用
func (s *UserService) LoginByCode(email, code string) (*model.User, error) {
	// 1.根据邮箱获取用户记录，判断用户是否存在
	user, err := s.UserRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 2.校验验证码
	latest, err := s.VerifyCodeRepo.GetLatestByEmailAndType(email, "login")
	if err != nil {
		return nil, err
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

	// 3.标记验证码已使用
	_ = s.VerifyCodeRepo.MarkUsed(latest.Id)

	return user, nil
}
