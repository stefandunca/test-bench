package first

import (
	"context"
	"crypto/ecdsa"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"

	"github.com/stretchr/testify/require"
)

const anvilDefaultEthBalance = 10000

func TestConvertAddresses(t *testing.T) {
	strAddress := "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"

	address := common.HexToAddress(strAddress)
	require.Equal(t, strAddress, address.Hex())
}

func TestGetBalanceLastBlock(t *testing.T) {
	client, _, anvil, tearDown := testClientWithBlocks()
	defer tearDown()

	addresses, err := anvil.AvailableAddresses()
	require.NoError(t, err)
	balance, err := client.BalanceAt(context.Background(), addresses[0], nil)
	require.NoError(t, err)

	tenThousandsEthAsWei := ethToWei(anvilDefaultEthBalance)
	require.Equal(t, 0, tenThousandsEthAsWei.Cmp(balance), "balance should be 10000 ETH")
}

func TestGetBalanceFirstBlock(t *testing.T) {
	client, _, anvil, tearDown := testClientWithBlocks()
	defer tearDown()

	addresses, err := anvil.AvailableAddresses()
	require.NoError(t, err)
	balance, err := client.BalanceAt(context.Background(), addresses[0], big.NewInt(0))
	require.NoError(t, err)

	require.Equalf(t, 10000.0, balanceToEther(balance), "balance should be 10000 ETH")
}

func TestHeaderByNumberLast(t *testing.T) {
	client, _, _, tearDown := testClientWithBlocks()
	defer tearDown()

	lastHeader, err := client.HeaderByNumber(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, lastHeader)
	require.Greater(t, lastHeader.Number.Cmp(big.NewInt(0)), 0, "last block should be greater than genesis block (0)")

	firstHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	require.NoError(t, err)
	require.NotNil(t, firstHeader)
	require.Equal(t, int64(0), firstHeader.Number.Int64())
	require.Greater(t, lastHeader.Time, firstHeader.Time, "last timestamp should be greater than first timestamp")
}

func TestBlockByNumber(t *testing.T) {
	client, _, _, tearDown := testClientWithBlocks()
	defer tearDown()

	lastBlock, err := client.BlockByNumber(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, lastBlock)
	require.Greater(t, lastBlock.Number().Cmp(big.NewInt(0)), 0, "last block should be greater than genesis block (0)")

	firstBlock, err := client.BlockByNumber(context.Background(), big.NewInt(0))
	require.NoError(t, err)
	require.NotNil(t, firstBlock)
	require.Equal(t, int64(0), firstBlock.Number().Int64())

	require.Greater(t, lastBlock.Time(), firstBlock.Time(), "last timestamp should be greater than first timestamp")
	require.Equal(t, 0, len(lastBlock.Transactions()), "no transactions expected, mined empty blocks")
	require.True(t, strings.Contains(strings.ToLower(lastBlock.Hash().Hex()), "0x"))
	// Anvil is a PoS network, so the difficulty is 0
	require.Equal(t, uint64(0), lastBlock.Difficulty().Uint64())
}

func TestGenerateNewWallet(t *testing.T) {
	_, testData, _, tearDown := testClient()
	defer tearDown()

	// Generate a random private key
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	require.Equal(t, privateKey.D.Cmp(big.NewInt(0)), 1, "private key should be greater than 0")

	privateKeyBytes := crypto.FromECDSA(privateKey)
	require.Equal(t, 32, len(privateKeyBytes), "Converted private key should be 32 bytes long")

	privateKeyHexStrWithPrefix := hexutil.Encode(privateKeyBytes)
	require.True(t, strings.Contains(privateKeyHexStrWithPrefix, "0x"), "Converted private key have 0x prefix")

	privateKeyStrNoPrefix := strings.TrimPrefix(privateKeyHexStrWithPrefix, "0x")
	require.False(t, strings.Contains(privateKeyStrNoPrefix, "0x"), "This is the private key which is used for signing transactions")

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok, "public key is of type *ecdsa.PublicKey")
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	require.Equal(t, 65, len(publicKeyBytes), "Converted public key should be 64 bytes long")

	publicKeyHexStrWithHexAndECPrefix := hexutil.Encode(publicKeyBytes)
	require.True(t, strings.Contains(publicKeyHexStrWithHexAndECPrefix, "0x04"), "Converted public key should have 0x04 prefix")

	publicKeyStrNoPrefix := strings.TrimPrefix(publicKeyHexStrWithHexAndECPrefix, "0x04")
	require.False(t, strings.Contains(publicKeyStrNoPrefix, "0x04"), "This is the public key which is used for generating public wallet addresses")

	// The public address is simply the Keccak-256 hash of the public key, and then we take the last 40 characters (20 bytes) and prefix it with 0x
	hash := sha3.NewLegacyKeccak256()
	// Strip down the byte of the 0x04 prefix before Keccak-256 hashing
	hash.Write(publicKeyBytes[1:])
	expectedPublicAddress := hexutil.Encode(hash.Sum(nil)[12:])

	strAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// Test that PubkeyToAddress has the same result ad the known process
	require.Equal(t, strings.ToUpper(expectedPublicAddress), strings.ToUpper(strAddress.Hex()))

	// Validate that the existing test accounts have a valid public address
	existingPrivateKeyHex := testData.PrivateKeys[0]
	existingPrivateKeyBytes, err := hexutil.Decode(existingPrivateKeyHex)
	require.NoError(t, err)
	require.Equal(t, 32, len(existingPrivateKeyBytes), "Converted existing private key should be 32 bytes long")

	// TODO: extract public key for existing account and validate that it matches the testData.Addresses
}

