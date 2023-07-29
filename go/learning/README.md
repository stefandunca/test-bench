# Info

Run tests in a loop

```sh
nodemon --ext "*.go" --exec 'go test -v ./go/learning/stuff/*.go 2>&1 | tee .dev/go-tests.log || exit 1'
```
