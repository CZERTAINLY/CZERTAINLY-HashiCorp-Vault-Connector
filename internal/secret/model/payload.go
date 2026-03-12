package model

import (
	"encoding/base64"
	"fmt"
)

func GetApiKeySecretContent(secret SecretContent) (string, error) {
	content, err := secret.AsApiKeySecretContent()
	if err != nil {
		return "", fmt.Errorf("unmarshalling SecretContent into ApiKeySecret failed: %w", err)
	}
	return content.Content, nil
}

func GetBasicAuthSecretContent(secret SecretContent) (username, password string, err error) {
	content, err := secret.AsBasicAuthSecretContent()
	if err != nil {
		err = fmt.Errorf("unmarshalling SecretContent into BasicAuthSecret failed: %w", err)
	}
	username = content.Username
	password = content.Password
	return
}

func GetGenericSecretContent(secret SecretContent) (string, error) {
	content, err := secret.AsGenericSecretContent()
	if err != nil {
		return "", fmt.Errorf("unmarshalling SecretContent into GenericSecret failed: %w", err)
	}
	return content.Content, nil
}

func GetJwtTokenSecretContent(secret SecretContent) (string, error) {
	content, err := secret.AsJwtTokenSecretContent()
	if err != nil {
		return "", fmt.Errorf("unmarshalling SecretContent into JwtTokenSecret failed: %w", err)
	}
	return content.Content, nil
}

func GetKeyStoreSecretContent(secret SecretContent) (keyStoreType KeyStoreType, content string, password string, err error) {
	ks, err := secret.AsKeyStoreSecretContent()
	if err != nil {
		err = fmt.Errorf("unmarshalling SecretContent into KeyStoreSecret failed: %w", err)
		return
	}
	content = ks.Content
	keyStoreType = ks.KeyStoreType
	password = ks.Password
	return
}

func GetKeyValueSecretContent(secret SecretContent) (map[string]any, error) {
	content, err := secret.AsKeyValueSecretContent()
	if err != nil {
		return nil, fmt.Errorf("unmarshalling SecretContent into KeyValueSecret failed: %w", err)
	}
	if len(content.Content) == 0 {
		return nil, fmt.Errorf("content of KeyValueSecret is empty")
	}
	return content.Content, nil
}

func GetPrivateKeySecretContent(secret SecretContent) ([]byte, error) {
	content, err := secret.AsPrivateKeySecretContent()
	if err != nil {
		return nil, fmt.Errorf("unmarshalling SecretContent into PrivateKeySecret failed: %w", err)
	}
	decoded, err := base64.StdEncoding.DecodeString(content.Content)
	if err != nil {
		return nil, fmt.Errorf("base64 decoding PrivateKeySecret content failed: %w", err)
	}
	return decoded, nil
}

func GetSecretKeySecretContent(secret SecretContent) (string, error) {
	content, err := secret.AsSecretKeySecretContent()
	if err != nil {
		return "", fmt.Errorf("unmarshalling SecretContent into SecretKeySecret failed: %w", err)
	}
	return content.Content, nil
}
