package repository

import(
	"aigo-coach/backend/model"
	"database/sql"
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
		INSERT INTO verify_codes (email, code, type, expire_time, used, apply_time)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.DB.Exec(query,
		code.Email,
		code.Code,
		code.Type,
		code.ExpiresAt,
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
			create_time,
			expire_time,
			apply_time,
			used
		FROM verify_codes
		WHERE email = ? AND type = ?
		ORDER BY create_time DESC
		LIMIT 1
	`

	var code model.VerifyCode
	err := r.DB.QueryRow(query, email, codeType).Scan(
		&code.Id,
		&code.Email,
		&code.Code,
		&code.Type,
		&code.CreatedAt,
		&code.ExpiresAt,
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
//IncrementApplyTimes 增加验证码的申请次数
func (r *VerifyCodeRepositorySQL) IncrementApplyTimes(id int64) error {
	query := `
		UPDATE verify_codes
		SET apply_time = apply_time + 1
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
