package main

import (
	"fmt"
	"net/http"
	"strings"
)

type router struct {
	contexts map[string]map[string]ContextFunc
}

func (r *router) HandleFunc(method, pattern string, cf ContextFunc) {
	fmt.Println("[router]::HandleFunc method:", method, ", pattern:", pattern, ", ContextFunc:", cf)
	m, ok := r.contexts[method]

	if ok == false {
		m = make(map[string]ContextFunc)
		r.contexts[method] = m
	}

	//http 메서드로 등록된 맵에 url 패턴과 핸들러 함수 등록
	m[pattern] = cf
}

func (r *router) getContextFunc() ContextFunc {
	fmt.Println("< router::getContextFunc")
	return func(c *Context) {

		fmt.Println("getContextFunc <- ", c)

		//http 메서드에 맞는 모든 handlers를 반복하면 요청 URL에 해당하는 handler를 찾음.
		for pattern, handler := range r.contexts[c.Request.Method] {
			if ok, params := match(pattern, c.Request.URL.Path); ok == true {
				for k, v := range params {
					c.Params[k] = v
					fmt.Println("match-params[", k, "]", v)
				}

				fmt.Println(">>>>match-params", len(params))
				//요청 URL에 해당하는 handler수행
				handler(c)
				return
			}
		} //for

		//요청 URL에 해당하는 handler를 찾지 못하면 NotFound 에러 처리
		http.NotFound(c.ResponseWriter, c.Request)
		return
	}

}

func match(pattern, path string) (bool, map[string]string) {

	fmt.Println("match pattern:", pattern, ", path:", path)

	//패턴과 패스가 정확히 일치하면 바로 true를 반환
	if pattern == path {
		return true, nil
	}

	_patterns := strings.Split(pattern, "/")
	_paths := strings.Split(path, "/")

	//패턴과 패스를 "/"로 구분한 후 부분 문자열 집합의 개수가 다르면 false를 반환
	if len(_patterns) != len(_paths) {
		return false, nil
	}

	//패턴에 일치하는 URL 매개변수를 담기 위한 params맵 생성
	_params := make(map[string]string)

	// "/"로 구분된 패턴/패스의 각 문자열을 하나씩 비교
	for i := 0; i < len(_patterns); i++ {
		switch {
		case _patterns[i] == _paths[i]:
			//패턴과 패스의 부분 문자열이 일치하면 바로 다음 루프 수행
		case len(_patterns[i]) > 0 && _patterns[i][0] == ':':
			//패턴이 ':'문자로 시작하면 params에 URL parmas를 담은 후 다음 루프 수행
			_params[_patterns[i][1:]] = _paths[i]
		default:
			//일치하는 경우가 없으면 false를 반환
			return false, nil
		} //switch
	} //for

	return true, _params
}
