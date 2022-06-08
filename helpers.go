package main

import (
	"github.com/gin-gonic/gin"
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

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}
	c.AbortWithStatusJSON(code, resp)
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
	hd.Salt = "new id hash salt data"
	hd.MinLength = 5
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{rand.Intn(1000), rand.Intn(1000), rand.Intn(1000), rand.Intn(1000)})
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
