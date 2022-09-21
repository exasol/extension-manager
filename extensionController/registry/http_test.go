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
	suite.registry = NewRegistry(suite.server.url())
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
