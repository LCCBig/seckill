package utils

import (
	"github.com/gin-gonic/gin"
)

/**
设置Cookie的值，并使其在指定时间内生效
https://www.kancloud.cn/lhj0702/sockstack_gin/1805369
*/
func DoSetCookie(ctx *gin.Context, cookieName string, cookieValue string, cookieMaxage int, secure bool) {
	//http://localhost:9090/home/tologin
	//fmt.Println("请求c.Request.URL.Host:", ctx.Request.URL.Host)                 // 没有任何输出
	//fmt.Println("请求c.Request.URL.Hostname():", ctx.Request.URL.Hostname())     // 没有任何输出
	//fmt.Println("请求c.Request.URL.Port():", ctx.Request.URL.Port())             // 没有任何输出
	//fmt.Println("请求c.Request.URL.String():", ctx.Request.URL.String())         // 输出: /home/tologin
	//fmt.Println("请求c.Request.URL.Scheme", ctx.Request.URL.Scheme)              // 没有任何输出
	//fmt.Println("请求c.Request.URL.RequestURI():", ctx.Request.URL.RequestURI()) // 输出: /home/tologin
	//fmt.Println("请求c.Request.Host:", ctx.Request.Host)                         // 输出: localhost:9090
	//fmt.Println("请求c.Request.RequestURI:", ctx.Request.RequestURI)             // 输出: /home/tologin

	//获取域名
	domain := ctx.Request.Host
	//name, value string, maxAge int, path, domain string, secure, httpOnly bool
	ctx.SetCookie(cookieName, cookieValue, cookieMaxage, "", domain, secure, false)
}
