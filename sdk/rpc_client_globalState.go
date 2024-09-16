package sdk

import (
	"encoding/json"
	"fmt"
)

type StateIdentifier struct {
	BlockHash string `json:"BlockHash,omitempty"`
}

func (c *RpcClient) QueryGlobalState(key string, BlockHash string) (GlobalStateResult, error) {
	resp, err := c.RpcCall("query_global_state", map[string]any{

		"key":  key,
		"path": []any{},
		"state_identifier": StateIdentifier{
			BlockHash: BlockHash,
		},
	})
	//log.Println(string(resp.Result))
	if err != nil {
		return GlobalStateResult{}, err
	}
	var result GlobalStateResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return GlobalStateResult{}, fmt.Errorf("failed to get result: %w", err)
	}
	//fmt.Println(*result.StoredValue.Withdraw)

	return result, nil
}

type GlobalStateResult struct {
	ApiVersion  string      `json:"api_version"`
	BlockHeader BlockHeader `json:"block_header"`
	StoredValue StoredValue `json:"stored_value"`
}
