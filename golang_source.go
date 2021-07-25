// Copyright (C) 2021  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
