package main

import (
	"github.com/speps/go-hashids/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int) string {
	return StringWithCharset(length, charset)
}

func generateJwtSecret() {
	jwtSecretKey = []byte(randomString(45))
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func generateNewHash(maxLength int) string {
	hd := hashids.NewData()
	hd.Salt = randomString(35)
	hd.MinLength = 5
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{rand.Intn(10000), rand.Intn(10000), rand.Intn(10000), rand.Intn(10000), rand.Intn(10000)})
	// Limit length to some defined max.
	i := 0
	for j := range e {
		if i == maxLength {
			return e[:j]
		}
		i++
	}
	return e
}

func unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}
