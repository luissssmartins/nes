package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"strings"
)

const maxROMSize = 8 << 20 

var romDir = func() string {
	u, err := user.Current()
	if err != nil {
		return ".nes/roms"
	}
	return path.Join(u.HomeDir, ".nes", "roms")
}()

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func downloadROM(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	name := path.Base(u.Path)
	if name == "." || name == "/" {
		name = u.Host
	}
	if !strings.HasSuffix(name, ".nes") {
		name += ".nes"
	}
	dst := path.Join(romDir, name)
	if _, err := os.Stat(dst); err == nil {
		return dst, nil
	}
	resp, err := http.Get(rawurl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %s", rawurl, resp.Status)
	}
	data, err := ioutil.ReadAll(io.LimitReader(resp.Body, maxROMSize+1))
	if err != nil {
		return "", err
	}
	if len(data) > maxROMSize {
		return "", fmt.Errorf("%s: file too large", rawurl)
	}
	if !bytes.HasPrefix(data, []byte("NES\x1a")) {
		return "", fmt.Errorf("%s: not an iNES (.nes) file", rawurl)
	}
	if err := os.MkdirAll(romDir, 0755); err != nil {
		return "", err
	}
	tmp := dst + ".part"
	if err := ioutil.WriteFile(tmp, data, 0644); err != nil {
		return "", err
	}
	return dst, os.Rename(tmp, dst)
}
