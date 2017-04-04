package backend

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type LocalFile struct {
	FileName string
	hashes   map[string]interface{}
}

func NewLocalFile(fileName string) (*LocalFile, error) {
	hashes, err := loadHashes(fileName)
	if err != nil {
		return nil, err
	}
	return &LocalFile{
		FileName: fileName,
		hashes:   hashes,
	}, nil
}

func (f *LocalFile) GetLease() (id string, err error) {
	fn := f.FileName + ".lock"
	if _, err := os.Stat(fn); !os.IsNotExist(err) {
		return "", err
	}
	file, err := os.Create(fn)
	defer file.Close()
	if err != nil {
		return "", err
	}
	return "ok", nil
}

func (f *LocalFile) ExtendLease(id string) (ok bool, err error) {
	_, err = f.GetLease()
	return err == nil, err
}

func (f *LocalFile) CancelLease() (err error) {
	err = saveHashes(f.FileName, f.hashes)
	if err != nil {
		return err
	}
	return os.Remove(f.FileName + ".lock")
}

func (f *LocalFile) IsProcessed(hash string) (bool, error) {
	_, ok := f.hashes[hash]
	return ok, nil
}

func (f *LocalFile) MarkProcessed(hash string) error {
	f.hashes[hash] = true
	return nil
}

func loadHashes(fileName string) (map[string]interface{}, error) {
	op := make(map[string]interface{})

	file, err := os.Open(fileName)
	defer file.Close()
	if os.IsNotExist(err) {
		return op, nil
	}
	if err != nil {
		return op, err
	}

	r := bufio.NewReader(file)

	var line string
	for {
		line, err = r.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line != "" {
			op[line] = true
		}
	}

	if err != io.EOF {
		return op, err
	}

	return op, nil
}

func saveHashes(fileName string, hashes map[string]interface{}) error {
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		return err
	}

	for k := range hashes {
		_, err := file.WriteString(k + "\n")
		if err != nil {
			return fmt.Errorf("failed to write hash to file: %v", err)
		}
	}
	return file.Sync()
}
