package shell

import (
	"bytes"
	"context"
	"fmt"
)

type BucketStat struct {
	Name        string
	BucketID    int64
	Ctime       string
	Policy      int32
	DataCount   int32
	ParityCount int32
	Encryption  int32
}

type Buckets struct {
	Method  string
	Buckets []BucketStat
}

func (bk BucketStat) String() string {
	return fmt.Sprintf(
		"Name: %s\n--BucketID: %d\n--Ctime: %s\n--Policy: %d\n--DataCount: %d\n--ParityCount: %d\n--Encryption: %d\n",
		bk.Name,
		bk.BucketID,
		bk.Ctime,
		bk.Policy,
		bk.DataCount,
		bk.ParityCount,
		bk.Encryption,
	)
}

func (bus Buckets) String() string {
	var str bytes.Buffer
	str.WriteString("Method: " + bus.Method + "\n")
	for _, buStat := range bus.Buckets {
		str.WriteString(buStat.String())
	}
	return str.String()
}

func (s *Shell) HeadBucket(BucketName string, options ...LfsOpts) (*Buckets, error) {
	var bks Buckets
	rb := s.Request("lfs/head_bucket", BucketName)
	for _, option := range options {
		option(rb)
	}

	if err := rb.Exec(context.Background(), &bks); err != nil {
		return nil, err
	}
	return &bks, nil
}

func (s *Shell) ListBuckets(options ...LfsOpts) (*Buckets, error) {
	var bks Buckets
	rb := s.Request("lfs/list_buckets")
	for _, option := range options {
		option(rb)
	}
	if err := rb.Exec(context.Background(), &bks); err != nil {
		return nil, err
	}
	return &bks, nil
}

func (s *Shell) CreateBucket(BucketName string, options ...LfsOpts) (*Buckets, error) {
	var bk Buckets
	rb := s.Request("lfs/create_bucket", BucketName)
	for _, option := range options {
		option(rb)
	}
	if err := rb.Exec(context.Background(), &bk); err != nil {
		return nil, err
	}
	return &bk, nil
}

func (s *Shell) DeleteBucket(BucketName string, options ...LfsOpts) (*Buckets, error) {
	var bk Buckets
	rb := s.Request("lfs/delete_bucket", BucketName)
	for _, option := range options {
		option(rb)
	}
	if err := rb.Exec(context.Background(), &bk); err != nil {
		return nil, err
	}
	return &bk, nil
}
