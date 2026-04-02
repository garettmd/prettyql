package web

import "embed"

//go:embed index.html
var StaticFiles embed.FS
