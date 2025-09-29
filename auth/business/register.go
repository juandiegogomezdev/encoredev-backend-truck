package business

import (
	"context"
	"fmt"
)

func (b *BusinessAuth) CheckUser(ctx context.Context, email string) (string, error) {
	exists, err := b.store.UserExistsByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if exists {
		return "", fmt.Errorf("user already exists")

	}

	token, err := b.tokenizer.GenerateConfirmRegisterToken(email)
	if err != nil {
		return "", err
	}

	// Send email with the token (using the mailer)
	go func() {
		b.mailer.Send(
			email, "Confirm your registration",
			fmt.Sprintf(`Abre el siguiente link: <a href="http://localhost:4000/static/confirm-register?token=%s"> Click here! </a>`, token))
	}()

	return token, nil
}
