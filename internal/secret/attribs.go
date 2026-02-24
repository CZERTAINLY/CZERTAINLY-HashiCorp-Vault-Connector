package secret

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"

	vcg "github.com/hashicorp/vault-client-go"
)

const (
	defaultRequestTimeout = 30 * time.Second
)

type Needs struct {
	address    string
	mount      string
	path       string
	reqTimeout time.Duration
	// TODO: rework this into a login structure with it's own methods so that we can have different auth methods
	token string
}

func ProcessAttrs(ctx context.Context, vaultAttrs, secretAttrs []sm.RequestAttribute) (Needs, error) {
	attrs := append(vaultAttrs, secretAttrs...)
	n := Needs{
		reqTimeout: defaultRequestTimeout,
	}
	for _, cpy := range attrs {
		// Invariant: Decision was made that we'll only accept v3 attributes
		attr, err := cpy.AsRequestAttributeV3()
		if err != nil {
			return n, fmt.Errorf("unmarshalling RequestAttribute into RequestAttributeV3 failed: %w", err)
		}

		switch attr.Uuid.String() {
		case SecretManagementVaultAddressUUID:
			if attr.ContentType != SecretManagementVaultAddressContentType {
				return n, fmt.Errorf("attribute %q has declared content type %q but received %q",
					SecretManagementVaultAddressUUID, SecretManagementVaultAddressContentType, attr.ContentType)
			}
			if len(*attr.Content) != 1 {
				return n, fmt.Errorf("attribute %q expects one content item, received: %d",
					SecretManagementVaultAddressUUID, len(*attr.Content))
			}
			strAttr, err := (*attr.Content)[0].AsStringAttributeContentV3()
			if err != nil {
				return n, fmt.Errorf("unmarshalling BaseAttributeContentDtoV3 into StringAttributeContentV3 failed for attribute %q: %w",
					SecretManagementVaultAddressUUID, err)
			}
			n.address = strAttr.Data

		case SecretManagementVaultRequestTimeoutSecondsUUID:
			if attr.ContentType != SecretManagementVaultRequestTimeoutSecondsContentType {
				return n, fmt.Errorf("attribute %q has declared content type %q but received %q",
					SecretManagementVaultRequestTimeoutSecondsUUID, SecretManagementVaultRequestTimeoutSecondsContentType, attr.ContentType)
			}
			if len(*attr.Content) != 1 {
				return n, fmt.Errorf("attribute %q expects one content item, received: %d",
					SecretManagementVaultRequestTimeoutSecondsUUID, len(*attr.Content))
			}
			intAttr, err := (*attr.Content)[0].AsIntegerAttributeContentV3()
			if err != nil {
				return n, fmt.Errorf("unmarshalling BaseAttributeContentDtoV3 into IntegerAttributeContentV3 failed for attribute %q: %w",
					SecretManagementVaultRequestTimeoutSecondsUUID, err)
			}
			n.reqTimeout = time.Duration(intAttr.Data) * time.Second

		case SecretManagementVaultMountUUID:
			if attr.ContentType != SecretManagementVaultMountContentType {
				return n, fmt.Errorf("attribute %q has declared content type %q but received %q",
					SecretManagementVaultMountUUID, SecretManagementVaultMountContentType, attr.ContentType)
			}
			if len(*attr.Content) != 1 {
				return n, fmt.Errorf("attribute %q expects one content item, received: %d",
					SecretManagementVaultMountUUID, len(*attr.Content))
			}
			strAttr, err := (*attr.Content)[0].AsStringAttributeContentV3()
			if err != nil {
				return n, fmt.Errorf("unmarshalling BaseAttributeContentDtoV3 into StringAttributeContentV3 failed for attribute %q: %w",
					SecretManagementVaultMountUUID, err)
			}
			n.mount = strAttr.Data

		case SecretManagementSecretPathUUID:
			if attr.ContentType != SecretManagementSecretPathContentType {
				return n, fmt.Errorf("attribute %q has declared content type %q but received %q",
					SecretManagementSecretPathUUID, SecretManagementSecretPathContentType, attr.ContentType)
			}
			if len(*attr.Content) != 1 {
				return n, fmt.Errorf("attribute %q expects one content item, received: %d",
					SecretManagementSecretPathUUID, len(*attr.Content))
			}
			strAttr, err := (*attr.Content)[0].AsStringAttributeContentV3()
			if err != nil {
				return n, fmt.Errorf("unmarshalling BaseAttributeContentDtoV3 into StringAttributeContentV3 failed for attribute %q: %w",
					SecretManagementSecretPathUUID, err)
			}
			n.path = strAttr.Data

		case SecretManagementVaultTokenUUID:
			if attr.ContentType != SecretManagementVaultTokenContentType {
				return n, fmt.Errorf("attribute %q has declared content type %q but received %q",
					SecretManagementVaultTokenUUID, SecretManagementVaultTokenContentType, attr.ContentType)
			}
			if len(*attr.Content) != 1 {
				return n, fmt.Errorf("attribute %q expects one content item, received: %d",
					SecretManagementVaultTokenUUID, len(*attr.Content))
			}
			strAttr, err := (*attr.Content)[0].AsStringAttributeContentV3()
			if err != nil {
				return n, fmt.Errorf("unmarshalling BaseAttributeContentDtoV3 into StringAttributeContentV3 failed for attribute %q: %w",
					SecretManagementVaultTokenUUID, err)
			}
			n.token = strAttr.Data

		default:
		}
	}
	// Invariants:
	// * mount MUST end with slash
	// * path MUST NOT end with slash
	if n.mount[len(n.mount)-1] != '/' {
		n.mount += "/"
	}
	n.path = strings.TrimSuffix(n.path, "/")

	if !n.satisfied() {
		return n, errors.New("required attributes missing")
	}
	return n, nil
}

func (n Needs) satisfied() bool {
	// TODO(improvement): could return an error with some specifics of what's missing etc...
	if n.address == "" || n.mount == "" || n.path == "" || n.token == "" {
		return false
	}
	return true
}

// TODO: after split, put this in the lower auth structure
func (n Needs) Client() (*vcg.Client, error) {
	client, err := vcg.New(
		vcg.WithAddress(n.address),
		vcg.WithRequestTimeout(n.reqTimeout),
	)
	if err != nil {
		return nil, err
	}
	if err := client.SetToken(n.token); err != nil {
		return nil, err
	}
	return client, nil
}
