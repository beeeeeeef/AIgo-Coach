package repository

import(
	"aigo-coach/backend/model"
	"database/sql"
	"time"
)

type VerifyCodeRepositorySQL struct {
	DB *sql.DB
}
func NewVerifyCodeRepositorySQL(db *sql.DB) *VerifyCodeRepositorySQL {
	return &VerifyCodeRepositorySQL{DB: db}
}
// Create 创建验证码记录
func(r*VerifyCodeRepositorySQL) Create(code *model.VerifyCode) error{
query := `
		INSERT INTO verify_codes (email, code, type, expired_at, used, apply_times)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.DB.Exec(query,
		code.Email,
		code.Code,
		code.Type,
		code.ExpiredAt,
		code.Used,
		code.ApplyTimes,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	code.Id = id
	return nil
}

// GetLatestByEmailAndType 获取指定邮箱和类型的最新验证码
func (r *VerifyCodeRepositorySQL) GetLatestByEmailAndType(email string, codeType string) (*model.VerifyCode, error) {
	query := `
		SELECT 
			id,
			email,
			code,
			type,
			created_at,
			expired_at,
			apply_times,
			used
		FROM verify_codes
		WHERE email = ? AND type = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	var code model.VerifyCode
	err := r.DB.QueryRow(query, email, codeType).Scan(
		&code.Id,
		&code.Email,
		&code.Code,
		&code.Type,
		&code.CreatedAt,
		&code.ExpiredAt,
		&code.ApplyTimes,
		&code.Used,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &code, nil
}
// RefreshCode 刷新验证码（更新 code、过期时间，apply_times+1，重置 used）
func (r *VerifyCodeRepositorySQL) RefreshCode(id int64, newCode string, expiredAt time.Time) error {
	query := `
		UPDATE verify_codes
		SET code = ?, expired_at = ?, apply_times = apply_times + 1, used = false
		WHERE id = ?
	`
	_, err := r.DB.Exec(query, newCode, expiredAt, id)
	return err
}

//IncrementApplyTimes 增加验证码的申请次数
func (r *VerifyCodeRepositorySQL) IncrementApplyTimes(id int64) error {
	query := `
		UPDATE verify_codes
		SET apply_times = apply_times + 1
		WHERE id = ?
	`

	_, err := r.DB.Exec(query, id)
	return err
}
// MarkUsed 将验证码标记为已使用
func (r *VerifyCodeRepositorySQL) MarkUsed(id int64) error {
	query := `
		UPDATE verify_codes
		SET used = true
		WHERE id = ?
	`
	_, err := r.DB.Exec(query, id)
	return err
}
