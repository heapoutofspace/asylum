## 1. Switch cache mounts to named volumes

- [x] 1.1 Replace bind mount loop in `appendVolumes` with named volume mounts (`--mount type=volume,src=<cname>-cache-<tool>,dst=<path>`)
- [x] 1.2 Remove `os.MkdirAll` calls for `~/.asylum/cache/<container>/` host directories
- [x] 1.3 Update container tests to expect named volume mounts instead of bind mounts

## 2. Verification

- [x] 2.1 Run `go test ./...` and `go vet ./...`
- [x] 2.2 Manual test: start container, verify cache volumes created (`docker volume ls | grep cache`)
- [x] 2.3 Manual test: `asylum --cleanup` removes cache volumes
