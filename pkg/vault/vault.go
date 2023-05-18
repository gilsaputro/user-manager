package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

const secret_path = "secret/data/config"

// Client is a list dependencies for vault package
type Client struct {
	vault *api.Client
}

// VaultMethod is list method for vault package
type VaultMethod interface {
	GetConfig() (map[string]string, error)
}

// NewVaultClient func to init vault and return the VaultMethod interface
func NewVaultClient(token, address string) (VaultMethod, error) {
	if len(token) <= 0 {
		return nil, fmt.Errorf("Error: Vault Token is invalid")
	}

	cfg := api.DefaultConfig()
	cfg.Address = address

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	client.SetToken(token)
	return &Client{
		vault: client,
	}, nil
}

// GetConfig is func to get config on secret
func (c *Client) GetConfig() (map[string]string, error) {
	// Get Secret in Vault
	res, err := readSecretFromPath(c.vault)
	if err != nil {
		fmt.Println("Error reading secret: ", err)
		return nil, err
	}
	// Parse as map string
	data := res.Data["data"].(map[string]interface{})
	secretMap := make(map[string]string)
	for key, value := range data {
		secretMap[key] = value.(string)
	}
	return secretMap, nil
}

func readSecretFromPath(vault *api.Client) (*api.Secret, error) {
	return vault.Logical().Read(secret_path)
}
