package repository

import "aigo-coach/backend/model"

type UserRepository interface {
	GetByEmail(email string) (*model.User, error)
	GetByID(id int64) (*model.User, error)
	Create(user *model.User) error
}

type VerifyCodeRepository interface {
	Create(code *model.VerifyCode) error
	GetLatestByEmailAndType(email string, codeType string) (*model.VerifyCode, error)
	IncrementApplyTimes(id int64) error
	MarkUsed(id int64) error
}
