package secret

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"

	vcg "github.com/hashicorp/vault-client-go"
	vcgSchema "github.com/hashicorp/vault-client-go/schema"
)

const (
	defaultRequestTimeout     = 30 * time.Second
	errstrDeclaredContentType = "attribute %q has declared content type %q but received %q"
	errstrContentNil          = "attribute %q has empty (nil) content"
	errstrContentLenNotOne    = "attribute %q expects one content item, received: %d"
)

func NewNeeds(k8sToken *string) Needs {
	return Needs{
		k8sToken:   k8sToken,
		reqTimeout: defaultRequestTimeout,
	}
}

type Needs struct {
	// do we have kubernetes Vault token stored in pod?
	k8sToken *string

	// common attribute values
	address    string
	mount      string
	path       string
	reqTimeout time.Duration
	credType   string

	// credential type specific attribute values
	roleID, roleSecret string
	role, jwt          string
}

func (n *Needs) Process(ctx context.Context, vaultAttrs, secretAttrs *[]sm.RequestAttribute) error {
	attrs := []sm.RequestAttribute{}
	if vaultAttrs != nil {
		attrs = append(attrs, *vaultAttrs...)
	}
	if secretAttrs != nil {
		attrs = append(attrs, *secretAttrs...)
	}

	for _, cpy := range attrs {
		// Invariant: Decision was made that we'll only accept v3 attributes
		var err error
		var attr sm.RequestAttributeV3

		if attr, err = cpy.AsRequestAttributeV3(); err != nil {
			return fmt.Errorf("unmarshalling RequestAttribute into RequestAttributeV3 failed: %w", err)
		}

		switch attr.Uuid.String() {
		case vaultManagementRole.Uuid:
			if n.role, err = strContentTypeDataAttrSingle(vaultManagementRole, attr); err != nil {
				return err
			}

		case vaultManagementJwt.Uuid:
			if n.jwt, err = resourceSecretContentTypeDataAttrSingle(vaultManagementJwt, attr); err != nil {
				return err
			}

		case vaultManagementCredentialType.Uuid:
			if n.credType, err = strContentTypeDataAttrSingle(vaultManagementCredentialType, attr); err != nil {
				return err
			}

		case vaultManagementRoleID.Uuid:
			if n.roleID, err = resourceSecretContentTypeDataAttrSingle(vaultManagementRoleID, attr); err != nil {
				return err
			}

		case vaultManagementRoleSecret.Uuid:
			if n.roleSecret, err = resourceSecretContentTypeDataAttrSingle(vaultManagementRoleSecret, attr); err != nil {
				return err
			}

		case vaultManagementURI.Uuid:
			if n.address, err = strContentTypeDataAttrSingle(vaultManagementURI, attr); err != nil {
				return err
			}

		case vaultManagementRequestTmout.Uuid:
			var i int
			if i, err = intContentTypeDataAttrSingle(vaultManagementRequestTmout, attr); err != nil {
				return err
			}
			n.reqTimeout = time.Duration(i) * time.Second

		case vaultManagementMount.Uuid:
			if n.mount, err = strContentTypeDataAttrSingle(vaultManagementMount, attr); err != nil {
				return err
			}

		case vaultManagementPath.Uuid:
			if n.path, err = strContentTypeDataAttrSingle(vaultManagementPath, attr); err != nil {
				return err
			}

		default:
			slog.Debug("Unknown RequestAttributeV3 encountered.", slog.String("uuid", attr.Uuid.String()), slog.String("name", attr.Name))
		}
	}

	return nil
}

func strContentTypeDataAttrSingle(ptrn sm.DataAttributeV3, recv sm.RequestAttributeV3) (string, error) {
	if recv.ContentType != ptrn.ContentType {
		return "", fmt.Errorf(errstrDeclaredContentType, ptrn.Uuid, ptrn.ContentType, recv.ContentType)
	}
	if recv.Content == nil {
		return "", fmt.Errorf(errstrContentNil, ptrn.Uuid)
	}
	if len(*recv.Content) != 1 {
		return "", fmt.Errorf(errstrContentLenNotOne, ptrn.Uuid, len(*recv.Content))
	}
	strAttr, err := (*recv.Content)[0].AsStringAttributeContentV3()
	if err != nil {
		return "", fmt.Errorf("unmarshalling BaseAttributeContentDtoV3 into StringAttributeContentV3 failed for attribute %q: %w", ptrn.Uuid, err)
	}
	return strAttr.Data, nil
}

