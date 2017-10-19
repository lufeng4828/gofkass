package gofkass

import (
	"os"
	"log"
	"strconv"
	"net/http"
	"github.com/imdario/mergo"
	"github.com/jochasinga/requests"
	"fmt"
	"bytes"
	"strings"
)

type FaasQueue struct {
	Name   string
	Kwargs map[string]interface{}
	Text   string
	Type   string
	Pipe   func(*requests.Response, map[string]interface{}) map[string]interface{}
	Func   func(name string, text string, kwargs map[string]interface{}, httpclient ...http.Client) (*requests.Response, error)
}

type Faas struct {
	Queue   []*FaasQueue
	Service string
	Port    int
}

func NewFaas(faasArgs ...interface{}) *Faas {
	faasService := ""
	faasPort := ""
	if len(faasArgs) == 2 {
		faasService, _ = faasArgs[0].(string)
		faasPort, _ = faasArgs[1].(string)
	} else {
		faasService = os.Getenv("FAAS_SERVICE")
		faasPort = os.Getenv("FAAS_PORT")
	}

	if len(faasService) == 0 || len(faasPort) == 0 {
		log.Println("fass、port must specify a value")
		return nil
	}
	faas := new(Faas)
	faas.Queue = make([]*FaasQueue, 0)
	faas.Service = faasService
	port, _ := strconv.Atoi(faasPort)
	faas.Port = port
	return faas
}

func (c *Faas) call(name string, text string, kwargs map[string]interface{}, httpclient ...http.Client) (*requests.Response, error) {
	url := fmt.Sprintf("http://%s:%d/function/%s", c.Service, c.Port, name)
	bodyType := "application/json"
	buf := new(bytes.Buffer)
	if len(text) > 0 {
		buf.Write([]byte(text))
		bodyType = "application/x-www-form-urlencoded"
	} else {
		buf.Write(I2bytes(kwargs))
	}
	var client http.Client
	if len(httpclient) == 0 {
		client = http.Client{
			Transport: &http.Transport{
				DisableKeepAlives: true,
			},
		}
	} else {
		client = httpclient[0]
	}
	request, _ := http.NewRequest("POST", url, buf)
	request.Header.Set("Content-Type", bodyType)
	resp, err := client.Do(request)
	if err != nil {
		return &requests.Response{}, err
	}
	response := requests.Response{Response: resp}
	return &response, err
}

func (c *Faas) Serv(name string, text string, kwargs map[string]interface{}) *Faas {
	name = strings.Replace(name, ".", "/", -1)
	c.Queue = append(c.Queue, &FaasQueue{Name: name, Text: text, Kwargs: kwargs, Type: "func", Func: c.call})
	return c
}

func (c *Faas) Pipe(f func(*requests.Response, map[string]interface{}) map[string]interface{}, kwargs map[string]interface{}) *Faas {
	if len(c.Queue) == 0 || c.Queue[len(c.Queue)-1].Type != "func" {
		log.Println("Pipe must be called after a Serv")
		return nil
	}
	c.Queue = append(c.Queue, &FaasQueue{Pipe: f, Kwargs: kwargs, Type: "pipe"})
	return c
}

func (c *Faas) Get(httpclient ...http.Client) (*requests.Response, error) {
	var result interface{}
	var err error
	for _, item := range c.Queue {
		switch item.Type {
		case "func":
			kwargs := make(map[string]interface{})
			kv, ok := result.(map[string]interface{})
			if ok {
				kwargs = kv
			}
			if item.Kwargs != nil {
				if err := mergo.Merge(&kwargs, item.Kwargs); err != nil {
					println("merge Demand map error:", err.Error())
				}
			}
			result, err = item.Func(item.Name, item.Text, kwargs, httpclient...)
			if err != nil {
				log.Println(err.Error())
			}
		case "pipe":
			kv, ok := result.(*requests.Response)
			if ok {
				result = item.Pipe(kv, item.Kwargs)
			} else {
				log.Println(err.Error())
			}
		}
	}
	kv, _ := result.(*requests.Response)
	return kv, nil
}
