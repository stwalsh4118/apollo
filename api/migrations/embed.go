package migrations

import "embed"

// Files stores SQL migrations embedded at build time.
//
//go:embed *.sql
var Files embed.FS
