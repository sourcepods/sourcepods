package gitssh

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/opentracing/opentracing-go"

	"github.com/gliderlabs/ssh"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/sourcepods/sourcepods/pkg/storage"
)

type Service interface {
	Handler() ssh.Handler
	ListenAndServe() error
	Shutdown(context.Context) error
}

// NewSSHServer returns a *grpc.Server serving SSH
//  is no `hostKeyPath` is given, random hostkeys will be generated...
func NewSSHServer(addr string, hostKeyPath string, logger log.Logger, cli *storage.Client) *ssh.Server {
	s := &ssh.Server{
		Addr:    addr,
		Handler: logHandler(mainHandler(cli), logger),
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

func mainHandler(cli *storage.Client) ssh.Handler {
	return func(s ssh.Session) {
		cmd := s.Command()
		if len(cmd) < 1 {
			fmt.Fprintf(s, "Welcome to SourcePods, %s\n", s.User())
			return
		}
		switch cmd[0] {
		case "git", "git-upload-pack", "git-receive-pack":
			storageHandler(cli, s)
		default:
			fmt.Fprintf(s, "unknown command given\n")
			s.Exit(1)
		}
	}
}

// NOTE: Windows sucks... sends "git upload-pack 'path/to/repo.git'" instead of "git-upload-pack 'path/to/repo.git'"
func windowsSucks(command []string) []string {
	if command[0] == "git" {
		command[0] = fmt.Sprintf("%s-%s", command[0], command[1])
		return append(command[0:1], command[2:]...)
	}
	return command
}

func storageHandler(cli *storage.Client, s ssh.Session) {
	ctx := s.Context().Value("span-ctx").(context.Context)
	span, ctx := opentracing.StartSpanFromContext(ctx, "ssh.StorageHandler")
	defer span.Finish()

	command := windowsSucks(s.Command())

	id := command[1]

	span.SetTag("repo_path", id)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	switch command[0] {
	case "git-upload-pack":
		ec, err := cli.UploadPack(ctx, id, s, s, s.Stderr())
		if err != nil {
			logger := s.Context().Value("logger").(log.Logger)
			level.Error(logger).Log(
				"msg", "upload-pack failed",
				"err", err.Error(),
			)
			s.Exit(1)
		}
		s.Exit(int(ec))
	case "git-receive-pack":
		ec, err := cli.ReceivePack(ctx, id, s, s, s.Stderr())
		if err != nil {
			logger := s.Context().Value("logger").(log.Logger)
			level.Error(logger).Log(
				"msg", "recieve-pack failed",
				"err", err.Error(),
			)
			s.Exit(1)
		}
		s.Exit(int(ec))
	default:
		fmt.Fprintf(s, "unknown command given\n")
		s.Exit(1)
	}
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
