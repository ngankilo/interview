package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type JWT struct {
	privateKey []byte
	publicKey  []byte
}

func GenToken(loginId, sub, prvKey string, timeout int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": loginId,
		"sub": sub,
		"exp": timeout,
	})

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(prvKey))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse private key")
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to sign token")
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func NewJWT() (JWT, error) {
	v := viper.New()
	v.AutomaticEnv()
	prvKey := v.GetString("JWT_PRIVATE_KEY")
	pubKey := v.GetString("JWT_PUBLIC_KEY")
	if prvKey == "" || pubKey == "" {
		return JWT{}, fmt.Errorf("missing JWT keys")
	}
	return JWT{
		privateKey: []byte(prvKey),
		publicKey:  []byte(pubKey),
	}, nil
}

func (j JWT) Validate(c *gin.Context) (jwt.MapClaims, string, error) {
	// if !ratelimit.Allow(c.ClientIP()) {
	// 	return nil, "", fmt.Errorf("rate limit exceeded")
	// }
	var auth string
	if len(c.GetHeader("Authorization")) > 0 {
		auth = c.GetHeader("Authorization")
	} else if len(c.Query("token")) > 0 {
		auth = c.Query("token")
	} else if len(c.Query("access_token")) > 0 {
		auth = c.Query("access_token")
	} else {
		return nil, "", fmt.Errorf("no token provided")
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse public key")
		return nil, token, fmt.Errorf("failed to parse public key: %w", err)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		log.Error().Err(err).Str("token", token).Msg("Failed to parse token")
		return nil, token, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		log.Warn().Str("token", token).Msg("Invalid token")
		return nil, token, fmt.Errorf("invalid token")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			log.Warn().Str("token", token).Msg("Token expired")
			return nil, token, fmt.Errorf("token expired")
		}
	}

	return claims, token, nil
}