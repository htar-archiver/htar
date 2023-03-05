package color

import (
  "github.com/gookit/color"
)

var (
  ArchiveValid = color.Style{color.FgGreen, color.OpBold}
  Error = color.Style{color.FgRed, color.OpBold}
  Partition = color.Style{color.FgCyan, color.OpBold}
  FileGroupSize = color.Style{color.BgGray}
  FileGroupFiles = color.Style{color.FgGray}
)
