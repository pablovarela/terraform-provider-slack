//go:build tools
// +build tools

package main

import (
	_ "github.com/bflad/tfproviderdocs"
	_ "github.com/katbyte/terrafmt"
	_ "github.com/ysmood/golangci-lint"
)
