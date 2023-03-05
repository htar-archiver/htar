package partition

import (
  "errors"
  "fmt"
  "github.com/c2h5oh/datasize"
  . "htar/pkg/core"
)

type LinearPartitioner struct {
  Attributes
  AllowSplit bool
}

func (p *LinearPartitioner) MakePartitions(groups []FileGroup) ([]Partition, error) {
    parts := make([]Partition, 0)
    split := false
    part := Partition{}
    empty := true

    for _, g := range groups {

      if g.TotalSize > p.MaxPartionSize {
        // requires splitting
        if !p.AllowSplit {
          return nil, fmt.Errorf("file group %q (%s) is too large to fit in partition without splitting",
            g.Name, datasize.ByteSize(g.TotalSize).HumanReadable())
        }      
        split = true
      }

      if empty && len(g.Files) > 0 {
        empty = false
      }

      var divides []FileGroup
      if !split {
        // add whole group
        divides = []FileGroup{g}
      } else {
        divides = make([]FileGroup, 0)
        divide := FileGroup{Name: g.Name}
        remaining := p.MaxPartionSize - part.TotalSize

        for _, f := range g.Files {
          if f.Size > p.MaxPartionSize {
            return nil, fmt.Errorf("file %q (%s) is too large to fit in partition",
              f.Path, datasize.ByteSize(f.Size).HumanReadable())
          }
          if divide.TotalSize + f.Size > remaining {
            // next divide
            divides = append(divides, divide)
            divide = FileGroup{Name: g.Name}
            remaining = p.MaxPartionSize
          }
          divide.Files = append(divide.Files, f)
          divide.TotalSize += f.Size
        }
        divides = append(divides, divide)
      }

      for _, d := range divides {
        if part.TotalSize + d.TotalSize > p.MaxPartionSize {
          // next partition
          parts = append(parts, part)
          part = Partition{}
        }
        part.Groups = append(part.Groups, d)
        part.TotalFiles += len(d.Files)
        part.TotalSize += d.TotalSize
      }
    }

    if empty {
      return nil, errors.New("no files in any source")
    }

    parts = append(parts, part)
    return parts, nil
}
