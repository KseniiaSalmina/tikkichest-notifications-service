package validation

import (
	"fmt"

	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/api/models"
	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/api/response_errors"
)

func User(user models.User) error {
	if err := UserID(user.ID); err != nil {
		return fmt.Errorf("invalid user data: %w", err)
	}

	if err := username(user.Username); err != nil {
		return fmt.Errorf("invalid user data: %w", err)
	}

	return nil
}

func UserID(id int) error {
	if id <= 0 {
		return response_errors.InvalidUserIDErr
	}

	return nil
}

func username(username string) error {
	if username == "" {
		return response_errors.InvalidUsernameErr
	}

	return nil
}
