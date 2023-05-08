package api

import (
	"context"
	"log"
	"time"

	"github.com/turbot/go-kit/types"

	"github.com/machinebox/graphql"
)

// Wiz API Client
type Client struct {
	Token   *string
	Graphql *graphql.Client
}

func CreateClient(ctx context.Context, config ClientConfig) (*Client, error) {
	return &Client{
		Token:   config.Token,
		Graphql: graphql.NewClient(types.StringValue(config.Url)),
	}, nil
}

func (client *Client) doRequest(req *graphql.Request, responseData interface{}) error {
	return client.DoRequest(req, responseData)
}

// execute graphql request
func (client *Client) DoRequest(req *graphql.Request, responseData interface{}) error {
	// set header fields
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+*client.Token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	start := time.Now()
	if err := client.Graphql.Run(ctx, req, &responseData); err != nil {
		return err
	}
	log.Println("graphql.time", time.Since(start).Milliseconds())
	return nil
}
