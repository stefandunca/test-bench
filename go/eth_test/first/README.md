# Dev info

```bash
python3 first/test.py
```

or to debug the output form `anvil`

```bash
anvil --port 8545

# Run tests in the loop
nodemon --ext "*.go" --exec 'sh -c "go test -v ./*.go" || exit 1'

# Run specific tests in the loop
nodemon --ext "*.go" --exec 'sh -c "go test -v ./*.go -run <TestNameOrFilter>" || exit 1'
```

## Docs

Examples

- Intro to the `newTx``: https://github.com/MariusVanDerWijden/web3go
- More examples: https://geth.ethereum.org/docs/developers/dapp-developer/native