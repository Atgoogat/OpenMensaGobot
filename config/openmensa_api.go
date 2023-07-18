package config

import "github.com/Atgoogat/openmensarobot/openmensa"

func GetOpenmensaApi() openmensa.OpenmensaApi {
	return openmensa.NewOpenmensaApi(openmensa.OPENMENSA_API_ENDPOINT)
}
