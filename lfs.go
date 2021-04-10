package shell

import (
	"bytes"
	"context"
	"strconv"
)

type UserPrivMessage struct {
	Address string
	Sk      string
}

type StringList struct {
	ChildLists []string
}

func (fl StringList) String() string {
	var buffer bytes.Buffer
	for i := 0; i < len(fl.ChildLists); i++ {
		buffer.WriteString(fl.ChildLists[i])
		buffer.WriteString("\n")
	}
	return buffer.String()
}

type IntList struct {
	ChildLists []int
}

func (fl IntList) String() string {
	var buffer bytes.Buffer
	for i := 0; i < len(fl.ChildLists); i++ {
		buffer.WriteString(strconv.Itoa(fl.ChildLists[i]))
		buffer.WriteString("\n")
	}
	return buffer.String()
}

type LfsOpts = func(*RequestBuilder) error

func SetAddress(addr string) LfsOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("address", addr)
		return nil
	}
}

func SetOp(key, value string) LfsOpts {
	return func(rb *RequestBuilder) error {
		rb.Option(key, value)
		return nil
	}
}

func SetObjectName(objectName string) LfsOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("objectname", objectName)
		return nil
	}
}

func SetPrefixFilter(prefix string) LfsOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("prefix", prefix)
		return nil
	}
}

func SetPolicy(policy int) LfsOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("policy", policy)
		return nil
	}
}

func SetDataCount(dataCount int) LfsOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("datacount", dataCount)
		return nil
	}
}

func SetParityCount(parityCount int) LfsOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("paritycount", parityCount)
		return nil
	}
}

func SetCrypto(crypto bool) LfsOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("encryption", crypto)
		return nil
	}
}

func SetAvailTime(enabled bool) LfsOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("availTime", enabled)
		return nil
	}
}

func SetSecretKey(sk string) LfsOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("secretekey", sk)
		return nil
	}
}
func SetPassword(pwd string) LfsOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("password", pwd)
		return nil
	}
}

func ForceFlush(enabled bool) LfsOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("force", enabled)
		return nil
	}
}

func UseErasureCodeOrMulRep(enabled bool) LfsOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("policy", enabled)
		return nil
	}
}

type PeerState struct {
	PeerID    string
	Connected bool
}

func (ps PeerState) String() string {
	if ps.Connected {
		return ps.PeerID + " connected"
	}
	return ps.PeerID + " unconnected"
}

type PeerList struct {
	Peers []PeerState
}

func (pl PeerList) String() string {
	var res string
	for _, ps := range pl.Peers {
		res += ps.String() + "\n"
	}
	return res
}

type queryEvent struct {
	ID    string
	Extra string
}

type GetBlockResult struct {
	IsExist bool
}

type BlockStat struct {
	Key  string
	Size int
}

func (s *Shell) CreateUser(options ...LfsOpts) (*UserPrivMessage, error) {
	var user UserPrivMessage
	rb := s.Request("create")
	for _, option := range options {
		option(rb)
	}

	if err := rb.Exec(context.Background(), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Shell) StartUser(address string, options ...LfsOpts) error {
	var res StringList
	rb := s.Request("lfs/start", address)
	for _, option := range options {
		option(rb)
	}
	if err := rb.Exec(context.Background(), &res); err != nil {
		return err
	}
	return nil
}

func (s *Shell) Fsync(options ...LfsOpts) error {
	var res StringList
	rb := s.Request("lfs/fsync")
	for _, option := range options {
		option(rb)
	}

	if err := rb.Exec(context.Background(), &res); err != nil {
		return err
	}
	return nil
}

func (s *Shell) ShowStorage(options ...LfsOpts) (string, error) {
	var res string
	rb := s.Request("lfs/show_storage")
	for _, option := range options {
		option(rb)
	}

	if err := rb.Exec(context.Background(), &res); err != nil {
		return "", err
	}
	return res, nil
}

func (s *Shell) ListKeepers(options ...LfsOpts) (*PeerList, error) {
	var res *PeerList
	rb := s.Request("lfs/list_keepers")
	for _, option := range options {
		option(rb)
	}

	if err := rb.Exec(context.Background(), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Shell) DeleteFrom(key, to string, options ...LfsOpts) (string, error) {
	var res string
	rb := s.Request("dht/deletefrom", key, to)
	for _, option := range options {
		option(rb)
	}

	if err := rb.Exec(context.Background(), &res); err != nil {
		return "", err
	}
	return res, nil
}

func (s *Shell) GetFrom(key, id string, options ...LfsOpts) (*queryEvent, error) {
	var res *queryEvent
	rb := s.Request("dht/getfrom", key, id)
	for _, option := range options {
		option(rb)
	}

	if err := rb.Exec(context.Background(), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Shell) GetBlockFrom(key, id string, options ...LfsOpts) (string, error) {
	var res string
	rb := s.Request("block/getfrom", key, id)
	for _, option := range options {
		option(rb)
	}

	if err := rb.Exec(context.Background(), &res); err != nil {
		return "", err
	}
	return res, nil
}

func (s *Shell) Kill(addr string, options ...LfsOpts) (*StringList, error) {
	var strlist StringList
	rb := s.Request("lfs/kill", addr)
	for _, option := range options {
		option(rb)
	}
	if err := rb.Exec(context.Background(), &strlist); err != nil {
		return nil, err
	}
	return &strlist, nil
}
