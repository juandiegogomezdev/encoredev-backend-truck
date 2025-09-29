package authstatic

import (
	"fmt"

	"encore.dev/middleware"
)

//encore:middleware target=tag:static
func StaticFilesMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	fmt.Println("Serving static file:", req.Data().Path)
	resp := next(req)

	return resp
}

// //encore:middleware
// func AuthMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
// 	// Obtener la cookie del request
// 	cookie, err := req.Data().Headers.Get("Cookie")
// 	// if err != nil {
// 	// 	rlog.Info("Cookie no encontrada", "error", err)
// 	// 	return middleware.Response{
// 	// 		HTTPStatus: http.StatusUnauthorized,
// 	// 		Err:        ErrMissingCookie,
// 	// 	}
// 	// }

// 	// // Validar la cookie
// 	// cookieData, err := validateCookie(cookie.Value)
// 	// if err != nil {
// 	// 	rlog.Info("Cookie inválida", "error", err)
// 	// 	return middleware.Response{
// 	// 		HTTPStatus: http.StatusUnauthorized,
// 	// 		Err:        err,
// 	// 	}
// 	// }

// 	// // Agregar los datos del usuario al contexto
// 	// ctx := context.WithValue(req.Context(), "user_data", cookieData)
// 	// req = req.WithContext(ctx)

// 	// // Continuar con el siguiente middleware/handler
// 	return next(req)
// }
