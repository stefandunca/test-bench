# Dev info

```bash
python3 first/test.py

# Or

ganache-cli --networkId 1337 -m "much repair shock carbon improve miss forget sock include bullet interest solution"

# Run tests in the loop with data file
nodemon --ext "*.go" --exec 'TEST_DATA_FILE="/Users/stefan/proj/test-bench/go/eth_test/first/test_data.json" sh -c "go test -v ./*.go" || exit 1'

# Run tests in the loop without data file
nodemon --ext "*.go" --exec 'sh -c "go test -v ./*.go" || exit 1'

# Run specific tests in the loop
nodemon --ext "*.go" --exec 'sh -c "go test -v ./*.go -run <TestNameOrFilter>" || exit 1'
```

## Docs

Examples

- Intro to the `newTx``: https://github.com/MariusVanDerWijden/web3go
- More examples: https://geth.ethereum.org/docs/developers/dapp-developer/native