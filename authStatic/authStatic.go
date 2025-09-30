package authstatic

import (
	"embed"
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

//encore:api public raw path=/static/org-select tag:static
func (s *AuthStaticService) ServeOrgSelect(w http.ResponseWriter, req *http.Request) {
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
		claims, err := s.tokenizer.ParseBaseClaims(token)
		// Delete the token if is invalid
		if err != nil {
			utils.DeleteDefaultCookieOptions(w, "auth_token")
		} else {
			tokenType = claims.TokenType
		}

	}

	switch tokenType {

	}

	if err != nil {
		http.Redirect(w, req, "/static/login", http.StatusSeeOther)
		return
	}

	s.staticHandler.ServeHTTP(w, req)
}

//encore:api public raw path=/static/*path tag:static
func (s *AuthStaticService) ServeStatic(w http.ResponseWriter, req *http.Request) {
	//
	s.staticHandler.ServeHTTP(w, req)
}

//encore:api public path=/static/*path
