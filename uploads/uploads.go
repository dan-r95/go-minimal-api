package uploads

import (
	"context"
	"errors"
	"fmt"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/gcerrors"
	"net/url"
	"path/filepath"
	"strconv"
)

type (
	// BlobStore is the representation of the file persistence
	BlobStore struct {
		base  string
		store blob.BucketURLOpener
		ctx   context.Context
	}

	Store struct {
		Save     func() error
		Close    func()
		MimeType string
		Length   uint64
	}
)

// NewStorage creates a new blob store
func NewStorage(ctx context.Context) (*BlobStore, error) {
	//adding  a  cloud store is possible
	opener :=
		&fileblob.URLOpener{
			Options: fileblob.Options{
				CreateDir: true,
				Metadata:  fileblob.MetadataDontWrite,
			},
		}
	path := filepath.Join("data")
	return &BlobStore{ctx: ctx, store: opener, base: filepath.Clean(path)}, nil
}

// GetImage returns the image with name and size from the storage. if specific size
// cannot be found, it tries to resize the image and store it.
// if the image has not been found at all, returns an error
func (bs *BlobStore) GetImage(name string, size string) (b []byte, err error) {
	var buck *blob.Bucket
	if buck, err = bs.getBucket(name); err == nil {
		defer func() { _ = buck.Close() }()
		b, err = buck.ReadAll(bs.ctx, size)
		//TODO: find if original picture was there
		if gcerrors.Code(err) == gcerrors.NotFound {
			err = errors.New("not found")
		}
	}
	//TODO: also save the image format somewhere
	//TODO: if size not there yet, resize it and store
	//ResizeImage()
	//err = bs.StoreImage(name, size, b)
	return
}

// StoreRawImage reads the image contents from the filer header and on success
// stores it in the blob store
func (bs *BlobStore) StoreRawImage(name string, size string, req *ResourceRequest) error {
	data, _, err := ReadImageAndValidate(req.File)
	if err != nil {
		return err
	}

	return bs.StoreImage(name, "raw", data)
}

//StoreImage writes the image to the blob store. directory syntax is "imagename/size".
// original upload is called "imagename/raw"
func (bs *BlobStore) StoreImage(name string, size string, data []byte) error {
	buck, err := bs.getBucket(name)
	if err != nil {
		return err
	}

	defer func() { _ = buck.Close() }()
	return buck.WriteAll(bs.ctx, fmt.Sprintf("%s/%s", name, size), data, nil)
}

func (bs *BlobStore) keyString(key uint64) string {
	return strconv.FormatUint(key, 10)
}

func (bs *BlobStore) getBucket(name string) (*blob.Bucket, error) {
	return bs.bucket(name)
}

func (bs *BlobStore) bucket(seg ...string) (*blob.Bucket, error) {
	path := filepath.Join(seg...)
	path = filepath.Join(bs.base, path)
	return bs.store.OpenBucketURL(bs.ctx, &url.URL{Path: path})
}
