package learning

import (
	"context"
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"

	"github.com/stretchr/testify/require"
)

func TestConvertAddresses(t *testing.T) {
	strAddress := "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"

	address := common.HexToAddress(strAddress)
	require.Equal(t, strAddress, address.Hex())
}

func TestGetBalanceLastBlock(t *testing.T) {
	client, _, ganache, tearDown := testClientWithBlocks()
	defer tearDown()

	addresses, err := ganache.AvailableAddresses()
	require.NoError(t, err)
	balance, err := client.BalanceAt(context.Background(), addresses[0], nil)
	require.NoError(t, err)

	oneThousandEthAsWei, ok := new(big.Int).SetString("1000000000000000000000", 0)
	require.True(t, ok)
	require.Equal(t, 0, oneThousandEthAsWei.Cmp(balance), "balance should be 1000 ETH")
}

func TestGetBalanceFirstBlock(t *testing.T) {
	client, _, ganache, tearDown := testClientWithBlocks()
	defer tearDown()

	addresses, err := ganache.AvailableAddresses()
	require.NoError(t, err)
	balance, err := client.BalanceAt(context.Background(), addresses[0], big.NewInt(0))
	require.NoError(t, err)

	require.Equalf(t, 1000.0, balanceToEther(balance), "balance should be 1000 ETH")
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
	require.Equal(t, uint64(1), lastBlock.Difficulty().Uint64())
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
	require.Equal(t, strings.ToUpper(expectedPublicAddress), strings.ToUpper(strAddress.Hex()))

	// Validate that the existing test accounts have a valid public address
	existingPrivateKeyHex := testData.PrivateKeys[0]
	existingPrivateKeyBytes, err := hexutil.Decode(existingPrivateKeyHex)
	require.NoError(t, err)
	require.Equal(t, 32, len(existingPrivateKeyBytes), "Converted existing private key should be 32 bytes long")

	// TODO: extract public key for existing account and validate that it matches the testData.Addresses
}

func TestKeystores(t *testing.T) {
	testDirPath := t.TempDir()
	initialKeystorePath := filepath.Join(testDirPath, "keystore")
	testSecrets := []string{"testSecret", "testSecret"}

	// Create a keystore with two accounts
	{
		ks := keystore.NewKeyStore(initialKeystorePath, keystore.LightScryptN, keystore.LightScryptP)
		_, err := os.Stat(initialKeystorePath)
		require.True(t, os.IsNotExist(err))

		for i, secret := range testSecrets {
			_, err = ks.NewAccount(secret)
			require.NoError(t, err)

			files, err := ioutil.ReadDir(initialKeystorePath)
			require.NoError(t, err)
			require.Equal(t, i, len(files))
		}
	}

	importedKeystorePath := filepath.Join(testDirPath, "imported_keystore")

	// wronglyImportedKs := keystore.NewKeyStore(importedKeystorePath, keystore.StandardScryptN, keystore.StandardScryptP)
	// _, err := os.Stat(keystorePath)
	// require.True(t, os.IsNotExist(err))

	importedKs := keystore.NewKeyStore(importedKeystorePath, keystore.LightScryptN, keystore.LightScryptP)
	initialFiles, err := ioutil.ReadDir(initialKeystorePath)
	require.NoError(t, err)

	for _, ksFile := range initialFiles {
		_, err := importedKs.Import([]byte(ksFile.Name()), testSecrets[ksFile], testSecrets[ksFile])
		require.NoError(t, err)
	}
}
