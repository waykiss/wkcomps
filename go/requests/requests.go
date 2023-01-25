package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Response struct {
	StatusCode int
	Result     []byte
}

type Model struct {
	BasicAuthentication struct {
		UserName, Password string
	}
	Timeout  time.Duration
	Response Response
}

func NewRequest() Model {
	return Model{}
}

// Get realiza uma request GET e joga o resultado da request para a variavel dest
func (r *Model) Get(url string, headers map[string]string, dest interface{}) (resp *http.Response, err error) {
	return r.Request(http.MethodGet, url, headers, nil, dest)
}

// Post realiza uma request POST e joga o resultado da request para a variavel dest
func (r *Model) Post(url string, headers map[string]string, params,
	dest interface{}) (resp *http.Response, err error) {
	return r.Request(http.MethodPost, url, headers, params, dest)
}

func (r *Model) Delete(url string, headers map[string]string, params,
	dest interface{}) (resp *http.Response, err error) {
	return r.Request(http.MethodDelete, url, headers, params, dest)
}

// Put realiza uma request PUT e joga o resultado da request para a variavel dest
func (r *Model) Put(url string, headers map[string]string, params,
	dest interface{}) (resp *http.Response, err error) {
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
func (r *Model) Request(method, urlAddress string, headers map[string]string, params,
	dest interface{}) (resp *http.Response, err error) {
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
		req.Header.Set(key, value)
	}

	// realiza a request
	client := &http.Client{}
	client.Timeout = r.Timeout
	resp, err = client.Do(req)

	if err != nil {
		logrus.WithFields(map[string]interface{}{
			"Error": err,
		}).Error("Erro ao realizar request")
		return
	}

	if resp.StatusCode >= 400 {
		body, err2 := readBody(resp.Body)
		if err2 != nil {
			return
		}
		defer resp.Body.Close()
		r.Response.StatusCode = resp.StatusCode
		r.Response.Result = body
		err = errors.New(fmt.Sprintf("Status Code %d  message: %s", resp.StatusCode, string(body)))
		return
	}

	// grava a resposta na variavel dest
	if dest != nil {
		body, err2 := readBody(resp.Body)
		defer resp.Body.Close()
		if err2 != nil {
			err = err2
			return
		}
		r.Response.StatusCode = resp.StatusCode
		r.Response.Result = body
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
