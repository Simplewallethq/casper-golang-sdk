package sdk

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client = NewRpcClient("http://65.21.238.180:7777/rpc")

func TestRpcClient_GetLatestBlock(t *testing.T) {
	_, err := client.GetLatestBlock()

	if err != nil {
		t.Errorf("can't get latest block")
	}
}

func TestRpcClient_GetDeploy(t *testing.T) {
	hash := "1dfdf144eb0422eae3076cd8a17e55089010a133c6555c881ac0b9e2714a1605"
	res, err := client.GetDeploy(hash)
	if err != nil {
		t.Errorf("can't get deploy:%s", err)
	}
	for _, arg := range res.Deploy.Session.Transfer.Args {
		fmt.Println(arg.Type)
		if arg.ArgValue.ClType.ClTypeClass != nil {
			if arg.ArgValue.ClType.ClTypeClass.Option != nil {
				fmt.Println(*arg.ArgValue.ClType.ClTypeClass.Option)
			}
			if arg.ArgValue.ClType.ClTypeClass.ByteArray != nil {
				fmt.Println(arg.ArgValue.ClType.ClTypeClass.ByteArray)
			}
		}
		if arg.ArgValue.ClType.String != nil {
			fmt.Println(*arg.ArgValue.ClType.String)
		}
		if arg.ArgValue.Parsed != nil {
			if arg.ArgValue.Parsed.Integer != nil {
				fmt.Println(*arg.ArgValue.Parsed.Integer)
			}
			if arg.ArgValue.Parsed.String != nil {
				fmt.Println(*arg.ArgValue.Parsed.String)
			}
		}
		if arg.ArgValue.Bytes != nil {
			fmt.Println(*arg.ArgValue.Bytes)
		}
		fmt.Println()
	}

	assert.Equal(t, "amount", res.Deploy.Session.Transfer.Args[0].Type)
	assert.Equal(t, "U512", *res.Deploy.Session.Transfer.Args[0].ArgValue.ClType.String)
	assert.Equal(t, "2500000000", *res.Deploy.Session.Transfer.Args[0].ArgValue.Parsed.String)
	assert.Equal(t, "0400f90295", *res.Deploy.Session.Transfer.Args[0].ArgValue.Bytes)

	assert.Equal(t, "target", res.Deploy.Session.Transfer.Args[1].Type)
	if res.Deploy.Session.Transfer.Args[1].ArgValue.ClType.ClTypeClass.ByteArray == nil {
		t.Errorf("cltype is not byte array")
	}

	assert.Equal(t, "462a2c0cdb4d438e8d04087b2a32081d58946546274fb6b5046ce95e050b78f6", *res.Deploy.Session.Transfer.Args[1].ArgValue.Parsed.String)

	assert.Equal(t, "id", res.Deploy.Session.Transfer.Args[2].Type)
	assert.Equal(t, "U64", *res.Deploy.Session.Transfer.Args[2].ArgValue.ClType.ClTypeClass.Option)
	assert.Nil(t, nil, res.Deploy.Session.Transfer.Args[2].ArgValue.Parsed)
	hash_unstake := "36477d92494ed1c0091d74bdee1536785900f1f0b8ebf4a40730531526ebb36f"
	_, err = client.GetDeploy(hash_unstake)
	if err != nil {
		t.Errorf("can't get deploy:%s", err)
	}

}

func TestRpcClient_GetDeploy2(t *testing.T) {
	hash := "1dfdf144eb0422eae3076cd8a17e55089010a133c6555c881ac0b9e2714a1605"
	res, err := client.GetDeploy(hash)
	if err != nil {
		t.Errorf("can't get deploy:%s", err)
	}
	for _, arg := range res.Deploy.Session.Transfer.Args {
		fmt.Println(arg.Type)
	}

}

func TestRpcClient_GetBlockState(t *testing.T) {
	stateRootHash := "c0eb76e0c3c7a928a0cb43e82eb4fad683d9ad626bcd3b7835a466c0587b0fff"
	key := "account-hash-a9efd010c7cee2245b5bad77e70d9beb73c8776cbe4698b2d8fdf6c8433d5ba0"
	path := []string{"special_value"}
	_, err := client.GetStateItem(stateRootHash, key, path)

	if err != nil {
		t.Errorf("can't get block state")
	}
}

func TestRpcClient_GetAccountBalance(t *testing.T) {
	stateRootHash := "c0eb76e0c3c7a928a0cb43e82eb4fad683d9ad626bcd3b7835a466c0587b0fff"
	key := "account-hash-a9efd010c7cee2245b5bad77e70d9beb73c8776cbe4698b2d8fdf6c8433d5ba0"

	balanceUref := client.GetAccountMainPurseURef(key)

	_, err := client.GetAccountBalance(stateRootHash, balanceUref)

	if err != nil {
		t.Errorf("can't get account balance")
	}
}

func TestRpcClient_GetAccountInfo(t *testing.T) {
	pubkey := "020237037ff4845669e59d3e7698e7d58eb97ca378960ac57478a86a6a3535460292"
	bad_pubkey := "020237037ff4845669e59d3e7698e7d58ee97ca378960ac57478a86a6a3535460292"

	block := "a705d5cf0bca0ec2cc0ffceb2913900669b1234c230340086304989d67dde7d7"

	_, err := client.GetAccountInfo(bad_pubkey, block)
	var rpcerror *RpcError

	if err != nil {
		if errors.As(err, &rpcerror) {
			if rpcerror.Code != -32003 {
				t.Errorf("bad error code")
			}
		}
	}
	_, err = client.GetAccountInfo(pubkey, block)
	if err != nil {
		t.Errorf("can't get account info")
	}
	//t.Errorf("can't get account info")
	//log.Println("Running test...")
	//log.Println(err)
}

