package archive

import (
  "strings"
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestWrongVersion(t *testing.T) {
  json := strings.NewReader("{\"_version\":\"wrong\"}")
  meta := &PartitionMeta{}
  
  err := meta.Decode(json)
  assert.NotNil(t, err)
  assert.Regexp(t, "^expected meta data version \".+\" but file contains \"wrong\"$", err.Error())
}

func TestMissingVersion(t *testing.T) {
  json := strings.NewReader("{}")
  meta := &PartitionMeta{}
  
  err := meta.Decode(json)
  assert.NotNil(t, err)
  assert.Regexp(t, "^expected meta data version \".+\" but file contains \"\"$", err.Error())
}
