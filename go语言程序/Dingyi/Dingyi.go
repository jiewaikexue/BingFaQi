package Dingyi

import (
	"fmt"
	"github.com/iunary/fakeuseragent"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"os"
	"reflect"
)

/*
ProjectName string            | 项目名称  -> 日志名称
ProxyDir string               | 代理     -> 代理内容
ReadTxtDir string             | 账号信息相关
TemplateLine string           | 单行账号格式

type TextInfo struct		  | 单行账号信息的结构体       -> 一定需要重新写
LenTemplate int				  | 结构体内部字段个数         -> 只需要自定义这个
WhoIsMapKey string            | map唯一的 主键


AllTxtInfoMap   map[string]*TextInfo  | map 存储所有 WhoIsMapKey 为key的完整的TextInfo  -> 最后做索引使用
Data  []*TextInfo                     | 所有账号信息 只为了给 http/https请求提供构造基础信息
Time    					  |      定时器倒计时



Proxies []Proxy_struct        |   代理池

type Proxy_struct struct {    |   单个代理节点
	Addr string
	//Available  bool
	UnUseTimes int          // 默认访问次数,每一次修改
	Mutex      sync.RWMutex // 读写锁,
}

--------------------------------
File_log  初始化
Console_Log 初始化              -> 用函数就可以了


*/

//===============结构体,每个

var WhoIsMapKey = "Address"
var WhoIsKeyToFindMasterINReq = "buyerAddress"

type TextInfo struct {
	Master    string
	Address   string
	PublicKey string
}

// =====================================
// var File_LOG = &zap.Logger{}    // 日志输出文件
var Console_LOG = &zap.Logger{} //换个地方 这里不好用
// var Err = errors.New("")        //错误
var GrayLogUseTCP bool = true
var GrayLogAddress string = "127.0.0.1:12201" //新增graylog输入的话 需要再docker配置文件内进行端口映射
var GrayLogShowName string = "ordzaar"

// ======================================
// 所有的requests使用那个API
var API string = "https://api-mainnet.magiceden.io/v2/ord/launchpad/psbt/minting"

// 主机的地址,所有请求去那个地址
var ClientAddr string = "api-mainnet.magiceden.io"

// 外循环循环次数,控制这所有的resquests请求内容,总共轮训几次
var WaiXunHuan int = 5

// 倒计时开始时间
var TimeString string = "2024-04-17 03:31:01"

// 定时器使用的时区
var TimeLocationString string = "Asia/Shanghai"

// 最后main运行结束,等待几秒好让全部的异步写入文件
const WaitTimAtEnd = 20

// 单个请求重试次数
const RequestTimes = 2

// 超时时间
const EveryReqTimeout = 5

