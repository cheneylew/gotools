package tool

import (
	"io/ioutil"
	"os"
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