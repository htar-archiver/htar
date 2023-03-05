package partition

import "testing"

func TestMakePartitions(t *testing.T) {
  groups := makeTestGroups()
  linear := &LinearPartitioner{MaxPartionSize: 99999}
  parts, err := linear.MakePartitions(groups)
  if err != nil {
    t.Fail()
  }
  if len(parts) != 1 {
    t.Fail()
  }
}
