package builder

import (
	"bufio"
	"bytes"
	"os"
	"path"
	"path/filepath"
)

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadFileLines(f *os.File) ([]string, error) {
	result := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result, nil
}

func BuildBuffer(lines []string) []byte {
	var buffer bytes.Buffer
	for _, line := range lines {
		buffer.Write([]byte(line))
		buffer.Write([]byte("\n"))
	}
	return buffer.Bytes()
}

func CreateAndWrite(pth, filename string, lines []string) error {
	fup, err := os.Create(path.Join(pth, filename))
	if err != nil {
		panic(err)
	}
	defer fup.Close()
	if _, err := fup.Write(BuildBuffer(lines)); err != nil {
		panic(err)
	}
	return nil
}
