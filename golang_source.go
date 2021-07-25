package elsa

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"shanhu.io/misc/errcode"
)

type golangSource struct {
	Version string
	SHA256  string
}

func createFileFrom(f string, r io.Reader) error {
	out, err := os.Create(f)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, r); err != nil {
		return err
	}

	return out.Sync()
}

func fileSHA256(f string) (string, error) {
	file, err := os.Open(f)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (s *golangSource) downloadTo(f string) error {
	name := fmt.Sprintf("go%s.src.tar.gz", s.Version)
	u := &url.URL{
		Scheme: "https",
		Host:   "dl.google.com",
		Path:   path.Join("/go", name),
	}

	req := &http.Request{
		Method: "GET",
		URL:    u,
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return errcode.Annotate(err, "get golang source")
	}
	defer resp.Body.Close()

	if err := createFileFrom(f, resp.Body); err != nil {
		return errcode.Annotate(err, "save golang source")
	}

	h, err := fileSHA256(f)
	if err != nil {
		return errcode.Annotate(err, "compute hash")
	}

	if h != s.SHA256 {
		return fmt.Errorf("incorrect sha256, want %s, got %s", s.SHA256, h)
	}
	return nil
}
