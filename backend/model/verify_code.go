package model
import "time"
type VerifyCode struct{
    Email string `db:"email"`
	Code string `db:"code"`
	Type string `db:"type"`
    ApplyTimes int `db:"apply_times"`
	ExpiresAt time.Time `db:"expires_at"`
    Used bool `db:"used"`
	CreatedAt time.Time `db:"created_at"`
	Id int64 `db:"id"`
	
}

