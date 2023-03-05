package partition

import (
  . "htar/pkg/core"
)

func makeTestGroups() []FileGroup {
  return []FileGroup {
    FileGroup {
      Name: "Test1",
      Files: []FileEntry {
        { Path: "/test1/a.bin", Size: 1024 },
        { Path: "/test1/b.bin", Size: 1024 },
        { Path: "/test1/c.bin", Size: 2048 },
      },
      TotalSize: 4096,
    },
    FileGroup {
      Name: "Test2",
      Files: []FileEntry {
        { Path: "/test2/a.bin", Size: 1024 },
        { Path: "/test2/b.bin", Size: 8192 },
        { Path: "/test2/c.bin", Size: 8192 },
      },
      TotalSize: 17408,
    },
  };
}
