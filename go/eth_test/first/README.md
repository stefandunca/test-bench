# Dev info

to debug the output form `anvil`

```bash
anvil --port 8545

# Run tests in the loop
nodemon --ext "*.go" --exec 'sh -c "go test -v ./*.go" || exit 1'

# Run specific tests in the loop
nodemon --ext "*.go" --exec 'sh -c "go test -v ./*.go -run <TestNameOrFilter>" || exit 1'
```
