package vault

import (
	"errors"
	"net/http"

	vcg "github.com/hashicorp/vault-client-go"
)

var (
	ErrAlreadyExists   = errors.New("already exists")
	ErrForbidden       = errors.New("forbidden")
	ErrNotFound        = errors.New("not found")
	ErrNotDeclaredType = errors.New("secret not the declared type")
)

// Assumption: err != nil
func toPkgErr(err error) error {
	switch {
	case vcg.IsErrorStatus(err, http.StatusForbidden):
		return ErrForbidden
	case vcg.IsErrorStatus(err, http.StatusNotFound):
		return ErrNotFound
	}
	return err
}
