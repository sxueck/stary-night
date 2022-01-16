package spider

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:95.0) Gecko/20100101 Firefox/95.0"

var bp = NewBytePoolCap(50, 50*1024, 50*1024)

type requestOptions struct {
	timeout   time.Duration
	data      string
	headers   map[string]string
	userAgent string
}

type Option struct {
	apply func(option *requestOptions)
}

type LengthReader struct {
	Source io.ReadCloser
	Length int
}

func defaultRequestOptions() *requestOptions {
	return &requestOptions{
		timeout: 5,
		data:    "",
		headers: map[string]string{
			"User-Agent":      UserAgent,
			"Accept-Encoding": "gzip",
			"Accept":          "*/*",
			"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
			"Connection":      "close",
		},
		userAgent: UserAgent,
	}
}

func WithData(data string) *Option {
	return &Option{
		apply: func(option *requestOptions) {
			option.data = data
		},
	}
}

// WithAppendHeaders : you can use any key-value pair
func WithAppendHeaders(kv map[string]string, overwrite bool) *Option {
	return &Option{
		apply: func(option *requestOptions) {
			if overwrite {
				option.headers = kv
			} else {
				for k, v := range kv {
					option.headers[k] = v
				}
			}
		},
	}
}

func WithObjectTimeout(timeout int) *Option {
	return &Option{
		apply: func(option *requestOptions) {
			option.timeout = time.Duration(timeout) * time.Microsecond
		},
	}
}

func exBody(res *http.Response) []byte {
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	return bs
}

func HttpRequestReader(method string, url string, options ...*Option) (*http.Response, error) {
	reqOpts := defaultRequestOptions()
	for _, opt := range options {
		opt.apply(reqOpts)
	}

	req, err := http.NewRequest(method, url, strings.NewReader(reqOpts.data))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize http request: %s", err)
	}

	for key, value := range reqOpts.headers {
		req.Header.Add(key, value)
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ERROR : %s", err)
	}

	return res, nil
}

func HttpRequestToBytes(method string, url string, options ...*Option) (*[]byte, error) {
	res, err := HttpRequestReader(method, url, options...)
	if err != nil {
		return nil, err
	}
	singleBody := exBody(res)

	bodyBuff := bp.Get()
	defer bp.Put(bodyBuff)

	bodyBuff = singleBody

	var unData *[]byte
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		unData, err = GzipUnCompress(bodyBuff)
		if err != nil {
			return nil, err
		}
	default: // identity
		unData = &bodyBuff
	}

	//fmt.Printf("%p\n", unData)

	return unData, nil
}

func GzipUnCompress(data []byte) (*[]byte, error) {
	b := new(bytes.Buffer)
	_ = binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	defer r.Close()
	if err != nil {
		log.Printf("[ParseGzip] NewReader error: %v, maybe data is ungzip\n", err)
		return nil, err
	} else {
		unData, err := ioutil.ReadAll(r)

		// Because a larger buffer area was requested
		// but there was no more data behind it
		// that's why the EOF error
		if err != nil && !strings.Contains(err.Error(), "EOF") {
			log.Printf("[ParseGzip]  ioutil.ReadAll error: %v\n", err)
			return nil, err
		}

		return &unData, nil
	}
}

func SurvivalChecks(url string) bool {
	res, err := HttpRequestReader("GET", url, WithObjectTimeout(400))
	if err != nil {
		log.Printf("%s inaccessible : %s", url, err)
		return false
	}

	if res.StatusCode != http.StatusOK {
		return false
	}

	return true
}
