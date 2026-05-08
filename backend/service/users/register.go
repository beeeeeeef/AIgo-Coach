package users

import (
    "aigo-coach/backend/repository"
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
	// 4.生成验证码，10 分钟后过期，保存到数据库
	
}