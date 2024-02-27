package vault

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/db"
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

func getLoginMethod(authority db.AuthorityInstance) LoginMethod {
	return AppRoleLogin{
		RoleId:   "60370fcc-1c96-6ca7-ea41-d92736def91a",
		SecretId: "6a6c7b86-5551-4849-e119-812da8086fcc",
	}
}

func GetClient(authority db.AuthorityInstance) (*vault.Client, error) {
	client, err := vault.New(
		vault.WithAddress("https://vault.czertainly.online:443"), // prepare a client with the given base address
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	return getLoginMethod(authority).Login(client)
}


