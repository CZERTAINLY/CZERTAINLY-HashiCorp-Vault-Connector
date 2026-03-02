package vault

import (
	"context"
	"encoding/base64"
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
		content, err := secret.AsApiKeySecretContent()
		if err != nil {
			return nil, fmt.Errorf("unmarshalling SecretContent into ApiKeySecret failed: %w", err)
		}
		decoded, err := base64.StdEncoding.DecodeString(content.Content)
		if err != nil {
			return nil, fmt.Errorf("base64 decoding ApiKeySecret content failed: %w", err)
		}
		return map[string]any{
			ContentKey: string(decoded),
		}, nil

	case sm.BasicAuth:
		content, err := secret.AsBasicAuthSecretContent()
		if err != nil {
			return nil, fmt.Errorf("unmarshalling SecretContent into BasicAuthSecret failed: %w", err)
		}
		return map[string]any{
			UsernameKey: content.Username,
			PasswordKey: content.Password,
		}, nil

	case sm.Generic:
		content, err := secret.AsGenericSecretContent()
		if err != nil {
			return nil, fmt.Errorf("unmarshalling SecretContent into GenericSecret failed: %w", err)
		}
		decoded, err := base64.StdEncoding.DecodeString(content.Content)
		if err != nil {
			return nil, fmt.Errorf("base64 decoding GenericSecret content failed: %w", err)
		}
		return map[string]any{
			ContentKey: string(decoded),
		}, nil

	case sm.JwtToken:
		content, err := secret.AsJwtTokenSecretContent()
		if err != nil {
			return nil, fmt.Errorf("unmarshalling SecretContent into JwtTokenSecret failed: %w", err)
		}
		decoded, err := base64.StdEncoding.DecodeString(content.Content)
		if err != nil {
			return nil, fmt.Errorf("base64 decoding JwtTokenSecret content failed: %w", err)
		}
		return map[string]any{
			ContentKey: string(decoded),
		}, nil

	case sm.KeyStore:
		content, err := secret.AsKeyStoreSecretContent()
		if err != nil {
			return nil, fmt.Errorf("unmarshalling SecretContent into KeyStoreSecret failed: %w", err)
		}
		decodedContent, err := base64.StdEncoding.DecodeString(content.Content)
		if err != nil {
			return nil, fmt.Errorf("base64 decoding KeyStoreSecret content failed: %w", err)
		}
		return map[string]any{
			KeyStoreTypeKey: content.KeyStoreType,
			ContentKey:      string(decodedContent),
			PasswordKey:     content.Password,
		}, nil

	case sm.KeyValue:
		content, err := secret.AsKeyValueSecretContent()
		if err != nil {
			return nil, fmt.Errorf("unmarshalling SecretContent into KeyValueSecret failed: %w", err)
		}
		if len(content.Content) == 0 {
			return nil, fmt.Errorf("content of KeyValueSecret is empty")
		}
		return content.Content, nil

	case sm.PrivateKey:
		content, err := secret.AsPrivateKeySecretContent()
		if err != nil {
			return nil, fmt.Errorf("unmarshalling SecretContent into PrivateKeySecret failed: %w", err)
		}
		decoded, err := base64.StdEncoding.DecodeString(content.Content)
		if err != nil {
			return nil, fmt.Errorf("base64 decoding PrivateKeySecret content failed: %w", err)
		}
		if !isPemFormat(decoded) {
			return nil, fmt.Errorf("%w: not PEM format", ErrNotDeclaredType)
		}
		return map[string]any{
			ContentKey: string(decoded),
		}, nil

	case sm.SecretKey:
		content, err := secret.AsSecretKeySecretContent()
		if err != nil {
			return nil, fmt.Errorf("unmarshalling SecretContent into SecretKeySecret failed: %w", err)
		}
		decoded, err := base64.StdEncoding.DecodeString(content.Content)
		if err != nil {
			return nil, fmt.Errorf("base64 decoding SecretKeySecret content failed: %w", err)
		}
		return map[string]any{
			ContentKey: string(decoded),
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
