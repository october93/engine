# Engine 

Engine is the monolith web backend powering October, a visual and pseudonymous social network.

![Preview](https://raw.githubusercontent.com/october93/engine/master/preview.png?token=AAFBT2IRGCUPZOHAW4BQGHS5IE2D6)

## Setup

### Requirements

1. Install Go
2. Install Docker
3. Install third party dependencies

```
go get github.com/vektah/gorunpkg
go get github.com/vektah/dataloaden

```

Docker Compose starts all necessary services Engine depends on and the run scripts builds and launches Engine.

```
docker-compose up
./scripts/run.sh
```
