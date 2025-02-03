package executor

import (
	"fmt"
	"log"
	"github.com/docker/docker/client"
)

var Client *client.Client

func GetNewExecutorClient() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	Client = cli
	return nil
}

func GetClient() (*client.Client, error) {
	if Client != nil {
		return Client, nil
	}

	return nil, fmt.Errorf("executor - client not initialized properly !! ")
}

func CloseClient() {

	err := RemoveAllContainers(Client)
	if err != nil {
		log.Println("error : CloseClient() failed to stop containers : ", err)
	}
	err = Client.Close()
	if err != nil {
		log.Println("error : CloseClient() :  failed to close client : ", err)
	}

}
