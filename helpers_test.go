package main

import (
	"fmt"
	"testing"
)

// Test randomString function
// Also tests StringWithCharset function
func TestRandomString(t *testing.T) {
	randMsg := randomString(20)
	randMsgTwo := randomString(20)
	if len(randMsg) < 20 {
		t.Fatalf(`randomString(20) = %q length fail`, randMsg)
	}
	if randMsg == randMsgTwo {
		t.Fatalf(`randomString == randomString %q %q`, randMsg, randMsgTwo)
	}
}

// Test unique function
func TestUnique(t *testing.T) {
	listWithDupes := []string{"abc", "cde", "efg", "efg", "abc", "cde"}
	listWithOutDupes := []string{"abc", "cde", "efg"}
	processedList := unique(listWithDupes)

	if len(listWithOutDupes) != len(processedList) {
		t.Fatalf(`processed list %q != randomString %q`, processedList, listWithOutDupes)
	}
	for i, v := range processedList {
		if v != listWithOutDupes[i] {
			t.Fatalf(`slices dont match %q %q`, v, listWithOutDupes[i])
		}
	}
}

// Test generateNewHash function
func TestGenerateHash(t *testing.T) {
	hashOne := generateNewHash(8)
	hashTwo := generateNewHash(8)
	if len(hashOne) > 8 {
		t.Fatalf(`generateNewHash(8) = %q length fail`, hashOne)
	}
	if hashOne == hashTwo {
		t.Fatalf(`generateNewHash == generateNewHash %q %q`, hashOne, hashTwo)
	}
	var hashList []string
	// Generate 10000 unique hashes.
	// Note, this has been tested at a couple million without issue, just takes a long time.
	for i := 1; i < 10000; i++ {
		hashList = append(hashList, generateNewHash(8))
	}
	// Remove any duplicate hashes. There should be none.
	processedList := unique(hashList)
	if len(hashList) != len(processedList) {
		t.Fatalf(`processed list %q != randomString %q`, processedList, hashList)
	}
	for i, v := range processedList {
		if v != hashList[i] {
			t.Fatalf(`slices dont match %q %q`, v, hashList[i])
		}
	}
}

// Test generateJwtSecret function
func TestGenerateJwtSecret(t *testing.T) {
	defaultValue := []byte("my_secret_key")
	generateJwtSecret()
	if string(defaultValue) == string(jwtSecretKey) {
		t.Fatalf(`jwtSecretKey not updated: %q`, jwtSecretKey)
	}
	if len(jwtSecretKey) < 45 {
		t.Fatalf(`jwtSecretKey length fail: %q`, jwtSecretKey)
	}
}

// Tests hashAndSalt function.
func TestHashAndSalt(t *testing.T) {
	randPass := randomString(20)
	hashedPassword := hashAndSalt([]byte(randPass))
	if randPass == hashedPassword {
		t.Fatalf(`password hash fail: %q | randPass: %q`, hashedPassword, randPass)
	}
}

// Tests comparePasswords function.
// Also tests hashAndSalt function.
func TestComparePasswords(t *testing.T) {
	fmt.Println("This test intentionally causes an error message (hashedPassword is not the hash of the given password)")
	testPassOne := "SomeNotRealPassword"
	testPassTwo := "SomeNotRealPassword"
	testPassThree := "ADifferentPassword"

	hashedPasswordOne := hashAndSalt([]byte(testPassOne))
	hashedPasswordTwo := hashAndSalt([]byte(testPassTwo))
	hashedPasswordThree := hashAndSalt([]byte(testPassThree))
	if !comparePasswords(hashedPasswordOne, []byte(testPassOne)) {
		t.Fatalf(`password compare fail: %q | randPass: %q`, hashedPasswordOne, testPassOne)
	}
	if !comparePasswords(hashedPasswordTwo, []byte(testPassTwo)) {
		t.Fatalf(`password compare fail: %q | randPass: %q`, hashedPasswordTwo, testPassTwo)
	}
	if !comparePasswords(hashedPasswordTwo, []byte(testPassOne)) {
		t.Fatalf(`password compare fail: %q | randPass: %q`, hashedPasswordTwo, testPassOne)
	}
	if !comparePasswords(hashedPasswordThree, []byte(testPassThree)) {
		t.Fatalf(`password compare fail: %q | randPass: %q`, hashedPasswordThree, testPassThree)
	}
	if comparePasswords(hashedPasswordThree, []byte(testPassOne)) {
		t.Fatalf(`password compare fail: %q | randPass: %q`, hashedPasswordTwo, testPassOne)
	}

}