func TestRpcClient_GetAccountBalanceByKeypair(t *testing.T) {
	stateRootHash := "c0eb76e0c3c7a928a0cb43e82eb4fad683d9ad626bcd3b7835a466c0587b0fff"
	stateRootHashNew, err := client.GetStateRootHash(stateRootHash)
	if err != nil {
		return
	}
	path := []string{"special_value"}
	_, err = client.GetStateItem(stateRootHashNew.StateRootHash, sourceKeyPair.AccountHash(), path)
	if err != nil {
		_, err := client.GetAccountBalanceByKeypair(stateRootHashNew.StateRootHash, sourceKeyPair)

		if err != nil {
			t.Errorf("can't get account balance")
		}
	}
}

func TestRpcClient_GetBlockByHeight(t *testing.T) {
	_, err := client.GetBlockByHeight(1034)

	if err != nil {
		t.Errorf("can't get block by height")
	}
}

func TestRpcClient_GetBlockTransfersByHeight(t *testing.T) {
	_, err := client.GetBlockByHeight(1034)

	if err != nil {
		t.Errorf("can't get block transfers by height")
	}
}

func TestRpcClient_GetBlockByHash(t *testing.T) {
	_, err := client.GetBlockByHash("")

	if err != nil {
		t.Errorf("can't get block by hash")
	}
}

func TestRpcClient_GetBlockTransfersByHash(t *testing.T) {
	_, err := client.GetBlockTransfersByHash("")

	if err != nil {
		t.Errorf("can't get block transfers by hash")
	}
}

func TestRpcClient_GetLatestBlockTransfers(t *testing.T) {
	_, err := client.GetLatestBlockTransfers()

	if err != nil {
		t.Errorf("can't get latest block transfers")
	}
}

func TestRpcClient_GetValidator(t *testing.T) {
	_, err := client.GetValidator()

	if err != nil {
		t.Errorf("can't get validator")
	}
}

func TestRpcClient_GetStatus(t *testing.T) {
	_, err := client.GetStatus()

	if err != nil {
		t.Errorf("can't get status")
	}
}

func TestRpcClient_GetPeers(t *testing.T) {
	_, err := client.GetPeers()

	if err != nil {
		t.Errorf("can't get peers")
	}
}

func TestRpcClient_GetEraInfo(t *testing.T) {
	_, err := client.GetEraInfo(1639836)
	if err != nil {
		t.Errorf("can't get era info")
	}
}

// make sure your account has balance
func TestRpcClient_PutDeploy(t *testing.T) {
	deploy := NewTransferToUniqAddress(*source, UniqAddress{
		PublicKey:  dest,
		TransferId: 10,
	}, big.NewInt(3000000000), big.NewInt(10000), "casper-test", "")

	assert.True(t, deploy.ValidateDeploy())
	deploy.SignDeploy(sourceKeyPair)

	result, err := client.PutDeploy(*deploy)

	if !assert.NoError(t, err) {
		t.Errorf("error : %v", err)
		return
	}

	assert.Equal(t, hex.EncodeToString(deploy.Hash), result.Hash)
}

func TestRpcClient_QueryGlobalState(t *testing.T) {
	key := "withdraw-4ac2caf8adfba7087888bde7a426ba179393cd68cfe3e7cb722126cb9452a643"
	key_deploy := "deploy-d40427fea9b2de7c8af0ca6c8a2e10bb15c44e41274e0e70a025d566afda339f"
	hash_deploy := "4cd6ed78cf13bac8af85afbadf76937936ee992244dc8ec49bec3072d21aad82"
	hash := "8027736250ac033f81074ab53920a5bb9e94dfa0c4fa3c4b22e974771d28f3ec"
	hash2 := "4d43badd4b0ed2f1c2ec977d84602125367b2a6d9bca1a3dfa08eaf1c5f90fe5"
	res1, err := client.QueryGlobalState(key, hash)
	if err != nil {
		t.Errorf("can't get block state %v", err)
	}
	if res1.StoredValue.Withdraw != nil {
		fmt.Println("withdraw", res1.StoredValue.Withdraw)
		if res1.StoredValue.Withdraw.Amount != nil {
			fmt.Println("amount", res1.StoredValue.Withdraw.Amount)
		}
	}

	res2, err := client.QueryGlobalState(key, hash2)
	if err != nil {
		t.Errorf("can't get block state %v", err)
	}
	if res2.StoredValue.Withdraw != nil {
		fmt.Println("withdraw", res2.StoredValue.Withdraw)
		if res2.StoredValue.Withdraw.Amount != nil {
			fmt.Println("amount", res2.StoredValue.Withdraw.Amount)
		}
	}
	res3, err := client.QueryGlobalState(key_deploy, hash_deploy)
	if err != nil {
		t.Errorf("can't get block state %v", err)
	}
	if res3.StoredValue.DeployInfo != nil {
		fmt.Println("deploy", res3.StoredValue.DeployInfo)
	}

}
