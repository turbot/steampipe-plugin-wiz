package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Vanta API Client
type Client struct {
	Token   *string
	Graphql *graphql.Client
}

func CreateClient(ctx context.Context, config ClientConfig) (*Client, error) {
	return &Client{
		Token:   config.ApiToken,
		Graphql: graphql.NewClient("https://api.us8.app.wiz.io/graphql"),
	}, nil
}
