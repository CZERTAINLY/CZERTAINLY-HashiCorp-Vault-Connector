package vault

import (
	"context"
	"errors"
	"net/http"

	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"

	vcg "github.com/hashicorp/vault-client-go"
	vcgSchema "github.com/hashicorp/vault-client-go/schema"
)

func (m *Manager) Create(ctx context.Context, client *vcg.Client, mount, path string, secret sm.SecretContent) error {
	u := lockRef(mount, path)
	m.locks.Lock(u)
	defer m.locks.Unlock(u)

	v, err := DetectKVVersion(ctx, client, mount)
	if err != nil {
		return err
	}

	payload, err := ToPayload(ctx, secret)
	if err != nil {
		return err
	}

	switch v {
	case KVVersionV1:
		_, err = client.Secrets.KvV1Read(ctx, path, vcg.WithMountPath(mount))
		switch {
		case err == nil:
			return ErrAlreadyExists

		case !vcg.IsErrorStatus(err, http.StatusNotFound):
			return toPkgErr(err)
		}

		_, err := client.Secrets.KvV1Write(ctx, path, payload, vcg.WithMountPath(mount))
		if err != nil {
			return toPkgErr(err)
		}
		return nil

	case KVVersionV2:
		_, err := client.Secrets.KvV2Write(ctx, path, vcgSchema.KvV2WriteRequest{
			Data: payload,
			Options: map[string]any{
				"cas": 0,
			},
		}, vcg.WithMountPath(mount))
		if err != nil {
			if vcg.IsErrorStatus(err, http.StatusBadRequest) {
				return ErrAlreadyExists
			}
			return toPkgErr(err)
		}
	default:
		return errors.New("unknown kv engine version")
	}

	return nil
}