func intContentTypeDataAttrSingle(ptrn sm.DataAttributeV3, recv sm.RequestAttributeV3) (int, error) {
	if recv.ContentType != ptrn.ContentType {
		return 0, fmt.Errorf(errstrDeclaredContentType, ptrn.Uuid, ptrn.ContentType, recv.ContentType)
	}
	if recv.Content == nil {
		return 0, fmt.Errorf(errstrContentNil, ptrn.Uuid)
	}
	if len(*recv.Content) != 1 {
		return 0, fmt.Errorf(errstrContentLenNotOne, ptrn.Uuid, len(*recv.Content))
	}
	intAttr, err := (*recv.Content)[0].AsIntegerAttributeContentV3()
	if err != nil {
		return 0, fmt.Errorf("unmarshalling BaseAttributeContentDtoV3 into IntegerAttributeContentV3 failed for attribute %q: %w", ptrn.Uuid, err)
	}
	return int(intAttr.Data), nil
}

func resourceSecretContentTypeDataAttrSingle(ptrn sm.DataAttributeV3, recv sm.RequestAttributeV3) (string, error) {
	if recv.ContentType != ptrn.ContentType {
		return "", fmt.Errorf(errstrDeclaredContentType, ptrn.Uuid, ptrn.ContentType, recv.ContentType)
	}
	if recv.Content == nil {
		return "", fmt.Errorf(errstrContentNil, ptrn.Uuid)
	}
	if len(*recv.Content) != 1 {
		return "", fmt.Errorf(errstrContentLenNotOne, ptrn.Uuid, len(*recv.Content))
	}
	resourceAttr, err := (*recv.Content)[0].AsResourceObjectContent()
	if err != nil {
		return "", fmt.Errorf("unmarshalling BaseAttributeContentDtoV3 into ResourceObjectContent failed for attribute %q: %w", recv.Uuid, err)
	}
	if resourceAttr.ContentType != sm.AttributeContentTypeResource {
		return "", fmt.Errorf("content item has wrong content type %q, expected %q, attribute %q", resourceAttr.ContentType, sm.AttributeContentTypeResource, recv.Uuid)
	}

	secret, err := resourceAttr.Data.AsResourceSecretContentData()
	if err != nil {
		return "", fmt.Errorf("unmarshalling data of ResourceObjectContent failed for attribute %q: %w", recv.Uuid, err)
	}
	if secret.Resource != sm.Secrets {
		return "", fmt.Errorf("resource field has wrong value for attribute %q, expected %q, got %q", recv.Uuid, sm.Secrets, secret.Resource)
	}
	if secret.Content == nil {
		return "", fmt.Errorf("content field is nil, attribute %q", recv.Uuid)
	}

	secretType, err := secret.Content.Discriminator()
	if err != nil {
		return "", fmt.Errorf("unmarshalling discriminator field for SecretContent failed, attribute %q: %w", recv.Uuid, err)
	}
	switch sm.SecretType(secretType) {
	case sm.ApiKey:
		apiKeyContent, err := sm.GetApiKeySecretContent(*secret.Content)
		if err != nil {
			return "", fmt.Errorf("attribute %q: %w", recv.Uuid, err)
		}
		return apiKeyContent, nil

	case sm.BasicAuth:
		_, password, err := sm.GetBasicAuthSecretContent(*secret.Content)
		if err != nil {
			return "", fmt.Errorf("attribute %q: %w", recv.Uuid, err)
		}
		return password, nil

	case sm.Generic:
		content, err := sm.GetGenericSecretContent(*secret.Content)
		if err != nil {
			return "", fmt.Errorf("attribute %q: %w", recv.Uuid, err)
		}
		return content, nil

	case sm.JwtToken:
		content, err := sm.GetJwtTokenSecretContent(*secret.Content)
		if err != nil {
			return "", fmt.Errorf("attribute %q: %w", recv.Uuid, err)
		}
		return string(content), nil

	case sm.KeyValue:
		content, err := sm.GetKeyValueSecretContent(*secret.Content)
		if err != nil {
			return "", fmt.Errorf("attribute %q: %w", recv.Uuid, err)
		}
		var value any
		for _, value = range content {
			break
		}
		s, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("attribute %q, value in %q is not a string", recv.Uuid, secretType)
		}

		return s, nil

	case sm.SecretKey:
		decoded, err := sm.GetSecretKeySecretContent(*secret.Content)
		if err != nil {
			return "", fmt.Errorf("attribute %q: %w", recv.Uuid, err)
		}
		return string(decoded), nil

	}

	return "", fmt.Errorf("attribute %q, unexpected secret type %q", recv.Uuid, secretType)
}

