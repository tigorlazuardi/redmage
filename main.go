package main

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/tigorlazuardi/redmage/app"
)

//go:embed public
var embeddedStatic embed.FS
var osDirStatic = os.DirFS("public")

func main() {
	var static fs.FS
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	if isGoRun {
		static = osDirStatic
	} else {
		var err error
		static, err = fs.Sub(embeddedStatic, "public")
		if err != nil {
			panic(err)
		}
	}
	if err := app.Start(static); err != nil {
		log.Fatal(err)
	}
}
