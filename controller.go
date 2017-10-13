package gofkass

import (
	"os"
	"fmt"
	"time"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type BaseController struct {
	beego.Controller
	AppPath string
	JsonData *simplejson.Json
}

func (c *BaseController) jsonResult(out interface{}) {
	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}

func (c *BaseController) IsPost() bool {
	return c.Ctx.Request.Method == "POST"
}

func (c *BaseController) Stringfy(data string) {
	c.Ctx.ResponseWriter.Write([]byte(data))
}

func (b *BaseController) Jint(key string, def ...int) int {
	if !b.Jexist(b.JsonData, key) && len(def) > 0 {
		return def[0]
	}
	value, err := b.JsonData.Get(key).Int()
	if err != nil {
		if len(def) > 0 {
			return def[0]
		}
	}
	return value
}

func (b *BaseController) Jexist(data *simplejson.Json, key string) bool {
	m, err := data.Map()
	if err == nil {
		if _, ok := m[key]; ok {
			return true
		}
	}
	return false
}

func (b *BaseController) Jfloat(key string, def ...float64) float64 {
	if !b.Jexist(b.JsonData, key) && len(def) > 0 {
		return def[0]
	}
	value, err := b.JsonData.Get(key).Float64()

	if err != nil {
		if len(def) > 0 {
			return def[0]
		}
	}
	return value
}

func (b *BaseController) Jbool(key string, def ...bool) bool {
	if !b.Jexist(b.JsonData, key) && len(def) > 0 {
		return def[0]
	}
	value, err := b.JsonData.Get(key).Bool()
	if err != nil {
		if len(def) > 0 {
			return def[0]
		}
	}
	return value
}

func (b *BaseController) Jstring(key string, def ...string) string {
	value, err := b.JsonData.Get(key).String()
	if err != nil || len(value) == 0 {
		if len(def) > 0 {
			return def[0]
		}
	}
	return value
}

func (b *BaseController) Jmap(key string, def ...map[string]interface{}) map[string]interface{} {
	if m, err := b.JsonData.Get(key).Map(); err == nil {
		return m
	}
	if len(def) > 0{
		return def[0]
	}
	return map[string]interface{}{}
}

func (b *BaseController) JstringArray(key string, def ...[]string) []string {
	value, err := b.JsonData.Get(key).StringArray()
	if err != nil || len(value) == 0 {
		if len(def) > 0 {
			return def[0]
		}
	}
	return value
}

func (b *BaseController) JTime(key string, def ...time.Time) time.Time {
	if !b.Jexist(b.JsonData, key) && len(def) > 0 {
		return def[0]
	}
	value, err := b.JsonData.Get(key).String()
	if err != nil {
		if len(def) > 0 {
			return def[0]
		}
	}
	time_, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
	if err != nil {
		return time.Time{}
	} else {
		return time_
	}
}

func (c *BaseController) Jsonify(data interface{}, code int, message string, success bool) {
	out := make(map[string]interface{})
	out["code"] = code
	out["data"] = data
	out["message"] = message
	out["success"] = success
	if code == 200 || code == 401 || code == 403 {
		c.Ctx.ResponseWriter.Header().Add("Content-Type", "application/json; charset=utf-8")
		c.Ctx.ResponseWriter.WriteHeader(code)
	}
	c.jsonResult(out)
}

func (c *BaseController) I2str(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json err:", err)
		return ""
	}
	return string(b)
}

func (c *BaseController) PrepareJson() bool {
	c.JsonData = c.GetJson()
	return c.JsonData != nil
}

func (c *BaseController) GetJson() *simplejson.Json {
	data, err := simplejson.NewJson(c.Ctx.Input.RequestBody)
	if err != nil {
		println(err.Error())
		return nil
	}
	return data
}

func (c *BaseController) Check() {
	c.Stringfy("ok")
}

func (c *BaseController) Version() {
	version := os.Getenv("FAAS_VERSION")
	c.Stringfy(string(version))
}

func (c *BaseController) Desc() {
	fileName := filepath.Join(beego.AppPath, "desc.yaml")
	file, err := os.Open(fileName)
	result := make(map[string]interface{})
	if err != nil {
		c.jsonResult(result)
	}
	content, _ := ioutil.ReadAll(file)
	yaml.Unmarshal(content, &result)
	c.jsonResult(result)
}