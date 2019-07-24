# package store

The `store` package provides the persistence layer. It is subdivided into persistence for the graph database and persistence for the domain models (`model` package) which are concerned with the user-facing part of the product.

## Rationale

Posting content on the October network or reacting to content has a mulitude of effects on the algorithmic level. One operation could potentially affect all the edges a single node is connected to. Assuming a graph with 1,000 nodes where every node is connect to every other node (complete graph) there are 1000 * (1000 - 1) / 2 = 499500 edges.

Triggering 500,000 writes on a database just for a single operation will quite quickly lead to an overload on the instance hosting the database. Instead we use a snapshot approach: the in-memort graph is persisted to disk. Individual nodes and their edges are written to snapahot files using gob for speed and type persistence.

Since writing snapshots on every update is expensive as well, we would use [**Command Logging**](https://docs.voltdb.com/UsingVoltDB/ChapCmdLog.php). Operations (commands) mutating the state are written to a log file. Every time an operation was succesful, a new entry is appended to the log. The log entry contains all the necessary data to reproduce the operation fully. When a snapshot is loaded, all operations are replayed from the command log.

Since the log entries only affect the node in question and each node stores its edges redundantly there are no problems related to concurrent execution. Due to this fact, there is also an absolute ordering on the log entries by time.
