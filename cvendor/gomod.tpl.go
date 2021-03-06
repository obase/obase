package cvendor

import "text/template"

var go_mod = template.Must(template.New("go_mod").Parse(`module github.com/obase

go 1.12

require (
	github.com/obase/api {{.Version}}
	github.com/obase/apix {{.Version}}
	github.com/obase/center {{.Version}}
	github.com/obase/clickhouse {{.Version}}
	github.com/obase/conf {{.Version}}
    github.com/obase/grpcx {{.Version}}
    github.com/obase/httpgw {{.Version}}
    github.com/obase/httpx {{.Version}}
    github.com/obase/kafka {{.Version}}
    github.com/obase/log {{.Version}}
	github.com/obase/mongo {{.Version}}
	github.com/obase/mongodb {{.Version}}
	github.com/obase/mysql {{.Version}}
	github.com/obase/redis {{.Version}}
)
`))
