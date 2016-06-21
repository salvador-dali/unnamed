package auth

import (
	"../../unnamed/config"
	"../../unnamed/structs"
	"crypto/rand"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
	"time"
)

// TODO read about these parameters
const (
	hashN      = 32768
	hashR      = 8
	hashP      = 1
	hashKeyLen = 32
)

func CreateJWT(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(time.Hour * 24 * time.Duration(config.Cfg.ExpDays)).Unix(),
	})

	return token.SignedString(config.Cfg.Secret)
}

func ValidateJWT(jwtToken string) (structs.JwtToken, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return config.Cfg.Secret, nil
	})

	var jwtJson structs.JwtToken
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		jwtJson.UserId = int(claims["id"].(float64))
		jwtJson.Exp = int(claims["exp"].(float64))
		return jwtJson, nil
	}

	return structs.JwtToken{}, err
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, config.Cfg.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		// error means that the system's system's random number generator does not have randomness
		return nil, err
	}
	return salt, nil
}

func PasswordHash(password string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(password), salt, hashN, hashR, hashP, hashKeyLen)
}
