package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func (c *RpcClient) GetDeploy(hash string) (DeployResult, error) {
	resp, err := c.RpcCall("info_get_deploy", map[string]string{
		"deploy_hash": hash,
	})
	if err != nil {
		return DeployResult{}, err
	}

	var result DeployResult
	err = json.Unmarshal(resp.Result, &result)

	if err != nil {
		return DeployResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

type DeployResult struct {
	Deploy           JsonDeploy            `json:"deploy"`
	ExecutionResults []JsonExecutionResult `json:"execution_results"`
}
type JsonDeploy struct {
	Hash      string             `json:"hash"`
	Header    JsonDeployHeader   `json:"header"`
	Session   *JsonDeploySession `json:"session,omitempty"`
	Approvals []JsonApproval     `json:"approvals"`
}
type JsonDeployHeader struct {
	Account      string    `json:"account"`
	Timestamp    time.Time `json:"timestamp"`
	TTL          string    `json:"ttl"`
	GasPrice     int       `json:"gas_price"`
	BodyHash     string    `json:"body_hash"`
	Dependencies []string  `json:"dependencies"`
	ChainName    string    `json:"chain_name"`
}
type JsonApproval struct {
	Signer    string `json:"signer"`
	Signature string `json:"signature"`
}
type JsonExecutionResult struct {
	BlockHash string          `json:"block_hash"`
	Result    ExecutionResult `json:"result"`
}
type ExecutionResult struct {
	Success      *SuccessExecutionResult `json:"success"`
	ErrorMessage *string                 `json:"error_message,omitempty"`
}
type SuccessExecutionResult struct {
	Transfers []string         `json:"transfers"`
	Cost      string           `json:"cost"`
	Effect    JsonDeployEffect `json:"effect"`
}
type JsonDeployEffect struct {
	Operations []any                  `json:"operations"`
	Transforms []JsonDeployTransforms `json:"transforms"`
}
type JsonDeployTransforms struct {
	Key       string              `json:"key"`
	Transform JsonDeployTransform `json:"transform"`
}

type JsonDeployTransform struct {
	TransformString *string
	WriteWithdraw   *JsonDeployWriteWithdraw
}

func (t *JsonDeployTransform) UnmarshalJSON(data []byte) error {
	var content any
	err := json.Unmarshal(data, &content)
	if err != nil {
		return err
	}
	switch v := content.(type) {
	case string:
		t.TransformString = &v
	case map[string]interface{}:
		if withdraw, ok := v["WriteWithdraw"]; ok {
			switch withdrawSlice := withdraw.(type) {
			case []interface{}:
				if len(withdrawSlice) == 1 {
					withdrawMap := withdrawSlice[0]
					switch withdrawValue := withdrawMap.(type) {
					case map[string]interface{}:
						var withdrawRes JsonDeployWriteWithdraw
						if amount, ok := withdrawValue["amount"]; ok {
							switch v := amount.(type) {
							case string:
								withdrawRes.Amount = v
							}
						}
						if era, ok := withdrawValue["era_of_creation"]; ok {
							switch v := era.(type) {
							case int:
								withdrawRes.EraOfCreation = v
							}
						}
						if bondingPurse, ok := withdrawValue["bonding_purse"]; ok {
							switch v := bondingPurse.(type) {
							case string:
								withdrawRes.BondingPurse = v
							}
						}
						if unbonderPublicKey, ok := withdrawValue["unbonder_public_key"]; ok {
							switch v := unbonderPublicKey.(type) {
							case string:
								withdrawRes.UnbonderPublicKey = v
							}
						}
						if validatorPublicKey, ok := withdrawValue["validator_public_key"]; ok {
							switch v := validatorPublicKey.(type) {
							case string:
								withdrawRes.ValidatorPublicKey = v
							}
						}
						t.WriteWithdraw = &withdrawRes
					}
				}
			}

		}

	}
	return nil
}

type JsonDeployWriteWithdraw struct {
	BondingPurse       string `json:"bonding_purse"`
	ValidatorPublicKey string `json:"validator_public_key"`
	UnbonderPublicKey  string `json:"unbonder_public_key"`
	EraOfCreation      int    `json:"era_of_creation"`
	Amount             string `json:"amount"`
}

// session
type JsonDeploySession struct {
	Transfer             *JsonDeployTransfer             `json:"Transfer,omitempty"`
	StoredContractByHash *JsonDeployStoredContractByHash `json:"StoredContractByHash,omitempty"`
}
type JsonDeployTransfer struct {
	Args []JsonDeployTransferArg `json:"args"`
}
type JsonDeployStoredContractByHash struct {
	Hash       string                  `json:"hash"`
	EntryPoint string                  `json:"entry_point"`
	Args       []JsonDeployTransferArg `json:"args"`
}
type JsonDeployTransferArg struct {
	Type     string
	ArgValue *JsonDeployArgClass
}
type JsonDeployArgClass struct {
	ClType JsonDeployClTypeUnion `json:"cl_type"`
	Bytes  *string               `json:"bytes,omitempty"`
	Parsed *JsonDeployParsed     `json:"parsed"`
}

type JsonDeployClTypeClass struct {
	Option    *string `json:"Option,omitempty"`
	List      any     `json:"List,omitempty"`
	ByteArray any     `json:"ByteArray,omitempty"`
	Result    any     `json:"Result,omitempty"`
	Map       any     `json:"Map,omitempty"`
	Tuple1    any     `json:"Tuple1,omitempty"`
	Tuple2    any     `json:"Tuple2,omitempty"`
	Tuple3    any     `json:"Tuple3,omitempty"`
}

type JsonDeployClTypeUnion struct {
	ClTypeClass *JsonDeployClTypeClass `json:"-"`
	String      *string                `json:"-"`
}

type JsonDeployParsed struct {
	Integer *int64
	String  *string
}

func (p *JsonDeployParsed) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err == nil {
		p.String = &s
		return nil
	}
	var i int64
	err = json.Unmarshal(data, &i)
	if err == nil {
		p.Integer = &i
		return nil
	}
	return errors.New("cannot unmarshal")
}

func (ac *JsonDeployClTypeUnion) UnmarshalJSON(data []byte) error {
	var c JsonDeployClTypeClass
	err := json.Unmarshal(data, &c)
	if err == nil {
		ac.ClTypeClass = &c
		return nil
	}
	var s string
	err = json.Unmarshal(data, &s)
	if err == nil {
		ac.String = &s
		return nil
	}

	return err

}

func (ta *JsonDeployTransferArg) UnmarshalJSON(data []byte) error {
	var arg []json.RawMessage
	err := json.Unmarshal(data, &arg)
	if err != nil {
		return err
	}
	if len(arg) < 2 {
		return errors.New("args out of range")
	}
	var typ string
	err = json.Unmarshal(arg[0], &typ)
	if err != nil {
		return err
	}
	ta.Type = typ
	var val JsonDeployArgClass
	err = json.Unmarshal(arg[1], &val)
	if err != nil {
		return err
	}
	ta.ArgValue = &val
	return nil
}
