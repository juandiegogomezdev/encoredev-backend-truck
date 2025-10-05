package appService

import (
	"context"
	"net/http"

	"encore.app/pkg/utils"
	"encore.dev/beta/errs"
)

// //encore:api private method=POST path=/sessions/org-select
// func (s *ServiceApp) CreateOrgSelectSession(ctx context.Context, req RequestCreateOrgSelectSession) (responseCreateOrgSelectSession, error) {
// 	orgSessionToken, err := s.b.CreateOrgSelectSession(ctx, req.UserID, req.DeviceInfo)
// 	if err != nil {
// 		return responseCreateOrgSelectSession{}, err
// 	}
// 	return responseCreateOrgSelectSession{
// 		OrgSelectSessionToken: orgSessionToken,
// 	}, nil
// }

// type RequestCreateOrgSelectSession struct {
// 	UserID     uuid.UUID `json:"user_id"`
// 	DeviceInfo string    `json:"device_info"`
// }

// type responseCreateOrgSelectSession struct {
// 	OrgSelectSessionToken string `json:"org_select_session_token"`
// }

// //encore:api private method=POST path=/session/membership
// func (s *ServiceApp) CreateMembershipSession(ctx context.Context, req RequestCreateMembershipSession) (responseCreateMembershipSession, error) {
// 	membershipSessionToken, err := s.b.CreateMembershipSession(ctx, req.MembershipID, req.SessionID)
// 	if err != nil {
// 		return responseCreateMembershipSession{}, err
// 	}
// 	return responseCreateMembershipSession{
// 		MembershipSessionToken: membershipSessionToken,
// 	}, nil
// }

// type RequestCreateMembershipSession struct {
// 	MembershipID uuid.UUID `json:"membership_id"`
// 	SessionID    uuid.UUID `json:"session_id"`
// }

// type responseCreateMembershipSession struct {
// 	MembershipSessionToken string `json:"membership_session_token"`
// }

//encore:api public method=POST path=/session/refresh
func (s *ServiceApp) RefreshToken(ctx context.Context) error {
	return &errs.Error{
		Code:    errs.Unimplemented,
		Message: "Not implemented",
	}
}

//encore:api public method=POST path=/session/delete/web
func (s *ServiceApp) DeleteWebUserSession(ctx context.Context, req requestDeleteSessionWeb) (responseDeleteSessionWeb, error) {
	// Generate the expired cookie to delete the cookie in the browser
	deleteCookie := utils.DeleteDefaultCookieOptions("auth_token")

	// Delete the session in a goroutine to not block the response
	go func() {
		// Check if the cookie is valid
		if req.SessionCookie == nil || req.SessionCookie.Value == "" {
			return
		}
		s.b.DeleteUserSession(ctx, req.SessionCookie.Value)
	}()

	return responseDeleteSessionWeb{
		SessionCookie: deleteCookie,
	}, nil

}

// Request struct to get the cookie from the request
type requestDeleteSessionWeb struct {
	SessionCookie *http.Cookie `cookie:"auth_token"`
}

// Set the expired cookie in the response header
type responseDeleteSessionWeb struct {
	SessionCookie string `header:"Set-Cookie"`
}

//encore:api public method=POST path=/session/delete/mobile
func (s *ServiceApp) DeleteMobileSession(ctx context.Context, req requestDeleteSessionMobile) (responseDeleteMobileSession, error) {
	if req.Authorization == "" {
		return responseDeleteMobileSession{}, &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "No se puede cerrar sesión.",
		}
	}

	return responseDeleteMobileSession{
		Success: true,
	}, nil
}

type requestDeleteSessionMobile struct {
	Authorization string `header:"Authorization"`
}

type responseDeleteMobileSession struct {
	Success bool `json:"success"`
}
