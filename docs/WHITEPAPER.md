# GO-KARAI

Reference client implementation and Karai network spec, written in Go.

### What is Karai?

Karai is a polyforest directed acyclic graph scaling solution and application substrate for modern cryptoasset networks and distributed applications.

### Purpose

The go-karai node provides a means for a user to establish, link and interact with Karai transaction channels.

The role of the go-karai node is to receive, arrange, and append transactions to the transaction channel graph. Channels are discovered by storing a pointer reference on TurtleCoin which functions as a hyperlink between the two networks. Once a user or application decodes the pointer and has connected to a Karai transaction channel, the user or application then becomes a participating member of the p2p operations of the channel.

## Technical Notes

P2P peer-addressing is handled by libp2p, with content-addressing provided by IPFS. The P2P connection is created with TLS, QUIC, and NAT Traversal enabled by default. Golang was chosen for this implementation out of preference for a fast and easy to learn compiled language that could produce a cross platform binary without the use of proprietary software and development environments.

### Transactions

Transactions are raw JSON documents similar to this example:

```json
{
    "tx_type": 2, // tx_type 2 is a normal transaction
    "tx_hash": "5f0433c98576b29a09a8e43e82fdb1f7df1a4895b1c05bddc10c717630c38ba6",
    "tx_prev": "26b858334a5d3a5e28a21af762aed33f822ec57d45bdf40368408f9a0bbad08c",
    "tx_data": { "tx_slot": 3 }
}
```

In this example a basic transaction has 4 fields, the transaction type, the current transaction hash, the previous transaction hash, and any extra data included in the transaction - in this case the transaction contains a JSON object record of its original graph slot location.

### Milestones

Governance changes are passed by the channel coordinator in what are known as 'milestones' which function similar to fork heights in conventional networks like TurtleCoin. These changes can reconfigure the way a transaction channel arranges transactions, responds to bursts of traffic, or even update the ownership of the channel. A channel is always created with 1 Root transaction, with the first milestone transaction always following immediately after.

Milestones are raw JSON documents that modify governance variables for the channel similar to this example:

```json
{
    "tx_type": 1, // tx_type 1 is a milestone transaction
    "tx_hash": "5477e3d072e206770cb2530e116ca2338b24537f113a5ec9d421cfa32a596ab3",
    "tx_prev": "96b94244d07d046eb7a153d44e3a0db41a52290436ccbfc924b4b5ae206ddb84",
    "tx_data": {
        "chan_milestone": 0,
        "channel_params": {
            "name": "karai-transaction-channel",
            "public": true,
            "n_txspread": 1,
            "n_interval": 30
        },
        "coord_params": {
            "majorSemver": 0,
            "minorSemver": 2,
            "patchSemver": 1,
            "address": "TRTLuxMpUNTBqfWshbc65E7Yqx17rQpHQZC5HyANaTL37AQm2fRsNDXG37jxPhXXa5NMJVLFJpQa9iQn9Se87VNuWwPHWScoZLY"
        },
        "peer_data": {
            "ip4": "123.34.45.56",
            "ifps_peer_id": "QmfMAVSMyw4T9Mu3y8hM4phLWbM5NhdYuq5HRmKy8kX3SD",
            "ipfs_pub_key": "CAASpgIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDI3gvSHQ/V3o7wWLp+KLw8w4k74JGF7+lxPAK0Z6SAp2CELvr+FJfflcIAnOna5NekFj3oZhgI3sTAMRixn802S+OUmBuFrtdxd8SX1PjwCmdzm+xTWU8IdrZbxzeHY/n4i34ZyOEybdWEvR4oExplxTk9mnZmKvZvIH3lCQIbfkhoJFTB4D4R5KG5YcEQ6/2hLvzdoMyUcVZRf7dRxWUyoXRdE5810tsCBECrRzLX9nWERP/ki4elvJlDQYU5bHUazZy4tbl9kEbP28gjm9XGYKxjAWyXG+uMZoohCujSNN3SQzo/5zE4VWzi4LC01ourl8xR9pd5HhzH1oKcYBoZAgMBAAE="
        }
    }
}
```

