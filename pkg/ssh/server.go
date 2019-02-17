package ssh

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/go-kit/kit/log"
	"github.com/sourcepods/sourcepods/pkg/ssh/mux"
	"github.com/sourcepods/sourcepods/pkg/storage"
)

// NewServer returns a *grpc.Server serving SSH
//  if no `hostKeyPath` is given, random hostkeys will be generated...
func NewServer(addr, hostKeyPath string, logger log.Logger, cli *storage.Client) *ssh.Server {
	gs := &gitStorage{client: cli}

	m := mux.New()
	m.Use(mux.RecoverWare(logger))
	m.Use(tracingHandler())
	m.Use(logHandler(logger))

	m.AddHandler("^git[ -]receive-pack ([0-9a-f/-]+)$", "ssh.Handler.ReceivePack", mux.HandlerFunc(gs.ReceivePack))
	m.AddHandler("^git[ -]upload-pack ([0-9a-f/-]+)$", "ssh.Handler.UploadPack", mux.HandlerFunc(gs.UploadPack))

	s := &ssh.Server{
		Addr:    addr,
		Handler: m.Handle(),
		PublicKeyHandler: func(ssh.Context, ssh.PublicKey) bool {
			// TODO: This needs to be implemented :D
			return true
		},
	}
	if len(hostKeyPath) != 0 {
		opts, err := loadHostKeys(hostKeyPath)
		if err != nil {
			panic(err)
		}
		for _, opt := range opts {
			s.SetOption(opt)
		}
	}

	return s
}

func loadHostKeys(dir string) ([]ssh.Option, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("ReadDir: %v", err)
	}

	var res []ssh.Option

	for _, fi := range files {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(fi.Name(), ".pub") {
			continue
		}
		fullDir := filepath.Join(dir, fi.Name())
		res = append(res, ssh.HostKeyFile(fullDir))
	}
	return res, nil
}
