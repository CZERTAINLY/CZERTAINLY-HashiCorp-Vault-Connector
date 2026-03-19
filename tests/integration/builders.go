//go:build integration

package integration

import "encoding/json"

// RequestBuilder constructs request payloads for Secret Provider endpoints.
type RequestBuilder struct {
	vaultURL   string
	roleID     string
	secretID   string
	mount      string
	secretPath string
}

// NewRequestBuilder creates a builder with Vault connection details.
// mount should include trailing slash (e.g. "kv-test/").
func NewRequestBuilder(vaultURL, roleID, secretID, mount string) *RequestBuilder {
	return &RequestBuilder{
		vaultURL: vaultURL,
		roleID:   roleID,
		secretID: secretID,
		mount:    mount,
	}
}

// WithSecretPath sets an optional relative secret path appended to the mount.
func (b *RequestBuilder) WithSecretPath(path string) *RequestBuilder {
	b.secretPath = path
	return b
}

// BuildCheckConnectionBody produces the bare []RequestAttribute JSON array
// expected by POST /v1/secretProvider/vaults.
// (Other endpoints wrap attributes in a DTO struct; this one does not.)
func (b *RequestBuilder) BuildCheckConnectionBody() []byte {
	attrs := b.vaultAttrs()
	data, _ := json.Marshal(attrs)
	return data
}

// BuildCreateSecretBody produces the CreateSecretRequestDto JSON body.
func (b *RequestBuilder) BuildCreateSecretBody(name string, secret map[string]any) []byte {
	body := map[string]any{
		"name":                   name,
		"secret":                 secret,
		"vaultAttributes":        b.vaultAttrs(),
		"vaultProfileAttributes": b.vaultProfileAttrs(),
	}
	if sa := b.secretAttrs(); sa != nil {
		body["secretAttributes"] = sa
	}
	data, _ := json.Marshal(body)
	return data
}

// BuildReadSecretBody produces the SecretRequestDto JSON body for reading a secret.
func (b *RequestBuilder) BuildReadSecretBody(name, secretType string) []byte {
	body := map[string]any{
		"name":                   name,
		"type":                   secretType,
		"vaultAttributes":        b.vaultAttrs(),
		"vaultProfileAttributes": b.vaultProfileAttrs(),
	}
	if sa := b.secretAttrs(); sa != nil {
		body["secretAttributes"] = sa
	}
	data, _ := json.Marshal(body)
	return data
}

// BuildUpdateSecretBody produces the UpdateSecretRequestDto JSON body.
func (b *RequestBuilder) BuildUpdateSecretBody(name string, secret map[string]any) []byte {
	body := map[string]any{
		"name":                   name,
		"secret":                 secret,
		"vaultAttributes":        b.vaultAttrs(),
		"vaultProfileAttributes": b.vaultProfileAttrs(),
	}
	if sa := b.secretAttrs(); sa != nil {
		body["secretAttributes"] = sa
	}
	data, _ := json.Marshal(body)
	return data
}

// BuildDeleteSecretBody produces the SecretRequestDto JSON body for deleting a secret.
func (b *RequestBuilder) BuildDeleteSecretBody(name string) []byte {
	body := map[string]any{
		"name":                   name,
		"type":                   "generic",
		"vaultAttributes":        b.vaultAttrs(),
		"vaultProfileAttributes": b.vaultProfileAttrs(),
	}
	if sa := b.secretAttrs(); sa != nil {
		body["secretAttributes"] = sa
	}
	data, _ := json.Marshal(body)
	return data
}

// --- Secret content helpers ---

// GenericSecret produces the SecretContent JSON for a generic secret.
func GenericSecret(value string) map[string]any {
	return map[string]any{"type": "generic", "content": value}
}

// BasicAuthSecret produces the SecretContent JSON for basic auth credentials.
func BasicAuthSecret(username, password string) map[string]any {
	return map[string]any{"type": "basicAuth", "username": username, "password": password}
}

// ApiKeySecret produces the SecretContent JSON for an API key.
// The field name is "content" — matching ApiKeySecretContent.Content in the generated model.
func ApiKeySecret(key string) map[string]any {
	return map[string]any{"type": "apiKey", "content": key}
}

