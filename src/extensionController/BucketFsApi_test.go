package extensionController

import (
	"context"
	"fmt"
	"time"

	"testing"

	"github.com/exasol/extension-manager/src/integrationTesting"
	"github.com/stretchr/testify/suite"
)

const DEFAULT_BUCKET_NAME = "default"

type BucketFsAPISuite struct {
	suite.Suite
	exasol integrationTesting.DbTestSetup
}

func TestBucketFsApiSuite(t *testing.T) {
	suite.Run(t, new(BucketFsAPISuite))
}

func (suite *BucketFsAPISuite) SetupSuite() {
	suite.exasol = *integrationTesting.StartDbSetup(&suite.Suite)
}

func (suite *BucketFsAPISuite) TeardownSuite() {
	suite.exasol.StopDb()
}

func (suite *BucketFsAPISuite) BeforeTest(suiteName, testName string) {
	suite.exasol.CreateConnection()
	suite.T().Cleanup(func() {
		suite.exasol.CloseConnection()
	})
}

func (suite *BucketFsAPISuite) TestListBuckets() {
	bfsAPI := suite.createBucketFs()
	result, err := bfsAPI.ListBuckets(context.Background(), suite.exasol.GetConnection())
	suite.NoError(err)
	suite.Contains(result, DEFAULT_BUCKET_NAME)
}

func (suite *BucketFsAPISuite) TestListFiles() {
	bfsAPI := suite.createBucketFs()
	fileName := fmt.Sprintf("myFile-%d.txt", time.Now().Unix())
	suite.NoError(suite.exasol.Exasol.UploadStringContent("12345", fileName))
	suite.T().Cleanup(func() {
		suite.NoError(suite.exasol.Exasol.DeleteFile(fileName))
	})
	result, err := bfsAPI.ListFiles(context.Background(), suite.exasol.GetConnection(), DEFAULT_BUCKET_NAME)
	suite.NoError(err)
	suite.Contains(result, BfsFile{Name: fileName, Size: 5})
}

func (suite *BucketFsAPISuite) createBucketFs() BucketFsAPI {
	return CreateBucketFsAPI()
}
