package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	caddyfile = `
0.0.0.0:3000 {
    proxy / localhost:3010
    proxy /api localhost:3020 {
        without /api
    }
}
`
)

func ActionDevSetup(c *cli.Context) error {
	if err := os.MkdirAll("./dev", 0755); err != nil {
		return errors.Wrap(err, "failed to create ./dev/ for development")
	}
	log.Println("Created ./dev/")

	if err := setupCaddy(); err != nil {
		return err
	}

	return nil
}

func setupCaddy() error {
	if err := downloadCaddy(); err != nil {
		return errors.Wrap(err, "failed to download caddy")
	}
	log.Println("Downloaded ./dev/caddy.zip")

	if err := extractCaddy(); err != nil {
		return errors.Wrap(err, "failed to extract caddy")
	}
	log.Println("Extracted ./dev/caddy.zip to ./dev/caddy")

	// Create Caddyfile with contents if it not exist
	if err := createCaddyfile(); err != nil {
		return errors.Wrap(err, "failed to create ./dev/Caddyfile")
	}
	log.Println("Created ./dev/Caddyfile")

	return nil
}

func downloadCaddy() error {
	caddyURL := fmt.Sprintf("https://caddyserver.com/download/%s/%s", runtime.GOOS, runtime.GOARCH)

	// Download & extract Caddy to ./dev/caddy if it not exist

	exist, err := exists("./dev/caddy")
	if err != nil {
		return err
	}
	if exist { // if it exist, don't do more
		return nil
	}

	exist, err = exists("./dev/caddy.zip")
	if err != nil {
		return err
	}
	if exist { // if it exist, don't do more
		return nil
	}

	out, err := os.Create("./dev/caddy.zip")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(caddyURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func extractCaddy() error {
	r, err := zip.OpenReader("./dev/caddy.zip")
	if err != nil {
		return err
	}

	var zippedCaddy *zip.File
	for _, file := range r.File {
		if file.Name == "caddy" {
			zippedCaddy = file
		}
	}

	fileReader, err := zippedCaddy.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	targetFile, err := os.OpenFile("./dev/caddy", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zippedCaddy.Mode())
	if err != nil {
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, fileReader); err != nil {
		return err
	}

	return nil
}

func createCaddyfile() error {
	exist, err := exists("./dev/Caddyfile")
	if err != nil {
		return err
	}
	if !exist {
		if err := ioutil.WriteFile("./dev/Caddyfile", []byte(strings.TrimSpace(caddyfile)), 0644); err != nil {
			return err
		}
	}

	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
