package remote

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ddliu/go-httpclient"
)

func http_get(url string, cid int64) bool {
	httpclient.Defaults(httpclient.Map{
		//		httpclient.OPT_USERAGENT: "my awsome httpclient",
		httpclient.OPT_TIMEOUT:        5,
		httpclient.OPT_FOLLOWLOCATION: 1,
	})

	res, err := httpclient.Get(url, map[string]string{})
	if err != nil {
		log.Println("error:", url, err)
		return false
	}
	body, err := res.ToString()
	//	fmt.Println(string(body))
	if len(body) > 200 {
		body = body[0:200]
	}
	log.Println(len(string(body)))
	log.Println(string(body))
	status := res.StatusCode
	log.Println(status)
	if status >= 400 {
		log.Println("notifyAdvertiser fail!!!", status, url)
		return false
	} else {
		log.Println("notifyAdvertiser scusess", status, url)
		return true
	}

}
func curl_get(url string, cid int64) bool {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)

	if err != nil {
		log.Println("error:", err)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//	fmt.Println(string(body))
	if len(body) > 200 {
		body = body[0:200]
	}
	log.Println(len(string(body)))
	log.Println(string(body))
	status := resp.StatusCode
	log.Println(status)
	if status >= 400 {
		return false
	} else {
		return true
	}

}

func HttpPost2(my_url string, postData string) bool {
	//HTTP POST请求
	log.Println("HttpPost2", my_url, postData)
	req, err := http.Post(my_url, "application/x-www-form-urlencoded", strings.NewReader(postData)) //这里定义链接和post的数据
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil || req.StatusCode > 200 {
		log.Println(err.Error())
		return false
	}
	if len(body) > 200 {
		body = body[0:200]
	}
	log.Println(len(string(body)))
	log.Println(string(body))
	return true
}
