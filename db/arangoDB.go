package db

import (
	arangodb "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"log"
)

func ConnectToArangoDB(endpointURL string, username string, password string) (arangodb.Client, error) {

	// Create an ArangoDB client
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{endpointURL},
	})
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
	}

	client, err := arangodb.NewClient(arangodb.ClientConfig{
		Connection:     conn,
		Authentication: arangodb.BasicAuthentication(username, password),
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	log.Println("Connected to ArangoDB")

	return client, nil

}
