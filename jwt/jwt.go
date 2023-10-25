package jwt

// 导入必要的库
import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	UserID   int    // 用户ID
	Username string // 用户名
	Email    string // 邮箱
	// 其他字段...
}

// 生成JWT令牌的函数
func GenerateToken(user User) (string, error) {
	// 创建一个新的令牌
	token := jwt.New(jwt.SigningMethodHS256)

	// 创建声明
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.UserID
	claims["username"] = user.Username
	claims["email"] = user.Email
	// 其他字段...

	// 签署令牌
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT 解析并验证 JWT 令牌
func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	// 解析 JWT 令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 在这里提供用于验证令牌签名的密钥
		return []byte("your-secret-key"), nil
	})

	if err != nil {
		fmt.Println("Failed to parse JWT:", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 从令牌中提取数据
		return claims, nil
	} else {
		fmt.Println("Token is not valid")
		return nil, fmt.Errorf("Token is not valid")
	}
}
