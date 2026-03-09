package vault

import (
	"context"
	"errors"

	vcg "github.com/hashicorp/vault-client-go"
)

func (m *Manager) Delete(ctx context.Context, client *vcg.Client, mount, path string) error {
	u := lockRef(mount, path)
	m.locks.Lock(u)
	defer m.locks.Unlock(u)

	var err error
	v, err := DetectKVVersion(ctx, client, mount)
	if err != nil {
		return err
	}

	switch v {
	case KVVersionV1:

		_, err = client.Secrets.KvV1Delete(ctx, path, vcg.WithMountPath(mount))

	case KVVersionV2:
		_, err = client.Secrets.KvV2DeleteMetadataAndAllVersions(ctx, path, vcg.WithMountPath(mount))

	default:
		return errors.New("unknown kv engine version")
	}

	if err != nil {
		return toPkgErr(err)
	}

	return nil
}
