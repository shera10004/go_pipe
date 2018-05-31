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

type result struct { //결과값을 저장할 구조체
	url  string //가져온 url
	name string //사용자 이름
}

const nilTag = "<nil>"

func fetch(url string) (*html.Node, error) {
	fmt.Println(">fetch url:", url)
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

func parseFollowing(gid int, doc *html.Node, urls chan string) <-chan string {
	name := make(chan string)

	go func(gid int) { //교착 상태가 되지 않도록 고루틴으로 실행
		count := 0

		var f func(*html.Node)

		f = func(n *html.Node) {
			//fmt.Println(" f func() start...", gid)
			findTarget := false
			if n.Type == html.ElementNode {

				//main page
				for _, at := range n.Attr {
					if at.Key == "class" && at.Val == "follow-list-name" {
						for _, at := range n.FirstChild.Attr {
							if at.Key == "title" {
								fmt.Println(">find name:", at.Val, ", gid:", gid)
								name <- at.Val

								findTarget = true
								break
							}
						} //for
					} //if

					if findTarget == true {
						isUrl := false
						for _, at := range n.FirstChild.FirstChild.Attr {
							if at.Key == "href" {
								furl := "https://github.com" + at.Val + "?tab=following"
								fmt.Println(">find href:", furl, ", gid:", gid)
								urls <- furl
								isUrl = true
								break
							}
						} //for
						if isUrl == false {
							fmt.Println("?????????? ", gid)
						}
						break
					} //if
				} //main for

				//sub page
				if findTarget == false {
					/*
						for _, at := range n.Attr {
							if at.Key == "class" && at.Val == "f4 link-gray-dark" {
								fmt.Println("n.Data", n.Data)

								for _, at := range n.FirstChild.Attr {
									if at.Key == "title" {
										fmt.Println(">find name:", at.Val)
										name <- at.Val

										findTarget = true
										break
									}
								} //for
							} //if

							if findTarget == true {
								for _, at := range n.FirstChild.FirstChild.Attr {
									if at.Key == "href" {
										furl := "https://github.com" + at.Val + "?/tab=following"
										fmt.Println(">find href:", furl)
										urls <- furl
										break
									}
								} //for
								break
							} //if
						} //sub for
						//*/
				}

			} //if

			if findTarget == false {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					//fmt.Println(">inn counter:", count)
					count++
					f(c) //재귀호출로 자식과 형제를 모두 탐색
				}
			} else {
				fmt.Println(">search end!", count)
			}

		} //end of f func()

		fmt.Println(">gen counter:", doc)
		if doc == nil {
			fmt.Println("nil?????")
			name <- nilTag
		} else {
			f(doc)
		}

	}(gid)

	return name
}

func crawl(gid int, url string, urls chan string, c chan<- result) {
	fmt.Println(">crawl inn", url)
	fetched.Lock()                   //맵은 뮤텍스로 보호
	if _, ok := fetched.m[url]; ok { //url 중복 처리 여부를 검사
		fmt.Println(">crawl out", url)
		fetched.Unlock()
		return
	}
	fetched.Unlock()

	doc, err := fetch(url) //url에서 파싱된 html데이터를 가져옴
	if err != nil {        //url을 가져오지 못했을 때
		fmt.Println(">crawl_err url:", url)
		go func(u string) { //교착 상태가 되지 않도록 고루틴을 생성
			urls <- u //채널에 url을 보냄
		}(url)
	}

	fetched.Lock()
	fetched.m[url] = err //가져온 url은 맵에 url과 에러 값 저장
	fetched.Unlock()

	name := <-parseFollowing(gid, doc, urls) //사용자 정보 출력, 팔로인 url을 구함
	fmt.Println(">>>>> rgo_result name:", name, ", url:", url)
	c <- result{url, name} //가져온 url과 사용자 이름을 구조체 인스턴스로 생성하여 채널 c에 보냄

}

//실제 작업을 처리하는 worker함수
func worker(gid int, done <-chan struct{}, urls chan string, c chan<- result) {
	for url := range urls { //urls 채널에서 url을 가져옴
		select {
		case <-done: //채널이 닫히면 worker 함수를 빠져나옴
			return
		default:
			fmt.Println("---- worker_crawl ----", gid)
			crawl(gid, url, urls, c) //url 처리
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	urls := make(chan string)   // 작업을 요청할 채널
	done := make(chan struct{}) // 작업 고루틴에 정지 신호를 보낼 채널
	c := make(chan result)      // 결과값을 저장할 채널

	var wg sync.WaitGroup
	const numWorkers = 10
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(n int) {
			worker(n, done, urls, c)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait() // 고루틴이 끝날때 까지 대기
		close(c)  // 대기가 끝나면 결과값 채널을 닫음
	}()

	urls <- "https://github.com/pyrasis/following"
	limitCount := 100

	count := 0
	for r := range c { // 결과 채널에 값이 들어올 때까지 대기한 뒤 값을 가져옴
		fmt.Println("### result:", r.name)

		if r.name == nilTag {
			close(done)
			fmt.Println(">nil :", r.url)
			break
		}

		count++
		if count > limitCount {
			close(done) //done을 닫아서 worker 고루틴을 종료
			break
		}
	}
}
