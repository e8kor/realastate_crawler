package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	handler "github.com/openfaas/templates-sdk/go-http"
)

// Entry is domain associated crawled json
type Entry struct {
	Domain string            `json:"domain"`
	Data   []json.RawMessage `json:"data"`
}

func Handle(r handler.Request) (handler.Response, error) {
	query, err := url.ParseQuery(r.QueryString)
	if err != nil {
		panic(err)
	}

	var (
		response handler.Response
		urls     = query["url"]
	)

	if urls == nil {
		urls = append(urls, os.Getenv("SOURCE_URL"))
	}

	crawlerResponse, err := http.Get("http://localhost/otodom?url=" + strings.Join(urls[:], "&url="))
	if err != nil {
		panic(err)
	}

	crawlerRequest := fmt.Sprintf(`{
		"domain": "otodom"
		"data": %s
	}`, string(streamToByte(crawlerResponse.Body)))

	persistorResponse, err := http.Post("http://localhost/persistor", "application/json", bytes.NewBuffer([]byte(crawlerRequest)))
	if err != nil {
		panic(err)
	}
	response = handler.Response{
		Body:       streamToByte(persistorResponse.Body),
		StatusCode: persistorResponse.StatusCode,
		Header:     persistorResponse.Header,
	}
	return response, nil
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
