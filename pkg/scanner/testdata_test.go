package scanner

import (
  "bytes"
  "io/fs"
  "testing"
  "testing/fstest"
  "github.com/stretchr/testify/assert"
)

func makeFileBytes(lines int) []byte {
  return bytes.Repeat([]byte("This is unit testing file content\n"), lines)
}

func makeTestFS() fstest.MapFS {
  m := make(fstest.MapFS)
  m["var/pool/data/Documents/Notes.txt"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(25),
	}
  m["var/pool/data/Documents/2020/doc1.txt"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(50),
	}
  m["var/pool/data/Documents/2020/doc2.txt"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(100),
	}
  m["var/pool/data/Documents/2021/doc1.txt"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(100),
	}
  m["var/pool/data/Music/Artist/Track01.mp3"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(1000),
	}
  m["var/pool/data/Music/Artist/Track02.mp3"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(1000),
	}
  m["var/pool/data/Music/Artist/Track03.mp3"] = &fstest.MapFile{
		Mode: 0666,
    Data: makeFileBytes(1000),
	}
  return m
}

func TestMapFS(t *testing.T) {
  fsys := makeTestFS()
  files, err := fs.ReadDir(fsys, "var/pool/data")
  assert.Nil(t, err)
  assert.True(t, len(files) > 1)
}
