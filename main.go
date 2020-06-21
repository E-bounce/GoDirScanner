package main

import (
	"flag"
	"fmt"
	"github.com/ebounce/GodirScanner/api"
	"sync"
	"time"
)
var wg sync.WaitGroup
func main() {
	t := time.Now()
	dirScanner := api.GoScanner
	flag.StringVar(&dirScanner.Method,"m","HEAD","指定使用扫描的方法一般就用这两种(GET,HEAD)")
	flag.BoolVar(&dirScanner.Detail,"v",false,"是否详细输出(详细输出会输出所有http返回码的内容)")
	flag.IntVar(&dirScanner.RoutineNum,"t",3,"指定需要的协程数")
	flag.StringVar(&dirScanner.Domain,"u","http://127.0.0.1:8080","设定需要扫描的域名，eg: -u https://baidu.com,请输入完整地址")
	flag.StringVar(&dirScanner.Dictname,"f","/dict/default.txt","设定使用的字典")
	flag.Parse()
	go dirScanner.ReadDict(&wg)
	for i:=0;i<dirScanner.RoutineNum;i++{
		wg.Add(1)
		go dirScanner.Scan(&wg)
	}
	wg.Wait()
	fmt.Printf("扫描已完成 共耗时 %v\n",time.Since(t))
}
