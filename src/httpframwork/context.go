package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Context struct {
	Params map[string]interface{}

	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type ContextFunc func(*Context)

//templates 템플릿 객체를 보관하기 위한 map
var templates = map[string]*template.Template{}

func (c *Context) RenderTemplate(path string, v interface{}) {
	//path에 해당하는 템플릿이 있는지 확인
	t, ok := templates[path]
	if ok == false {
		t = template.Must(template.ParseFiles(filepath.Join(".", path)))
		templates[path] = t
	}

	// v 값을 템플릿 내부로 전달하여 만들어진 최종 결과를 c.ResponseWriter에 출력
	t.Execute(c.ResponseWriter, v)
}

func (c *Context) Redirect(url string) {
	log.Println("Redirect url:", url)
	http.Redirect(c.ResponseWriter, c.Request, url, http.StatusMovedPermanently)
}
