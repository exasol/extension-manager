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
	result, err := suite.listFiles()
	suite.NoError(err)
	suite.Empty(result)
}

func (suite *BucketFsAPISuite) TestListSingleFile() {
	fileName := fmt.Sprintf("myFile-%d.txt", time.Now().Unix())
	suite.uploadStringContent(fileName, "12345")
	result, err := suite.listFiles()
	suite.NoError(err)
	suite.Len(result, 1)
	suite.Equal([]BfsFile{{Name: fileName, Path: DEFAULT_BUCKET_PATH + fileName, Size: 5}}, result)
}

func (suite *BucketFsAPISuite) TestListFilesRecursively() {
	file1 := "file1"
	file2 := "dir1/file2"
	file3 := "dir2/file2"
	suite.uploadStringContent(file1, "1")
	suite.uploadStringContent(file2, "12")
	suite.uploadStringContent(file3, "123")
	result, err := suite.listFiles()
	suite.NoError(err)
	suite.Len(result, 3)
	suite.Equal([]BfsFile{
		{Name: "file2", Path: DEFAULT_BUCKET_PATH + file2, Size: 2},
		{Name: "file2", Path: DEFAULT_BUCKET_PATH + file3, Size: 3},
		{Name: "file1", Path: DEFAULT_BUCKET_PATH + file1, Size: 1}}, result)
}

func (suite *BucketFsAPISuite) TestFindAbsolutePathNoFileFound() {
	time.Sleep(3 * time.Second)
	result, err := suite.findAbsolutePath("no-such-file")
	suite.EqualError(err, `file "no-such-file" not found in BucketFS`)
	suite.Equal("", result)
}

func (suite *BucketFsAPISuite) TestFindAbsolutePathFileInRoot() {
	fileName := "file.txt"
	suite.uploadStringContent(fileName, "123")
	result, err := suite.findAbsolutePath(fileName)
	suite.NoError(err)
	suite.Equal("/buckets/bfsdefault/default/"+fileName, result)
}

func (suite *BucketFsAPISuite) TestFindAbsolutePathFileInSubDir() {
	suite.uploadStringContent("dir/file.txt", "123")
	result, err := suite.findAbsolutePath("file.txt")
	suite.NoError(err)
	suite.Equal("/buckets/bfsdefault/default/dir/file.txt", result)
}

func (suite *BucketFsAPISuite) TestFindAbsolutePathMultipleFiles() {
	suite.uploadStringContent("dirA/file.txt", "123")
	suite.uploadStringContent("dirB/file.txt", "98765")
	result, err := suite.findAbsolutePath("file.txt")
	suite.NoError(err)
	suite.Equal("/buckets/bfsdefault/default/dirA/file.txt", result)
}

func (suite *BucketFsAPISuite) listFiles() ([]BfsFile, error) {
	bfsClient := CreateBucketFsAPI(DEFAULT_BUCKET_PATH)
	return bfsClient.ListFiles(context.Background(), suite.exasol.GetConnection())
}

func (suite *BucketFsAPISuite) findAbsolutePath(fileName string) (string, error) {
	bfsClient := CreateBucketFsAPI(DEFAULT_BUCKET_PATH)
	return bfsClient.FindAbsolutePath(context.Background(), suite.exasol.GetConnection(), fileName)
}

func (suite *BucketFsAPISuite) uploadStringContent(fileName string, fileContent string) {
	err := suite.exasol.Exasol.UploadStringContent(fileContent, fileName)
	if err != nil {
		suite.FailNowf("Failed to upload file %q. Cause: %w", fileName, err)
	}
	suite.T().Cleanup(func() {
		suite.NoError(suite.exasol.Exasol.DeleteFile(fileName))
	})
}
