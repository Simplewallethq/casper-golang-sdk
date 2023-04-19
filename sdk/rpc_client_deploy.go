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
	Success      SuccessExecutionResult `json:"success"`
	ErrorMessage *string                `json:"error_message,omitempty"`
}
type SuccessExecutionResult struct {
	Transfers []string `json:"transfers"`
	Cost      string   `json:"cost"`
}

// session
type JsonDeploySession struct {
	Transfer *JsonDeployTransfer `json:"Transfer,omitempty"`
}
type JsonDeployTransfer struct {
	Args []JsonDeployTransferArg `json:"args"`
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
