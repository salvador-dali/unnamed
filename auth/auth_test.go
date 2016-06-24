package auth

import (
	"../config"
	"encoding/base64"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	config.Init()
	retCode := m.Run()
	os.Exit(retCode)
}

func TestCreateJWT(t *testing.T) {
	type jwtJson struct {
		Id  int
		Exp int
	}

	currentTime := int(time.Now().Unix())
	for _, v := range []int{6, 2, 1, 5, 8} {
		jwt, err := CreateJWT(v)
		if err != nil {
			t.Errorf("Expect correct jwt. Got %v", err)
		}

		if len(jwt) < 10 {
			t.Errorf("Jwt is too short: %v", jwt)
		}

		parts := strings.Split(jwt, ".")
		if len(parts) != 3 || parts[0] != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" {
			t.Errorf("Jwt does not consist of three parts or first part is not right: %v", jwt)
		}

		data, err := base64.RawStdEncoding.DecodeString(parts[1])
		if err != nil {
			t.Errorf("Second part is not a base64: %v, %v", parts[1], err)
		}

		var claim jwtJson
		if err := json.Unmarshal(data, &claim); err != nil {
			t.Errorf("Second part is not a json: %v, %v", data, err)
		}

		if claim.Id != v || claim.Exp <= currentTime {
			t.Errorf("Claim is not correct %v, %v", claim, v)
		}
	}
}

func TestGenerateSalt(t *testing.T) {
	for i := 0; i < 50; i++ {
		salt, err := GenerateSalt()
		if err != nil || len(salt) != config.Cfg.SaltLen {
			t.Errorf("Salt is not initialized: %v, %v", salt, err)
		}
	}
}

func TestPasswordHash(t *testing.T) {
	table := []struct {
		pwd  string
		salt string
		hash []byte
	}{
		{"password", "sdfasfqwer", []byte{116, 209, 127, 138, 223, 106, 95, 47, 96, 25, 167, 82, 126, 43, 18, 57, 133, 138, 187, 30, 123, 73, 89, 197, 126, 164, 121, 60, 230, 136, 123, 4}},
		{"1234fasdf", "safsadf", []byte{66, 178, 213, 127, 41, 231, 143, 44, 81, 222, 147, 163, 195, 184, 97, 228, 61, 19, 244, 212, 187, 34, 91, 2, 38, 36, 87, 93, 210, 93, 9, 142}},
		{"123asdf(q25L2sa", "asdfasd", []byte{140, 1, 57, 249, 186, 47, 41, 190, 67, 118, 90, 173, 208, 190, 71, 125, 224, 212, 61, 12, 100, 60, 67, 135, 221, 87, 190, 7, 197, 71, 228, 187}},
	}
	for _, v := range table {
		hash, err := PasswordHash(v.pwd, []byte(v.salt))
		if err != nil {
			t.Errorf("Expect to execute correctly: %v", err)
		}
		if !reflect.DeepEqual(hash, v.hash) {
			t.Errorf("Hashes do not match: %v, %v", hash, v.hash)
		}
	}

	table[0].hash[0] = 222
	table[1].hash[6] = 222
	table[2].hash[16] = 222
	for _, v := range table {
		hash, err := PasswordHash(v.pwd, []byte(v.salt))
		if err != nil {
			t.Errorf("Expect to execute correctly: %v", err)
		}
		if reflect.DeepEqual(hash, v.hash) {
			t.Errorf("Hashes should not match: %v, %v", hash, v.hash)
		}
	}
}

