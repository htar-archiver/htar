package partition

import (
  . "htar/pkg/core"
)

type Attributes struct {
  MaxPartionSize int64
}

func (p *Attributes) GetMaxSize() int64{
  return p.MaxPartionSize
}

type Partitioner interface{
  GetMaxSize() int64
  MakePartitions(groups []FileGroup) ([]Partition, error)
}
