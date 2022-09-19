package file

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gocarina/gocsv"
	"github.com/vincent-petithory/dataurl"
	"github.com/waykiss/wkcomps/str"
	"golang.org/x/text/encoding/charmap"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func CreateFile(file []byte, filePath string) (err error) {
	dir := filepath.Dir(filePath)
	if err = CreateDirIfNotExists(dir); err != nil {
		return
	}
	osFile, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer osFile.Close()
	_, err = osFile.Write(file)
	if err != nil {
		return
	}
	return
}

// CreateDirIfNotExists cria uma diretorio no caminho especificado caso o mesmo nao exista
func CreateDirIfNotExists(dir string) (err error) {
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModeDir|0755)
	}
	return
}

// Exists ehck if the file exists given the path
func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// DirectotyExists check if directoty exists given a path
func DirectotyExists(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}
	return false
}

// RemoveLastDirectoryFromPath remove o ultimo diretorio de um caminho especificado no parametro
func RemoveLastDirectoryFromPath(p string) string {
	r := strings.Split(filepath.SplitList(filepath.Clean(p))[0], string(filepath.Separator))
	s := strings.Join(r[:len(r)-1], string(filepath.Separator))
	return s
}

// CreateTextFile create a text fileUtil, not binary fileUtil
func CreateTextFile(filePath, content string) error {

	fileDir := filepath.Dir(filePath)
	if !Exists(fileDir) {
		_ = os.MkdirAll(fileDir, os.ModePerm)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// FileUrlToBase64 faz o download de um arquivo dado uma url e retorna o mesmo no formato base64
func FileUrlToBase64(url string) (result string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	fileData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = base64.StdEncoding.EncodeToString(fileData)
	return
}

// FileUrlToReader faz o download de um arquivo dado uma url e retorna um objeto do tipo `bytes.Reader`
func FileUrlToReader(url string) (result *bytes.Reader, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	fileData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = bytes.NewReader(fileData)
	return
}

// FileUrlToLocalfile faz o download de um arquivo dado uma url e salva no caminho passado como parametro
func FileUrlToLocalfile(url, localpath string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	fileData, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("falha ao realizar download. Arquivo %s Resposta : %s", url, fileData)
	}
	if err != nil {
		return
	}
	err = CreateFile(fileData, localpath)
	return
}

// FileToBase64 function to convert a local fileUtil into base64
func FileToBase64(filepath string) (result string, err error) {

	// Open fileUtil on disk.
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()
	// Read the fileUtil into byte slice.
	reader := bufio.NewReader(file)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	result = base64.StdEncoding.EncodeToString(content)
	return
}

// StringToBase64 converte string em base64
func StringToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// ByteToBase64 converte arquivo em byte para base64
func ByteToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Base64ToString converte string no formato base64 para string normal
func Base64ToString(s string) string {
	return string(Base64ToBytes(s))
}

// Base64ToBytes converte uma string em base64 para bytes
func Base64ToBytes(s string) []byte {
	s = GetBase64Content(s)
	b, _ := base64.StdEncoding.DecodeString(s)
	return b
}

// GetBase64Content essa funcao remove os dados sobre o base64 da string, quando converte para base 64, ele pode retornar
// com fragmento dizendo o tipo de arquivo e isso nao é válido para criar um arquivo
func GetBase64Content(b64 string) string {
	// se vier com a string de base64, retire e pegue somente o conteudo de base64
	idx := strings.LastIndex(b64, "base64,")
	if idx > 0 {
		b64 = b64[idx+7:]
	}
	return b64
}

// DecodeSO8859_1 funcao para decodificar bytes que possuem acentos, retorna com acento corretamente
func DecodeISO88591(b string) string {
	dec := charmap.ISO8859_1.NewDecoder()
	arBdest := make([]byte, len(b)*2)
	n, _, err := dec.Transform(arBdest, []byte(b), true)
	if err != nil {
		return ""
	}
	return string(arBdest[:n])
}

// CreateTempFile creates a temporary file and returns its path
// Use fileNameExt to determine the filename along with the extension. EX: file.pdf
func CreateTempFile(file []byte, fileNameExt string) (path string, err error) {
	tmpDir := GetUniqueTempDir()
	if err = CreateDirIfNotExists(tmpDir); err != nil {
		return
	}

	path = filepath.FromSlash(fmt.Sprintf("%s/%s", tmpDir, fileNameExt))
	err = CreateFile(file, path)
	return
}

// GetUniqueTempDir retorna um diretorio temporário(com o path compatível com o OS) adicionado de uma string que o deixa único
func GetUniqueTempDir() string {
	return filepath.FromSlash(
		fmt.Sprintf("%s/%s", os.TempDir(), str.RandString(20, string(str.RandStringCharsAllLettersAndNumbers))),
	)
}

// GetDir retorna o diretorio dado um caminho, ja retorna com o slash compativel com o sistema operacional
func GetDir(path string) (dir string) {
	dir = filepath.Dir(path)
	dir = filepath.FromSlash(dir)
	return
}

// IsBase64 verifica se uma string é base64 valida
func IsBase64(b64 string) bool {
	b64 = GetBase64Content(b64)
	_, err := base64.StdEncoding.DecodeString(b64)
	return err == nil
}

// FileBase64ToFile funcao para salvar uma string de um arquivo em base64 em um arquivo o disco e no caminho especificado
// no parametro `filename`
func FileBase64ToFile(b64 string, filename string) (err error) {
	b64 = GetBase64Content(b64)
	dec, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return
	}
	dir := filepath.Dir(filename)
	err = CreateDirIfNotExists(dir)
	if err != nil {
		return
	}
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	if _, err = file.Write(dec); err != nil {
		return
	}
	if err = file.Sync(); err != nil {
		return
	}
	return
}