func (n *Needs) CommonCheck() error {
	switch {
	case strings.TrimSpace(n.address) == "":
		return fmt.Errorf("missing attribute uuid %q, name %q", vaultManagementURI.Uuid, vaultManagementURI.Name)
	case strings.TrimSpace(n.mount) == "":
		return fmt.Errorf("missing attribute uuid %q, name %q", vaultManagementMount.Uuid, vaultManagementMount.Name)

	case strings.TrimSpace(n.credType) == "":
		return fmt.Errorf("missing attribute uuid %q, name %q", vaultManagementCredentialType.Uuid, vaultManagementCredentialType.Name)
	}

	if err := n.credTypeCheck(); err != nil {
		return err
	}

	return nil
}

func (n Needs) credTypeCheck() error {
	switch n.credType {
	case credentialTypeAppRole.Data:
		if n.roleID == "" || n.roleSecret == "" {
			return fmt.Errorf("required attributes for credential type %q missing %s(%s), %s(%s)",
				n.credType, vaultManagementRoleID.Uuid, vaultManagementRoleID.Name,
				vaultManagementRoleSecret.Uuid, vaultManagementRoleSecret.Name)
		}
	case credentialTypeJwt.Data:
		if n.role == "" || n.jwt == "" {
			return fmt.Errorf("required attributes for credential type %q missing %s(%s), %s(%s)",
				n.credType, vaultManagementRole.Uuid, vaultManagementRole.Name,
				vaultManagementJwt.Uuid, vaultManagementJwt.Name)
		}
	case credentialTypeK8s.Data:
		if n.k8sToken != nil {
			if n.role == "" {
				return fmt.Errorf("required attributes for credential type %q missing %q(%s)",
					n.credType, vaultManagementRole.Uuid, vaultManagementRole.Name)
			}
		} else {
			return fmt.Errorf("unknown credential type %q", n.credType)
		}

	default:
		return fmt.Errorf("unknown credential type %q", n.credType)
	}
	return nil
}

func (n Needs) ConnectionCheck() error {
	switch {
	case strings.TrimSpace(n.address) == "":
		return fmt.Errorf("missing attribute uuid %q, name %q", vaultManagementURI.Uuid, vaultManagementURI.Name)
	case strings.TrimSpace(n.credType) == "":
		return fmt.Errorf("missing attribute uuid %q, name %q", vaultManagementCredentialType.Uuid, vaultManagementCredentialType.Name)
	}

	if err := n.credTypeCheck(); err != nil {
		return err
	}

	return nil
}

func (n Needs) Client(ctx context.Context) (*vcg.Client, error) {

	client, err := vcg.New(
		vcg.WithAddress(n.address),
		vcg.WithRequestTimeout(n.reqTimeout),
	)
	if err != nil {
		return nil, err
	}

	switch n.credType {
	case credentialTypeAppRole.Data:
		resp, err := client.Auth.AppRoleLogin(
			ctx,
			vcgSchema.AppRoleLoginRequest{
				RoleId:   n.roleID,
				SecretId: n.roleSecret,
			},
		)
		if err != nil {
			return nil, err
		}
		if err := client.SetToken(resp.Auth.ClientToken); err != nil {
			return nil, err
		}

	case credentialTypeJwt.Data:
		resp, err := client.Auth.JwtLogin(
			ctx,
			vcgSchema.JwtLoginRequest{
				Jwt:  n.jwt,
				Role: n.role,
			},
		)
		if err != nil {
			return nil, err
		}

		if err := client.SetToken(resp.Auth.ClientToken); err != nil {
			return nil, err
		}

	case credentialTypeK8s.Data:
		if n.k8sToken == nil {
			return nil, fmt.Errorf("unknown credential type %q", credentialTypeK8s.Data)
		}
		resp, err := client.Auth.KubernetesLogin(ctx, vcgSchema.KubernetesLoginRequest{
			Jwt:  *n.k8sToken,
			Role: n.role,
		})
		if err != nil {
			return nil, err
		}
		if err := client.SetToken(resp.Auth.ClientToken); err != nil {
			return nil, err
		}
	}

	return client, nil
}
