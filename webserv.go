package gofkass

import (
	"time"
	"github.com/bitly/go-simplejson"
	"github.com/lufeng4828/beego"
	"github.com/lufeng4828/beego/httplib"
	"fmt"
	"net/url"
	"strings"
	"os"
)

func init() {
	httplib.SetDefaultSetting(
		httplib.BeegoHTTPSettings{
			UserAgent:        "beegoServer",
			ConnectTimeout:   600 * time.Second,
			ReadWriteTimeout: 1800 * time.Second,
			Gzip:             true,
			DumpBody:         true,
		})
}


type WebServ struct {
	AppName   string
	ApiUrl    string
	SecretId  string
	Signature string
	Token     string
}

func NewWebServ(appName string, config ...string) *WebServ {
	webServ := new(WebServ)
	webServ.AppName = appName
	if len(config) == 2 {
		webServ.SecretId = config[0]
		webServ.Signature = config[1]
		webServ.ApiUrl = fmt.Sprintf("http://%s-app:8080", appName)
	}else if len(config) == 3{
		webServ.SecretId = config[0]
		webServ.Signature = config[1]
		webServ.ApiUrl = config[2]
	}else{
		webServ.SecretId = os.Getenv("API_SECRECT_ID")
		webServ.Signature = os.Getenv("API_SIGNATURE")
		webServ.ApiUrl = fmt.Sprintf("http://%s-app:8080", appName)
	}

	return webServ
}

func (c *WebServ) Bind(appName string, apiUrl string, secretId string, signature string) {
	c.AppName = appName
	c.ApiUrl = apiUrl
	c.SecretId = secretId
	c.Signature = signature
}

func (c *WebServ) SetToken(token string) {
	c.Token = token
}

func (c *WebServ) ToJson(data string) *simplejson.Json {
	if json, ok := simplejson.NewJson([]byte(data)); ok == nil {
		return json
	}
	return nil
}

func (c *WebServ) Call(resource string, contentType string, body map[string]interface{}, method ... string) (*simplejson.Json, error) {
	resource = strings.Replace(resource, ".", "/", -1)
	var res string
	var err error
	url_ := fmt.Sprintf("%s/resource/v1/%s.json", c.ApiUrl, resource)
	if len(method) == 0 || method[0] == "POST" {
		req := httplib.Post(url_)
		if len(c.SecretId) > 0 && len(c.Signature) > 0 {
			req.Header("x-secretid", c.SecretId)
			req.Header("x-signature", c.Signature)
		} else {
			if len(c.Token) > 0 {
				req.Header("x-token", c.Token)
			}
		}

		if contentType == "application/json" {
			req.Header("Content-Type", "application/json")
			req.JSONBody(body)
		} else {
			for k, v := range body {
				req.Param(k, v.(string))
			}
		}
		res, err = req.String()
	} else {
		query := url.Values{}
		for key, value := range body {
			if str, err := I2String(value); err == nil {
				query.Add(key, str)
			}

		}
		url_ = fmt.Sprintf("%s?%s", url_, query.Encode())
		req := httplib.Get(url_)
		if len(c.SecretId) > 0 && len(c.Signature) > 0 {
			req.Header("x-secretid", c.SecretId)
			req.Header("x-signature", c.Signature)
		} else {
			if len(c.Token) > 0 {
				req.Header("x-token", c.Token)
			}
		}
		res, err = req.String()
	}
	result := c.ToJson(res)
	if result == nil{
		beego.Error(res)
	}
	return result, err
}

func (c *WebServ) TextPost(resource string, body map[string]interface{}) (*simplejson.Json, error) {
	return c.Call(resource, "", body, "POST")
}

func (c *WebServ) JsonPost(resource string, body map[string]interface{}) (*simplejson.Json, error) {
	return c.Call(resource, "application/json", body, "POST")
}

func (c *WebServ) Get(resource string, body map[string]interface{}) (*simplejson.Json, error) {
	return c.Call(resource, "application/json", body, "GET")
}