package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserClaims struct {
	jwt.StandardClaims
	SessionID int64
}

func (u *UserClaims) Valid() error {
	if !u.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("token has expired")
	}

	if u.SessionID == 0 {
		return fmt.Errorf("invalid session ID")
	}

	return nil
}

var key []byte

func main() {
	for i := 1; i <= 64; i++ {
		key = append(key, byte(i))
	}
	pass := "123456789"

	hashedPass, err := hashPassword(pass)
	if err != nil {
		panic(err)
	}

	if err := comparePassword(pass, hashedPass); err != nil {
		log.Println("Not logged in")
		return
	}
	log.Println("Logged in!!")

	// HMAC
	sign, err := signMessage([]byte("Hello"))
	if err != nil {
		fmt.Errorf("signMessage failed to sign: %w", err)
	}
	bool, err := checkSig([]byte("Hello"), sign)
	if err != nil {
		fmt.Errorf("check sign failed: %w", err)
		return
	}
	fmt.Println("Sign verified", bool)
}

func hashPassword(password string) ([]byte, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error while generating bcrypt hash from password: %w", err)
	}
	return bs, nil
}

func comparePassword(password string, hashedPass []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPass, []byte(password))
}

func signMessage(msg []byte) ([]byte, error) {
	hash := hmac.New(sha512.New, key)
	_, err := hash.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("error in signMessage while hashing message: %w", err)
	}
	signature := hash.Sum(nil)
	return signature, nil
}

func checkSig(msg, sig []byte) (bool, error) {
	newSig, err := signMessage(msg)
	if err != nil {
		return false, fmt.Errorf("error in checkSig while getting signature of message: %w", err)
	}
	same := hmac.Equal(newSig, sig)
	return same, nil
}
