package requests

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Response struct {
	StatusCode   int
	HttpResponse *http.Response
	ResponseBody []byte
}

type Request struct {
	BasicAuthentication struct {
		UserName, Password string
	}
	Timeout time.Duration
}

func NewRequest() Request {
	return Request{}
}

// Get realiza uma request GET e joga o resultado da request para a variavel dest
func (r *Request) Get(url string, headers map[string]string, dest interface{}) (resp *Response, err error) {
	return r.Request(http.MethodGet, url, headers, nil, dest)
}

// Post realiza uma request POST e joga o resultado da request para a variavel dest
func (r *Request) Post(url string, headers map[string]string, params,
	dest interface{}) (resp *Response, err error) {
	return r.Request(http.MethodPost, url, headers, params, dest)
}

func (r *Request) Delete(url string, headers map[string]string, params,
	dest interface{}) (resp *Response, err error) {
	return r.Request(http.MethodDelete, url, headers, params, dest)
}

// Put realiza uma request PUT e joga o resultado da request para a variavel dest
func (r *Request) Put(url string, headers map[string]string, params,
	dest interface{}) (resp *Response, err error) {
	return r.Request(http.MethodPut, url, headers, params, dest)
}

func readBody(body io.ReadCloser) (resp []byte, err error) {
	resp, err = ioutil.ReadAll(body)
	if err != nil {
		logrus.WithFields(map[string]interface{}{
			"Error": err,
		}).Error("Erro ao ler resposta(body)")
	}
	return
}

// Request realiza uma request e joga o resultado da request para o parametro `dest`
func (r *Request) Request(method, urlAddress string, headers map[string]string, params,
	dest interface{}) (resp *Response, err error) {
	logrus.WithFields(map[string]interface{}{
		"Method": method,
		"Url":    urlAddress,
	}).Info("Realizando request")

	var payload io.Reader
	postFormValues, isPostForm := params.(url.Values)

	switch isPostForm {
	case true:
		payload = strings.NewReader(postFormValues.Encode())
	default:
		if !isPostForm && params != nil {
			paramsBytes, err2 := json.Marshal(params)
			if err2 != nil {
				err = err2
				return
			}
			payload = bytes.NewBuffer(paramsBytes)
		}
	}

	// cria o objeto request
	req, err := http.NewRequest(method, urlAddress, payload)
	if err != nil {
		logrus.WithFields(map[string]interface{}{
			"Error": err,
		}).Error("Erro ao criar objeto request")
		return
	}

	if r.BasicAuthentication.UserName != "" {
		req.SetBasicAuth(r.BasicAuthentication.UserName, r.BasicAuthentication.Password)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	resp = &Response{}
	// realiza a request
	client := &http.Client{}
	client.Timeout = r.Timeout
	resp.HttpResponse, err = client.Do(req)

	if err != nil {
		logrus.WithFields(map[string]interface{}{
			"Error": err,
		}).Error("Erro ao realizar request")
		return
	}

	// faz o parse da request
	body, err2 := readBody(resp.HttpResponse.Body)
	if err2 != nil {
		return
	}
	defer resp.HttpResponse.Body.Close()
	resp.ResponseBody = body

	// grava a resposta na variavel dest
	if dest != nil && resp.HttpResponse.StatusCode == http.StatusOK || resp.HttpResponse.StatusCode == http.StatusCreated {
		err = json.Unmarshal(body, dest)
		if err != nil {
			logrus.WithFields(map[string]interface{}{
				"Error": err,
			}).Error("Erro ao fazer unmarshal de json da request")
			return
		}
	}

	return
}
