package auth

import (
	"../../unnamed/config"
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

func ParseJWT(jwtToken string) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return config.Cfg.Secret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
	}
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
