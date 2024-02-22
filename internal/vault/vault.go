package vault

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

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
		log.Fatal(err)
		return nil, err
	}
	fmt.Println(resp.Auth.ClientToken)
	if err := client.SetToken(resp.Auth.ClientToken); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return client, nil
}

type LoginWithToken struct {
	Token string
}

func (l LoginWithToken) Login(client *vault.Client) (*vault.Client, error) {
	return nil, nil
}

type APIClientConfig struct {
	LoginMethod LoginMethod
	EndpointURL string
}

func GetAPIClient(config APIClientConfig) (*vault.Client, error) {
	// prepare a client with the given base address
	client, err := vault.New(
		vault.WithAddress(config.EndpointURL),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	return config.LoginMethod.Login(client)

}
