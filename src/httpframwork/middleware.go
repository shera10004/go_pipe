package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

type Middleware func(next ContextFunc) ContextFunc

func logHandler(next ContextFunc) ContextFunc {
	fmt.Println("<logHandler")
	return func(c *Context) {
		//next(c)를 실행하기 전에 현재 시간을 기록
		t := time.Now()

		// 다음 핸들러 수행
		next(c)

		// 웹 요청 정보와 전체 소요 시간을 로그로 남김
		log.Printf("logHandler - [%s] %q %v\n",
			c.Request.Method,
			c.Request.URL.String(),
			time.Now().Sub(t))
	}
}

func recoverHandler(next ContextFunc) ContextFunc {
	fmt.Println("<recoverHandler")
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(c.ResponseWriter,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()
		next(c)
	}
}

func parseFormHandler(next ContextFunc) ContextFunc {
	fmt.Println("<parseFormHandler")
	return func(c *Context) {
		c.Request.ParseForm()
		fmt.Println(c.Request.PostForm)
		for k, v := range c.Request.PostForm {
			if len(v) > 0 {
				c.Params[k] = v[0]
				log.Println("parseForm Param[", k, "]:", v[0])
			}
		} //for
		next(c)
	}
}

func parseJSONBodyHandler(next ContextFunc) ContextFunc {
	fmt.Println("<parseJSONBodyHandler")
	return func(c *Context) {
		var m map[string]interface{}
		if json.NewDecoder(c.Request.Body).Decode(&m); len(m) > 0 {
			for k, v := range m {
				c.Params[k] = v
				log.Println("parseJSON Param[", k, "]:", v)
			} //for
		}
		next(c)
	}
}

func staticHandler(next ContextFunc) ContextFunc {
	fmt.Println("<staticHandler")
	var (
		dir       = http.Dir("")
		indexFile = "index.html"
	)
	return func(c *Context) {
		fmt.Println("staticHandler <--- ", c)
		//http 메서드가 GET이나 HEAD가 아니면 바로 다음 핸들러 수행
		if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			next(c)
			return
		}

		file := c.Request.URL.Path
		f, err := dir.Open(file)
		if err != nil {
			// URL 경로에 해당하는 파일 열기에 실패하면 다음 핸들러 수행
			next(c)
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			// 파일의 상태가 정상이 아니면 바로 다음 핸들러 수행
			next(c)
			return
		}

		// URL 경로가 디렉터리면 indexFile을 사용
		if fi.IsDir() {
			//디렉터리 경로를 URL로 사용하면 경로 끝에 "/"를 붙여야함
			if !strings.HasSuffix(c.Request.URL.Path, "/") {
				http.Redirect(c.ResponseWriter, c.Request, c.Request.URL.Path+"/", http.StatusFound)
				return
			}

			//디렉터리를 가리키는 URL 경로에 indexFile 이름을 붙여서 전체 파일 경로 생성
			file = path.Join(file, indexFile)

			//indexFile 열기 시도
			f, err = dir.Open(file)
			if err != nil {
				next(c)
				return
			}
			defer f.Close()

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				//indexFile 상태가 정상이 아니면 바로 다음 핸들러 수행
				next(c)
				return
			}
		}
		//file의 내용 전달(next 핸들러로 제어권을 넘기지 않고 요청 처리를 종료함)
		http.ServeContent(c.ResponseWriter, c.Request, file, fi.ModTime(), f)
	}
}

func (c *Context) RenderJson(v interface{}) {
	//HTTP Status를 StatusOK로 지정
	c.ResponseWriter.WriteHeader(http.StatusOK)
	//Content-Type을 application/json으로 지정
	c.ResponseWriter.Header().Set("Content-Type", "application/json,charset=utf-8")

	//v 값을 json으로 출력
	if err := json.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
		//에러 발생시 RenderErr메서드 호출
		c.RenderErr(http.StatusInternalServerError, err)
	}
}

func (c *Context) RenderXml(v interface{}) {

	c.ResponseWriter.WriteHeader(http.StatusOK)
	//Content-Type을 application/xml으로 지정
	c.ResponseWriter.Header().Set("Content-Type", "application/xml,charset=utf-8")

	if err := xml.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
		c.RenderErr(http.StatusInternalServerError, err)
	}
}

func (c *Context) RenderErr(code int, err error) {
	if err != nil {
		if code > 0 {
			//정상적인 code를 전달하면 HTTP Status를 해당 code로 지정
			http.Error(c.ResponseWriter, http.StatusText(code), code)
		} else {
			//정상적인 code가 아니면 HTTP Status를 StatusInternalServerError로 지정
			defaultErr := http.StatusInternalServerError
			http.Error(c.ResponseWriter, http.StatusText(defaultErr), defaultErr)
		}
	}
}
