package utils

import (
	"fmt"
	"net/http"
	"time"
)

// Options for setting cookies.
// type CookieOptions struct {
// 	Name     string
// 	Value    string
// 	Path     string
// 	HttpOnly bool
// }

// Default cookie options.
func DefaultCookieOptions(name, value string) string {

	// Generate cookie options
	return fmt.Sprintf("%s=%s; Path=/; HttpOnly; SameSite=Strict", name, value)
	// return CookieOptions{
	// 	Name:     "auth_token",
	// 	Value:    value,
	// 	Path:     "/",
	// 	HttpOnly: true,
	// 	// Secure: true
	// }
}

// This delete the cookie created with DefaultCookieOptions
func DeleteDefaultCookieOptions(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

type CookieOptionsDelete struct {
	Name string
	Path string
}
