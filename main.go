package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/containers/buildah"
	"github.com/containers/storage/pkg/unshare"
	"github.com/joyrex2001/kubedock/cmd"
)

//go:embed README.md
var readme string

//go:embed config.md
var config string

//go:embed LICENSE
var license string

func main() {
	if buildah.InitReexec() {
		return
	}
    fmt.Println("After init exec")
	fmt.Println(os.Geteuid())
	unshare.MaybeReexecUsingUserNamespace(false)
	cmd.README = readme
	cmd.LICENSE = license
	cmd.CONFIG = config
	cmd.Execute()
}
