package migrations

import "embed"

// FS holds SQL migrations embedded for desktop builds and environments without a migrations directory.
//
//go:embed *.sql
var FS embed.FS
