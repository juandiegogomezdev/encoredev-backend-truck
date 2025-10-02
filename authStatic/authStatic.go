package authstatic

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"encore.app/pkg/myjwt"
	"encore.app/pkg/utils"
)

var secrets struct {
	JWT_SECRET_KEY string
}

//encore:service
type AuthStaticService struct {
	tokenizer     myjwt.JWTTokenizer
	staticHandler http.Handler
}

func initAuthStaticService() (*AuthStaticService, error) {
	tokenizer := myjwt.NewJWTTokenizer(secrets.JWT_SECRET_KEY)
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.FS(assets)))
	return &AuthStaticService{tokenizer: tokenizer, staticHandler: staticHandler}, nil
}

//go:embed dist/*
var dist embed.FS

var assets, _ = fs.Sub(dist, "dist")

//encore:api public raw path=/static/*path tag:static
func (s *AuthStaticService) ServeStaticFiles(w http.ResponseWriter, req *http.Request) {
	//
	// Verificacion diferentes al archivo endpoint de abajo
	url := req.URL.Path
	if strings.HasSuffix(url, ".js") || strings.HasSuffix(url, ".css") {
		s.staticHandler.ServeHTTP(w, req)
		return
	}

	// Get token from cookie
	cookie, err := req.Cookie("auth_token")
	token := ""
	if err == nil {
		token = cookie.Value
	}

	// Identify if token is valid and type of access
	tokenType := myjwt.TokenTypeUnknown
	if token != "" {
		claims, tokenStatus := s.tokenizer.ParseBaseClaims(token)
		switch tokenStatus {
		case myjwt.TokenStatusExpired:
			// Refresh token if is expired and is in the database
			// TODO: Implement refresh token
		case myjwt.TokenStatusInvalid:
			// Delete the token if is invalid
			expiredCookie := utils.DeleteDefaultCookieOptions2("auth_token")
			http.SetCookie(w, expiredCookie)
		case myjwt.TokenStatusValid:
			// Set the token type
			tokenType = claims.TokenType
		}
	}

	switch tokenType {
	case myjwt.TokenTypeUnknown:
		if strings.HasSuffix(url, "org-select/") || strings.HasSuffix(url, "app/") {
			http.Redirect(w, req, "http://localhost:4000/static/login", http.StatusSeeOther)
			return
		}
	case myjwt.TokenTypeOrgSelect:
		fmt.Println("url:", url)
		if !strings.HasSuffix(url, "org-select/") {
			http.Redirect(w, req,
				"http://localhost:4000/static/org-select/", http.StatusFound)
			return
		}
	case myjwt.TokenTypeMembership:
		if !strings.HasSuffix(url, "app/") {
			http.Redirect(w, req,
				"http://localhost:4000/static/app/", http.StatusFound)

			return
		}
	}

	s.staticHandler.ServeHTTP(w, req)
}
