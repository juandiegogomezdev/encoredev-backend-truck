package business

import (
	"math/rand"

	"encore.app/auth/store"
	"encore.app/pkg/myjwt"
	"encore.app/pkg/resendmailer"
	"golang.org/x/crypto/bcrypt"
)

type BusinessAuth struct {
	store     *store.AuthStore
	tokenizer myjwt.JWTTokenizer
	mailer    resendmailer.ResendMailer
}

func NewAuthBusiness(store *store.AuthStore, tokenizer myjwt.JWTTokenizer, mailer resendmailer.ResendMailer) *BusinessAuth {
	return &BusinessAuth{store: store, tokenizer: tokenizer, mailer: mailer}
}

// Hash the password using bcrypt.
func GenerateHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Validate if the provided password matches the stored hash.
func (b *BusinessAuth) validatePassword(storedHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	return err == nil
}

// Generate a code with n digits
func (b *BusinessAuth) generateCodeLogin(n int) string {
	numbers := "0123456789"
	code := make([]byte, n)
	for i := range code {
		code[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(code)
}
