package Proxy

import (
	"2024_4_15/Dingyi"
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

// 定义代理
//var Proxies []Proxy_struct

// 结构体的构造函数，设置默认失效次数
func NewProxyStructWithTimes(addr string, times int) Dingyi.Proxy_struct {
	return Dingyi.Proxy_struct{
		Addr: addr,
		//UnUseTimes: times,
		//Available:  true,
		//Mutex: sync.RWMutex{},
	}
}

// 默认10次失效后放弃
func NewProxyStruct(addr string) Dingyi.Proxy_struct {
	return Dingyi.Proxy_struct{
		Addr: addr,
		//UnUseTimes: 500,
		//Available:  true,
		//Mutex: sync.RWMutex{},
	}
}

func ReadInProxyPool(dir string) []Dingyi.Proxy_struct {
	file, err := os.OpenFile(dir, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println("代理文件打开失败:", err)
	}
	// 随手开,随手关 -> 后面自动关闭
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("代理文件关闭失败")
		}
	}(file)
	//Proxies 为代理池
	proxies := make([]Dingyi.Proxy_struct, 0)
	myscanner := bufio.NewScanner(file)
	for myscanner.Scan() {
		// 从文件中一行行读取ip+port
		Ip := myscanner.Text()
		proxies = append(proxies, NewProxyStruct(Ip))

	}
	if err := myscanner.Err(); err != nil {
		fmt.Println("代理文件读取错误:", err)
	}

	return proxies
}

// 代理更换策略:到达某一个特定的重试次数就直接 直接扔一个错误, ==> 抛弃
// 从代理列表中随机选择一个代理地址
// 没有进行代理的检测
func SelectRandomProxy() *Dingyi.Proxy_struct {
	proxies := *(Dingyi.GetProxyPool())
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tmp_proxy := proxies[r.Intn(len(proxies))]
	return &tmp_proxy
}

/*
代理检测模块
*/
func PingTargetUrlUseAll(TargetUrl string, Proxies *[]Dingyi.Proxy_struct) {
	for _, proxy := range *Proxies {
		fmt.Println("检测目标: %s   使用ping命令. 代理地址: %s", TargetUrl, proxy.Addr)
		if pingWithProxy(TargetUrl, proxy.Addr) {
			fmt.Println("代理: %s 有效", proxy.Addr)
		} else {
			fmt.Println("代理 %s  失效", proxy.Addr)
		}
	}
}

/*
代理检测模块
ping某个代理
*/
func pingWithProxy(TargetUrl, proxyAddr string) bool {
	// 设置代理
	proxy, err := net.Dial("tcp", proxyAddr)
	if err != nil {
		fmt.Println("Error connecting to Proxy:", err)
		return false
	}
	defer proxy.Close()

	// 设置Ping超时
	timeout := time.Duration(7 * time.Second)
	conn, err := net.DialTimeout("tcp", TargetUrl, timeout)
	if err != nil {
		fmt.Println("Error pinging target via Proxy:", err)
		return false
	}
	defer conn.Close()
	return true
}
