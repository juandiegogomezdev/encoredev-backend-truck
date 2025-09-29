package utils

// Options for setting cookies.
type CookieOptions struct {
	Name     string
	Value    string
	Path     string
	HttpOnly bool
}

// Default cookie options.
func DefaultCookieOptions(name, value string) CookieOptions {
	return CookieOptions{
		Name:     "auth_token",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		// Secure: true
	}
}
