package utils

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
)

func DecodeJSON(r io.Reader, dst interface{}) error {
	if dst == nil {
		return errors.New("DecodeJSON: dest is nil")
	}
	return json.NewDecoder(r).Decode(dst)
}

func DecodeJSONFromFile(filepath string, dst interface{}) error {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return DecodeJSON(bytes.NewReader(buf), dst)
}