// TODO: TestGenerateHDWallet using https://github.com/ethereum/go-ethereum/blob/master/accounts/hd.go

// A keystore is a file containing an encrypted wallet private key. Keystores in go-ethereum can only contain one wallet key pair per file.
func TestKeystores(t *testing.T) {
	testDirPath := t.TempDir()
	initialKeystorePath := filepath.Join(testDirPath, "keystore")
	testSecrets := []string{"testSecret1", "testSecret2"}
	accToSecretIndex := make(map[common.Address]int, 2)

	// Create a keystore
	ks := keystore.NewKeyStore(initialKeystorePath, keystore.LightScryptN, keystore.LightScryptP)
	_, err := os.Stat(initialKeystorePath)
	// Keystore creation on't touch the filesystem, yet
	require.True(t, os.IsNotExist(err))
	for i, secret := range testSecrets {
		// Generate a new wallet
		acc, err := ks.NewAccount(secret)
		require.NoError(t, err)
		accToSecretIndex[acc.Address] = i

		// NewAccount should have created a new keystore file
		files, err := ioutil.ReadDir(initialKeystorePath)
		require.NoError(t, err)
		require.Equal(t, i+1, len(files))
	}

	// Import an existing keystore into a new keystore
	importedKeystorePath := filepath.Join(testDirPath, "imported_keystore")

	importedKs := keystore.NewKeyStore(importedKeystorePath, keystore.LightScryptN, keystore.LightScryptP)

	originalAccounts := ks.Accounts()

	// Read all files from the initial keystore and validate that there is a matching account in the original keystore with the same file path
	initialFiles, err := ioutil.ReadDir(initialKeystorePath)
	require.NoError(t, err)
	for i, ksFile := range initialFiles {
		// Read keystore file
		filePath := filepath.Join(initialKeystorePath, ksFile.Name())
		jsonContent, err := ioutil.ReadFile(filePath)
		require.NoError(t, err)

		var acc accounts.Account
		for _, acc = range originalAccounts {
			if acc.URL.Path == filePath {
				break
			}
		}
		require.Equal(t, acc.URL.Path, filePath, "existing file found in the original keystore")

		// Import keystore content in the new keystore
		importedAcc, err := importedKs.Import(jsonContent, testSecrets[i], testSecrets[i])
		require.NoError(t, err)
		require.True(t, ks.HasAddress(importedAcc.Address))

		err = ks.Delete(acc, testSecrets[accToSecretIndex[acc.Address]])
		require.NoError(t, err)
		_, err = os.Stat(filePath)
		require.True(t, os.IsNotExist(err))
	}
}

func TestAddressIsValid(t *testing.T) {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	require.True(t, re.MatchString("0x323b5d4c32345ced77393b3530b1eed0f346429d"), "Address valid")
	require.False(t, re.MatchString("0xZYXb5d4c32345ced77393b3530b1eed0f346429d"), "Address NOT valid")
}

func TestAddressIsFromASmartContract(t *testing.T) {
	client, _, anvil, tearDown := testClient()
	defer tearDown()

	addresses, err := anvil.AvailableAddresses()
	require.NoError(t, err)

	bytecode, err := client.CodeAt(context.Background(), addresses[0], nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
	}

	isContract := len(bytecode) > 0
	require.False(t, isContract, "Address is not a smart contract")
}
