package partition

import (
  . "htar/pkg/core"
)

type Partitioner interface{
  MakePartitions(groups []FileGroup) []Partition
}
