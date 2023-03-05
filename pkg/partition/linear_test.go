package partition

import (
  "testing"
  "reflect"
  "github.com/stretchr/testify/assert"
)

func TestMakeSinglePartition(t *testing.T) {
  groups := makeTestGroups()
  linear := &LinearPartitioner{MaxPartionSize: 99999}
  parts, err := linear.MakePartitions(groups)
  assert.Nil(t, err)
  assert.Equal(t, len(parts), 1)
  assert.True(t, reflect.DeepEqual(parts[0].Groups, groups))
}

func TestMakePartitions(t *testing.T) {
  groups := makeTestGroups()
  linear := &LinearPartitioner{MaxPartionSize: 20480}
  parts, err := linear.MakePartitions(groups)
  assert.Nil(t, err)
  assert.Equal(t, 2, len(parts))
  assert.True(t, reflect.DeepEqual(parts[0].Groups[0], groups[0]))
  assert.True(t, reflect.DeepEqual(parts[1].Groups[0], groups[1]))
}

func TestMakePartitionsAllowSplit(t *testing.T) {
  groups := makeTestGroups()
  linear := &LinearPartitioner{MaxPartionSize: 16384, AllowSplit: true}
  parts, err := linear.MakePartitions(groups)
  assert.Nil(t, err)
  assert.Equal(t, 2, len(parts))
  assert.Equal(t, 2, len(parts[0].Groups))
  assert.Equal(t, "Test1", parts[0].Groups[0].Name)
  assert.Equal(t, "Test2", parts[0].Groups[1].Name)
  assert.Equal(t, 1, len(parts[1].Groups))
  assert.Equal(t, "Test2", parts[1].Groups[0].Name)
}

func TestGroupTooLarge(t *testing.T) {
  groups := makeTestGroups()
  linear := &LinearPartitioner{MaxPartionSize: 4096}
  parts, err := linear.MakePartitions(groups)
  assert.Nil(t, parts)
  assert.EqualError(t, err, "file group \"Test2\" (17.0 KB) is too large to fit in partition without splitting")
}

func TestFileTooLarge(t *testing.T) {
  groups := makeTestGroups()
  linear := &LinearPartitioner{MaxPartionSize: 4096, AllowSplit: true}
  parts, err := linear.MakePartitions(groups)
  assert.Nil(t, parts)
  assert.EqualError(t, err, "file \"/test2/b.bin\" (8.0 KB) is too large to fit in partition")
}
