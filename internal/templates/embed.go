package templates

import "embed"

//go:embed all:presets all:shared all:layers
var FS embed.FS
