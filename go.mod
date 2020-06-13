module github.com/my-cargonaut/cargonaut

// +heroku goVersion go1.14
// +heroku install ./cmd/cargonaut

go 1.14

require (
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/cors v1.1.1
	github.com/go-chi/render v1.0.1
	github.com/golangci/golangci-lint v1.27.0
	github.com/jmoiron/sqlx v1.2.1-0.20190826204134-d7d95172beb5
	github.com/lib/pq v1.7.0
	github.com/peterbourgon/ff/v2 v2.0.0
	github.com/rakyll/statik v0.1.7
	github.com/rubenv/sql-migrate v0.0.0-20200429072036-ae26b214fa43
	github.com/satori/go.uuid v1.2.0
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9
	gotest.tools/gotestsum v0.4.2
)
