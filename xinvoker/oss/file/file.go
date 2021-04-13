package file

import (
	"errors"
	"github.com/coder2z/g-server/xinvoker/oss/standard"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
	isDelete bool
	bucket   string
	cdnName  string
}

func NewOss(cdnName string, bucket string, isDelete bool) (client standard.Oss, err error) {
	client = &Client{
		isDelete: isDelete,
		bucket:   bucket,
		cdnName:  cdnName,
	}
	return
}

func (c *Client) PutObject(dstPath string, reader io.Reader, options ...standard.Option) error {
	dstPath = c.bucket + "/" + dstPath
	err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
	if err != nil {
		return err
	}

	fileByte, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dstPath, fileByte, os.ModePerm)
}

func (c *Client) PutObjectFromFile(dstPath, srcPath string, options ...standard.Option) (err error) {
	dstPath = c.bucket + "/" + dstPath
	err = os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
	if err != nil {
		return
	}
	var b []byte
	b, err = ioutil.ReadFile(srcPath)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(dstPath, b, os.ModePerm)
	if err != nil {
		return
	}

	if c.isDelete {
		err = os.Remove(srcPath)
	}
	return
}

func (c *Client) GetObject(dstPath string, options ...standard.Option) (output []byte, err error) {
	return ioutil.ReadFile(dstPath)
}

func (c *Client) GetObjectToFile(dstPath, srcPath string, options ...standard.Option) error {
	return errors.New("file type is not supported GetObjectToFile method")
}

func (c *Client) DeleteObject(dstPath string) (err error) {
	err = os.Remove(strings.TrimLeft(dstPath, "/"))
	return
}

func (c *Client) DeleteObjects(dstPaths []string, options ...standard.Option) (output standard.DeleteObjectsResult, err error) {
	for _, filePath := range dstPaths {
		err1 := os.Remove(filePath)
		if err1 != nil {
			if err != nil {
				err = errors.New(err.Error() + ", err is " + err1.Error())
			} else {
				err = err1
			}
		}
	}
	return
}

func (c *Client) IsObjectExist(dstPath string) (bool, error) {
	_, err := os.Lstat(dstPath)
	return !os.IsNotExist(err), nil
}

func (c *Client) ListObjects(options ...standard.Option) (standard.ListObjectsResult, error) {
	return standard.ListObjectsResult{}, errors.New("file type is not supported ListObjects method")
}

func (c *Client) SignURL(dstPath string, method string, expiredInSec int64, options ...standard.Option) (resp string, err error) {
	resp = c.cdnName + dstPath
	return
}
