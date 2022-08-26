package restAPI

import (
	"github.com/Nightapes/go-rest/pkg/openapi"
)

var authentication = map[string][]string{BasicAuth: {}, BearerAuth: {}}

type exaPath struct {
	*openapi.PathBuilder
}

func newPathWithDbQueryParams() *exaPath {
	path := &exaPath{getV1PublicBasePath(openapi.NewPathBuilder())}
	path.WithQueryParameter("dbHost", openapi.STRING, "Exasol database hostname", true)
	path.WithQueryParameter("dbPort", openapi.INTEGER, "Exasol database port number", true)
	return path
}

func getV1PublicBasePath(builder *openapi.PathBuilder) *openapi.PathBuilder {
	return builder.Add("api").Add("v1")
}
