package secret

import (
	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"
)

// Vault specific attribute definitions
var (
	vaultManagementURI = sm.DataAttributeV3{
		Uuid:          "ffd606d5-5fd0-4425-9a5a-29c2713ce18d",
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          "data_vault_management_uri",
		ContentType:   sm.AttributeContentTypeString,
		Description:   ptrStr("Vault URI should be in the following format: `http(s)://<vault-url>:<port>`."),
		Properties: sm.DataAttributeProperties{
			Label:    "Vault URI",
			Visible:  true,
			Required: true,
		},
	}
	vaultManagementRequestTmout = sm.DataAttributeV3{
		Uuid:          "4494de6b-7c33-44d2-8609-c3a561f5e3f1",
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          "data_vault_management_request_timeout",
		ContentType:   sm.AttributeContentTypeInteger,
		Description:   ptrStr("Request timeout in seconds applied to each Vault request."),
		Properties: sm.DataAttributeProperties{
			Label:    "Individual Vault request timeout",
			Visible:  true,
			Required: false,
		},
	}
	vaultManagementMount = sm.DataAttributeV3{
		Uuid:          "11541b02-6752-4651-8df3-86bed296af78",
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          "data_vault_management_mount",
		ContentType:   sm.AttributeContentTypeString,
		Description:   ptrStr("Vault mount."),
		Properties: sm.DataAttributeProperties{
			Label:    "Vault mount",
			Visible:  true,
			Required: true,
		},
	}
	vaultManagementCredentialType = sm.DataAttributeV3{
		Uuid:          "f461a9ab-7a99-4b41-b190-d0338e833064",
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          "data_vault_management_credential_type",
		ContentType:   sm.AttributeContentTypeString,
		Description:   ptrStr("List of available Vault authentication methods."),
		Properties: sm.DataAttributeProperties{
			Label:       "Please select an authentication method",
			Visible:     true,
			Required:    true,
			ReadOnly:    false,
			List:        true,
			MultiSelect: false,
		},
	}
	vaultManagementPath = sm.DataAttributeV3{
		Uuid:          "19c0493b-1eb3-4d20-9394-610f63078109",
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          "data_vault_management_secret_path",
		Type:          sm.Data,
		ContentType:   sm.AttributeContentTypeString,
		Description:   ptrStr("Path of secret in Vault without trailing slash."),
		Properties: sm.DataAttributeProperties{
			Label:    "Secret Path",
			Visible:  true,
			Required: false,
		},
	}
)

// credential type specific attribute definitions
var (
	credentialTypeAppRole = sm.StringAttributeContentV3{
		Reference: ptrStr("AppRole"),
		Data:      "approle",
	}
	credentialTypeJwt = sm.StringAttributeContentV3{
		Reference: ptrStr("JWT/OIDC"),
		Data:      "jwt",
	}
	credentialTypeK8s = sm.StringAttributeContentV3{
		Reference: ptrStr("Kubernetes"),
		Data:      "kubernetes",
	}
	vaultManagementRoleID = sm.DataAttributeV3{
		Uuid:          "302af8ad-0c4d-4eb2-9add-2d4a894c6b32",
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          "data_vault_management_role_id",
		ContentType:   sm.AttributeContentTypeResource,
		Description:   ptrStr("Vault AppRole ID."),
		Properties: sm.DataAttributeProperties{
			Label:    "AppRole ID",
			Visible:  true,
			Required: true,
		},
	}
	vaultManagementRoleSecret = sm.DataAttributeV3{
		Uuid:          "f8ee975c-aad5-4b48-bac9-46daa1a9a689",
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          "data_vault_management_role_secret",
		ContentType:   sm.AttributeContentTypeResource,
		Description:   ptrStr("Vault AppRole Secret."),
		Properties: sm.DataAttributeProperties{
			Label:    "AppRole Secret",
			Visible:  true,
			Required: true,
		},
	}
	vaultManagementRole = sm.DataAttributeV3{
		Uuid:          "e869cab0-80c9-44a0-900b-51791827edeb",
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          "data_vault_management_role",
		ContentType:   sm.AttributeContentTypeString,
		Description:   ptrStr("Vault Role."),
		Properties: sm.DataAttributeProperties{
			Label:    "Vault role",
			Visible:  true,
			Required: true,
		},
	}
	vaultManagementJwt = sm.DataAttributeV3{
		Uuid:          "94ec433b-9c50-4fcd-aabd-ab8da204d5db",
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          "data_vault_management_jwt",
		ContentType:   sm.AttributeContentTypeResource,
		Description:   ptrStr("Vault JWT."),
		Properties: sm.DataAttributeProperties{
			Label:    "JWT/OIDC",
			Visible:  true,
			Required: true,
		},
	}
)

const (
	VaultManagementCredentialGroupUUID = "2371992e-e074-4128-a53a-a877d6e548c6"
	VaultManagementCredentialGroupName = "group_vault_management_credential"
)
