package auth

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type CognitoClient struct {
	AppClientID string
	*cip.Client
}

func Init() *CognitoClient {
	// Load the shared AWS config (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return &CognitoClient{
		os.Getenv("COGNITO_APP_CLIENT_ID"),
		cip.NewFromConfig(cfg),
	}
}
