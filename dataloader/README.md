# Dataloader

* Concept: [dataloaders](https://github.com/facebook/dataloader)
* Go library: [dataloaden](https://github.com/vektah/dataloaden)

### Adding a dataloader

1. Add appropriate generation header in `loaders.go` (see examples in the file already)
1. `go generate dataloader/loaders.go`
1. Add new loader to the loader constructor function