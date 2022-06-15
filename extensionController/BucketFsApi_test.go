package extensionController

import (
	"testing"

	"github.com/exasol/extension-manager/integrationTesting"
	"github.com/stretchr/testify/suite"
)

type BucketFsAPISuite struct {
	integrationTesting.IntegrationTestSuite
}

func TestBucketFsApiSuite(t *testing.T) {
	suite.Run(t, new(BucketFsAPISuite))
}

func (suite *BucketFsAPISuite) TestListBuckets() {
	connectionWithNoAutocommit, err := suite.Exasol.CreateConnectionWithConfig(false)
	suite.NoError(err)
	defer func() { suite.NoError(connectionWithNoAutocommit.Close()) }()
	bfsAPI := CreateBucketFsAPI(connectionWithNoAutocommit)
	result, err := bfsAPI.ListBuckets()
	suite.NoError(err)
	suite.Assert().Contains(result, "default")
}

func (suite *BucketFsAPISuite) TestListFiles() {
	connectionWithNoAutocommit, err := suite.Exasol.CreateConnectionWithConfig(false)
	suite.NoError(err)
	defer func() { suite.NoError(connectionWithNoAutocommit.Close()) }()
	bfsAPI := CreateBucketFsAPI(connectionWithNoAutocommit)
	suite.NoError(suite.Exasol.UploadStringContent("12345", "myFile.txt"))
	defer func() { suite.NoError(suite.Exasol.DeleteFile("myFile.txt")) }()
	result, err := bfsAPI.ListFiles("default")
	suite.NoError(err)
	suite.Assert().Contains(result, BfsFile{Name: "myFile.txt", Size: 5})
}
