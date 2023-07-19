package bfs

import (
	"context"
	"fmt"
	"time"

	"testing"

	"github.com/exasol/extension-manager/pkg/integrationTesting"
	"github.com/stretchr/testify/suite"
)

const DEFAULT_BUCKET_PATH = "/buckets/bfsdefault/default/"

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

/* [utest -> dsn~extension-components~1] */
func (suite *BucketFsAPISuite) TestListEmptyDir() {
	bfsAPI := suite.createBucketFs()
	result, err := bfsAPI.ListFiles(context.Background(), suite.exasol.GetConnection())
	suite.NoError(err)
	suite.Empty(result)
}

func (suite *BucketFsAPISuite) TestListSingleFile() {
	bfsAPI := suite.createBucketFs()
	fileName := fmt.Sprintf("myFile-%d.txt", time.Now().Unix())
	suite.NoError(suite.exasol.Exasol.UploadStringContent("12345", fileName))
	suite.T().Cleanup(func() {
		suite.NoError(suite.exasol.Exasol.DeleteFile(fileName))
	})
	result, err := bfsAPI.ListFiles(context.Background(), suite.exasol.GetConnection())
	suite.NoError(err)
	suite.Len(result, 1)
	suite.Equal([]BfsFile{{Name: fileName, Path: DEFAULT_BUCKET_PATH + fileName, Size: 5}}, result)
}

func (suite *BucketFsAPISuite) TestListFilesRecursively() {
	bfsAPI := suite.createBucketFs()
	file1 := "file1"
	file2 := "dir1/file2"
	file3 := "dir2/file2"
	suite.NoError(suite.exasol.Exasol.UploadStringContent("1", file1))
	suite.NoError(suite.exasol.Exasol.UploadStringContent("12", file2))
	suite.NoError(suite.exasol.Exasol.UploadStringContent("123", file3))
	suite.T().Cleanup(func() {
		suite.NoError(suite.exasol.Exasol.DeleteFile(file1))
		suite.NoError(suite.exasol.Exasol.DeleteFile(file2))
		suite.NoError(suite.exasol.Exasol.DeleteFile(file3))
	})
	result, err := bfsAPI.ListFiles(context.Background(), suite.exasol.GetConnection())
	suite.NoError(err)
	suite.Len(result, 3)
	suite.Equal([]BfsFile{
		{Name: "file2", Path: DEFAULT_BUCKET_PATH + file2, Size: 2},
		{Name: "file2", Path: DEFAULT_BUCKET_PATH + file3, Size: 3},
		{Name: "file1", Path: DEFAULT_BUCKET_PATH + file1, Size: 1}}, result)
}

func (suite *BucketFsAPISuite) createBucketFs() BucketFsAPI {
	return CreateBucketFsAPI(DEFAULT_BUCKET_PATH)
}
