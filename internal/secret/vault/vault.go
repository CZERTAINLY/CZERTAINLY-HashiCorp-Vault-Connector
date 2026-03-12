package vault

import (
	"context"
	"encoding/pem"
	"fmt"

	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"

	vcg "github.com/hashicorp/vault-client-go"
	"github.com/moby/locker/v2"
)

type Manager struct {
	locks locker.RWMutexMap[string]
}

func New() *Manager {
	return &Manager{}
}

func (m *Manager) ConnCheck(ctx context.Context, client *vcg.Client) error {
	if _, err := client.System.MountsListSecretsEngines(ctx); err != nil {
		return toPkgErr(err)
	}

	return nil
}

func ToPayload(ctx context.Context, secret sm.SecretContent) (map[string]any, error) {
	secretType, err := secret.Discriminator()
	if err != nil {
		return nil, fmt.Errorf("unmarshalling discriminator field for SecretContent failed: %w", err)
	}

	switch sm.SecretType(secretType) {
	case sm.ApiKey:
		apiKeyContent, err := sm.GetApiKeySecretContent(secret)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			ContentKey: apiKeyContent,
		}, nil

	case sm.BasicAuth:
		username, password, err := sm.GetBasicAuthSecretContent(secret)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			UsernameKey: username,
			PasswordKey: password,
		}, nil

	case sm.Generic:
		content, err := sm.GetGenericSecretContent(secret)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			ContentKey: content,
		}, nil

	case sm.JwtToken:
		content, err := sm.GetJwtTokenSecretContent(secret)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			ContentKey: string(content),
		}, nil

	case sm.KeyStore:
		keyStoreType, content, password, err := sm.GetKeyStoreSecretContent(secret)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			KeyStoreTypeKey: keyStoreType,
			ContentKey:      string(content),
			PasswordKey:     password,
		}, nil

	case sm.KeyValue:
		content, err := sm.GetKeyValueSecretContent(secret)
		if err != nil {
			return nil, err
		}
		return content, nil

	case sm.PrivateKey:
		decoded, err := sm.GetPrivateKeySecretContent(secret)
		if err != nil {
			return nil, err
		}
		if !isPemFormat(decoded) {
			return nil, fmt.Errorf("%w: not PEM format", ErrNotDeclaredType)
		}
		return map[string]any{
			ContentKey: string(decoded),
		}, nil

	case sm.SecretKey:
		secretKeyContent, err := sm.GetSecretKeySecretContent(secret)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			ContentKey: secretKeyContent,
		}, nil
	}
	return nil, fmt.Errorf("unknown secret type %q", secretType)
}

func isPemFormat(pk []byte) bool {
	block, _ := pem.Decode(pk)
	return block != nil
}

func lockRef(mount, path string) string {
	return fmt.Sprintf("%s%s", mount, path)
}
