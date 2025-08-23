package static

import _ "embed"

//go:embed index.css
var IndexCSS []byte

//go:embed index.html
var IndexHTML []byte
