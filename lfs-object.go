package shell

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	files "github.com/ipfs/go-ipfs-files"
)

type ObjectStat struct {
	Name           string
	Size           int64
	MD5            string
	Ctime          string
	Dir            bool
	LatestChalTime string
}

type Objects struct {
	Method  string
	Objects []ObjectStat
}

var (
	errLfsServiceNotReady   = errors.New("lfs service not ready")
	errGroupServiceNotReady = errors.New("group service not ready")
)

func (ob ObjectStat) String() string {
	FloatStorage := float64(ob.Size)
	var OutStorage string
	if FloatStorage < 1024 && FloatStorage > 0 {
		OutStorage = fmt.Sprintf("%.2f", FloatStorage) + "B"
	} else if FloatStorage < 1048576 && FloatStorage >= 1024 {
		OutStorage = fmt.Sprintf("%.2f", FloatStorage/1024) + "KB"
	} else if FloatStorage < 1073741824 && FloatStorage >= 1048576 {
		OutStorage = fmt.Sprintf("%.2f", FloatStorage/1048576) + "MB"
	} else {
		OutStorage = fmt.Sprintf("%.2f", FloatStorage/1073741824) + "GB"
	}
	return fmt.Sprintf(
		"ObjectName: %s\n--ObjectSize: %s\n--MD5: %s\n--Ctime: %s\n--Dir: %t\n--LatestChalTime: %s\n",
		ob.Name,
		OutStorage,
		ob.MD5,
		ob.Ctime,
		ob.Dir,
		ob.LatestChalTime,
	)
}

func (obs Objects) String() string {
	var str bytes.Buffer
	str.WriteString("Method: " + obs.Method + "\n")
	for _, obStat := range obs.Objects {
		str.WriteString(obStat.String())
	}
	return str.String()
}

func (s *Shell) HeadObject(ObjectName, BucketName string, options ...LfsOpts) (*Objects, error) {
	var objs Objects
	rb := s.Request("lfs/head_object", BucketName, ObjectName)
	for _, option := range options {
		option(rb)
	}

	if err := rb.Exec(context.Background(), &objs); err != nil {
		return nil, err
	}
	return &objs, nil
}

func (s *Shell) GetObject(ObjectName, BucketName string, options ...LfsOpts) (io.ReadCloser, error) {
	var err error
	rb := s.Request("lfs/get_object", BucketName, ObjectName)
	for _, option := range options {
		option(rb)
	}
	resp, err := rb.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.Output, err
}

func (s *Shell) GenShare(ObjectName, BucketName string, options ...LfsOpts) (string, error) {
	var slink string
	rb := s.Request("lfs/gen_share", BucketName, ObjectName)
	for _, option := range options {
		option(rb)
	}
	if err := rb.Exec(context.Background(), &slink); err != nil {
		return "", err
	}
	return slink, nil
}

func (s *Shell) GetShare(shareLink, outputName string, options ...LfsOpts) (io.ReadCloser, error) {
	rb := s.Request("lfs/get_share", shareLink, outputName)
	for _, option := range options {
		option(rb)
	}
	resp, err := rb.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.Output, err
}

func (s *Shell) GetObjectToFile(ObjectName, BucketName, outPath string, options ...LfsOpts) error {
	var err error
	var p string
	rootExists := true
	rootIsDir := false
	if stat, err := os.Stat(outPath); err != nil && os.IsNotExist(err) {
		rootExists = false
	} else if err != nil {
		return err
	} else if stat.IsDir() {
		rootIsDir = true
	}
	if rootIsDir == true {
		p = path.Join(outPath, ObjectName)
	} else if rootExists == false {
		p = outPath
	} else {
		return errors.New("The outpath already has file: " + ObjectName)
	}
	var file *os.File
	if _, err := os.Stat(p); err != nil && os.IsNotExist(err) {
		file, err = os.Create(p)
		if err != nil {
			return err
		}
	} else {
		return errors.New("The outpath already has file: " + ObjectName)
	}

	defer file.Close()

	rb := s.Request("lfs/get_object", BucketName, ObjectName)
	for _, option := range options {
		option(rb)
	}
	resp, err := rb.Send(context.Background())
	if err != nil {
		return err
	}
	written, err := io.Copy(file, resp.Output)
	if err != nil {
		fmt.Println("Download", ObjectName, " err", err)
	}
	fmt.Println("Download", ObjectName, "finish, write data", written)
	return err
}

func (s *Shell) ListObjects(BucketName string, options ...LfsOpts) (*Objects, error) {
	var objs Objects
	rb := s.Request("lfs/list_objects", BucketName)
	for _, option := range options {
		option(rb)
	}

	if err := rb.Exec(context.Background(), &objs); err != nil {
		return nil, err
	}
	return &objs, nil
}

func (s *Shell) PutObject(r io.Reader, ObjectName, BucketName string, options ...LfsOpts) (*Objects, error) {
	fr := files.NewReaderFile(r)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", fr)})
	fileReader := files.NewMultiFileReader(slf, true)
	var objs Objects
	rb := s.Request("lfs/put_object", BucketName)
	for _, option := range options {
		option(rb)
	}
	rb.Option("objectname", ObjectName)
	rb = rb.Body(fileReader)
	if err := rb.Exec(context.Background(), &objs); err != nil {
		return nil, err
	}
	return &objs, nil
}

func (s *Shell) DeleteObject(ObjectName, BucketName string, options ...LfsOpts) (*Objects, error) {
	var objs Objects
	rb := s.Request("lfs/delete_object", BucketName)
	for _, option := range options {
		option(rb)
	}

	if err := rb.Exec(context.Background(), &objs); err != nil {
		return nil, err
	}
	return &objs, nil
}
