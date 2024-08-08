Minimal reproduction of an issue that hit Warden Prorotocol testnet
(buenavista) during the upgrade to v0.4.0.

Part of a migration of the database was inserting new key-value pairs into the
chain state, in a non-deterministic order.

Even if the end data is the same, the shape of the IAVL tree is different,
leading to different hashes and finally, to the nefarious CONSENSUS FAILURE!!!
apphash mismatch errors.


## Run this repo

This Go program creates a new go-leveldb database (`test.db`) and inserts a
bunch of key-value pairs into it, taken from a Go map. Iterating over a map in
Go is not deterministic, so the order of the key-value pairs is not.

The program then commits the IAVL tree to the database, and prints the hash of
the resulting tree.

```bash
❯ rm -rf test.db/ && go run .
Hash d23ce4210812193b4caa7ce9e393ffbd40888af24a8373856f2858662d2c6c14 Version 1

❯ rm -rf test.db/ && go run .
Hash ad64363b6eabef8191ab7de4d4d999a5ba99bd4bfed1f9c6c9e6a47c67b25df4 Version 1

❯ rm -rf test.db/ && go run .
Hash a5dd0899756ff398416d6b4c455f627c74017b967d484cdf06bf8ffa311b517b Version 1

❯ rm -rf test.db/ && go run .
Hash 7a0d38e100eee064cb98908f1f33f01872376eeee6187e41bc76dfe8fe72706a Version 1

❯ rm -rf test.db/ && go run .
Hash ef9b8d260aded3af5a17faac5a8106b44eb06c77113f25723734579e05ff08d6 Version 1
```


## Identifying this issue

A tool called `iaviewer`, developed by Ethan Frey and part of the IAVL library
can be used to inspect the IAVL tree.

More instructions can be found in its README:
https://github.com/cosmos/iavl/tree/master/cmd/iaviewer.

In this example, we can use it like this:

```
❯ rm -rf test.db/ && go run . && mv test.db test-1.db

❯ rm -rf test.db/ && go run . && mv test.db test-2.db

❯ go run github.com/cosmos/iavl/cmd/iaviewer@v1.2.0 data test-1.db/ "" > dump1.txt

❯ go run github.com/cosmos/iavl/cmd/iaviewer@v1.2.0 data test-2.db/ "" > dump2.txt

❯ diff dump1.txt dump2.txt
2003c2003
< Hash: EF9B8D260ADED3AF5A17FAAC5A8106B44EB06C77113F25723734579E05FF08D6
---
> Hash: 143A56967C6D21ECF4B9B748E03F60D3E8CC816F8F76DAB711B71D04B4DC0BBF
```

The only difference between the two dumps is the hash of the IAVL tree. This
means that the data is the same, but the shape of the IAVL tree is different.
