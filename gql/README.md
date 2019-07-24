# GraphQL

This API is currently only used for admin. [Learn GraphQL](https://graphql.org/learn/)

## Changing the API

1. Make appropriate changes to schema.graphql
1. [conditional] Add any predefined types you are using in the schema in `types.json`
1. `go generate gql/graphql.go` to generate the new interface (located in `generated.go`)
1. Update and implement any interfaces changes from `generated.go`

[Read more about `gqlgen`](https://gqlgen.com/)

## Performance

Because of GraphQLs heavily nested nature, retrieving data must be loosely coupled. However, this could
potentially result in multiplying the query count by `n` every nest. To alleviate this, [dataloaders](https://github.com/facebook/dataloader)
are used to batch requests.

We use [dataloaden](https://github.com/vektah/dataloaden) to build our dataloaders using codegen in go.

See [dataloader](/dataloader) package for more specifics.
