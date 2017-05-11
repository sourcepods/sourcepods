package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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
	if err := os.MkdirAll("./dev/sessions/", 0755); err != nil {
		return errors.Wrap(err, "failed to create ./dev/sessions/ for development")
	}
	log.Println("Created ./dev/")

	if err := setupNodeModules(); err != nil {
		return err
	}

	if err := setupCaddy(); err != nil {
		return err
	}

	return nil
}

func setupNodeModules() error {
	yarn, err := exec.LookPath("yarn")
	if err != nil {
		return err
	}

	cmd := exec.Command(yarn)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func setupCaddy() error {
	url := fmt.Sprintf("https://caddyserver.com/download/%s/%s", runtime.GOOS, runtime.GOARCH)
	archive := ""

	switch runtime.GOOS {
	case "darwin":
		archive = "./dev/caddy.zip"
	default:
		archive = "./dev/caddy.tar.gz"
	}

	if err := downloadCaddy(url, archive); err != nil {
		return errors.Wrap(err, "failed to download caddy")
	}
	log.Println("Downloaded ./dev/caddy.zip")

	if err := extractCaddy(archive); err != nil {
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

func downloadCaddy(url, archive string) error {

	// Download & extract Caddy to ./dev/caddy if it not exist

	exist, err := exists("./dev/caddy")
	if err != nil {
		return err
	}
	if exist { // if it exist, don't do more
		return nil
	}

	exist, err = exists(archive)
	if err != nil {
		return err
	}
	if exist { // if it exist, don't do more
		return nil
	}

	out, err := os.Create(archive)
	if err != nil {
		return err
	}
	defer out.Close()

	log.Println("Downloading", url)
	resp, err := http.Get(url)
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

func extractCaddy(archive string) error {
	switch runtime.GOOS {
	case "darwin":
		r, err := zip.OpenReader(archive)
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
	default:
		archiveFile, err := os.Open(archive)
		if err != nil {
			return err
		}
		defer archiveFile.Close()

		gzipREader, err := gzip.NewReader(archiveFile)
		if err != nil {
			return err
		}

		tarReader := tar.NewReader(gzipREader)
		for {
			header, err := tarReader.Next()

			if err == io.EOF {
				break
			}

			if err != nil {
				return err
			}

			if header.Name == "caddy" && header.Typeflag == tar.TypeReg {
				f, err := os.Create("./dev/caddy")
				if err != nil {
					return err
				}
				defer f.Close()

				if err := f.Chmod(0744); err != nil {
					return err
				}

				if _, err := io.Copy(f, tarReader); err != nil {
					return err
				}
				return nil
			}
		}
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
