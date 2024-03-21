package vault

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/db"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/logger"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

var log = logger.Get()

const DEFAULT_K8S_TOKEN_PATH = "/var/run/secrets/kubernetes.io/serviceaccount/token"

type LoginMethod interface {
	Login(client *vault.Client) (*vault.Client, error)
}

type AppRoleLogin struct {
	RoleId   string
	SecretId string
}

func (l AppRoleLogin) Login(client *vault.Client) (*vault.Client, error) {
	ctx := context.Background()
	resp, err := client.Auth.AppRoleLogin(
		ctx,
		schema.AppRoleLoginRequest{
			RoleId:   l.RoleId,
			SecretId: l.SecretId,
		},
		//vault.WithMountPath("my/approle/path"), // optional, defaults to "approle"
	)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	//fmt.Println(resp.Auth.ClientToken)
	if err := client.SetToken(resp.Auth.ClientToken); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return client, nil
}

type LoginWithToken struct{}

func (l LoginWithToken) Login(client *vault.Client) (*vault.Client, error) {
	ctx := context.Background()
	token, err := os.ReadFile(DEFAULT_K8S_TOKEN_PATH) // Replace with your actual file path
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	jwt := string(token)
	authInfo, err := client.Auth.JwtLogin(ctx, schema.JwtLoginRequest{
		Jwt:  jwt,
		Role: "czertainly-role",
	}, vault.WithMountPath("jwt"))
	if err != nil {
		return nil, fmt.Errorf("unable to log in with JWT auth: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after JWT login")
	}

	err = client.SetToken(authInfo.Auth.ClientToken)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type LoginWithK8sToken struct {
}

func (l LoginWithK8sToken) Login(client *vault.Client) (*vault.Client, error) {
	ctx := context.Background()
	token, err := os.ReadFile(DEFAULT_K8S_TOKEN_PATH) // Replace with your actual file path
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	jwt := string(token)
	authInfo, err := client.Auth.KubernetesLogin(ctx, schema.KubernetesLoginRequest{
		Jwt:  jwt,
		Role: "czertainly-role",
	}, vault.WithMountPath("kubernetes"))
	if err != nil {
		return nil, fmt.Errorf("unable to log in with Kubernetes auth: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after K8s login")
	}

	err = client.SetToken(authInfo.Auth.ClientToken)
	if err != nil {
		return nil, err
	}
	return client, nil
}
func getLoginMethod(authority db.AuthorityInstance) LoginMethod {
	switch authority.CredentialType {
	case model.JWTOIDC_CRED:
		return LoginWithToken{}
	case model.KUBERNETES_CRED:
		return LoginWithK8sToken{}
	case model.APPROLE_CRED:
		return AppRoleLogin{
			RoleId:   authority.RoleId,
			SecretId: authority.RoleSecret,
		}

	}
	return nil
}

func GetClient(authority db.AuthorityInstance) (*vault.Client, error) {
	client, err := vault.New(
		vault.WithAddress(authority.URL),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		log.Error(err.Error())
	}
	return getLoginMethod(authority).Login(client)
}
