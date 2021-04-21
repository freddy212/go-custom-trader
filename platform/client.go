package platform

import (
	coinbasepro "github.com/freddy212/go-coinbasepro"
)

var clientInstance *coinbasepro.Client

func GetClientInstance() *coinbasepro.Client {
	if clientInstance == nil {
		clientInstance = coinbasepro.NewClient()
	}
	return clientInstance
}
