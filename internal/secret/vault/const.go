package vault

// Constants for string values of keys written to and expected from Vault secrets
const (
	UsernameKey     = "username"
	PasswordKey     = "password"
	ContentKey      = "content"
	KeyStoreTypeKey = "key-store-type"
)

// KVVersion is an enum for Vault KeyValue engine versions
type KVVersion int

const (
	KVVersionV1 KVVersion = iota
	KVVersionV2
)

func (v KVVersion) String() string {
	switch v {
	case KVVersionV1:
		return "v1"
	case KVVersionV2:
		return "v2"
	}
	return "unknown"
}
