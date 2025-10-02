package authstatic

//encore:middleware target=tag:static
// func (s *AuthStaticService) StaticFilesMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
// 	// fmt.Println("Serving static file:", req.Data())
// 	data := req.Data()

// 	// If files end (.js, .css). In this case pass to the handler/endpoint function
// 	path := data.Path
// 	if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") {
// 		return next(req)
// 	}

// 	// Get the cookie from the request header
// 	cookie := data.Headers.Get("Cookie")
// 	parseCookies

// 	// If the file url end with /org-select
// 	endPath := strings.TrimPrefix(path, "/static/")

// 	switch endPath {
// 	case "org-select/":
// 	case "confirm-login/":
// 	case ""
// 	}

// 	if strings.HasPrefix(endPath, "org-select") {
// 		// Check if has access_token cookie
// 		return next(req)
// 	}

// 	// Check if has the auth_token cookie
// 	token := data.Headers.Get("Cookie")
// 	if token == "" {
// 		fmt.Println("No auth_token cookie found")
// 	} else {
// 		fmt.Println("auth_token cookie found:", token)
// 	}

// 	resp := next(req)

// 	return resp
// }
