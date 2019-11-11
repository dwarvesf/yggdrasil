package toolkit

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

func NewVaultClient() (*api.Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: os.Getenv("VAULT_ADDR"),
	})
	if err != nil {
		return nil, err
	}
	client.SetToken(os.Getenv("VAULT_TOKEN"))
	return client, nil
}

func GetVaultValueFromKey(v *api.Client, k string) (string, error) {
	val, err := v.Logical().Read(os.Getenv("VAULT_PATH"))
	if err != nil {
		return "", err
	}
	m := val.Data["data"].(map[string]interface{})
	return fmt.Sprintf("%v", m[k]), nil
}