// ======================================
func DiyEveryProject(T *TextInfo, API string) *fasthttp.Request {
	tmp := fasthttp.AcquireRequest()
	if tmp == nil {
		os.Exit(1)
	}
	All_url := API + "?collectionSymbol=runesterminal&price=1" + "&buyerAddress=" + T.Address + "&buyerTokenReceiveAddress=" + T.Address + "&buyerPublicKey=" + T.PublicKey + "&feerateTier=fastestFee"
	//fmt.Println(All_url)
	tmp.SetRequestURI(All_url)
	tmp.Header.SetMethod(fasthttp.MethodGet)

	//tmp.Header.Set("Cookie", "__cf_bm=DCF5.6O2sAtkwBv1nYoMeRz804Ga3vNpv5msqCw7frE-1713301443-1.0.1.1-4UCaLkpFZ2kiHE4xZIiNYXs3.QhEPKbEMw8zA1tNdTpPRK01bbLlx1YWmpriU3ag5ucUgwkIKNn1km0Tdv0eYQ; _cfuvid=zhQioO83NMvrtecJtVuYVVllbwuBCRGCurg.9wqA5g4-1713301443101-0.0.1.1-604800000; __cf_bm=GSvjXnH1sxhn3hMVXD.bfNUtisRHbXNix37_sXCunRY-1713329411-1.0.1.1-Ua9mESmCFiNRuTHIgNZIUmYnpW0A5Prgol31peIbFAbPdyR9bPhLOxH2OzxCPZ.TIF5Owd9az8TLDbIko.RYHA; _cfuvid=VHQ0NSNR7KF8WVhysmhCiiXsXoleZgAtBqRekTSPsTM-1713329411326-0.0.1.1-604800000")
	tmp.Header.SetHost(ClientAddr)
	//randomAgent := fakeuseragent.RandomUserAgent()
	//tmp.Header.SetUserAgent(randomAgent)
	//tmp.Header.Set("Accept-Encoding", "gzip, deflate, br")
	//tmp.Header.Set("Connection", "Keep-Alive")
	//tmp.Header.Set("Accept", "application/json")

	//Host: api-mainnet.magiceden.io
	//Connection: keep-alive
	//Pragma: no-cache
	//	Cache-Control: no-cache
	//	sec-ch-ua: "Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"
	//	sec-ch-ua-mobile: ?0
	//	sec-ch-ua-platform: "Windows"
	//	Upgrade-Insecure-Requests: 1
	//	User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36
	//Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
	//	Sec-Fetch-Site: none
	//	Sec-Fetch-Mode: navigate
	//	Sec-Fetch-User: ?1
	//	Sec-Fetch-Dest: document
	//	Accept-Encoding: gzip, deflate, br, zstd
	//	Accept-Language: zh-CN,zh;q=0.9
	//	Cookie: cf_clearance=PNMrJ09OaPeGMBk4Uc5pBy6wMy_w71MSqVopABiI0_I-1713287258-1.0.1.1-.5.QwpS6tsl54flt1H8qUTM9gKiutQqzK5wpPur23YV5UcRRT0jsxh_Klusn.9TOtYMD8XRL1DhhGW6SB1LavA; rl_user_id=RudderEncrypt%3AU2FsdGVkX18yrSO5pzhX1Slq3xj0PYT3Ro%2BLeDwDCzA%3D; rl_trait=RudderEncrypt%3AU2FsdGVkX1%2FDnIKpK5pu2KbJFMgc4tcakcQAWST%2BjdE%3D; rl_group_id=RudderEncrypt%3AU2FsdGVkX19WwIBEvRMbak0NVN5%2BeIhjA9dGkwpSp7E%3D; rl_group_trait=RudderEncrypt%3AU2FsdGVkX1%2FE6uEWWWPe9DmlH72FZGRgTBFwGj68Ddw%3D; rl_anonymous_id=RudderEncrypt%3AU2FsdGVkX19fn9upYKvr4i34BQxdhIKRa0k4RXgiGgoSRlFQUUhM9fA9j8cii%2BLdL1%2FgOxRC2jPQvRgsTTVnaA%3D%3D; rl_page_init_referrer=RudderEncrypt%3AU2FsdGVkX19gqF4oNptwIhcpNKRm9kmjoV1laoS4oH4%3D; rl_page_init_referring_domain=RudderEncrypt%3AU2FsdGVkX19SSIwwS4TL%2BbZVti1wMlzbj9EQd9BVbSQ%3D; rs_ga=GA1.1.e7d7e8a1-0d9b-4b33-bab7-08f3259c950f; rs_ga_8BCG117VGT=GS1.1.1713287242717.1.0.1713287244.60.0.0; rl_session=RudderEncrypt%3AU2FsdGVkX19ikPPduelRK97axUh0MZGqaYkrQZf283Ge%2F946NZ4Phxn0%2FkkMUAkVprlaXzcPW60%2ByFKZ4IyGs1pNzSeZp4PqdNnjR3UBn7zk8T81UN9vTV3EBGEw9Vcc68bF%2FDE4NEGBDxACHPsf5A%3D%3D; intercom-id-htawnd0o=083063de-4c1a-4faf-a5e5-b37189425526; intercom-session-htawnd0o=; intercom-device-id-htawnd0o=bf931941-0ead-48f3-bf85-89229878aa3e; _cfuvid=_Iu86PLUqMYzZCzMn2Gw0pwO46H6C1EP4SmDtb3Dxik-1713330152438-0.0.1.1-604800000; __cf_bm=gRB.7Z1B_fkeUpmqn4FU97tBN6LT57MDuYQ4iO5w7s4-1713333104-1.0.1.1-cZcSKD_u82l_038UWrbvdMSsKR9Zjt7KAMzEAfq2G1aoKP6f.Ejq9o7_m469bsWlc15BK2ZOCABJyQo02Eet6Q

	tmp.Header.Set("Connection", "keep-alive")
	tmp.Header.Set("Cache-Control", "no-cache")
	//tmp.Header.Set("sec-ch-ua",{""\"Google Chrome\"","v":"123", "\"Not:A-Brand\"","v":"8", "\"Chromium\"","v":"123\"},)
	tmp.Header.Set("sec-ch-ua", "\"Google Chrome\""+","+"v:123"+","+"\"Not:A-Brand\""+","+"v:8"+","+"\"Chromium\""+","+"v:123\"")
	tmp.Header.Set("sec-ch-ua-mobile", "?0")
	tmp.Header.Set("Upgrade-Insecure-Requests", "1")

	//tmp.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	//chromeAgent := fakeuseragent.GetUserAgent(fakeuseragent.BrowserChrome) + "(Accepts JSON only)"
	randomAgent := fakeuseragent.RandomUserAgent() + "(Accepts JSON only)"
	tmp.Header.SetUserAgent(randomAgent)

	tmp.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//tmp.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	tmp.Header.Set("Sec-Fetch-Site", "none")
	tmp.Header.Set("Sec-Fetch-Mode", "navigate")
	tmp.Header.Set("Sec-Fetch-User", "?1")
	tmp.Header.Set("Sec-Fetch-Dest", "document")
	tmp.Header.Set("Accept", "application/json")
	tmp.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	tmp.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	tmp.Header.Set("Cookie",
		"cf_clearance=PNMrJ09OaPeGMBk4Uc5pBy6wMy_w71MSqVopABiI0_I-1713287258-1.0.1.1-.5.QwpS6tsl54flt1H8qUTM9gKiutQqzK5wpPur23YV5UcRRT0jsxh_Klusn.9TOtYMD8XRL1DhhGW6SB1LavA; rl_user_id=RudderEncrypt%3AU2FsdGVkX18yrSO5pzhX1Slq3xj0PYT3Ro%2BLeDwDCzA%3D; rl_trait=RudderEncrypt%3AU2FsdGVkX1%2FDnIKpK5pu2KbJFMgc4tcakcQAWST%2BjdE%3D; rl_group_id=RudderEncrypt%3AU2FsdGVkX19WwIBEvRMbak0NVN5%2BeIhjA9dGkwpSp7E%3D; rl_group_trait=RudderEncrypt%3AU2FsdGVkX1%2FE6uEWWWPe9DmlH72FZGRgTBFwGj68Ddw%3D; rl_anonymous_id=RudderEncrypt%3AU2FsdGVkX19fn9upYKvr4i34BQxdhIKRa0k4RXgiGgoSRlFQUUhM9fA9j8cii%2BLdL1%2FgOxRC2jPQvRgsTTVnaA%3D%3D; rl_page_init_referrer=RudderEncrypt%3AU2FsdGVkX19gqF4oNptwIhcpNKRm9kmjoV1laoS4oH4%3D; rl_page_init_referring_domain=RudderEncrypt%3AU2FsdGVkX19SSIwwS4TL%2BbZVti1wMlzbj9EQd9BVbSQ%3D; rs_ga=GA1.1.e7d7e8a1-0d9b-4b33-bab7-08f3259c950f; rs_ga_8BCG117VGT=GS1.1.1713287242717.1.0.1713287244.60.0.0; rl_session=RudderEncrypt%3AU2FsdGVkX19ikPPduelRK97axUh0MZGqaYkrQZf283Ge%2F946NZ4Phxn0%2FkkMUAkVprlaXzcPW60%2ByFKZ4IyGs1pNzSeZp4PqdNnjR3UBn7zk8T81UN9vTV3EBGEw9Vcc68bF%2FDE4NEGBDxACHPsf5A%3D%3D; intercom-id-htawnd0o=083063de-4c1a-4faf-a5e5-b37189425526; intercom-session-htawnd0o=; intercom-device-id-htawnd0o=bf931941-0ead-48f3-bf85-89229878aa3e; _cfuvid=_Iu86PLUqMYzZCzMn2Gw0pwO46H6C1EP4SmDtb3Dxik-1713330152438-0.0.1.1-604800000; __cf_bm=gRB.7Z1B_fkeUpmqn4FU97tBN6LT57MDuYQ4iO5w7s4-1713333104-1.0.1.1-cZcSKD_u82l_038UWrbvdMSsKR9Zjt7KAMzEAfq2G1aoKP6f.Ejq9o7_m469bsWlc15BK2ZOCABJyQo02Eet6Q")

	//tmp.Header.SetContentType("utf-8")
	return tmp
}

