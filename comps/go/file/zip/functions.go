package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/waykiss/wkcomps/file"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ZipBase64ToBytesArray recebe o arquivo zip no formato base64 e retorna a lista de arquivos encontrados nele em bytes
func ZipBase64ToBytesArray(zipBase64 string) (files [][]byte, err error) {
	//base64 to reader
	zipString := file.Base64ToString(zipBase64)

	//convert to byte array
	data := []byte(zipString)

	//ready zip file from byte array
	zipa, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return
	}

	//lendo os arquivos contidos no zip
	for _, f := range zipa.File {
		rc, fileErr := f.Open()
		if fileErr != nil {
			err = fileErr
			return
		}
		defer func() {
			if err := rc.Close(); err != nil {
				return
			}
		}()
		fileBytes, errIo := ioutil.ReadAll(rc)
		if errIo != nil {
			err = errIo
			return
		}

		files = append(files, fileBytes)
	}
	return
}

func NewZipFile() (z zipFile) {
	z = zipFile{}
	z.File = new(bytes.Buffer)
	z.Zip = zip.NewWriter(z.File)
	return
}

type zipFile struct {
	File *bytes.Buffer
	Zip  *zip.Writer
}

func (z *zipFile) AddFile(filePath string, fileBytes []byte) (err error) {
	f, err := z.Zip.Create(filePath)
	if err != nil {
		return
	}
	_, err = f.Write(fileBytes)
	return
}

// Unzip descompacta um arquivo zip para o destino especificado no parametro dest
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
