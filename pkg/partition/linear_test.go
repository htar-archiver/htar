package partition

import (
  . "htar/pkg/core"

  "testing"
  "github.com/stretchr/testify/assert"
)

func TestMakeSinglePartition(t *testing.T) {
  groups := makeTestGroups()
  linear := &LinearPartitioner{Attributes: Attributes{MaxPartionSize: 99999}}
  parts, err := linear.MakePartitions(groups)
  assert.Nil(t, err)
  assert.Equal(t, 1, len(parts))
  assert.Equal(t, 6, parts[0].TotalFiles)
  assert.Equal(t, parts[0].Groups, groups)
}

func TestMakePartitions(t *testing.T) {
  groups := makeTestGroups()
  linear := &LinearPartitioner{Attributes: Attributes{MaxPartionSize: 20480}}
  parts, err := linear.MakePartitions(groups)
  assert.Nil(t, err)
  assert.Equal(t, 2, len(parts))
  assert.Equal(t, parts[0].Groups[0], groups[0])
  assert.Equal(t, parts[1].Groups[0], groups[1])
}

func TestMakePartitionsAllowSplit(t *testing.T) {
  groups := makeTestGroups()
  linear := &LinearPartitioner{Attributes: Attributes{MaxPartionSize: 16384}, AllowSplit: true}
  parts, err := linear.MakePartitions(groups)
  assert.Nil(t, err)
  assert.Equal(t, 2, len(parts))
  assert.Equal(t, 2, len(parts[0].Groups))
  assert.Equal(t, "Test1", parts[0].Groups[0].Name)
  assert.Equal(t, "Test2", parts[0].Groups[1].Name)
  assert.Equal(t, 1, len(parts[1].Groups))
  assert.Equal(t, "Test2", parts[1].Groups[0].Name)
  assert.Equal(t, len(parts[1].Groups[0].Files), parts[1].TotalFiles)
}

func TestGroupTooLarge(t *testing.T) {
  groups := makeTestGroups()
  linear := &LinearPartitioner{Attributes: Attributes{MaxPartionSize: 4096}}
  parts, err := linear.MakePartitions(groups)
  assert.Nil(t, parts)
  assert.EqualError(t, err, "file group \"Test2\" (17.0 KB) is too large to fit in partition without splitting")
}

func TestFileTooLarge(t *testing.T) {
  groups := makeTestGroups()
  linear := &LinearPartitioner{Attributes: Attributes{MaxPartionSize: 4096}, AllowSplit: true}
  parts, err := linear.MakePartitions(groups)
  assert.Nil(t, parts)
  assert.EqualError(t, err, "file \"/test2/b.bin\" (8.0 KB) is too large to fit in partition")
}

func TestNoFilesInSource(t *testing.T) {
  groups := make([]FileGroup, 0)
  linear := &LinearPartitioner{}
  parts, err := linear.MakePartitions(groups)
  assert.Nil(t, parts)
  assert.EqualError(t, err, "no files in any source")
}
