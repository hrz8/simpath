package templates

import "embed"

//go:embed *.html partials/*.html
var TemplatesFS embed.FS
