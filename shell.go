// package shell implements a remote API interface for a running ipfs daemon
package shell

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	gohttp "net/http"
	"os"
	"path"
	"strings"
	"time"

	files "github.com/ipfs/go-ipfs-files"
	homedir "github.com/mitchellh/go-homedir"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

const (
	DefaultPathName = ".mefs"
	DefaultPathRoot = "~/" + DefaultPathName
	DefaultApiFile  = "api"
	EnvDir          = "MEFS_PATH"
)

type Shell struct {
	url     string
	httpcli gohttp.Client
}

func NewLocalShell() *Shell {
	baseDir := os.Getenv(EnvDir)
	if baseDir == "" {
		baseDir = DefaultPathRoot
	}

	baseDir, err := homedir.Expand(baseDir)
	if err != nil {
		return nil
	}

	apiFile := path.Join(baseDir, DefaultApiFile)

	if _, err := os.Stat(apiFile); err != nil {
		return nil
	}

	api, err := ioutil.ReadFile(apiFile)
	if err != nil {
		return nil
	}

	return NewShell(strings.TrimSpace(string(api)))
}

func NewShell(url string) *Shell {
	c := &gohttp.Client{
		Transport: &gohttp.Transport{
			Proxy:             gohttp.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	}

	return NewShellWithClient(url, c)
}

func NewShellWithClient(url string, c *gohttp.Client) *Shell {
	if a, err := ma.NewMultiaddr(url); err == nil {
		_, host, err := manet.DialArgs(a)
		if err == nil {
			url = host
		}
	}
	var sh Shell
	sh.url = url
	sh.httpcli = *c
	// We don't support redirects.
	sh.httpcli.CheckRedirect = func(_ *gohttp.Request, _ []*gohttp.Request) error {
		return fmt.Errorf("unexpected redirect")
	}
	return &sh
}

func (s *Shell) SetTimeout(d time.Duration) {
	s.httpcli.Timeout = d
}

func (s *Shell) Request(command string, args ...string) *RequestBuilder {
	return &RequestBuilder{
		command: command,
		args:    args,
		shell:   s,
	}
}

type IdOutput struct {
	NetworkAddr  string
	AccountAddr  string
	PublicKey    string
	Addresses    []string
	AgentVersion string
}

// ID gets information about a given peer.  Arguments:
//
// peer: peer.ID of the node to look up.  If no peer is specified,
//   return information about the local peer.
func (s *Shell) ID(peer ...string) (*IdOutput, error) {
	if len(peer) > 1 {
		return nil, fmt.Errorf("Too many peer arguments")
	}

	var out IdOutput
	if err := s.Request("id", peer...).Exec(context.Background(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

const (
	DirectPin    = "direct"
	RecursivePin = "recursive"
	IndirectPin  = "indirect"
)

type PeerInfo struct {
	Addrs []string
	ID    string
}

func (s *Shell) FindPeer(peer string) (*PeerInfo, error) {
	var peers struct{ Responses []PeerInfo }
	err := s.Request("dht/findpeer", peer).Exec(context.Background(), &peers)
	if err != nil {
		return nil, err
	}
	if len(peers.Responses) == 0 {
		return nil, errors.New("peer not found")
	}
	return &peers.Responses[0], nil
}

func (s *Shell) ResolvePath(path string) (string, error) {
	var out struct {
		Path string
	}
	err := s.Request("resolve", path).Exec(context.Background(), &out)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(out.Path, "/ipfs/"), nil
}

// returns ipfs version and commit sha
func (s *Shell) Version() (string, string, error) {
	ver := struct {
		Version string
		Commit  string
	}{}

	if err := s.Request("version").Exec(context.Background(), &ver); err != nil {
		return "", "", err
	}
	return ver.Version, ver.Commit, nil
}

func (s *Shell) IsUp() bool {
	_, _, err := s.Version()
	return err == nil
}

func (s *Shell) BlockStat(path string) (string, int, error) {
	var inf struct {
		Key  string
		Size int
	}

	if err := s.Request("block/stat", path).Exec(context.Background(), &inf); err != nil {
		return "", 0, err
	}
	return inf.Key, inf.Size, nil
}

func (s *Shell) BlockGet(path string) ([]byte, error) {
	resp, err := s.Request("block/get", path).Send(context.Background())
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	if resp.Error != nil {
		return nil, resp.Error
	}

	return ioutil.ReadAll(resp.Output)
}

func (s *Shell) BlockPut(block []byte, format, mhtype string, mhlen int) (string, error) {
	var out struct {
		Key string
	}

	fr := files.NewBytesFile(block)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", fr)})
	fileReader := files.NewMultiFileReader(slf, true)

	return out.Key, s.Request("block/put").
		Option("mhtype", mhtype).
		Option("format", format).
		Option("mhlen", mhlen).
		Body(fileReader).
		Exec(context.Background(), &out)
}

type SwarmStreamInfo struct {
	Protocol string
}

type SwarmConnInfo struct {
	Addr    string
	Peer    string
	Latency string
	Muxer   string
	Streams []SwarmStreamInfo
}

type SwarmConnInfos struct {
	Peers []SwarmConnInfo
}

// SwarmPeers gets all the swarm peers
func (s *Shell) SwarmPeers(ctx context.Context) (*SwarmConnInfos, error) {
	v := &SwarmConnInfos{}
	err := s.Request("swarm/peers").Exec(ctx, &v)
	return v, err
}

type swarmConnection struct {
	Strings []string
}

// SwarmConnect opens a swarm connection to a specific address.
func (s *Shell) SwarmConnect(ctx context.Context, addr ...string) error {
	var conn *swarmConnection
	err := s.Request("swarm/connect").
		Arguments(addr...).
		Exec(ctx, &conn)
	return err
}
