package gofkass

import (
	"testing"
	"github.com/jochasinga/requests"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFaas(t *testing.T) {

	Convey("Subject: Test Faas\n", t, func() {
		faas := NewFaas("faas.test.foo.com", "80")
		res, err := faas.Serv("cron.check", "", nil).Get()

		Convey("1.<echo.check> Should Be 'ok'", func() {
			So(res.String(), ShouldEqual, "ok")
		})

		faas = NewFaas("faas.test.foo.com", "80")
		res, err = faas.Serv("echo", "test", nil).Get()

		Convey("2.<echo> Should Be 'test'", func() {
			dict, err := FromJson(string(res.JSON()))
			if err != nil{
				println(err)
			}
			text, _ := dict["text"]
			So(text, ShouldEqual, "test")
		})

		faas = NewFaas("faas.test.foo.com", "80")
		res, err = faas.Serv("echo", "", map[string]interface{}{"name": "foo"}).Get()

		Convey("3.<echo> Should Be 'test'", func() {
			dict, err := FromJson(string(res.JSON()))
			if err != nil{
				println(err)
			}
			text, _ := dict["text"].(string)
			kwargs, _ := dict["kwargs"].(map[string]interface{})
			name, _ := kwargs["name"].(string)
			So(text + name, ShouldEqual, "foo")
		})

		faas = NewFaas("faas.test.foo.com", "80")
		res, err = faas.Serv("echo", "", map[string]interface{}{"name": "foo"}).Pipe(func(r *requests.Response, kwargs map[string]interface{}) map[string]interface{}{
			dict, err := FromJson(string(r.JSON()))
			if err != nil{
				println(err)
			}

			kwargs1, _ := dict["kwargs"].(map[string]interface{})
			kwargs1["pipe"] = "pipe"
			return kwargs1
		}, nil).Serv("echo", "", nil).Get()

		Convey("4.<echo> Should Be 'test'", func() {
			dict, err := FromJson(string(res.JSON()))
			if err != nil{
				println(err)
			}
			text, _ := dict["text"].(string)
			kwargs, _ := dict["kwargs"].(map[string]interface{})
			name, _ := kwargs["name"].(string)
			pipe, _ := kwargs["pipe"].(string)
			So(text + name + pipe, ShouldEqual, "foopipe")
		})

		if err != nil{
			println(err.Error())
		}

	})
}