module github.com/exasol/extension-manager

go 1.19

require (
	github.com/dop251/goja v0.0.0-20230402114112-623f9dda9079
	github.com/dop251/goja_nodejs v0.0.0-20230322100729-2550c7b6c124
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/exasol/exasol-driver-go v0.4.7
	github.com/exasol/exasol-test-setup-abstraction-server/go-client v0.3.2
	github.com/go-chi/chi/v5 v5.0.8
	github.com/kinbiko/jsonassert v1.1.1
	github.com/swaggo/http-swagger v1.3.4
)

// Can't upgrade to latest version because of error
// no required module provides package github.com/getkin/kin-openapi/jsoninfo;
// Will be fixed in https://github.com/Nightapes/go-rest/issues/6
require github.com/getkin/kin-openapi v0.98.0 // indirect

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/spec v0.20.8 // indirect
	github.com/google/pprof v0.0.0-20230323073829-e72429f035bd // indirect
	github.com/iancoleman/orderedmap v0.2.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/swaggo/swag v1.8.12 // indirect
	golang.org/x/tools v0.7.0 // indirect
)

require (
	github.com/Nightapes/go-rest v0.3.2
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dlclark/regexp2 v1.8.1 // indirect
	github.com/exasol/error-reporting-go v0.1.1 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.12.0 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/leodido/go-urn v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/swaggo/files v1.0.1 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
