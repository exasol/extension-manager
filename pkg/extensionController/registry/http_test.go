package registry

import (
	"testing"

	"github.com/exasol/extension-manager/pkg/integrationTesting"
	"github.com/stretchr/testify/suite"
)

type HttpRegistrySuite struct {
	suite.Suite
	server   *integrationTesting.MockRegistryServer
	registry Registry
}

func TestBucketFsApiSuite(t *testing.T) {
	suite.Run(t, new(HttpRegistrySuite))
}

func (suite *HttpRegistrySuite) SetupSuite() {
	suite.server = integrationTesting.NewMockRegistryServer(&suite.Suite)
	suite.server.Start()
}

func (suite *HttpRegistrySuite) TeardownSuite() {
	suite.server.Close()
}

func (suite *HttpRegistrySuite) SetupTest() {
	suite.registry = NewRegistry(suite.server.IndexUrl())
	suite.server.Reset()
}

func (suite *HttpRegistrySuite) TestFindExtensions_noExtensionsAvailable() {
	suite.server.SetRegistryContent(`{}`)
	extensions, err := suite.registry.FindExtensions()
	suite.NoError(err)
	suite.Empty(extensions)
}

/* [itest -> dsn~extension-registry~1] */
/* [itest -> dsn~extension-definitions-storage~1] */
func (suite *HttpRegistrySuite) TestFindExtensions() {
	suite.server.SetRegistryContent(`{"extensions":[{"id": "ext1"},{"id": "ext2"},{"id": "ext3"}]}`)
	extensions, err := suite.registry.FindExtensions()
	suite.NoError(err)
	suite.Equal([]string{"ext1", "ext2", "ext3"}, extensions)
}

func (suite *HttpRegistrySuite) TestReadExtensionFailsWhenLoadingIndex() {
	suite.server.SetRegistryContent(`invalid`)
	content, err := suite.registry.ReadExtension("unknown-ext-id")
	suite.EqualError(err, "failed to decode registry content: invalid character 'i' looking for beginning of value")
	suite.Equal("", content)
}

func (suite *HttpRegistrySuite) TestReadExtensionFailsForUnknownExtension() {
	url := suite.server.BaseUrl() + "/ext1.js"
	suite.server.SetRegistryContent(`{"extensions":[{"id": "ext1", "url": "` + url + `"}]}`)
	content, err := suite.registry.ReadExtension("unknown-ext-id")
	suite.ErrorContains(err, `extension "unknown-ext-id" not found`)
	suite.Equal("", content)
}

func (suite *HttpRegistrySuite) TestReadExtensionFailsForFailedStatusCode() {
	url := suite.server.BaseUrl() + "/ext1.js"
	suite.server.SetRegistryContent(`{"extensions":[{"id": "ext1", "url": "` + url + `"}]}`)
	content, err := suite.registry.ReadExtension("ext1")
	suite.ErrorContains(err, `failed to load extension "ext1": registry at `+url+` returned status "404 Not Found"`)
	suite.Equal("", content)
}

func (suite *HttpRegistrySuite) TestReadExtensionSucceeds() {
	url := suite.server.BaseUrl() + "/ext1.js"
	suite.server.SetPathContent("/ext1.js", "ext-content")
	suite.server.SetRegistryContent(`{"extensions":[{"id": "ext1", "url": "` + url + `"}]}`)
	content, err := suite.registry.ReadExtension("ext1")
	suite.Nil(err)
	suite.Equal("ext-content", content)
}
