# package rpc

The `rpc` package provides access to the API of Engine across a network in the form of a method set.

## Rationale

The graph database is designed with a distributed system architecture in mind. Even though a single instance will suffice for the first thousand user or more there should be a clear path for distributing it across machines. With Remote Procedure Calls (RPC) methods can be invoked transparently without necessarily knowing whether the method is going to be called locally or remotely. Most of the logic goes through [`Node`](https://github.com/october93/engine/blob/master/graph/node.go) which gives reason to distribute computation by distributing `Node` instances over different machines.

The set of methods is defined as part of the `RPC` interface. Additionally, the `RPC` interface is also the way for the clients, consumer-facing and internal tools, to communicate with the backend. This includes operations executed directly on the graph as well as operations affecting the pure domain model, i.e. user and session management, card creation, etc.

## Structure

* `protocol` implements the transport layer 
* `client` provides an exported client implementation
* `server` provides the HTTP & WebSocket logic to host all the endpoints

## Outlook

The `RPC` interface is currently flat without any hierarchy. One advantage of RPC is to offer a structured way of accessing an API by utilizing different objects with different states. Currently, RPCs which require a user take a `nodeID` parameter in order to fetch the corresponding user object. Instead, the set of methods should be subdivided into different sets, e.g. user methods, admin methods, etc.
