package learning

// TODO: enable after finishing generateTransactions
// func TestTransactionsGetCount(t *testing.T) {
// 	client, _, ganache, tearDown := testClient()
// 	defer tearDown()

// 	addresses, err := ganache.AvailableAddresses()
// 	require.NoError(t, err)
// 	require.Greater(t, len(addresses), 3)

// 	err = generateTransactions(ganache,
// 		map[int][]TestTransaction{
// 			3: {
// 				{
// 					From:    addresses[0],
// 					To:      addresses[1],
// 					Balance: big.NewInt(100),
// 				},
// 				{
// 					From:    addresses[2],
// 					To:      addresses[3],
// 					Balance: big.NewInt(1000),
// 				},
// 			},
// 		})
// 	require.NoError(t, err)

// 	block, err := client.BlockByNumber(context.Background(), big.NewInt(3))
// 	require.NoError(t, err)

// 	count, err := client.TransactionCount(context.Background(), block.Hash())
// 	require.NoError(t, err)
// 	require.Equal(t, uint(2), count)
// }
