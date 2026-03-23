## 1. Dockerfile

- [x] 1.1 Change `mise install java@temurin-17 java@temurin-21 java@temurin-25` to `mise install java@17 java@21 java@25`
- [x] 1.2 Change `mise use --global java@temurin-21` to `mise use --global java@21`

## 2. Entrypoint

- [x] 2.1 Replace the `case` statement with a single `mise use --global java@"${ASYLUM_JAVA_VERSION}"` for all versions

## 3. Verification

- [x] 3.1 Run `go test ./...` and `go vet ./...`
- [x] 3.2 Manual test: build image with `--rebuild`, verify `java -version` shows Java 21
- [x] 3.3 Manual test: project with `.tool-versions` containing `java 25` — no "missing" warning
