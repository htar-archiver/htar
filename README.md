htar
====
A hashed, encoding only, partitioning implementation of tar in golang.

## Development

```bash
# run all unit tests
go test ./...
# test writing files to pipe
go build && ./htar pack --size 1TB --root /go/src/htar --pipe "mbuffer -R 10mb -o /dev/null" ..:1
```