// JwtTokenSecret produces the SecretContent JSON for a JWT token.
// The token is stored and returned as-is (no JWT signature validation by the connector).
func JwtTokenSecret(token string) map[string]any {
	return map[string]any{"type": "jwtToken", "content": token}
}

// SecretKeySecret produces the SecretContent JSON for a symmetric secret key.
// content should be a base64-encoded string.
func SecretKeySecret(content string) map[string]any {
	return map[string]any{"type": "secretKey", "content": content}
}

// PrivateKeySecret produces the SecretContent JSON for a PEM-encoded private key.
// pemBase64 must be the standard-encoding base64 of the PEM text (the HTTP layer decodes
// it before calling pem.Decode to validate structure). On read, the connector re-encodes
// the stored PEM as base64 in the response.
func PrivateKeySecret(pemBase64 string) map[string]any {
	return map[string]any{"type": "privateKey", "content": pemBase64}
}

// KeyValueSecret produces the SecretContent JSON for arbitrary key-value pairs.
// data is stored as-is in Vault with no schema validation.
func KeyValueSecret(data map[string]any) map[string]any {
	return map[string]any{"type": "keyValue", "content": data}
}

// KeyStoreSecret produces the SecretContent JSON for a JKS or PKCS12 keystore.
// content should be the keystore bytes base64-encoded.
// keyStoreType must be "JKS" or "PKCS12" (case-sensitive).
func KeyStoreSecret(content, password, keyStoreType string) map[string]any {
	return map[string]any{
		"type":         "keyStore",
		"content":      content,
		"password":     password,
		"keyStoreType": keyStoreType,
	}
}

// --- Internal attribute builders ---

// vaultAttrs returns the vault-instance attributes (URI + credential type + AppRole creds).
func (b *RequestBuilder) vaultAttrs() []map[string]any {
	return []map[string]any{
		stringAttr(AttrUUIDVaultURI, AttrNameVaultURI, b.vaultURL),
		stringAttr(AttrUUIDCredentialType, AttrNameCredentialType, CredTypeAppRole),
		resourceSecretAttr(AttrUUIDRoleID, AttrNameRoleID, b.roleID),
		resourceSecretAttr(AttrUUIDRoleSecret, AttrNameRoleSecret, b.secretID),
	}
}

// vaultProfileAttrs returns the vault-profile attributes (mount point).
func (b *RequestBuilder) vaultProfileAttrs() []map[string]any {
	return []map[string]any{
		stringAttr(AttrUUIDMount, AttrNameMount, b.mount),
	}
}

// secretAttrs returns the optional secret-level attributes (relative path).
func (b *RequestBuilder) secretAttrs() []map[string]any {
	if b.secretPath == "" {
		return nil
	}
	return []map[string]any{
		stringAttr(AttrUUIDSecretPath, AttrNameSecretPath, b.secretPath),
	}
}

// stringAttr produces a RequestAttributeV3 JSON object for a string-typed attribute.
func stringAttr(uuid, name, value string) map[string]any {
	return map[string]any{
		"uuid":        uuid,
		"name":        name,
		"contentType": "string",
		"version":     "v3",
		"content": []map[string]any{
			{"data": value},
		},
	}
}

// resourceSecretAttr produces a RequestAttributeV3 JSON object for a resource-typed attribute
// (used for roleID and roleSecret). The value is wrapped as an apiKey-typed secret.
//
// The connector's needs.go calls resourceSecretContentTypeDataAttrSingle() which expects:
// - contentType = "resource" on the attribute
// - content[0] is a ResourceObjectContent with contentType="resource"
// - data.resource = "secrets"
// - data.content has type discriminator "apiKey" and the actual value in "content"
//   (matching ApiKeySecretContent.Content json:"content" in the generated model)
func resourceSecretAttr(uuid, name, value string) map[string]any {
	return map[string]any{
		"uuid":        uuid,
		"name":        name,
		"contentType": "resource",
		"version":     "v3",
		"content": []map[string]any{
			{
				"contentType": "resource",
				"data": map[string]any{
					"resource": "secrets",
					"name":     "test-credential",
					"uuid":     "00000000-0000-0000-0000-000000000001",
					"content": map[string]any{
						"type":    "apiKey",
						"content": value,
					},
				},
			},
		},
	}
}
