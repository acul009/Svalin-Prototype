package permissions

import (
	"context"
	"fmt"
	"github.com/rahn-it/svalin/config"
	"github.com/rahn-it/svalin/ent"
	"github.com/rahn-it/svalin/ent/user"
	"github.com/rahn-it/svalin/pki"
)

var ErrPermissionDenied = PermissionDeniedError{}

type PermissionDeniedError struct {
	PublicKey *pki.PublicKey
	Reason    string
}

func (e PermissionDeniedError) Error() string {
	return fmt.Sprintf("permission denied: %s", e.Reason)
}

func (e PermissionDeniedError) Is(target error) bool {
	_, ok := target.(PermissionDeniedError)
	return ok
}

func MayStartCommand(sender *pki.PublicKey, command string) error {
	isRoot, err := pki.Root.MatchesKey(sender)
	if err != nil {
		return fmt.Errorf("failed to check if public key is CA: %w", err)
	}

	if isRoot {
		return nil
	}

	switch command {

	case "verify-certificate-chain":
		return nil

	default:
		db := config.DB()

		encoded := sender.Base64Encode()

		_, err = db.User.Query().Where(user.PublicKeyEQ(encoded)).Only(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				return PermissionDeniedError{
					PublicKey: sender,
					Reason:    "requested sender is not a user",
				}
			}
			return fmt.Errorf("failed to query user: %w", err)
		}

		return nil
	}

}
