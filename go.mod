module github.com/my-cargonaut/cargonaut

// +heroku goVersion go1.14
// +heroku install ./cmd/cargonaut

go 1.14

require (
	github.com/SermoDigital/jose v0.9.2-0.20180104203859-803625baeddc
	github.com/fatih/structtag v1.2.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/cors v1.1.1
	github.com/go-chi/jwtauth v4.0.4+incompatible
	github.com/go-chi/render v1.0.1
	github.com/golangci/golangci-lint v1.27.0
	github.com/gomodule/redigo v1.8.2
	github.com/jmoiron/sqlx v1.2.1-0.20190826204134-d7d95172beb5
	github.com/lib/pq v1.7.0
	github.com/manifoldco/promptui v0.7.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/peterbourgon/ff/v2 v2.0.0
	github.com/rakyll/statik v0.1.7
	github.com/rubenv/sql-migrate v0.0.0-20200616145509-8d140a17f351
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9
	gotest.tools/gotestsum v0.5.0
)
