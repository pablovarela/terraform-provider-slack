//go:build tools
// +build tools

package tools

import (
	_ "github.com/bflad/tfproviderdocs"
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/gordonklaus/ineffassign"
	_ "github.com/katbyte/terrafmt"
)
