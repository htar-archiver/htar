package scanner

import (
  "bytes"
  "io/fs"
  "math"
  "testing"
  "testing/fstest"
  "github.com/stretchr/testify/assert"
)

var testFileContent = []byte("This is unit testing file content\n")

func makeFileBytes(sizeKb int) []byte {
  size := sizeKb * 1024
  lines := int(math.Ceil(float64(size) / float64(len(testFileContent))))
  return bytes.Repeat(testFileContent, lines)[:size]
}

func makeTestFS() fstest.MapFS {
  m := make(fstest.MapFS)
  m["var/pool/data/Documents/Notes.txt"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(2),
	}
  m["var/pool/data/Documents/2020/doc1.txt"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(5),
	}
  m["var/pool/data/Documents/2020/doc2.txt"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(10),
	}
  m["var/pool/data/Documents/2021/doc1.txt"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(11),
	}
  m["var/pool/data/Music/Artist/Track01.mp3"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(120),
	}
  m["var/pool/data/Music/Artist/Track02.mp3"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(50),
	}
  m["var/pool/data/Music/Artist/Track03.mp3"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(80),
	}
  return m
}

func TestMapFS(t *testing.T) {
  fsys := makeTestFS()
  files, err := fs.ReadDir(fsys, "var/pool/data")
  assert.Nil(t, err)
  assert.True(t, len(files) > 1)
}
