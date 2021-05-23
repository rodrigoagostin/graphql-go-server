package graph

import "github.com/rodrigagostin/graphql-server/graph/model"

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	videos []*model.Video
}
