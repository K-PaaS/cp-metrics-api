package action

import (
	"fmt"
)

func GetProviderSecretData(clusterId string) (map[string]interface{}, error) {
	iaasAuth, err := GetSecretInfo(clusterId)
	if err != nil {
		fmt.Println("GetSecretInfo(providerName) :: ", err)
	}

	return iaasAuth, err
}
