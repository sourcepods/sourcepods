package ssh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/pkg/errors"
)

type apiStorage struct {
	client *http.Client
}

type lfsPayload struct {
	Href    string `json:"href"`
	Headers struct {
		Authorization string `json:"Authorization"`
	} `json:"headers"`
	ExpiresAt string `json:"expires_at"`
}

func (s *apiStorage) LFSUpload(ctx context.Context, sess ssh.Session) error {
	path, err := parsePath(ctx)
	if err != nil {
		return err
	}

	// Token is fetched based on the PK...
	token := "foobar"

	payload := lfsPayload{
		Href:      fmt.Sprintf("https://example.com/%s/info/lfs/", path),
		ExpiresAt: time.Now().Add(5 * time.Minute).Format(time.RFC3339),
	}
	payload.Headers.Authorization = fmt.Sprintf("Bearer %s", token)

	b, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	_, err = sess.Write(b)
	return err
}

func (s *apiStorage) LFSDownload(ctx context.Context, sess ssh.Session) error {
	// NOTE: Cheating...
	return s.LFSUpload(ctx, sess)
}
