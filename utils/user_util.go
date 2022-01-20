package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"seckill/dbs"
	"seckill/models"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type SyncWriteBuffer struct {
	//锁
	BufferLock sync.Mutex
	//缓冲区
	Buffer *bufio.Writer
}

func (sw *SyncWriteBuffer) Write(b []byte) (nn int, err error) {
	defer sw.BufferLock.Unlock()
	sw.BufferLock.Lock()
	return sw.Buffer.Write(b)
}

/**
向数据库中插入
*/
func CreateUser() {
	//获取Mysql连接
	clinet := dbs.GetMysqlClinet()
	wg := sync.WaitGroup{}

	//设置原子类型，防止并发问题
	var num int32 = 31

	for i := 31; i <= 5031; i++ {
		wg.Add(1)
		go func() {
			for {
				usernameNum := atomic.LoadInt32(&num)
				if !atomic.CompareAndSwapInt32(&num, usernameNum, usernameNum+1) {
					fmt.Println("cas失败！！！")
					continue
				}
				var user models.User
				user.UserId = int(usernameNum)
				user.UserName = "1800000" + fmt.Sprintf("%04d", usernameNum)
				user.Password = "80997b0bea687adf917aa98d4e763f3e"
				user.Salt = "xxyyzz"
				user.NikeName = "test" + strconv.FormatInt(int64(usernameNum), 10)
				user.InsertUser(clinet)
				wg.Done()
				//退出循环，要不将后面这段user插入部分放for 外边也行
				break
			}
		}()
	}
	wg.Wait()
	fmt.Println("插入完成")
}

/**
发送http请求用户登陆
*/
func UserLogin(startNum int) {
	time.Sleep(time.Second * 10)
	fmt.Println("登陆开始！！！！！！！！")
	//获取表表中最大的id
	clinet := dbs.GetMysqlClinet()
	row := clinet.QueryRow("SELECT MAX(id) FROM t_user")
	var idMax int
	err := row.Scan(&idMax)
	if err != nil {
		fmt.Println("查询出错")
		panic(err)
	}
	//id Channal 装入ID
	var idChan chan int
	idChan = make(chan int, 1000)

	//同步机制
	wg := sync.WaitGroup{}

	go func() {
		for i := startNum; i <= idMax; i++ {
			idChan <- i
		}
		close(idChan)
	}()

	//记录成功次数
	var successful int32 = 0

	//写入文件
	file, err := os.OpenFile("/Users/gaea/Desktop/config.txt", os.O_RDWR, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
		return
	}

	sw := SyncWriteBuffer{sync.Mutex{}, bufio.NewWriter(file)}

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for id := range idChan {
				//通过id查询数据
				var username string
				queryRow := clinet.QueryRow("SELECT username FROM t_user WHERE id = ?", id)
				err := queryRow.Scan(&username)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(username)

				//http 客户端
				httpClient := &http.Client{}
				//配置form参数
				urlMap := url.Values{}
				urlMap.Add("username", username)
				urlMap.Add("password", "569efca813bd9fc07037334b1f1608cc")
				parms := ioutil.NopCloser(strings.NewReader(urlMap.Encode())) //把form数据编下码
				//发送登陆请求
				request, err := http.NewRequest("POST", "http://localhost:9090/user/login", parms)
				if err != nil {
					panic(err)
					return
				}
				request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				//获取response
				resp, err := httpClient.Do(request)

				if err != nil {
					panic(err)
					return
				}

				//获取 cookie
				cookie := resp.Cookies()[0]
				fmt.Println(cookie.Value)
				atomic.AddInt32(&successful, 1)
				//关闭流，防止超时
				resp.Body.Close()

				//并发写入缓冲区 ,害怕缓冲区满 可以一定条件就Flush()写入文件，再次调用Write()写在再Flush()
				_, Err := sw.Write([]byte(username + "," + cookie.Value + "\n"))
				if Err != nil {
					panic(Err)
				}
			}
		}()
	}
	wg.Wait()
	sw.Buffer.Flush()
	fmt.Printf("成功次数：%d\n", successful)
}
