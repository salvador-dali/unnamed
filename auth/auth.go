package auth

import (
	"../config"
	"../misc"
	"crypto/rand"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
	"strings"
	"time"
)

// TODO read about these parameters
const (
	hashN      = 32768
	hashR      = 8
	hashP      = 1
	hashKeyLen = 32
)

func CreateJWT(userId int, verified bool) (string, error) {
	claims := jwt.MapClaims{
		"id":  userId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * time.Duration(config.Cfg.ExpDays)).Unix(),
	}

	if !verified {
		claims["unverified"] = 1
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.Cfg.Secret)
}

func ValidateJWT(jwtToken string) (misc.JwtToken, error) {
	if strings.Count(jwtToken, ".") != 2 {
		return misc.JwtToken{}, errors.New("Not a JWT token")
	}

	if len(jwtToken) <= 90 {
		return misc.JwtToken{}, errors.New("JWT token is too short")
	}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return config.Cfg.Secret, nil
	})

	var jwtJson misc.JwtToken
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if len(claims) < 3 || len(claims) > 4 {
			return misc.JwtToken{}, errors.New("Token with wrong number of claims")
		}
		_, ok1 := claims["id"]
		_, ok2 := claims["iat"]
		_, ok3 := claims["exp"]
		if !(ok1 && ok2 && ok3) {
			return misc.JwtToken{}, errors.New("One of the required claims is missing")
		}

		// required claims
		jwtJson.UserId = int(claims["id"].(float64))
		jwtJson.Iat = int(claims["iat"].(float64))
		jwtJson.Exp = int(claims["exp"].(float64))

		// optional claims
		if _, ok := claims["unverified"]; ok {
			jwtJson.Verified = false
		} else {
			jwtJson.Verified = true
		}

		return jwtJson, nil
	}

	return misc.JwtToken{}, err
}

func ExtendJWT(jwtToken string) (string, error) {
	token, err := ValidateJWT(jwtToken)
	if err != nil {
		return "", err
	}

	return CreateJWT(token.UserId, true)
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
