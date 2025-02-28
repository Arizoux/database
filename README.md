What this Project is about:

This project is a simple file-based Key-Value-Database, which is based on the guide book "Build your own database from scratch in Go part 2" by James Smith. 

It will contain the following:
- B+Tree with an own node format, inserting / deletion functions with copy-on-write for crash resistance, as well as tests for the B+Tree
- Key-Value store with a copy-on-write B+Tree backed by a file
- A free-list for space management of disk pages, so that used pages in a KV-store are recycled
- Concurrency transactions

More information soon to follow.
