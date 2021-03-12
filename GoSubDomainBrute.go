package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/panjf2000/ants/v2"
	"net"

	"net/http"
	"os"
	"sync"
	"time"
)

var(
	help bool
	target string
	dict string
	targets []string
	rate int

)


var c *http.Client = &http.Client{
	Transport: &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*2)
			if err != nil {
				return nil, err
			}
			return c, nil

		},
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * 2,
	},
}


func init(){
	flag.BoolVar(&help,"h, --help",false,"help, 帮助命令")
	flag.StringVar(&target,"target","","目标url,例如 -target example.com")
	flag.StringVar(&dict,"dict","","字典文件")
	flag.IntVar(&rate,"rate",100,"速率,速率越高误差越大")
	flag.Usage = usage
	flag.Parse()

}

func usage(){
	fmt.Fprintf(color.Output,color.HiCyanString(`Go语言子域名爆破工具
Options:
`))
	flag.PrintDefaults()
}


func Brute(target string) bool{
	resp,err := c.Get(target)
	if resp != nil {
		defer resp.Body.Close()
		return true
	}
	if err != nil{
		return false
	}
	return false
}


func main(){
	start := time.Now()
	dictfile, err := os.OpenFile(dict,os.O_RDONLY,1)
	if err != nil {
		fmt.Fprintln(color.Output,time.Now().Format("2006-01-02 15:04:05"),color.RedString("[ERROR]")+":",err)
	}
	dicts := bufio.NewScanner(dictfile)
	for dicts.Scan() {
		tmp := dicts.Text()
		targettmp := "http://" + tmp + "." + target
		targets = append(targets, targettmp)
	}
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(rate, func(i interface{}) {
		if (Brute(i.(string))){
			fmt.Fprintln(color.Output,time.Now().Format("2006-01-02 15:04:05"),color.HiCyanString(i.(string)))
		}
		wg.Done()
	})
	for _,v := range targets{
		wg.Add(1)
		_ = p.Invoke(v)
	}
	wg.Wait()
	end := time.Since(start)
	fmt.Fprintln(color.Output,color.HiGreenString("共用时:",end))

}
