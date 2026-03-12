package vault

import (
	"context"
	"errors"

	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"

	vcg "github.com/hashicorp/vault-client-go"
	vcgSchema "github.com/hashicorp/vault-client-go/schema"
)

func (m *Manager) Update(ctx context.Context, client *vcg.Client, mount, path string, secret sm.SecretContent) (sm.SecretType, error) {
	u := lockRef(mount, path)
	m.locks.Lock(u)
	defer m.locks.Unlock(u)

	v, err := DetectKVVersion(ctx, client, mount)
	if err != nil {
		return sm.SecretType(""), err
	}

	payload, secretType, err := ToPayload(ctx, secret)
	if err != nil {
		return secretType, err
	}

	switch v {
	case KVVersionV1:

		_, err = client.Secrets.KvV1Read(ctx, path, vcg.WithMountPath(mount))
		if err != nil {
			return secretType, toPkgErr(err)
		}

		_, err := client.Secrets.KvV1Write(ctx, path, payload, vcg.WithMountPath(mount))
		if err != nil {
			return secretType, toPkgErr(err)
		}
		return secretType, nil

	case KVVersionV2:
		_, err = client.Secrets.KvV2Read(ctx, path, vcg.WithMountPath(mount))
		if err != nil {
			return secretType, toPkgErr(err)
		}

		_, err := client.Secrets.KvV2Write(ctx, path, vcgSchema.KvV2WriteRequest{Data: payload}, vcg.WithMountPath(mount))
		if err != nil {
			return secretType, toPkgErr(err)
		}
		return secretType, nil
	}
	return secretType, errors.New("unknown kv engine version")
}
