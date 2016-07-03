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
		Iat int
		Exp int
	}

	currentTime := int(time.Now().Unix())
	for _, v := range []int{6, 2, 1, 5, 8} {
		jwt, err := CreateJWT(v, true)
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

		if claim.Id != v || claim.Exp <= currentTime || claim.Iat != currentTime {
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
		{"password", "sdfasfqwer", []byte{116, 209, 127, 138, 223, 106, 95, 47, 96, 25, 167, 82, 126, 43, 18, 57, 133, 138, 187, 30, 123, 73, 89, 197, 126, 164, 121, 60, 230, 136, 123, 4, 246, 244, 97, 111, 47, 97, 100, 84, 58, 14, 252, 205, 188, 186, 102, 168, 140, 217, 137, 16, 229, 178, 70, 23, 6, 195, 189, 207, 73, 130, 221, 198, 52, 147, 149, 58, 167, 204, 237, 28, 163, 196, 211, 5, 222, 232, 143, 65, 229, 174, 154, 90, 27, 43, 252, 156, 200, 34, 102, 173, 223, 182, 117, 248, 95, 139, 247, 121, 195, 255, 3, 123, 240, 27, 239, 119, 227, 142, 212, 199, 112, 159, 231, 238, 38, 38, 226, 148, 247, 16, 109, 50, 152, 57, 161, 47}},
		{"1234fasdf", "safsadf", []byte{66, 178, 213, 127, 41, 231, 143, 44, 81, 222, 147, 163, 195, 184, 97, 228, 61, 19, 244, 212, 187, 34, 91, 2, 38, 36, 87, 93, 210, 93, 9, 142, 83, 197, 224, 33, 81, 111, 203, 44, 159, 194, 18, 166, 53, 206, 122, 247, 180, 11, 128, 117, 206, 60, 174, 240, 255, 233, 209, 251, 234, 61, 227, 127, 91, 112, 214, 103, 202, 231, 67, 16, 183, 42, 232, 109, 252, 188, 45, 171, 190, 55, 85, 248, 125, 66, 228, 28, 67, 159, 8, 130, 74, 235, 237, 201, 183, 1, 27, 44, 168, 54, 237, 77, 246, 138, 47, 219, 218, 147, 70, 192, 93, 66, 167, 211, 155, 195, 254, 138, 172, 215, 235, 253, 73, 53, 198, 182}},
		{"123asdf(q25L2sa", "asdfasd", []byte{140, 1, 57, 249, 186, 47, 41, 190, 67, 118, 90, 173, 208, 190, 71, 125, 224, 212, 61, 12, 100, 60, 67, 135, 221, 87, 190, 7, 197, 71, 228, 187, 27, 100, 99, 100, 0, 146, 185, 58, 177, 161, 146, 67, 106, 58, 139, 16, 35, 90, 17, 24, 243, 31, 166, 89, 44, 115, 213, 121, 75, 139, 134, 241, 71, 221, 139, 78, 58, 242, 238, 52, 120, 184, 182, 64, 193, 151, 104, 24, 246, 101, 179, 139, 88, 37, 15, 163, 10, 178, 79, 152, 99, 118, 253, 47, 135, 113, 18, 14, 172, 162, 159, 99, 183, 186, 244, 6, 199, 245, 142, 113, 209, 58, 192, 51, 210, 28, 96, 7, 108, 54, 211, 22, 54, 90, 169, 130}},
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
		iat int
		exp int
		jwt string
	}{
		{1, 1466833211, 1639633211, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjF9.LYey3jgBd70QYjygbZvoPqXGJHj90nZ8VUm2yeVlVVo"},
		{2, 1466833211, 1639633211, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjJ9.dbBN08ZNdGhKbPhFRSccRWvMgSxSTjlM3wC7K2oz3_M"},
		{3, 1466833211, 1639633211, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjN9.WF7GGKA2XB3Th5lztqseW1fixf9XApTYpwDhcvq_sDw"},
		{4, 1466833211, 1639633211, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjR9.cI2Ie6KDVQhWk1VRuP_UzE1HpKFfyT0jgTe9J2g7pJA"},
		{5, 1466833211, 1639633211, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjV9.bmmgOyeN700onUcVfJcFT4dn5XyNY7rdUfpYDhlfdOc"},
		{6, 1466833211, 1639633211, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjZ9.vqP4oem2PeQpzBBC2enSXYrKg2xDcPa8iXcJToSmWHs"},
		{7, 1466833211, 1639633211, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjd9.KlfEaHwqWLMGVA9MUIu_z8oSNaXbioJ6_mgftlbWpeI"},
		{8, 1466833211, 1639633211, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjh9.0u293Hl2-cJawLI1JlEcE1fYBB6yrkMvKUiGHy61-2A"},
		{9, 1466833211, 1639633211, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjl9.JS9Xc135ndkunTa2oKess5KCX4WVCcvAkI7bVsV4YVo"},
	}
	for _, v := range tableSuccess {
		jwtJson, err := ValidateJWT(v.jwt)
		if err != nil || jwtJson.UserId != v.id || jwtJson.Exp != v.exp || jwtJson.Iat != v.iat {
			t.Errorf("Expect a correct unexpired token. Got %v, %v", jwtJson, err)
		}
	}

	tableSuccess[0].jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyOTI2MTQsImlkIjoiMTIzIn0.-MW6PUpl8PastumA03R-i0ZuW21aG_3ZQnruSHS6aBo"
	tableSuccess[1].jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyOTI4MzIsImlkIjoxfQ.E3KRJgFfpKHgexw13grm9-neaXrlb7sLjk5Q9XsBeRY"
	tableSuccess[2].jwt += "a"
	tableFail := []struct {
		id  int
		iat int
		exp int
		jwt string
	}{
		tableSuccess[0],
		tableSuccess[1],
		tableSuccess[2],
		{1, 0, 1639292831, "eyJ0eXAiOiJKV1QifQ.eyJleHAiOjE2MzkyOTI4MzEsImlkIjoxfQ.E3KRJgFfpKHgexw13grm9-neaXrlb7sLjk5Q9XsBeRY"},
		{1, 0, 1639292831, "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyOTI4MzEsImlkIjoxfQ.E3KRJgFfpKHgexw13grm9-neaXrlb7sLjk5Q9XsBeRY"},
		{1, 0, 1639292831, "eyJhbGciOiIiLCJ0eXAiOiJKV1QifQ.eyJleHAiOjE2MzkyOTI4MzEsImlkIjoxfQ.E3KRJgFfpKHgexw13grm9-neaXrlb7sLjk5Q9XsBeRY"},
		{0, 0, 0, "asfasdfasdf"},
		{0, 0, 0, "wrong.token"},
		{0, 0, 0, "wrong.token.asdf"},
		{0, 0, 0, ""},
		{0, 0, 0, "..........."},
		{0, 0, 0, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyOTI4MzEsImlkIjoxfQ.E3KRJgFfpKHgexw13grm9-neaXrlb7sLjk5Q9XsBeRY"},
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

func TestExtendJWT(t *testing.T) {
	tableSuccess := []struct {
		id  int
		jwt string
	}{
		{1, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjF9.LYey3jgBd70QYjygbZvoPqXGJHj90nZ8VUm2yeVlVVo"},
		{2, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjJ9.dbBN08ZNdGhKbPhFRSccRWvMgSxSTjlM3wC7K2oz3_M"},
		{3, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjN9.WF7GGKA2XB3Th5lztqseW1fixf9XApTYpwDhcvq_sDw"},
		{4, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjR9.cI2Ie6KDVQhWk1VRuP_UzE1HpKFfyT0jgTe9J2g7pJA"},
	}
	for _, v := range tableSuccess {
		currentTime := int(time.Now().Unix())
		jwtJson, err := ExtendJWT(v.jwt)
		if err != nil || len(jwtJson) < 100 {
			t.Errorf("Expect a correct unexpired token. Got %v, %v", jwtJson, err)
		}

		data, err := ValidateJWT(jwtJson)
		if data.UserId != v.id || data.Iat != currentTime {
			t.Errorf("Expect %v, %v. Got %v, %v", v.id, currentTime, data.UserId, data.Iat)
		}
	}

	tableFail := []struct {
		id  int
		jwt string
	}{
		{1, tableSuccess[0].jwt + "a"},
		{2, tableSuccess[0].jwt + "a"},
		{3, tableSuccess[0].jwt + "a"},
		{4, tableSuccess[0].jwt + "a"},
	}
	for _, v := range tableFail {
		jwtJson, err := ExtendJWT(v.jwt)
		if err == nil || jwtJson != "" {
			t.Errorf("Expect failure. Got %v, %v", err, jwtJson)
		}
	}
}
