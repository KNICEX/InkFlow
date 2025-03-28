package ioc

import client "github.com/gorse-io/gorse-go"

func InitGorseCli() *client.GorseClient {
	return client.NewGorseClient("localhost:8080", "")
}
