package bfs_test

import (
	"context"
	"fmt"
	"time"

	"testing"

	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
	"github.com/exasol/extension-manager/pkg/integrationTesting"
	"github.com/stretchr/testify/suite"
)

const DEFAULT_BUCKET_PATH = "/buckets/bfsdefault/default/"

type BucketFsAPIITestSuite struct {
	suite.Suite
	exasol integrationTesting.DbTestSetup
}

func TestBucketFsApiITestSuite(t *testing.T) {
	suite.Run(t, new(BucketFsAPIITestSuite))
}

func (suite *BucketFsAPIITestSuite) SetupSuite() {
	suite.exasol = *integrationTesting.StartDbSetup(&suite.Suite)
}

func (suite *BucketFsAPIITestSuite) TeardownSuite() {
	suite.exasol.StopDb()
}

func (suite *BucketFsAPIITestSuite) BeforeTest(suiteName, testName string) {
	suite.exasol.CreateConnection()
	suite.T().Cleanup(func() {
		suite.exasol.CloseConnection()
	})
}

/* [utest -> dsn~extension-components~1]. */
func (suite *BucketFsAPIITestSuite) TestListEmptyDir() {
	result, err := suite.listFiles()
	suite.NoError(err)
	suite.Empty(result)
}

/* [itest -> dsn~configure-bucketfs-path~1]. */
func (suite *BucketFsAPIITestSuite) TestListSingleFile() {
	fileName := fmt.Sprintf("myFile-%d.txt", time.Now().Unix())
	suite.uploadStringContent(fileName, "12345")
	result, err := suite.listFiles()
	suite.NoError(err)
	suite.Len(result, 1)
	suite.Equal([]bfs.BfsFile{{Name: fileName, Path: DEFAULT_BUCKET_PATH + fileName, Size: 5}}, result)
}

func (suite *BucketFsAPIITestSuite) TestListFilesRecursively() {
	file1 := "file1"
	file2 := "dir1/file2"
	file3 := "dir2/file2"
	suite.uploadStringContent(file1, "1")
	suite.uploadStringContent(file2, "12")
	suite.uploadStringContent(file3, "123")
	result, err := suite.listFiles()
	suite.NoError(err)
	suite.Len(result, 3)
	suite.Equal([]bfs.BfsFile{
		{Name: "file2", Path: DEFAULT_BUCKET_PATH + file2, Size: 2},
		{Name: "file2", Path: DEFAULT_BUCKET_PATH + file3, Size: 3},
		{Name: "file1", Path: DEFAULT_BUCKET_PATH + file1, Size: 1}}, result)
}

func (suite *BucketFsAPIITestSuite) TestFindAbsolutePathNoFileFound() {
	time.Sleep(3 * time.Second)
	result, err := suite.findAbsolutePath("no-such-file")
	suite.EqualError(err, `file "no-such-file" not found in BucketFS`)
	suite.Equal("", result)
}

/* [itest -> dsn~configure-bucketfs-path~1] */
/* [itest -> dsn~resolving-files-in-bucketfs~1]. */
func (suite *BucketFsAPIITestSuite) TestFindAbsolutePathFileInRoot() {
	fileName := "file01.txt"
	suite.uploadStringContent(fileName, "123")
	result, err := suite.findAbsolutePath(fileName)
	suite.NoError(err)
	suite.Equal("/buckets/bfsdefault/default/"+fileName, result)
}

func (suite *BucketFsAPIITestSuite) TestFindAbsolutePathFileInSubDir() {
	suite.uploadStringContent("dir/file02.txt", "123")
	result, err := suite.findAbsolutePath("file02.txt")
	suite.NoError(err)
	suite.Equal("/buckets/bfsdefault/default/dir/file02.txt", result)
}

func (suite *BucketFsAPIITestSuite) TestFindAbsolutePathMultipleFiles() {
	suite.uploadStringContent("dirA/file03.txt", "123")
	suite.uploadStringContent("dirB/file03.txt", "98765")
	result, err := suite.findAbsolutePath("file03.txt")
	suite.NoError(err)
	suite.Equal("/buckets/bfsdefault/default/dirA/file03.txt", result)
}

func (suite *BucketFsAPIITestSuite) TestFindAbsolutePathMultipleFilesFirstFileDoesNotMatch() {
	suite.uploadStringContent("dirA/file05.txt", "123")
	suite.uploadStringContent("dirB/file06.txt", "98765")
	result, err := suite.findAbsolutePath("file06.txt")
	suite.NoError(err)
	suite.Equal("/buckets/bfsdefault/default/dirB/file06.txt", result)
}

func (suite *BucketFsAPIITestSuite) listFiles() ([]bfs.BfsFile, error) {
	bfsClient := suite.createBucketFsClient()
	defer bfsClient.Close()
	return bfsClient.ListFiles()
}

func (suite *BucketFsAPIITestSuite) findAbsolutePath(fileName string) (string, error) {
	bfsClient := suite.createBucketFsClient()
	defer bfsClient.Close()
	return bfsClient.FindAbsolutePath(fileName)
}

func (suite *BucketFsAPIITestSuite) createBucketFsClient() bfs.BucketFsAPI {
	bfsClient, err := bfs.CreateBucketFsAPI(DEFAULT_BUCKET_PATH, context.Background(), suite.exasol.GetConnection())
	if err != nil {
		suite.FailNow("Creating BFS API failed: " + err.Error())
	}
	return bfsClient
}

func (suite *BucketFsAPIITestSuite) uploadStringContent(fileName string, fileContent string) {
	err := suite.exasol.Exasol.UploadStringContent(fileContent, fileName)
	if err != nil {
		suite.FailNowf("Failed to upload file %q. Cause: %w", fileName, err)
	}
	suite.T().Cleanup(func() {
		suite.NoError(suite.exasol.Exasol.DeleteFile(fileName))
	})
}
