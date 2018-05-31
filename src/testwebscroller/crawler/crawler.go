package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"

	"golang.org/x/net/html"
)

var fetched = struct {
	m          map[string]error //중복 검사를 위한 url과 error값 저장
	sync.Mutex                  //뮤텍스 임베딩
}{m: make(map[string]error)} //변수를 선언하면서 이름이 없는 구조체를 정의하고 초기값을 생성하여 대입

func fetch(url string) (*html.Node, error) {
	res, err := http.Get(url) //url에서 html데이터를 가져옴
	if err != nil {
		log.Println(err)
		return nil, err
	}

	doc, err := html.Parse(res.Body) //res.Body를 넣으면 파싱된 데이터가 리턴됨
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return doc, nil
}

func parseFollowing(doc *html.Node) []string {
	var urls = make([]string, 0)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" { // img 태그
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "avatar left" { //class가 gravatar left인 요소
					for _, a := range n.Attr {
						if a.Key == "alt" {
							fmt.Println(a.Val) //사용자 이름 출력
							break
						}
					}
				}

				if a.Key == "class" && a.Val == "avatar" {
					user := n.Parent.Attr[0].Val //부모 요소의 첫 번째 속성(href)
					//사용자 이름으로 팔로잉 url조합
					urls = append(urls, "https://github.com"+user+"/following")
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c) //재귀호출로 자식과 형제를 모두 탐색
		}
	}

	f(doc)

	return urls
}

func crawl(url string) {
	fetched.Lock()                   //맵은 뮤텍스로 보호
	if _, ok := fetched.m[url]; ok { //url 중복 처리 여부를 검사
		fetched.Unlock()
		return
	}
	fetched.Unlock()

	doc, err := fetch(url) //url에서 파싱된 html데이터를 가져옴

	fetched.Lock()
	fetched.m[url] = err //가져온 url은 맵에 url과 에러 값 저장
	fetched.Unlock()

	urls := parseFollowing(doc) //사용자 정보 출력, 팔로인 url을 구함

	done := make(chan bool)
	for _, u := range urls { //url 개수만큼
		go func(url string) { //고루틴 생성
			crawl(url) //재귀호출
			done <- true
		}(u)
	}

	for i := 0; i < len(urls); i++ {
		<-done //고루틴이 모두 실행될 때까지 대기
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	crawl("https://github.com/pyrasis/following") //크롤링 시작

}
