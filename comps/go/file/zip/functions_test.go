package zip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestZipBase64ToBytesArrayFail testa com entradas de arquivos nao sendo zip
func TestZipBase64ToBytesArrayFail(t *testing.T) {
	arquivoBase64Txt := "dGVzdGU="
	arquivoBase64TxtPrefixoErrado := "data:image/png;base64,dGVzdGU="
	arquivoBase64TxtPrefixoErrado = "data:application/zip;base64,dGVzdGU="
	values := []struct {
		valueBase64 string
		expected    string
	}{
		{arquivoBase64Txt, "teste"},
		{arquivoBase64TxtPrefixoErrado, "teste"},
		{arquivoBase64TxtPrefixoErrado, "teste"},
	}
	for _, v := range values {
		b, err := ZipBase64ToBytesArray(v.valueBase64)
		assert.NotNil(t, err.Error(), "zip: not a valid zip file", "Deveria da erro, arquivo zip incorreto!")
		assert.NotNil(t, v.expected, b, "deveria ter retornado  vazio")
	}
}

// TestZipBase64ToBytesArraySuccess testa arquivo correto com diferentes prefixo
func TestZipBase64ToBytesArraySuccess(t *testing.T) {
	arquivoBase64Zip := "UEsDBBQAAAAIAMRKTVMPSbTmBwAAAAUAAAAJAAAAdGVzdGUudHh0K0ktLkkFAFBLAQIfABQAAAAIAMRKTVMPSbTmBwAAAAUAAAAJACQAAAAAAAAAIAAAAAAAAAB0ZXN0ZS50eHQKACAAAAAAAAEAGAC+AazxLMDXAb4BrPEswNcBvgGs8SzA1wFQSwUGAAAAAAEAAQBbAAAALgAAAAAA"
	arquivoBase64ZipComPrefixoIncorreto := "data:image/png;base64,UEsDBBQAAAAIAMRKTVMPSbTmBwAAAAUAAAAJAAAAdGVzdGUudHh0K0ktLkkFAFBLAQIfABQAAAAIAMRKTVMPSbTmBwAAAAUAAAAJACQAAAAAAAAAIAAAAAAAAAB0ZXN0ZS50eHQKACAAAAAAAAEAGAC+AazxLMDXAb4BrPEswNcBvgGs8SzA1wFQSwUGAAAAAAEAAQBbAAAALgAAAAAA"
	arquivoBase64ZipComPrefixoIncorreto = "data:application/zip;base64,UEsDBBQAAAAIAMRKTVMPSbTmBwAAAAUAAAAJAAAAdGVzdGUudHh0K0ktLkkFAFBLAQIfABQAAAAIAMRKTVMPSbTmBwAAAAUAAAAJACQAAAAAAAAAIAAAAAAAAAB0ZXN0ZS50eHQKACAAAAAAAAEAGAC+AazxLMDXAb4BrPEswNcBvgGs8SzA1wFQSwUGAAAAAAEAAQBbAAAALgAAAAAA"
	values := []struct {
		valueBase64 string
		expected    string
	}{
		{arquivoBase64Zip, "teste"},
		{arquivoBase64ZipComPrefixoIncorreto, "teste"},
		{arquivoBase64ZipComPrefixoIncorreto, "teste"},
	}
	for _, v := range values {
		b, err := ZipBase64ToBytesArray(v.valueBase64)
		assert.Nil(t, err, "Nao deveria da erro, arquivo zip correto!")
		assert.Equal(t, v.expected, string(b[0]), "deveria ter retornado o conteudo do arquivo")
	}
}
