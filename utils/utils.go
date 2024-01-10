package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	other = Others{
		mt:       new(sync.RWMutex),
		Services: make([]string, 0),
	}
)

type Others struct {
	Services []string
	mt       *sync.RWMutex
}

func OtherService() []string {
	other := os.Getenv("OTHERS")
	otherservice := strings.Split(other, ",")

	return otherservice
}

func init() {
	other.mt.Lock()
	defer other.mt.Unlock()
	other.Services = OtherService()
}

type RequestUrl struct {
	Method   string
	Url      string
	Body     string
	Response Response
}

type Response struct {
	Status int
}

func JsonConvert(data map[string]string) string {
	strByte, err := json.Marshal(data)
	if err != nil {

		return ""
	}

	return string(strByte)
}

func Cache(data map[string]string, statusCode int, path string) [][]byte {
	bodys := make([][]byte, 0)
	other.mt.RLock()
	services := other.Services
	other.mt.RUnlock()
	localhost := os.Getenv("HOSTPORT")

	for _, v := range services {
		if localhost == v {
			continue
		}
		request := RequestUrl{
			Method: "GET",
			Url:    fmt.Sprintf("http://%s/%s", v, path),
			Response: Response{
				Status: statusCode,
			},
		}

		if len(data) != 0 {
			str := JsonConvert(data)
			if str == "" {
				return nil
			}

			request.Body = str
		}

		bodys = append(bodys, request.UrlRequest())
	}

	return bodys
}

func (r RequestUrl) UrlRequest() []byte {
	payload := strings.NewReader(r.Body)

	client := &http.Client{}

	req, err := http.NewRequest(r.Method, r.Url, payload)
	if err != nil {
		log.Println("cant make request err:", err.Error())

		return nil
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Println("failed to call API err:", err.Error())

		return nil
	}

	if res.StatusCode != r.Response.Status {
		log.Println("status code not matching ", res.StatusCode)

		return nil
	}

	body, _ := io.ReadAll(res.Body)

	fmt.Println("request as required code: ", res.StatusCode)

	return body
}

func ToMap(data []byte) map[string]string {
	smap := make(map[string]string)
	_ = json.Unmarshal(data, &smap)

	return smap
}
