package pkg

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"path"
	"time"

	minio "github.com/minio/minio-go"
)

//KutoCDN CDN实例
type KutoCDN struct {
	s3       *minio.Client
	endpoint string
	timeout  time.Duration
}

//NewCDN 新建CDN实例
func NewCDN(addr string, accessKeyID string, secretAccessKey string, t time.Duration) *KutoCDN {
	mc, err := minio.New(addr, accessKeyID, secretAccessKey, true)
	if err != nil {
		log.Fatal(err)
	}

	return &KutoCDN{
		s3:       mc,
		endpoint: addr,
		timeout:  t,
	}
}

//Put 上传CDN
func (cdn *KutoCDN) Put(bucket, object string, reader io.Reader, size int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := cdn.s3.PutObjectWithContext(ctx, bucket, object, reader, size, minio.PutObjectOptions{
		ContentType: mime.TypeByExtension(path.Ext(object)),
	})

	if err != nil {
		return err
	}

	return nil
}

//GenerateURL 生成url链接
func (cdn *KutoCDN) GenerateURL(bucket, object string) string {
	return fmt.Sprintf("https://%s/%s/%s", cdn.endpoint, bucket, object)
}
