# Engine 

Engine is the monolith web backend powering October, a visual and pseudonymous social network designed for the attention economy. Read more about it [here](https://github.com/october93/october).

## Requirements
- Go
- Docker

## Setup

Docker Compose starts all necessary services Engine depends on and the run scripts builds and launches Engine.

```
go get github.com/vektah/gorunpkg
go get github.com/vektah/dataloaden
docker-compose up
./scripts/run.sh
```
