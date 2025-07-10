package pkg

import (
	"fmt"
	"math/big"
	"strings"

	"testing"

	"github.com/ethereum/go-ethereum/params"
)

var (
	// parameters used for signing transactions
	chainID  = big.NewInt(1074)
	gasPrice = big.NewInt(30 * params.GWei)

	// would be nice to use a networkID that's different from chainID,
	// but some clients don't support the distinction properly.
	networkID = big.NewInt(1074)
)

var files = map[string]string{
	"/genesis.json": "./init/genesis.json",
}

type testSpec struct {
	Name  string
	About string
	Run   func(*TestEnv)
}

var tests = []testSpec{
	// HTTP RPC tests.
	// {Name: "http/BalanceAndNonceAt", Run: balanceAndNonceAtTest}, // pass
	// {Name: "http/CanonicalChain", Run: canonicalChainTest},
	// {Name: "http/CodeAt", Run: CodeAtTest}, // pass
	// {Name: "http/ContractDeployment", Run: deployContractTest},
	// {Name: "http/ContractDeploymentOutOfGas", Run: deployContractOutOfGasTest},
	// {Name: "http/EstimateGas", Run: estimateGasTest},
	// {Name: "http/GenesisBlockByHash", Run: genesisBlockByHashTest},
	// {Name: "http/GenesisBlockByNumber", Run: genesisBlockByNumberTest},
	// {Name: "http/GenesisHeaderByHash", Run: genesisHeaderByHashTest},
	// {Name: "http/GenesisHeaderByNumber", Run: genesisHeaderByNumberTest},
	// {Name: "http/Receipt", Run: receiptTest},
	// {Name: "http/SyncProgress", Run: syncProgressTest},
	// {Name: "http/TransactionCount", Run: transactionCountTest},
	// {Name: "http/TransactionInBlock", Run: transactionInBlockTest},
	// {Name: "http/TransactionReceipt", Run: TransactionReceiptTest},

	// HTTP ABI tests.
	// {Name: "http/ABICall", Run: callContractTest},
	// {Name: "http/ABITransact", Run: transactContractTest},

	// // WebSocket RPC tests.
	// {Name: "ws/BalanceAndNonceAt", Run: balanceAndNonceAtTest},
	// {Name: "ws/CanonicalChain", Run: canonicalChainTest},
	// {Name: "ws/CodeAt", Run: CodeAtTest},
	// {Name: "ws/ContractDeployment", Run: deployContractTest},
	// {Name: "ws/ContractDeploymentOutOfGas", Run: deployContractOutOfGasTest},
	// {Name: "ws/EstimateGas", Run: estimateGasTest},
	// {Name: "ws/GenesisBlockByHash", Run: genesisBlockByHashTest},
	// {Name: "ws/GenesisBlockByNumber", Run: genesisBlockByNumberTest},
	// {Name: "ws/GenesisHeaderByHash", Run: genesisHeaderByHashTest},
	// {Name: "ws/GenesisHeaderByNumber", Run: genesisHeaderByNumberTest},
	// {Name: "ws/Receipt", Run: receiptTest},
	// {Name: "ws/SyncProgress", Run: syncProgressTest},
	// {Name: "ws/TransactionCount", Run: transactionCountTest},
	// {Name: "ws/TransactionInBlock", Run: transactionInBlockTest},
	// {Name: "ws/TransactionReceipt", Run: TransactionReceiptTest},

	// // WebSocket subscription tests.
	// {Name: "ws/NewHeadSubscription", Run: newHeadSubscriptionTest},
	// {Name: "ws/LogSubscription", Run: logSubscriptionTest},
	// {Name: "ws/TransactionInBlockSubscription", Run: transactionInBlockSubscriptionTest},

	// // WebSocket ABI tests.
	// {Name: "ws/ABICall", Run: callContractTest},
	// {Name: "ws/ABITransact", Run: transactContractTest},
}

// func runLESTests(t *testing.T, serverNode *hivesim.Client) {
// 	// Configure LES client.
// 	clientParams := clientEnv.Set("HIVE_NODETYPE", "light")
// 	// Disable mining.
// 	clientParams = clientParams.Set("HIVE_MINER", "")
// 	clientParams = clientParams.Set("HIVE_CLIQUE_PRIVATEKEY", "")

// 	enode, err := serverNode.EnodeURL()
// 	if err != nil {
// 		t.Fatal("can't get node peer-to-peer endpoint:", enode)
// 	}

// 	// // Sync all sink nodes against the source.
// 	// t.RunAllClients(hivesim.ClientTestSpec{
// 	// 	Role:        "eth1_les_client",
// 	// 	Name:        "CLIENT as LES client",
// 	// 	Description: "This runs the RPC tests against an LES client.",
// 	// 	Parameters:  clientParams,
// 	// 	Files:       files,
// 	// 	AlwaysRun:   true,
// 	// 	Run: func(t *testing.T, client *hivesim.Client) {
// 	// 		err := client.RPC().Call(nil, "admin_addPeer", enode)
// 	// 		if err != nil {
// 	// 			t.Fatalf("connection failed:", err)
// 	// 		}
// 	// 		waitSynced(client.RPC())
// 	// 		runAllTests(t, client, client.Type+" LES")
// 	// 	},
// 	// })
// }

// runAllTests runs the tests against a client instance.
// Most tests simply wait for tx inclusion in a block so we can run many tests concurrently.
func RunAllTests(t *testing.T, host string, clientName string) {
	vault := newVault()

	s := newSemaphore(len(tests))
	for _, test := range tests {
		test := test
		s.get()
		go func() {
			defer s.put()
			t.Run(fmt.Sprintf("%s (%s)", test.Name, clientName), func(t *testing.T) {
				switch test.Name[:strings.IndexByte(test.Name, '/')] {
				case "http":
					runHTTP(t, host, vault, test.Run)
				case "ws":
					runWS(t, host, vault, test.Run)
				default:
					panic("bad test prefix in name " + test.Name)
				}
			},
			)
		}()
	}
	s.drain()
}

type semaphore chan struct{}

func newSemaphore(n int) semaphore {
	s := make(semaphore, n)
	for i := 0; i < n; i++ {
		s <- struct{}{}
	}
	return s
}

func (s semaphore) get() { <-s }
func (s semaphore) put() { s <- struct{}{} }

func (s semaphore) drain() {
	for i := 0; i < cap(s); i++ {
		<-s
	}
}
