package assets

import (
	"embed"
	"os"

	"github.com/kettek/go-multipath/v2"
)

var FS multipath.FS

//go:embed tilemap.txt
var _embededFS embed.FS

func init() {
	FS.AddFS(os.DirFS("./"))
	FS.InsertFS(_embededFS, multipath.LastPriority)
}
