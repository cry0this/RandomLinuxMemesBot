package discord

import (
	"errors"
	"io"
	"net/http"
	"os"
)

func DownloadFile(fileName string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	if written == 0 {
		return errors.New("empty file")
	}

	return nil
}
