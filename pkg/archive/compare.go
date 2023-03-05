package archive

import(
  "bytes"
  "fmt"
)

func CompareHashes(expected Hashes, actual Hashes) []string {
  diff := make([]string, 0)

  for path, hash := range expected {
    other, hasPath := actual[path]
    if hasPath && bytes.Equal(other, hash) {
      continue
    }

    diff = append(diff, fmt.Sprintf("%s: checksum mismatch [expected=%x,actual=%x]", path, hash, other))
  }
  
  for path, _ := range actual {
    if _, hasPath := expected[path]; !hasPath {
      diff = append(diff, fmt.Sprintf("%s: unexpected file", path))
    }
  }

  if (len(diff) > 0) {
    return diff
  }
  return nil
}
