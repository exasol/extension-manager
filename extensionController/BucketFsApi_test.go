package extensionController

import (
	"context"
	"fmt"
	"time"

	"testing"

	"github.com/exasol/extension-manager/integrationTesting"
	"github.com/stretchr/testify/suite"
)

const DEFAULT_BUCKET_NAME = "default"

type BucketFsAPISuite struct {
	integrationTesting.IntegrationTestSuite
}

func TestBucketFsApiSuite(t *testing.T) {
	suite.Run(t, new(BucketFsAPISuite))
}

func (suite *BucketFsAPISuite) TestListBuckets() {
	bfsAPI := suite.createBucketFs()
	result, err := bfsAPI.ListBuckets()
	suite.NoError(err)
	suite.Contains(result, DEFAULT_BUCKET_NAME)
}

func (suite *BucketFsAPISuite) TestListFiles() {
	bfsAPI := suite.createBucketFs()
	fileName := fmt.Sprintf("myFile-%d.txt", time.Now().Unix())
	suite.NoError(suite.Exasol.UploadStringContent("12345", fileName))
	defer func() { suite.NoError(suite.Exasol.DeleteFile(fileName)) }()
	result, err := bfsAPI.ListFiles(DEFAULT_BUCKET_NAME)
	suite.NoError(err)
	suite.Contains(result, BfsFile{Name: fileName, Size: 5})
}

func (suite *BucketFsAPISuite) createBucketFs() BucketFsAPI {
	return CreateBucketFsAPI(context.Background(), suite.Connection)
}
