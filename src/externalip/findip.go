package externalip

import (
	"io"
	"net/http"
	"strings"
	"sync"
)

type findIp struct {
	ip  string
	err error
}

var lockObj sync.Mutex

func (f *findIp) Write(p []byte) (n int, err error) {
	f.ip = string(p)
	f.ip = strings.Trim(f.ip, "\n")
	return len(p), nil
}
func (f *findIp) externalIP() (string, error) {
	resp, err := http.Get("http://www.myexternalip.com/raw")
	if err != nil {
		f.ip = ""
		f.err = err
	} else {
		io.Copy(f, resp.Body)
	}

	defer resp.Body.Close()

	//fmt.Printf("%#v\n", resp)
	return f.ip, f.err
}

// 외부 IP를 반환 합니다.
func ExternalIP() (string, error) {
	lockObj.Lock()
	defer lockObj.Unlock()

	fw := findIp{}
	return fw.externalIP()
}