// GetTempPath retorna um diretorio temporário junto com o nome do arquivo passado como parametro
func GetTempPath(filename string) string {
	return filepath.FromSlash(fmt.Sprintf("%s/%s/%s",
		os.TempDir(),
		str.RandString(20, string(str.RandStringCharsAllLettersAndNumbers)),
		filename))
}

// FileBase64ToReader convert string in base64 to Reader interface
func FileBase64ToReader(contentBase64 string) *bytes.Reader {
	// check if has comma on string
	commaIdx := strings.IndexByte(contentBase64, ',')
	if commaIdx > 0 {
		contentBase64 = contentBase64[strings.IndexByte(contentBase64, ',')+1:]
	}
	dec, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(dec)
}

// FileToString convert text fileUtil to string
func FileToString(filename string) string {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	return string(file) // convert content to a 'string'
}

// CsvStringToMap Convert Csv stringUtil to Map
func CsvStringToMap(fileString string, delimiter rune, upperColumns bool) (r []map[string]string, err error) {
	lines := strings.Split(fileString, "\n")
	if upperColumns {
		lines[0] = strings.ToUpper(lines[0])
	}
	colunas := strings.Split(lines[0], string(delimiter))
	for i := 1; i < len(lines); i++ {
		campos := strings.Split(lines[i], string(delimiter))
		record := map[string]string{}
		for idx, c := range campos {
			record[colunas[idx]] = c
		}
		r = append(r, record)
	}
	return
}

// CsvStringToStruct converte um arquivo csv no formato de string para uma struct, essa struct passada deve ser uma lista
func CsvStringToStruct(fileString string, delimiter rune, quoted bool, dest interface{}) (err error) {
	r, err := getCsvFromString(fileString, delimiter, quoted)
	if err != nil {
		return
	}
	err = gocsv.Unmarshal(r, dest)
	return
}

// getCsvFromString retorna um objeto(io.Reader) para que seja feito o parser usando o pacote govsc
func getCsvFromString(fileString string, delimiter rune, quoted bool) (r io.Reader, err error) {
	tempFilepath := GetTempPath("csvfile.csv")
	if IsBase64(fileString) {
		if err = FileBase64ToFile(fileString, tempFilepath); err != nil {
			return
		}
	} else {
		if err = CreateFile([]byte(fileString), tempFilepath); err != nil {
			return
		}
	}
	defer func() { _ = os.Remove(tempFilepath) }()

	file, err := os.OpenFile(tempFilepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	// isso é necessário devido a acentuacao nas strings
	r = charmap.ISO8859_1.NewDecoder().Reader(file)
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		if quoted {
			r.LazyQuotes = true
			r.TrimLeadingSpace = true
		}
		r.Comma = delimiter
		r.FieldsPerRecord = -1

		return r
	})
	return
}

// GetFileListFromDirectory returns slice/list of the files within the path passed by parameter, specify if it will
// returns recursively or not
func GetFileListFromDirectory(directoryPath string, recursively bool) ([]string, error) {
	var filesList []string
	err := filepath.Walk(directoryPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !IsDirectory(path) {
				filesList = append(filesList, path)
			}
			return nil
		})
	return filesList, err
}

// getMineTypeFromReader get the file Mine type from string base64 format
func GetMineTypeFromBase64(strBase64 string) string {
	reader := bytes.NewReader([]byte(strBase64))
	return GetMineTypeFromReader(reader)
}

// getMineTypeFromReader get the Mine type from reader
func GetMineTypeFromReader(reader *bytes.Reader) string {
	mime, err := mimetype.DetectReader(reader)
	if err != nil {
		return "text/plain"
	}
	return mime.String()
}

// getMineTypeFromReader get the Mine type from reader
func GetFileExtensionFromBase64(strBase64 string) (string, error) {
	dataURL, err := dataurl.DecodeString(strBase64)
	if err != nil {
		return "", err
	}
	return dataURL.Subtype, nil
}

// GetFileNameWithExtensionFromPath retorna o nome do arquivo com a extensao dado um determinado path
func GetFileNameWithExtensionFromPath(path string) string {
	return filepath.Base(path)
}

// GetFileExtensionFromPath retorna a extensao do aruivo dado um determinado path
func GetFileExtensionFromPath(path string) string {
	return filepath.Ext(path)
}

// IsDirectory check if the path is a directory or not
func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// GetFileSizeFromBase64 get the file size of the base64 string, return the int value what represents the size in KB
func GetFileSizeFromBase64(fileBase64 string) int {

	l := len(fileBase64)

	// count how many trailing '=' there are (if any)
	eq := 0
	if l >= 2 {
		if fileBase64[l-1] == '=' {
			eq++
		}
		if fileBase64[l-2] == '=' {
			eq++
		}

		l -= eq
	}

	// basically:
	//
	// eq == 0 :    bits-wasted = 0
	// eq == 1 :    bits-wasted = 2
	// eq == 2 :    bits-wasted = 4

	// each base64 character = 6 bits

	// so orig length ==  (l*6 - eq*2) / 8

	return (l*3 - eq) / 4 / 1024
}

// removeXmlVersion remove a string da versao do xml, pois quando gera esta vindo essa string
func RemoveXmlVersion(xmlString string) string {
	return strings.ReplaceAll(xmlString, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>", "")
}