The go-karai program generates a pointer that links a Karai transaction channel to the TurtleCoin Network. These pointers can be scanned for on the TurtleCoin Network by looking for an encoded address in the `tx_data` field. Each pointer is generated in Hex format as well as Ascii format for human readable output. All pointers get stored and scanned as hex. Ascii pointers are for convenience only, but generally have the format of `ktx(ip:port)`.

Karai Transaction Channel pointers look like the pointer in this example:

```
Hex:    6b7478287b22f12c114729
Ascii:  ktx(123.34.241.44:4423)
```

## Coordinating Transactions

Karai designates the founding node of a transaction channel as the Channel Coordinator. The CC functions as an authority for timestamping and arranging transactions. The graph can be operated in two modes regarding graph construction:

-   In **'linear mode'** transactions arranged in the typical order follow one another in a single timeline (zero arrangement).
-   In **'wave mode'**, the transactions take on a new role depending on their position in a multipath time limited subgraph. (varied arrangement)

Karai differs from conventional networks in that there is no mining, and no empty space being created when there are no transactions to be appended. In wave mode, after channel creation or after a period of inactivity, a new transaction arriving will trigger a 'listen interval' that gathers subsequent transactions. This interval is a key variable in subgraph formation, but can be disregarded entirely in linear mode. While the interval is open, a transaction still can be determined to be valid thought its final position in the subgraph can be variable.

Note: At the time of this writing, positional transaction metadata is non-hashed data, meaning there is a degree of maleability with regards to tip selection and final transaction position in the subgraph.

## Subgraph Wave construction

A subgraph is a directed rooted tree made of type 2 transactions gathered in the same listen interval that have differing roles depending on their position.

The subgraph root is known as the 'wave tip'. This special role is assigned to the first transaction in a listen interval. Upon initiating the subgraph, this first transaction accumulates peer metadata for the peer ids that own a transaction in the subgraph.

As a bonus to scalability over linear mode configurations, the subgraph tip metadata also aids in non-linear edge-tracing of the graph by allowing scanning for owned transactions to occur from milestone to wavetip rather than traversing every transaction, decreasing the time for a peer to reach a fully synced state.

### Example Wave Construction

A typical rooted tree arrangement of transactions in the subgraph:

0. Individual transactions are timestamped and assigned slots in rows of the wave lattice.
1. (Wave tip [1]): `G[[0]1], tx (1A)` (First transaction enters slot A, stream width = 1) 9 tx remain
1. Wave row [2]: `G[[0]2], tx 2A->(1A)`, `tx 2B->(1A)` (Next 2 tx enter slots A+B, stream width = 2) 7 tx remain
1. Wave row [3]: `G[[0]3], tx 3A->(2A)`, `tx 3B->(2A & 2B)`, `tx 3C->(2B)` (Next 3 tx enter slots A+B+C, stream width = 3) 4 tx remain
1. Wave row [4]: `G[[0]4], tx 4A->(3A)`, `tx 4B->(3A & 3B)`, `tx 4C->(3B & 3C)`, `tx 4D->(3C)` (Final 4 tx enter slots A+B+C+D, stream width = 4) No tx remain

### Wave Lattice

Subgraph ordering is determined by a simple expanding lattice tip selection algorithm which was designed to account for transaction waves needing to receive and arrange an indeterminate number of vertices during a listen interval. In simple terms with this default lattice tip selection configuration, a wave tip is the parent node for subsequent transactions with the wave tip row being [1] in the lattice. Subsequent rows add +1 children per row until the listen interval has expired. All rows but the final row will have two child nodes per parent node.
Note: Lattice position at this time is a non-hashed element rendering position somewhat maleable, but this could change.

#### Graph Proofs

A graph proof is a modified Karai transaction channel pointer that can be stored on the TurtleCoin network that contains a channel address and a signed historical hash of their Karai transaction channel graph. This Graph Proof serves as a means to immutably notarize a Karai transaction channel at a certain point in history as untampered.

A graph proof inherits the elements of a typical channel pointer but includes a Channel Coordinator signed hash of a matching hash coordinate where the hash can be verified on the corresponding Karai transaction channel.
