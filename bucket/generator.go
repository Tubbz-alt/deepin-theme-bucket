package bucket

import (
	"fmt"
	"github.com/Iceyer/go-sdk/upyun/form"
	"io"
	"net/http"
)

var (
	BucketHost = "http://theme-store.b0.upaiyun.com"
)

type generator interface {
	GetURL(datatype, id string) (string, error)
	Get(datatype string, id string) (io.ReadCloser, error)
	Put(datatype string, r io.Reader) error
	putDir(rootPath string) error
}

type creator interface {
	Get(id string) (io.ReadCloser, error)
	GetURL(id string) (string, error)
	Put(id string, r io.Reader) error
}

type urlCreator struct {
	client      *http.Client
	urlTemplate string
}

func (uc *urlCreator) remotePath(id string) string {
	return fmt.Sprintf(uc.urlTemplate, id)
}

func (uc *urlCreator) GetURL(id string) (string, error) {
	return BucketHost + fmt.Sprintf(uc.urlTemplate, id), nil
}

func (uc *urlCreator) Get(id string) (io.ReadCloser, error) {
	url, _ := uc.GetURL(id)
	rsp, err := uc.client.Get(url)
	if nil != err {
		return nil, err
	}
	if http.StatusOK != rsp.StatusCode {
		return nil, fmt.Errorf(rsp.Status + url)
	}
	return rsp.Body, nil
}

var upform = form.NewUpForm("theme-store", &remoteSignature{})

func (uc *urlCreator) Put(id string, r io.Reader) error {
	err := upform.SlicePostData(r, uc.remotePath(id))
	return err
}
