package utils

import(
    "crypto/hmac"      // HMAC 签名算法
	"crypto/sha256"     // SHA-256 哈希
	"encoding/base64"   // Base64 编解码
	"encoding/json"     // JSON 序列化
	"errors"
	"fmt"
	"strings"           // 字符串操作（切割 token）
	"time"              // 过期时间处理
)

var (
	ErrTokenExpired = errors.New("token 已过期")
	ErrTokenInvalid = errors.New("token 无效")
	ErrTokenMalformed = errors.New("token 格式错误")
)

//headerClaims JWT 的头部(header)
// {"alg":"HS256","typ":"JWT"} 告诉接收方：我是用 HMAC-SHA256 签名的 JWT
type headerClaims struct {
	Alg string `json:"alg"` // 签名算法
	Typ string `json:"typ"` // 令牌类型
}

//payloadClaims JWT 的负载(payload)
// {"sub":42,"exp":1748700000} sub 是 userID，exp 是过期时间的 Unix 时间戳
type payloadClaims struct {
	Sub int64 `json:"sub"` // 用户 userID
	Exp int64 `json:"exp"` // 过期时间（Unix 时间戳） expiration
}

//base64UrlEncode 对数据进行 Base64 URL 安全编码
func base64UrlEncode(data []byte) string {
	encoded := base64.URLEncoding.EncodeToString(data)
	return encoded
}

//base64UrlDecode 对 Base64 URL 安全编码的数据进行解码
func base64UrlDecode(encoded string) ([]byte, error) {
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

// GenerateToken 生成 JWT token
// 参数
// userID: 用户 ID,存在 JWT 的 sub 字段中
// secretKey: 用于签名的密钥
// ttl: 过期时间，单位为时间段（如 time.Hour）
//
// 返回值
//   一段完整的 JWT 字符串，格式为 header.payload.signature
//   每一段都是 Base64URL 编码的
//
// 流程（自己想）
