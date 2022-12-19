package action

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"kube-metric-collector/common"
	"log"
	"os"
)

var (
	leasePath = "secret/data/cluster/"
	vaultIP   string
	vaultPort string
	roleId    string
	secretId  string
)

type Client struct {
	Token         string
	Client        *api.Client
	LeaseDuration int
}

type appRoleLogin struct {
	RoleID   string `json:"role_id"`
	SecretID string `json:"secret_id"`
}

func readProperties() {

	profile := "local"
	if len(os.Getenv("PROFILE")) > 0 {
		profile = os.Getenv("PROFILE")
	}

	if profile == "prod" {
		vaultIP = os.Getenv("VAULT_IP")
		vaultPort = os.Getenv("VAULT_PORT")
		roleId = os.Getenv("VAULT_ROLE_ID")
		secretId = os.Getenv("VAULT_SECRET_ID")
	} else {
		vaultIP = common.ConfInfo["vault.server.address"]
		vaultPort = common.ConfInfo["vault.server.port"]
		roleId = common.ConfInfo["vault.id"]
		secretId = common.ConfInfo["vault.secret.id"]
	}

	fmt.Printf("IP: %s  Port: %s role: xxxxxxxx  secret: xxxxxxxx\n", vaultIP, vaultPort)
}

func GetSecretInfo(clusterId string) (map[string]interface{}, error) {
	leasePath = "secret/data/cluster/" + clusterId
	//leasePath = leasePath + clusterId //secret/data/{cp-clusterId}
	log.Println("leasePath : ", leasePath)
	secret, err := ReadSecret(leasePath)

	if err != nil {
		//fmt.Println("ReadSecret(leasePath) :: ", err)
		return nil, err
	}

	if secret == nil {
		fmt.Println("no secret values in path")
		return nil, nil
	}

	return secret, nil
}

func ReadSecret(secretPath string) (map[string]interface{}, error) {
	client, err := NewVaultClient()
	if err != nil {
		log.Println("NewVaultClient() :: ", err)
		return nil, err
	}

	token, err := client.AuthUser()
	if err != nil {
		log.Println("client.AuthUser() :: ", err)
		return nil, err
	}

	client.Client.SetToken(token)
	secret, err := client.Client.Logical().Read(secretPath)
	if secret == nil {
		//log.Println("secret is nil(", secretPath, ")")
		err = fmt.Errorf("secret is nil(%s)", secretPath)
		return nil, err
	}
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		//log.Fatalf("Data type assertion filed: %#v", secret.Data["data"])
		return nil, err
	}

	if err != nil || secret == nil {
		return nil, err
	}

	return data, err

}

func (vaultClient *Client) AuthUser() (string, error) {
	request := vaultClient.Client.NewRequest("POST", "/v1/auth/approle/login")

	login := appRoleLogin{SecretID: secretId, RoleID: roleId}

	if err := request.SetJSONBody(login); err != nil {
		return "", err
	}

	resp, err := vaultClient.Client.RawRequest(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return "", err
	}

	return secret.Auth.ClientToken, nil
}

func NewVaultClient() (*Client, error) {
	readProperties()

	vaultClient := Client{}

	client, err := api.NewClient(&api.Config{
		Address: "http://" + vaultIP + ":" + vaultPort,
	})

	vaultClient.Client = client

	return &vaultClient, err
}
