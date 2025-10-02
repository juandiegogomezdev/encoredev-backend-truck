package business

// // GenerateOrgAccessToken creates a new session for the user and returns a JWT token for accessing the API.
// func (b *BusinessAuth) GenerateOrgAccessToken(ctx context.Context, userID uuid.UUID) (string, error) {
// 	// TODO get device info from request

// 	// Count active sessions for the user
// 	count, err := b.store.CountSessionsByUserID(ctx, userID)
// 	if err != nil {
// 		fmt.Println("Error counting user sessions:", err)
// 		return "", &errs.Error{
// 			Code:    errs.Internal,
// 			Message: "Error al generar acceso al usuario",
// 		}
// 	}

// 	// Limit to 5 active sessions
// 	if count >= 5 {
// 		return "", &errs.Error{
// 			Code:    errs.PermissionDenied,
// 			Message: "Has alcanzado el límite de sesiones activas. Cierra sesión en otros dispositivos para continuar.",
// 		}
// 	}

// 	// Generate token with userID
// 	// This is before saving the session to avoid save and then fail to generate the token

// 	// Create user session record
// 	newSession := sessions.RequestCreateSessionStruct{
// 		UserID:     userID,
// 		SessionID:  sesionId,
// 		DeviceInfo: "web",
// 		CreatedAt:  utils.GetCurrentTime(),
// 		ExpiresAt:  utils.GetExpiryTime(24 * time.Hour),
// 	}

// 	// Save session in the database
// 	err = b.store.CreateUserSession(ctx, newSession)
// 	if err != nil {
// 		fmt.Println("Error saving user session:", err)
// 		return "", &errs.Error{
// 			Code:    errs.Internal,
// 			Message: "Error al generar acceso al usuario",
// 		}
// 	}

// 	return generatedToken, nil
// }

// func (b *BusinessAuth) GenerateConfirmRegisterToken(newEmail string) (string, error) {
// 	claims, err := b.tokenizer.GenerateConfirmRegisterToken(newEmail)
// 	if err != nil {
// 		fmt.Println("Error generating confirm register token:", err)
// 		return "", &errs.Error{
// 			Code:    errs.Internal,
// 			Message: "Error al generar el token de confirmación",
// 		}
// 	}
// 	return claims, nil
// }
