package sessions

import (
	"context"
	"net/http"

	"encore.app/pkg/myjwt"
	"encore.app/sessions/business"
	"encore.app/sessions/store"
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
func (s *ServiceSessions) CreateOrgSelectSession(ctx context.Context, userID uuid.UUID, deviceInfo string) (string, error) {
	orgSessionToken, err := s.b.CreateOrgSelectSession(ctx, userID, deviceInfo)
	if err != nil {
		return "", err
	}
	return orgSessionToken, nil
}

//encore:api private method=POST path=/session/membership
func (s *ServiceSessions) CreateMembershipSession(ctx context.Context, membershipID, sessionID uuid.UUID) (string, error) {
	membershipSessionToken, err := s.b.CreateMembershipSession(ctx, membershipID, sessionID)
	if err != nil {
		return "", err
	}
	return membershipSessionToken, nil
}

//encore:api public method=DELETE path=/session
func (s *ServiceSessions) DeleteWebUserSession(ctx context.Context, req requestDeleteSession) error {
	// Check if the cookie is valid
	if req.CookieToken.Value == "" {
		return &http.ProtocolError{ErrorString: "No auth_token cookie provided"}
	}

}

type requestDeleteSession struct {
	CookieToken http.Cookie `cookie:"auth_token"`
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
