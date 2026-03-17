package secret

import (
	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"
)

// Vault specific attribute definitions
var (
	vaultManagementInfo = sm.InfoAttributeV3{
		Uuid:          "890470a6-8cdd-4af9-a344-4f409dda4a64",
		Version:       ptr(int32(3)),
		SchemaVersion: sm.V3,
		Name:          "info_vault_management_explanation",
		Description:   ptr("Create a new HashiCorp Vault instance configuration"),
		ContentType:   sm.AttributeContentTypeText,
		Properties: sm.InfoAttributeProperties{
			Label:   "HashiCorp Vault instance configuration",
			Visible: true,
		},
	}
	vaultManagementURI = sm.DataAttributeV3{
		Uuid:          "ffd606d5-5fd0-4425-9a5a-29c2713ce18d",
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          "data_vault_management_uri",
		ContentType:   sm.AttributeContentTypeString,
		Description:   ptr("Vault URL should be in the following format: `http(s)://<vault-url>:<port>`."),
		Properties: sm.DataAttributeProperties{
			Label:    "Vault URL",
			Visible:  true,
			Required: true,
		},
	}
	vaultManagementMount = sm.DataAttributeV3{
		Uuid:          "11541b02-6752-4651-8df3-86bed296af78",
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          "data_vault_management_mount",
		ContentType:   sm.AttributeContentTypeString,
		Description:   ptr("Vault mount point"),
		Properties: sm.DataAttributeProperties{
			Label:    "Vault mount point",
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
		Description:   ptr("List of available Vault authentication methods."),
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
		Description:   ptr("Path of secret in Vault without trailing slash."),
		Properties: sm.DataAttributeProperties{
			Label:    "Secret path prefix",
			Visible:  true,
			Required: false,
		},
	}
)

// credential type specific attribute definitions
var (
	credentialTypeAppRole = sm.StringAttributeContentV3{
		Reference: ptr("AppRole"),
		Data:      "approle",
	}
	credentialTypeJwt = sm.StringAttributeContentV3{
		Reference: ptr("JWT/OIDC"),
		Data:      "jwt",
	}
	credentialTypeK8s = sm.StringAttributeContentV3{
		Reference: ptr("Kubernetes"),
		Data:      "kubernetes",
	}
	vaultManagementRoleID = sm.DataAttributeV3{
		Version:     3,
		Uuid:        "302af8ad-0c4d-4eb2-9add-2d4a894c6b32",
		Name:        "data_vault_management_role_id",
		Description: ptr("Vault AppRole ID."),
		ContentType: sm.AttributeContentTypeResource,
		Properties: sm.DataAttributeProperties{
			Resource: ptr(sm.Secrets),
			Visible:  true,
			Label:    "AppRole ID",
			Required: true,
		},
		AttributeCallback: &sm.AttributeCallback{
			Mappings: []sm.AttributeCallbackMapping{
				{
					To:      "SECRET_TYPE.EQUALS",
					Targets: []sm.AttributeValueTarget{sm.Filter},
					Value:   []sm.SecretType{sm.Generic, sm.ApiKey},
				},
			},
		},
		SchemaVersion: sm.V3,
	}
	vaultManagementRoleSecret = sm.DataAttributeV3{
		Uuid:          "f8ee975c-aad5-4b48-bac9-46daa1a9a689",
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          "data_vault_management_role_secret",
		ContentType:   sm.AttributeContentTypeResource,
		Description:   ptr("Vault AppRole Secret."),
		Properties: sm.DataAttributeProperties{
			Resource: ptr(sm.Secrets),
			Label:    "AppRole Secret",
			Visible:  true,
			Required: true,
		},
		AttributeCallback: &sm.AttributeCallback{
			Mappings: []sm.AttributeCallbackMapping{
				{
					To:      "SECRET_TYPE.EQUALS",
					Value:   []sm.SecretType{sm.ApiKey, sm.Generic},
					Targets: []sm.AttributeValueTarget{sm.Filter},
				},
			},
		},
	}
	vaultManagementRole = sm.DataAttributeV3{
		Uuid:          "e869cab0-80c9-44a0-900b-51791827edeb",
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          "data_vault_management_role",
		ContentType:   sm.AttributeContentTypeString,
		Description:   ptr("Vault Role."),
		Properties: sm.DataAttributeProperties{
			Label:    "Vault role",
			Visible:  true,
			Required: true,
		},
	}
)

const (
	vaultInfoContentDescrConst = "### HashiCorp Vault instance configuration.\n\nProvide URL of your Vault and select one of the available authentication methods:\n-  **AppRole** - Use AppRole authentication method with Role ID and Secret ID\n-  **Kubernetes** - Use Kubernetes authentication method with Service Account Token (automatically taken from the environment)\n-  **JWT/OIDC** - Use JWT/OIDC authentication method with provided JWT token (automatically taken from the environment)\n\n**Vault mount point** - The mount point is the root level \"directory\" where a secrets engine is enabled in Vault\n\n**Secret path prefix** - Relative path that is prepended to each request. Useful if you don't want to have all the secrets in the root level"
)

const (
	VaultManagementCredentialGroupUUID = "2371992e-e074-4128-a53a-a877d6e548c6"
	VaultManagementCredentialGroupName = "group_vault_management_credential"
)
