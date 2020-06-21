package api

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GoDirScanner struct {
	Detail    bool
	RoutineNum int
	Domain    string
	Dictname   string
	Urls      chan string
	Client    http.Client
	Data string
	Method string
}

func errorHandle(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
var lock sync.Mutex
func showResult(url string, resp *http.Response, detail bool) {
	code := strconv.Itoa(resp.StatusCode)
	if detail{
		switch {
		case code[0:1] == "2":
			fmt.Printf("%c[4;40;32mURL: %v StatusCode: %v%c[0m \n", 0x1B,url, code,0x1B)
		case code[0:1] == "3":
			fmt.Printf("%c[4;40;34mURL: %v StatusCode: %v%c[0m \n", 0x1B,url, code,0x1B)
		case code[0:1] == "4":
			fmt.Printf("%c[0;40;33mURL: %v StatusCode: %v%c[0m \n", 0x1B,url, code,0x1B)
		case code[0:1] == "5":
			fmt.Printf("%c[0;40;31mURL: %v StatusCode: %v%c[0m \n", 0x1B,url, code,0x1B)
		}
	} else {
		switch {
		case code[0:1] == "2":
			fmt.Printf("%c[4;40;32mURL: %v StatusCode: %v%c[0m \n", 0x1B,url, code,0x1B)
		case code[0:1] == "3":
			fmt.Printf("%c[4;40;34mURL: %v StatusCode: %v%c[0m \n", 0x1B,url, code,0x1B)
		}
	}
}

func extendString(string1, string2 string) string {
	var builder strings.Builder
	builder.WriteString(string1)
	last_char := string1[len(string1):]
	first_char := string2[0:1]
	if last_char == "/" && first_char == "/" {
		string2 = string2[1:]
	}
	if last_char != "/" && first_char != "/" {
		builder.WriteString("/")
	}
	builder.WriteString(string2)
	return builder.String()
}

func GetPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = dir +"/"
	errorHandle(err)
	return strings.Replace(dir, "\\", "/", -1)
}

func (this *GoDirScanner) ReadDict(wg *sync.WaitGroup) {
	dir := GetPath()+this.Dictname
	file, err := os.Open(dir)
	errorHandle(err)
	defer file.Close()
	Scanner := bufio.NewScanner(file)
	this.Urls = make(chan string,5)
	for Scanner.Scan() {
		timer := time.NewTimer(time.Second*2)
		if Scanner.Text()=="" {
			continue
		}
		select {
			case this.Urls <- extendString(this.Domain,Scanner.Text()):
				continue
		case <- timer.C:
			if Scanner.Text()!="" {
				continue
			}
		}
	}
}

func (this *GoDirScanner) Scan(wg *sync.WaitGroup) {
	this.Client = http.Client{
		Timeout: time.Second * 10,
	}
	for {
		//time.Sleep(time.Millisecond*50)
		timeSet := time.NewTimer(time.Second*2)
		select {
			case url,ok := <-this.Urls:
				if !ok {
					wg.Done()
				}
				Request, err := http.NewRequest(this.Method, url, nil)
				if err != nil{
					continue
				}
				Request.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36")
				resp, err := this.Client.Do(Request)
				errorHandle(err)
				showResult(url, resp, this.Detail)
				errorHandle(resp.Body.Close())
				time.Sleep(time.Millisecond*3)
		case <-timeSet.C:
			defer func() {
				if r := recover(); r!=nil {}
			}()
			wg.Done()
			}
		}
	}

var GoScanner GoDirScanner
