package main

import (
	"2024_4_15/Dingyi"
	"2024_4_15/GrayLogWithZapByTCP"
	"2024_4_15/GrayLogWithZapByUDP"
	"2024_4_15/Proxy"
	"2024_4_15/ReadTxtInitResq"
	"2024_4_15/Time"
	"2024_4_15/request"
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"sync"
	"time"
)

var AllTxtInfoMap *map[string]*Dingyi.TextInfo
var Data = make([]*Dingyi.TextInfo, 0)

//import "Console"  "Dingyi.Console_LOG"

func sendRequests(requests *[]*fasthttp.Request, client *fasthttp.Client, wg *sync.WaitGroup) {
	defer wg.Wait() // 等待所有请求处理完成
	for _, req := range *requests {
		go func(req *fasthttp.Request) {
			defer wg.Done()
			startTime := time.Now() // 记录请求开始时间
			resp := fasthttp.AcquireResponse()
			defer fasthttp.ReleaseResponse(resp)
			if err := client.Do(req, resp); err != nil {
				Dingyi.Console_LOG.Error("Request failed", zap.Error(err))
				return
			}
			duration := time.Since(startTime)   // 计算请求耗时
			key := extractKeyFromRequest(req)   // 从请求中提取关键字段
			handleResponse(resp, duration, key) // 将响应、耗时和关键字段一起处理
		}(req)
	}
}

func extractKeyFromRequest(req *fasthttp.Request) string {
	// 确保提取的键与Map中的键完全匹配
	key := string(req.URI().QueryArgs().Peek(Dingyi.WhoIsKeyToFindMasterINReq))
	//fmt.Println("Extracted key:", key)
	return key
}

func handleResponse(resp *fasthttp.Response, duration time.Duration, key string) {
	body := resp.Body()
	statusCode := resp.StatusCode()
	info, exists := (*AllTxtInfoMap)[key]
	//fmt.Println(*AllTxtInfoMap)
	if statusCode != 200 {
		//Console_Log.Info("Request succeeded",
		//	zap.Int("status", statusCode),
		//	zap.Duration("duration", duration),
		//	zap.String("Master", fmt.Sprintf("%v", info.Master)),
		//	zap.ByteString("body", body)) // 错误信息

		Dingyi.Console_LOG.Info("Request succeeded",
			zap.Int("status", statusCode),
			zap.Duration("duration", duration),
			zap.String("Master", info.Master),
			zap.ByteString("body", body)) // 记录成功响应的body和额外信息
		return
	}

	// 从 AllTxtInfoMap 中获取与 key 相关的信息
	if !exists {
		Dingyi.Console_LOG.Error("Key not found in AllTxtInfoMap", zap.String("key", key))
		return
	}

	// 如果状态码为200，记录成功日志、响应体、耗时和额外信息
	Dingyi.Console_LOG.Info("Request succeeded",
		zap.Int("status", statusCode),
		zap.Duration("duration", duration),
		zap.String("Master", fmt.Sprintf("%v", info.Master)),
		zap.ByteString("body", body)) // 记录成功响应的body和额外信息
}

//var Console_LOG = &zap.Logger{}

func main() {
	defer sync_time() // 确保资源正确清理
	if Dingyi.GrayLogUseTCP {
		Dingyi.Console_LOG = GrayLogWithZapByTCP.GetZapLogger("info", true, Dingyi.GrayLogAddress)
	} else {
		Dingyi.Console_LOG = GrayLogWithZapByUDP.GetZapLogger("info", false, Dingyi.GrayLogAddress)
	}

	AllTxtInfoMapTemp, Data, err := ReadTxtInitResq.ReadTxtInitResq("main/test.txt", Dingyi.WhoIsMapKey)
	if err != nil {
		Dingyi.Console_LOG.Error("Failed to initialize resources", zap.Error(err))
		return
	}
	AllTxtInfoMap = AllTxtInfoMapTemp // 必须这样写,不然会被提前释放
	Dingyi.ProxyPool = Proxy.ReadInProxyPool("Proxy/proxy.txt")
	if len(*AllTxtInfoMap) == 0 {
		return
	}
	if err != nil {
		Dingyi.Console_LOG.Error("Failed to initialize resources", zap.Error(err))
		return
	}

	requests := request.MakeRequestsPackage(&Data, Dingyi.DiyEveryProject, len(Data), Dingyi.API)
	if requests == nil {
		Dingyi.Console_LOG.Error("No requests to process")
		return
	}

	timer, err := Time.NewScheduledTimerFromString(Dingyi.TimeString, Dingyi.TimeLocationString)

	if err != nil {
		Dingyi.Console_LOG.Error("Failed to create timer", zap.Error(err))
		return
	}
	defer timer.Stop()

	client := request.MakeClient(Dingyi.ClientAddr, false) //=============================> 是否使用代理
	if client == nil {
		Dingyi.Console_LOG.Error("Failed to create HTTP client")
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(*requests) * 100)
	<-timer.C // 等待定时器触发
	for i := 0; i < 100; i++ {
		go sendRequests(requests, client, wg)
	}

	wg.Wait() // 等待所有请求处理完成
	Dingyi.Console_LOG.Info("Timer triggered, starting request sending")

	Dingyi.Console_LOG.Info("All requests have been processed")
	return
}

func sync_time() {
	time.Sleep(Dingyi.WaitTimAtEnd * time.Second)
}
