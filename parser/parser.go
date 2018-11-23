package parser

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
)

//https://github.com/ethereum/wiki/wiki/JSON-RPC

// Call Call
func Call(host, method string, params interface{}) (*JSON2Response, error) {
	j := NewJSON2RequestBlank()
	j.Method = method
	j.Params = params
	j.ID = 1

	_method := "POST"

	address := host

	data, err := j.JSONString()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(_method, address, strings.NewReader(data))
	req.Header.Set("Content-type", "application/json")
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jResp := new(JSON2Response)

	err = json.Unmarshal(body, jResp)
	if err != nil {
		return nil, err
	}

	return jResp, nil
}

// ParseTokenTransfer ParseTokenTransfer
// 0xa9059cbb0000000000000000000000000cf0698955123303a9a36ce470552c8d10ee6198000000000000000000000000000000000000000000000001158e460913d00000
func ParseTokenTransfer(inputData string) (string, string, error) {

	if len(inputData) != 138 {
		return "", "0", errors.New("len is error")
	}

	if inputData[0:10] != "0xa9059cbb" {
		return "", "0", errors.New("MethodID is error")
	}

	_addr := inputData[34:74]

	_value := new(big.Int)
	_value, ok := _value.SetString(inputData[74:138], 16)
	if !ok {
		return "", "0", errors.New("MethodID is error")
	}

	return _addr, _value.String(), nil
}