//=================================================================================================
//=================================================================================================
//=================================================================================================
//=================================================================================================

// 使用反射设置结构体的值==================
// ================= 这个函数写死,不用动
func (t *TextInfo) SetValues(values []string) {
	// 获取结构体的反射值
	v := reflect.ValueOf(t).Elem()
	// 如果 values 的长度与结构体字段的数量不匹配，则输出错误信息并返回
	if len(values) != v.NumField() {
		fmt.Println("参数个数错误")
		return
	}

	// 遍历结构体的每个字段
	for i := 0; i < v.NumField(); i++ {
		// 获取字段的反射值
		fieldValue := v.Field(i)
		// 如果字段可以被设置，则将 values 中对应位置的值赋给该字段
		if fieldValue.CanSet() {
			fieldValue.SetString(values[i])
		}
	}
}

// ++++++++++++代理相关+++++++++++++++++++
type Proxy_struct struct {
	Addr string
	//Available  bool
	//UnUseTimes int // 默认访问次数,每一次修改
	//Mutex      sync.RWMutex // 读写锁,
}

// 定义一个函数,用来获取代理池
var ProxyPool []Proxy_struct = make([]Proxy_struct, 0)
var ProxyDir string = "Proxy/proxy.txt"

// 这个参数仅仅是给proxy文件里面的selectproxy函数使用
func GetProxyPool() *[]Proxy_struct {
	return &ProxyPool
}
