package tool

import (
	"io/ioutil"
	"os"
	"net/http"
	"fmt"
	"io"
)

func WriteFile(path string, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0777)
}

func ReadFile(path string) ([]byte, error) {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	return ioutil.ReadAll(fi)
}

func ReadFileString(path string) (string, error) {
	bytes, err := ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func DownloadFileToBytes(url string) ([]byte, error) {
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return bodyBytes, nil
	}

	return nil, fmt.Errorf("download file! http status not ok")
}

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
func DownloadFileToPath(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}