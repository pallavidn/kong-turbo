package discovery


import (
	"crypto/tls"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"fmt"
)

const (
	defaultTimeOut = time.Second * 5
)

type HttpClient struct {
	client   *http.Client
	username string
	password string

	addr string
}

func NewHttpClient(addr string) (*HttpClient, error) {
	//1. get http client
	client := &http.Client{
		Timeout: defaultTimeOut,
	}

	//2. check whether it is using ssl
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}

	uaddr, err := url.Parse(addr)
	if err != nil {
		glog.Errorf("Invalid url:%v, %v", addr, err)
		return nil, err
	}
	if uaddr.Scheme == "https" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	glog.V(2).Infof("Created new http client: %v", addr)
	fmt.Printf("Created new http client: %v", addr)

	//3. generate request
	c := &HttpClient{
		client: client,
		addr:   addr,
	}

	return c, nil
}

// SetUser set the login user/password for the prometheus client
func (c *HttpClient) SetUser(username, password string) {
	c.username = username
	c.password = password
}

func (c *HttpClient) generateRequest() (*http.Request, error) {
	req, err := http.NewRequest("POST", c.addr, nil)
	if err != nil {
		glog.Errorf("Failed to generate a http.request: %v", err)
		return nil, err
	}

	if len(c.username) > 0 {
		req.SetBasicAuth(c.username, c.password)
	}

	return req, nil
}

func (c *HttpClient) DoPost() (string, error) {
	//1. generate request
	req, err := c.generateRequest()
	if err != nil {
		glog.Errorf("Failed to generate http request: %v", err)
		return "", err
	}
	glog.V(3).Infof("Request URL: %v", req.URL)
	//2. send the request
	resp, err := c.client.Do(req)
	if err != nil {
		glog.Errorf("Failed to send http request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Failed to read response: %v", err)
		return "", err
	} else {
		glog.V(3).Infof("Response: %v", string(result))
	}

	return string(result), nil
}

func (c *HttpClient) DoGet(resourcePath string) ([]byte, error){
	//1. generate request
	requestURL := fmt.Sprintf("%s/%s", c.addr, resourcePath)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		glog.Errorf("Failed to generate a http.request: %v", err)
		return nil, err
	}

	glog.V(3).Infof("Request: %v", req.URL)
	//fmt.Printf("Request: %v\n", req.URL)
	//2. send the request
	resp, err := c.client.Do(req)
	if err != nil {
		glog.Errorf("Failed to send http request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var byteContent []byte
	byteContent, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(" Error in ioutil.ReadAll: %s", err)
	} else {
		//glog.V(3).Infof("Response: %v", string(byteContent))
		//fmt.Printf("Response: %v\n", string(byteContent))
	}

	return byteContent, nil
}
