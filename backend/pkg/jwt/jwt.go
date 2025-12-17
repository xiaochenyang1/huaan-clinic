package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"huaan-medical/pkg/config"
)

// TokenType Token类型
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims 自定义JWT声明
type Claims struct {
	UserID    int64     `json:"user_id"`
	OpenID    string    `json:"open_id,omitempty"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// AdminClaims 管理员JWT声明
type AdminClaims struct {
	AdminID  int64  `json:"admin_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// TokenPair Token对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // 过期时间（秒）
}

var jwtConfig *config.JWTConfig

// Init 初始化JWT配置
func Init(cfg *config.JWTConfig) {
	jwtConfig = cfg
}

// GenerateTokenPair 生成Token对
func GenerateTokenPair(userID int64, openID string) (*TokenPair, error) {
	now := time.Now()

	// 生成Access Token
	accessClaims := Claims{
		UserID:    userID,
		OpenID:    openID,
		TokenType: AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtConfig.AccessTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "huaan-medical",
			Subject:   "user",
		},
	}
	accessToken, err := generateToken(accessClaims)
	if err != nil {
		return nil, err
	}

	// 生成Refresh Token
	refreshClaims := Claims{
		UserID:    userID,
		OpenID:    openID,
		TokenType: RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtConfig.RefreshTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "huaan-medical",
			Subject:   "user",
		},
	}
	refreshToken, err := generateToken(refreshClaims)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(jwtConfig.AccessTokenExpire.Seconds()),
	}, nil
}

// GenerateAdminToken 生成管理员Token
func GenerateAdminToken(adminID int64, username, role string) (string, error) {
	now := time.Now()

	claims := AdminClaims{
		AdminID:  adminID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtConfig.AccessTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "huaan-medical",
			Subject:   "admin",
		},
	}

	return generateToken(claims)
}

// ParseToken 解析用户Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig.Secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token已过期")
		}
		return nil, errors.New("token无效")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token解析失败")
}

// ParseAdminToken 解析管理员Token
func ParseAdminToken(tokenString string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig.Secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token已过期")
		}
		return nil, errors.New("token无效")
	}

	if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token解析失败")
}

// RefreshAccessToken 刷新Access Token
func RefreshAccessToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := ParseToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != RefreshToken {
		return nil, errors.New("无效的refresh token")
	}

	return GenerateTokenPair(claims.UserID, claims.OpenID)
}

// generateToken 生成Token
func generateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.Secret))
}
