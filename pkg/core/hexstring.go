package core

import (
  "encoding/hex"
  "encoding/json"
)

type HexString []byte

func (xs HexString) String() string {
  return hex.EncodeToString([]byte(xs))
}

func (xs HexString) MarshalJSON() ([]byte, error) {
  s := hex.EncodeToString([]byte(xs))
	return json.Marshal(s)
}

func (xs *HexString) UnmarshalJSON(data []byte) error {
  var s string
  err := json.Unmarshal(data, &s)
  if err != nil {
    return err
  }

  b, err := hex.DecodeString(s)
  if err != nil {
    return err
  }

  *xs = HexString(b)
  return nil
}
