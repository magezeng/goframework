package Utils

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
)

func GzipDecode2(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer func() {
		_ = reader.Close()
	}()
	return ioutil.ReadAll(reader)
}

func ParseRequestParams(params interface{}) (string, *bytes.Reader, error) {
	if params == nil {
		return "", nil, errors.New("illegal parameter")
	}
	data, err := json.Marshal(params)
	if err != nil {
		return "", nil, errors.New("json convert string error")
	}
	jsonBody := string(data)
	binBody := bytes.NewReader(data)
	return jsonBody, binBody, nil
}
