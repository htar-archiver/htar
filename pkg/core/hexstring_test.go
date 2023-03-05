package core

import (
  "crypto/sha256"
  "encoding/json"
  "testing"

  "github.com/stretchr/testify/assert"
)

type hexTest struct {
  Value HexString `json:"value"`
}

func TestHexStringToString(t *testing.T) {
  xs := HexString([]byte("ABC"))
  assert.Equal(t, "414243", xs.String())
}

func TestHexStringEncode(t *testing.T) {
  sha := sha256.New()

  o := hexTest{
    Value: sha.Sum(nil),
  }

  b, err := json.Marshal(o)
  assert.Nil(t, err)

  // sha256 of an empty string
  assert.Equal(t, "{\"value\":\"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\"}", string(b))
}

func TestHexStringDecode(t *testing.T) {
  o1 := hexTest{
    Value: []byte("Hello World"),
  }

  b, err := json.Marshal(o1)
  assert.Nil(t, err)

  var o2 hexTest
  err = json.Unmarshal(b, &o2)
  assert.Nil(t, err)

  assert.Equal(t, o1, o2)
}
