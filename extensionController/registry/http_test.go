package registry

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type HttpRegistrySuite struct {
	suite.Suite
	server   *mockRegistryServer
	registry Registry
}

func TestBucketFsApiSuite(t *testing.T) {
	suite.Run(t, new(HttpRegistrySuite))
}

func (suite *HttpRegistrySuite) SetupSuite() {
	suite.server = newMockRegistryServer(&suite.Suite)
	suite.server.start()
}

func (suite *HttpRegistrySuite) TeardownSuite() {
	suite.server.close()
}

func (suite *HttpRegistrySuite) SetupTest() {
	suite.registry = NewRegistry(suite.server.indexUrl())
	suite.server.reset()
}

func (suite *HttpRegistrySuite) TestFindExtensions_noExtensionsAvailable() {
	suite.server.setRegistryContent(`{}`)
	extensions, err := suite.registry.FindExtensions()
	suite.NoError(err)
	suite.Empty(extensions)
}

func (suite *HttpRegistrySuite) TestFindExtensions() {
	suite.server.setRegistryContent(`{"extensions":[{"id": "ext1"},{"id": "ext2"},{"id": "ext3"}]}`)
	extensions, err := suite.registry.FindExtensions()
	suite.NoError(err)
	suite.Equal([]string{"ext1", "ext2", "ext3"}, extensions)
}

func (suite *HttpRegistrySuite) TestReadExtensionFailsWhenLoadingIndex() {
	suite.server.setRegistryContent(`invalid`)
	content, err := suite.registry.ReadExtension("unknown-ext-id")
	suite.EqualError(err, "failed to decode registry content: invalid character 'i' looking for beginning of value")
	suite.Equal("", content)
}

func (suite *HttpRegistrySuite) TestReadExtensionFailsForUnknownExtension() {
	url := suite.server.baseUrl() + "/ext1.js"
	suite.server.setRegistryContent(`{"extensions":[{"id": "ext1", "url": "` + url + `"}]}`)
	content, err := suite.registry.ReadExtension("unknown-ext-id")
	suite.ErrorContains(err, `extension "unknown-ext-id" not found`)
	suite.Equal("", content)
}

func (suite *HttpRegistrySuite) TestReadExtensionFailsForFailedStatusCode() {
	url := suite.server.baseUrl() + "/ext1.js"
	suite.server.setRegistryContent(`{"extensions":[{"id": "ext1", "url": "` + url + `"}]}`)
	content, err := suite.registry.ReadExtension("ext1")
	suite.ErrorContains(err, `failed to load extension "ext1": registry at `+url+` returned status "404 Not Found"`)
	suite.Equal("", content)
}

func (suite *HttpRegistrySuite) TestReadExtensionSucceeds() {
	url := suite.server.baseUrl() + "/ext1.js"
	suite.server.setPathContent("/ext1.js", "ext-content")
	suite.server.setRegistryContent(`{"extensions":[{"id": "ext1", "url": "` + url + `"}]}`)
	content, err := suite.registry.ReadExtension("ext1")
	suite.Nil(err)
	suite.Equal("ext-content", content)
}
