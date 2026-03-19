//go:build integration

package integration

import "time"

// Timeouts for WaitForCondition and test assertions
const (
	ShortTimeout  = 5 * time.Second
	MediumTimeout = 15 * time.Second
	LongTimeout   = 30 * time.Second
)

// Attribute UUIDs — copied from internal/secret/const.go (unexported vars, not importable)
const (
	AttrUUIDVaultURI       = "ffd606d5-5fd0-4425-9a5a-29c2713ce18d"
	AttrUUIDCredentialType = "f461a9ab-7a99-4b41-b190-d0338e833064"
	AttrUUIDRoleID         = "302af8ad-0c4d-4eb2-9add-2d4a894c6b32"
	AttrUUIDRoleSecret     = "f8ee975c-aad5-4b48-bac9-46daa1a9a689"
	AttrUUIDMount          = "11541b02-6752-4651-8df3-86bed296af78"
	AttrUUIDPathPrefix     = "19c0493b-1eb3-4d20-9394-610f63078109"
	AttrUUIDSecretPath     = "17e54346-3c10-4afe-b221-b4e0325c306d"
	AttrUUIDNamespace      = "b7755e40-3ad3-404b-af8d-55a8a1105213"
)

// Attribute names — must match the Name field in internal/secret/const.go
const (
	AttrNameVaultURI       = "data_vault_management_uri"
	AttrNameCredentialType = "data_vault_management_credential_type"
	AttrNameRoleID         = "data_vault_management_role_id"
	AttrNameRoleSecret     = "data_vault_management_role_secret"
	AttrNameMount          = "data_vault_management_profile_mount"
	AttrNamePathPrefix     = "data_vault_management_profile_secret_path_prefix"
	AttrNameSecretPath     = "data_secret_management_secret_path"
)

// Credential type values — must match internal/secret/const.go credentialTypeAppRole.Data
const (
	CredTypeAppRole = "approle"
)

// Expected UUIDs in GET /v1/secretProvider/vaults/attributes response
const (
	GroupUUIDVaultCredentials = "2371992e-e074-4128-a53a-a877d6e548c6"
)
