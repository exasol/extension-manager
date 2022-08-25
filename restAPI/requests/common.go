package requests

import (
	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/exasol/extension-manager/restAPI/core"
)

var authentication = map[string][]string{core.BasicAuth: {}, core.BearerAuth: {}}

type ExaPath struct {
	*openapi.PathBuilder
}

func NewPathWithDbQueryParams() *ExaPath {
	path := &ExaPath{getV1PublicBasePath(openapi.NewPathBuilder())}
	path.WithQueryParameter("dbHost", openapi.STRING, "Exasol database hostname", true)
	path.WithQueryParameter("dbPort", openapi.INTEGER, "Exasol database port number", true)
	return path
}

func getV1PublicBasePath(builder *openapi.PathBuilder) *openapi.PathBuilder {
	return builder.Add("api").Add("v1")
}
