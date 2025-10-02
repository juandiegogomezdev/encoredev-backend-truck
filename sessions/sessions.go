package sessions

import (
	"context"
	"net/http"

	"encore.app/pkg/myjwt"
	"encore.app/pkg/utils"
	"encore.app/sessions/business"
	"encore.app/sessions/store"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"encore.dev/types/uuid"
)

var sessionsDB = sqldb.NewDatabase("db_sessions", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

//encore:service
type ServiceSessions struct {
	b *business.BusinessSession
}

var secrets struct {
	JWT_SECRET_KEY string
	RESEND_API_KEY string
}

func initServiceSessions() (*ServiceSessions, error) {
	s := store.NewSessionStore(sessionsDB)

	tokenizer := myjwt.NewJWTTokenizer(secrets.JWT_SECRET_KEY)
	b := business.NewBusinessSession(s, tokenizer)

	return &ServiceSessions{b: b}, nil
}

//encore:api private method=POST path=/sessions/org-select
func (s *ServiceSessions) CreateOrgSelectSession(ctx context.Context, req RequestCreateOrgSelectSession) (responseCreateOrgSelectSession, error) {
	orgSessionToken, err := s.b.CreateOrgSelectSession(ctx, req.UserID, req.DeviceInfo)
	if err != nil {
		return responseCreateOrgSelectSession{}, err
	}
	return responseCreateOrgSelectSession{
		OrgSelectSessionToken: orgSessionToken,
	}, nil
}

type RequestCreateOrgSelectSession struct {
	UserID     uuid.UUID `json:"user_id"`
	DeviceInfo string    `json:"device_info"`
}

type responseCreateOrgSelectSession struct {
	OrgSelectSessionToken string `json:"org_select_session_token"`
}

//encore:api private method=POST path=/session/membership
func (s *ServiceSessions) CreateMembershipSession(ctx context.Context, req RequestCreateMembershipSession) (responseCreateMembershipSession, error) {
	membershipSessionToken, err := s.b.CreateMembershipSession(ctx, req.MembershipID, req.SessionID)
	if err != nil {
		return responseCreateMembershipSession{}, err
	}
	return responseCreateMembershipSession{
		MembershipSessionToken: membershipSessionToken,
	}, nil
}

type RequestCreateMembershipSession struct {
	MembershipID uuid.UUID `json:"membership_id"`
	SessionID    uuid.UUID `json:"session_id"`
}

type responseCreateMembershipSession struct {
	MembershipSessionToken string `json:"membership_session_token"`
}

//encore:api public method=POST path=/session/delete/web
func (s *ServiceSessions) DeleteWebUserSession(ctx context.Context, req requestDeleteSessionWeb) (responseDeleteSessionWeb, error) {
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
func (s *ServiceSessions) DeleteMobileSession(ctx context.Context, req requestDeleteSessionMobile) (responseDeleteMobileSession, error) {
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

// //encore:api private method=POST path=/sessions
// func (s *ServiceSessions) CreateSession(ctx context.Context, req RequestCreateSessionStruct) (ResponseCreateSessionStruct, error) {
// 	sessionID, err := s.b.CreateUserSession(ctx, req.UserID, req.DeviceInfo)
// 	if err != nil {
// 		return ResponseCreateSessionStruct{}, err
// 	}
// 	return ResponseCreateSessionStruct{
// 		success:   true,
// 		sessionID: sessionID,
// 	}, nil
// }

// type RequestCreateSessionStruct struct {
// 	UserID     uuid.UUID
// 	DeviceInfo string
// }

// type ResponseCreateSessionStruct struct {
// 	success   bool
// 	sessionID uuid.UUID
// }
