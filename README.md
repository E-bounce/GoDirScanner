# Project

自制利用`Golang`实现简单的目录扫描器

# Install

如果你的本地有`Golang`语言环境，可以直接

```go
go build main.go
```

生成二进制可执行文件，然后就可以在命令行里愉快的使用了。

# Help

使用`-h`参数能够查看简易的帮助:

```shell
./main -h
Usage of ./main:
  -f string
    	设定使用的字典 (default "./dict/default.txt")
  -m string
    	指定使用扫描的方法一般就用这两种(GET,HEAD) (default "HEAD")
  -t int
    	指定需要的协程数 (default 3)
  -u string
    	设定需要扫描的域名，eg: -u https://baidu.com,请输入完整地址 (default "http://127.0.0.1:8080")
  -v	是否详细输出(详细输出会输出所有http返回码的内容)
```

能够得到一些简易的帮助。

# Example

基本只需要指定扫描的`URL`，其他东西都不用管，默认使用的是`dirsearch`项目的字典，如果需要指定字典文件，可以把[fuzzDict](https://github.com/TheKingOfDuck/fuzzDicts0)项目里面的字典搬过来，指定即可，以本地`python开启的简单http服务器为例`

```shell
python3 -m http.server 8080
/main -u http://127.0.0.1:8080
```

默认字典跑完大约需要`10s`左右，由于怕扫的太快，所以源码中添加的休眠限制一下速度，如果有进一步的需求可直接删除休眠代码，`-t 20`时，大概只需要`4s`的时间扫完，

# About Code

所有的源码均在`api.go`里，基本使用一个`Go`结构体完成所有数据方面的制定~~(强行写成了面向对象的编程)~~，其他的地方主要就是`select`配合`time.Newtimer`进行超时无阻塞的并发操作了，主要原因是本地测试的时候，读字典的`GoRoutine`跑的没有拿数据的协程快，避免协程直接挂掉，才出此下策解决。