func TestValidateJWT(t *testing.T) {
	tableSuccess := []struct {
		id  int
		exp int
		jwt string
	}{
		{1, 1639292831, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyOTI4MzEsImlkIjoxfQ.E3KRJgFfpKHgexw13grm9-neaXrlb7sLjk5Q9XsBeRY"},
		{2, 1639555670, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU2NzAsImlkIjoyfQ.-o8iN6TXLqeyUR8bkJ3WCfDr7527BZ9aHY12qCfOCvE"},
		{3, 1639555719, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU3MTksImlkIjozfQ.Agi-2KpwE-J8B4wUwOz5n-5mcg8P9cUF9qqCwsL2USI"},
		{4, 1639555743, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU3NDMsImlkIjo0fQ.ceGmymRfiO2sv-WV-_7z63FePcdZ36wrQmugHtyI94g"},
		{5, 1639555774, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU3NzQsImlkIjo1fQ.FMx5hJQ-KdV1lCrOhP_UrKXhKvY1DfNeDzsnO2wlGwI"},
		{6, 1639555785, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU3ODUsImlkIjo2fQ.sTQ9HMqrpaP1R6tl7mgrCPjbr52-qWpensYB2IsoaNo"},
		{7, 1639555800, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU4MDAsImlkIjo3fQ.DhJpM75XmrvJet37OhEff0jN3ZBrpoBMbUoSOaCaqTM"},
		{8, 1639555811, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU4MTEsImlkIjo4fQ.vF0Vo_Mpha7FcYhu7BraRfJqsn8hMBednlFGTMumAhk"},
		{9, 1639555822, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU4MjIsImlkIjo5fQ.huTzZZ2ToM1wflgT42oirBRwnyTZbtAJZw6hm6-aJck"},
		{123, 1639292614, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyOTI2MTQsImlkIjoxMjN9.-MW6PUpl8PastumA03R-i0ZuW21aG_3ZQnruSHS6aBo"},
		{1341, 1639292852, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyOTI4NTIsImlkIjoxMzQxfQ.39L7Bl8ZzVg5N2E_X2sl-T43hlsNFloI8X1ZnMDySeA"},
	}
	for _, v := range tableSuccess {
		jwtJson, err := ValidateJWT(v.jwt)
		if err != nil || jwtJson.UserId != v.id || jwtJson.Exp != v.exp {
			t.Errorf("Expect a correct unexpired token. Got %v, %v", jwtJson, err)
		}
	}

	tableSuccess[0].jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyOTI2MTQsImlkIjoiMTIzIn0.-MW6PUpl8PastumA03R-i0ZuW21aG_3ZQnruSHS6aBo"
	tableSuccess[1].jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyOTI4MzIsImlkIjoxfQ.E3KRJgFfpKHgexw13grm9-neaXrlb7sLjk5Q9XsBeRY"
	tableSuccess[2].jwt += "a"
	tableFail := []struct {
		id  int
		exp int
		jwt string
	}{
		tableSuccess[0],
		tableSuccess[1],
		tableSuccess[2],
		{1, 1639292831, "eyJ0eXAiOiJKV1QifQ.eyJleHAiOjE2MzkyOTI4MzEsImlkIjoxfQ.E3KRJgFfpKHgexw13grm9-neaXrlb7sLjk5Q9XsBeRY"},
		{1, 1639292831, "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyOTI4MzEsImlkIjoxfQ.E3KRJgFfpKHgexw13grm9-neaXrlb7sLjk5Q9XsBeRY"},
		{1, 1639292831, "eyJhbGciOiIiLCJ0eXAiOiJKV1QifQ.eyJleHAiOjE2MzkyOTI4MzEsImlkIjoxfQ.E3KRJgFfpKHgexw13grm9-neaXrlb7sLjk5Q9XsBeRY"},
		{0, 0, "asfasdfasdf"},
		{0, 0, "wrong.token"},
		{0, 0, "wrong.token.asdf"},
		{0, 0, ""},
		{0, 0, "..........."},
	}
	for _, v := range tableFail {
		jwtJson, err := ValidateJWT(v.jwt)
		if err == nil || jwtJson.UserId != 0 || jwtJson.Exp != 0 {
			t.Errorf("Expect a wrong JWT. Got %v %v", jwtJson, err)
		}
	}

	expired := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjY0OTY3ODksImlkIjoxfQ.YHEkCkXzifjjk2lTbSsV3gsTtJOtWfE_S8T4tYMjvXs",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjY0OTY4MzUsImlkIjo0MzJ9.jg9Lw77-WrcYy9EzR9DVU0C7LR0B_IIPLdwn9r9qPDM",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjY0OTY4NjYsImlkIjozNDIyMX0.diuB5EdkeLv_k9bjvMonhy5mqwch4emkHx6p_2GPjfY",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjY0OTY4ODcsImlkIjoyfQ.01Eu08mRAi4WJmmFqqXsN15IsHiA7M73zrkydKzass0",
	}
	for _, v := range expired {
		jwtJson, err := ValidateJWT(v)
		if err == nil || jwtJson.Exp != 0 || jwtJson.UserId != 0 || err.Error() != "Token is expired" {
			t.Errorf("Token expired. Got %v %v", jwtJson, err)
		}
	}
}
