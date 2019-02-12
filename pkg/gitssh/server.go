package gitssh

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/opentracing/opentracing-go"

	"github.com/gliderlabs/ssh"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/sourcepods/sourcepods/pkg/storage"
)

// NewSSHServer returns a *grpc.Server serving SSH
func NewSSHServer(port int, hostKeyPath string, logger log.Logger, cli *storage.Client) *ssh.Server {
	s := &ssh.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mainHandler(cli, logger),
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

func sendPktLine(w io.Writer, format string, args ...interface{}) {
	pkt := fmt.Sprintf(format, args...)
	fmt.Fprintf(w, "%04X%s", 4+len(pkt), pkt)
}

func flushPkt(w io.Writer) {
	fmt.Fprint(w, "0000")
}

func readLine(r io.Reader) (string, error) {
	lenByte := make([]byte, 4)
	_, err := r.Read(lenByte)
	if err != nil {
		return "", err
	}
	n, err := strconv.ParseInt(string(lenByte[:]), 16, 16)
	if err != nil {
		return "", err
	}
	foo := make([]byte, n)
	n2, err := r.Read(foo)
	if err != nil {
		return "", err
	}
	if n2 != int(n) {
		return "", fmt.Errorf("invalid pkt-line: %d != %d", n2, n)
	}
	return string(foo), nil
}

func mainHandler(cli *storage.Client, logger log.Logger) ssh.Handler {
	return func(s ssh.Session) {
		defer s.Close()
		span, _ := opentracing.StartSpanFromContext(s.Context(), "ssh.MainHandler")
		span.SetTag("remote-addr", s.RemoteAddr().String())
		defer span.Finish()

		level.Info(logger).Log(
			"msg", "new connection",
			"user", s.User(),
			"command", fmt.Sprintf("%v", s.Command()),
		)
		cmd := s.Command()
		if len(cmd) < 1 {
			fmt.Fprintf(s, "Welcome to SourcePods, %s\n", s.User())
			return
		}
		switch cmd[0] {
		case "git", "git-upload-pack", "git-receive-pack":
			storageHandler(logger, cli, s)
		default:
			fmt.Fprintf(s, "unknown command given\n")
			s.Exit(1)
		}
	}
}

func storageHandler(logger log.Logger, cli *storage.Client, s ssh.Session) {
	command := s.Command()
	// NOTE: Windows sucks... sends "git upload-pack 'path/to/repo.git'" instead of "git-upload-pack 'path/to/repo.git'"
	if command[0] == "git" {
		command[0] = fmt.Sprintf("%s-%s", command[0], command[1])
		command = append(command[0:1], command[2:]...)
	}

	id := command[1]

	span, ctx := opentracing.StartSpanFromContext(s.Context(), "ssh.StorageHandler")
	span.SetTag("repo_path", id)
	defer span.Finish()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	switch command[0] {
	case "git-upload-pack":
		ec, err := cli.UploadPack(ctx, id, s, s, s.Stderr())
		if err != nil {
			level.Error(logger).Log(
				"msg", "upload-pack failed",
				"err", err.Error(),
			)
		}
		if ec != 0 {
			s.Exit(1)
		}
	case "git-receive-pack":
		ec, err := cli.ReceivePack(ctx, id, s, s, s.Stderr())
		if err != nil {
			level.Error(logger).Log(
				"msg", "recieve-pack failed",
				"err", err.Error(),
			)
		}
		if ec != 0 {
			s.Exit(1)
		}
	}
	s.Exit(0)
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
