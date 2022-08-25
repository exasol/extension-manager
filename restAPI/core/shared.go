package core

import (
	"github.com/Nightapes/go-rest/pkg/openapi"
)

const (
	TagExtension = "Extension"
	TagInstance  = "Instance"

	BearerAuth = "BearerAuth"
	BasicAuth  = "BasicAuth"
)

type ExaPath struct {
	*openapi.PathBuilder
}

// NewPublicPath /api/v1
func NewPathWithDbQueryParams() *ExaPath {
	path := &ExaPath{getV1PublicBasePath(openapi.NewPathBuilder())}
	path.WithQueryParameter("dbHost", openapi.STRING, "Exasol database hostname", true)
	path.WithQueryParameter("dbPort", openapi.INTEGER, "Exasol database port number", true)
	return path
}

func getV1PublicBasePath(builder *openapi.PathBuilder) *openapi.PathBuilder {
	return builder.
		Add("api").
		Add("v1")
}

func (e *ExaPath) GetInstalledExtensionsBasePath() *ExaPath {
	e.Add("installedExtensions")
	return e
}
