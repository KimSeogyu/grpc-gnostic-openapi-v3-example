package openapi

import (
	"embed"
	_ "embed"
)

//go:embed *
var StaticFiles embed.FS
