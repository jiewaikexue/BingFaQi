package request

import (
	"2024_4_15/Dingyi"
	"2024_4_15/Proxy"
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"net"
	"time"
)
import "github.com/valyala/fasthttp"

//type TextInfo Dingyi.TextInfo

//type EveryProjectRequestNeedMod func(T *TextInfo, API string) *fasthttp.Request

// 默认实现

func MakeRequestsPackage(Data *([]*Dingyi.TextInfo),
	ReqmakerNeedcreateeveryproject func(T *Dingyi.TextInfo, API string) *fasthttp.Request,
	Data_len int, API string) *[]*fasthttp.Request {
	requests := make([]*fasthttp.Request, Data_len)
	for i := 0; i < Data_len; i++ {
		requests[i] = ReqmakerNeedcreateeveryproject((*Data)[i], API)
	}
	if len(requests) <= 1 {
		Dingyi.Console_LOG.Info("MakeRequestsPackage err")
		return nil
	}
	return &requests
}

//func MakeClient(ClientAddr string, ProxyUse bool) *fasthttp.Client {
//	if !ProxyUse {
//		return &fasthttp.Client{
//			ReadTimeout:  5 * time.Second,
//			WriteTimeout: 5 * time.Second,
//		}
//	} else {
//		return &fasthttp.Client{
//			//Dial:         proxyregister(),
//			Dial: func(addr string) (net.Conn, error) {
//				proxyAddr := Proxy.SelectRandomProxy().Addr // 选择一个代理地址
//				//proxyAddr := "127.0.0.1:10809"
//				proxyConn, err := net.DialTimeout("tcp", proxyAddr, 10*time.Second)
//				if err != nil {
//					Dingyi.Console_LOG.Info("Failed to connect to proxy: " + err.Error())
//					return nil, err
//				}
//
//				// 发送 CONNECT 请求到代理服务器
//				fmt.Fprintf(proxyConn, "CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", addr, addr)
//				response := make([]byte, 4096)
//				n, err := proxyConn.Read(response)
//				if err != nil {
//					Dingyi.Console_LOG.Info("Failed to read response from proxy: " + err.Error())
//					return nil, err
//				}
//
//				// 简单检查代理响应是 200 连接建立
//				if !bytes.Contains(response[:n], []byte("200 Connection established")) {
//					Dingyi.Console_LOG.Info("Proxy failed to connect: " + string(response[:n]))
//					return nil, fmt.Errorf("proxy connection failed: %s", response[:n])
//				}
//				return proxyConn, nil
//			},
//			ReadTimeout:  5 * time.Second,
//			WriteTimeout: 5 * time.Second,
//		}
//	}
//}

func MakeClient(clientAddr string, proxyUse bool) *fasthttp.Client {
	//// 定义TLS配置
	//tlsConfig := &tls.Config{
	//	MinVersion:               tls.VersionTLS12, // 设置最小TLS版本为TLS 1.2
	//	PreferServerCipherSuites: true,             // 优先使用服务器密码套件
	//	InsecureSkipVerify:       false,            // 不跳过服务器证书验证
	//}

	if !proxyUse {
		return &fasthttp.Client{
			ReadTimeout:  20 * time.Second,
			WriteTimeout: 20 * time.Second,
			//TLSConfig:    tlsConfig,
		}
	} else {
		return &fasthttp.Client{
			Dial: func(addr string) (net.Conn, error) {
				proxyAddr := Proxy.SelectRandomProxy().Addr
				proxyConn, err := net.DialTimeout("tcp", proxyAddr, 10*time.Second)
				if err != nil {
					Dingyi.Console_LOG.Info("Failed to connect to proxy", zap.String("error", err.Error()))
					return nil, err
				}

				// 发送 CONNECT 请求到代理服务器
				fmt.Fprintf(proxyConn, "CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", addr, addr)
				response := make([]byte, 4096)
				n, err := proxyConn.Read(response)
				if err != nil {
					Dingyi.Console_LOG.Info("Failed to read response from proxy", zap.String("error", err.Error()))
					return nil, err
				}

				// 检查代理响应是否是 200 连接建立
				if !bytes.Contains(response[:n], []byte("200 Connection established")) {
					errMsg := fmt.Sprintf("Proxy connection failed: %s", response[:n])
					Dingyi.Console_LOG.Info("Proxy failed to connect", zap.String("response", errMsg))
					return nil, fmt.Errorf(errMsg)
				}
				return proxyConn, nil
			},
			ReadTimeout:  20 * time.Second,
			WriteTimeout: 20 * time.Second,
			//TLSConfig:    tlsConfig,
		}
	}
}

// 代理注册函数,在该函数内 先确定好到底使用哪一个代理
// 然后将该代理设计成一个闭包函数
// // 最后 让client注册
//
//	func proxyregister() fasthttp.DialFunc {
//		tmp_proxy := Proxy.SelectRandomProxy().Addr
//		return func(addr string) (net.Conn, error) {
//			use_proxy := tmp_proxy
//			proxyConn, err := net.DialTimeout("tcp", use_proxy, time.Second*10)
//			if err != nil {
//				return nil, err
//				Dingyi.Console_LOG.Info("proxy bind err")
//			}
//			//fmt.Fprintf(proxyConn, "CONNECT %s HTTP/1.1\r\n\r\n", addr)
//			// 需要在这里进行代理的判断,我不关心代理是谁
//			//Dingyi.Console_LOG.info("")
//			return proxyConn, nil
//		}
//	}

//===================== proxy代理注册机
//func proxyregister() fasthttp.DialFunc {
//	return func(addr string) (net.Conn, error) {
//		proxyAddr := Proxy.SelectRandomProxy().Addr // 选择一个代理地址
//		proxyConn, err := net.DialTimeout("tcp", proxyAddr, 10*time.Second)
//		if err != nil {
//			Dingyi.Console_LOG.Info("Failed to connect to proxy: " + err.Error())
//			return nil, err
//		}
//
//		// 发送 CONNECT 请求到代理服务器
//		fmt.Fprintf(proxyConn, "CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", addr, addr)
//		response := make([]byte, 4096)
//		n, err := proxyConn.Read(response)
//		if err != nil {
//			Dingyi.Console_LOG.Info("Failed to read response from proxy: " + err.Error())
//			return nil, err
//		}
//
//		// 简单检查代理响应是 200 连接建立
//		if !bytes.Contains(response[:n], []byte("200 Connection established")) {
//			Dingyi.Console_LOG.Info("Proxy failed to connect: " + string(response[:n]))
//			return nil, fmt.Errorf("proxy connection failed: %s", response[:n])
//		}
//		return proxyConn, nil
//	}
//}
