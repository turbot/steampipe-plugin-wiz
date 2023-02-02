package api

import (
	"context"

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
		Token:   config.ApiToken,
		Graphql: graphql.NewClient(types.StringValue(config.Url)),
	}, nil
}
