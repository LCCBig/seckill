package good

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"seckill/common"
	"seckill/controllers"
	"seckill/dbs"
	"seckill/models"
	"seckill/services"
	"seckill/utils"
	"seckill/vo"
	"strconv"
	"time"
)

func ToList(ctx *gin.Context) {
	//获取用户
	_, exists := ctx.Get("user")
	if !exists {
		//如果用户未登陆或登陆状态过期，则重定向到登陆界面
		ctx.Redirect(http.StatusTemporaryRedirect, "/home/tologin")
		return
	}
	goodsList := services.GetGoodList()
	ctx.HTML(http.StatusOK, "goodList.tmpl", gin.H{"goodsList": goodsList})
}

/**
进入商品详情界面
*/
func ToDetail(ctx *gin.Context) {
	//获取用户
	user, exists := ctx.Get("user")
	if !exists {
		//如果用户未登陆或登陆状态过期，则重定向到登陆界面
		ctx.Redirect(http.StatusTemporaryRedirect, "/home/tologin")
		return
	}
	idStirng := ctx.Param("secKillId")
	secKillId, err := strconv.Atoi(idStirng)
	if err != nil {
		//参数出错
		print(err)
	}
	goodVo := services.GetGood(secKillId)

	startDate := goodVo.StartDate
	endDate := goodVo.EndDate
	nowDate := time.Now()

	//秒杀状态
	secKillStatus := 0
	//剩余时间
	remainSeconds := 0

	if nowDate.Before(startDate) {
		//秒杀还未开始
		remainSeconds = startDate.Second() - nowDate.Second()
	} else if nowDate.After(endDate) {
		//秒杀已经结束
		secKillStatus = 2
		remainSeconds = -1
	} else {
		//秒杀进行中
		secKillStatus = 1
		remainSeconds = 0
	}
	ctx.HTML(http.StatusOK, "goodsDetail.tmpl", gin.H{
		"user":          user,
		"goods":         goodVo,
		"secKillStatus": secKillStatus,
		"remainSeconds": remainSeconds})
}

func GoodDetailInfo(ctx *gin.Context) {
	//获取用户
	user, exists := ctx.Get("user")
	if !exists {
		//如果用户未登陆或登陆状态过期，则重定向到登陆界面
		ctx.Redirect(http.StatusTemporaryRedirect, "/home/tologin")
	}
	idStirng := ctx.Param("secKillId")
	secKillId, err := strconv.Atoi(idStirng)
	if err != nil {
		//参数出错
		panic(err)
	}
	goodVo := services.GetGood(secKillId)

	startDate := goodVo.StartDate
	endDate := goodVo.EndDate
	nowDate := time.Now()

	//秒杀状态
	secKillStatus := 0
	//剩余时间
	remainSeconds := 0

	if nowDate.Before(startDate) {
		//秒杀还未开始
		remainSeconds = startDate.Second() - nowDate.Second()
	} else if nowDate.After(endDate) {
		//秒杀已经结束
		secKillStatus = 2
		remainSeconds = -1
	} else {
		//秒杀进行中
		secKillStatus = 1
		remainSeconds = 0
	}
	//装填数据
	var detail vo.GoodDetailVo
	detail.User, _ = user.(models.User)
	detail.Good = *goodVo
	detail.SecKillStatus = secKillStatus
	detail.RemainSeconds = remainSeconds

	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{})
	dataMap["obj"] = detail
	controllers.Response(ctx, http.StatusOK, "查询成功", dataMap)
}

/**
优化前：134
优化后：400左右
*/
func DoSecKill(ctx *gin.Context) {
	//获取用户
	userInterFace, exists := ctx.Get("user")
	if !exists {
		//如果用户未登陆或登陆状态过期，则重定向到登陆界面
		ctx.Redirect(http.StatusTemporaryRedirect, "/home/tologin")
		return
	}
	//goodsId
	idStirng := ctx.PostForm("secKillId")
	secKillId, err := strconv.Atoi(idStirng)
	if err != nil {
		//参数出错
		panic(err)
	}

	//使用内存标记判断是否重复抢购，减少redis的访问
	over := services.GetEmptyStock(secKillId)
	if over {
		controllers.Response(ctx, 500500, "商品已售罄", nil)
		return
	}

	user, ok := userInterFace.(models.User)
	if !ok {
		fmt.Println("类型不匹配")
	}
	//设置分布式锁，防止玩家重复抢购，setnx
	//获取redis连接
	redisClinet := dbs.GetRedisClinet()
	background := context.Background()

	//decr := redisClinet.Decr(background, "secKillGood:" + strconv.Itoa(secKillId))
	//redisClinet.Incr(background, "secKillGood:" + strconv.Itoa(secKillId))
	//decr = redisClinet.Decr(background, "secKillGood:" + strconv.Itoa(secKillId))
	//redisClinet.Incr(background, "secKillGood:" + strconv.Itoa(secKillId))
	//return

	boolCmd := redisClinet.SetNX(background, "userId:"+strconv.Itoa(user.UserId)+"_"+"secKillGood:"+strconv.Itoa(secKillId), "repeat", time.Duration(5)*time.Second)
	result, err := boolCmd.Result()
	if !result {
		controllers.Response(ctx, 500501, "该商品每人限购一件", nil)
		return
	}

	//进行秒杀
	decr := redisClinet.Decr(background, "secKillGood:"+strconv.Itoa(secKillId))
	if decr.Val() < 0 {
		//库存补偿
		redisClinet.Incr(background, "secKillGood:"+strconv.Itoa(secKillId))
		//不要尝试再这直接将内存进行标记，锁竞争会导致redis连接失效
		services.SetEmptyStock(secKillId)
		controllers.Response(ctx, 500500, "库存不足", nil)
		return
	}

	secKillOrderCache := services.GetSecKillOrderCache(user.UserId, secKillId)
	if secKillOrderCache != nil {
		//库存补偿
		redisClinet.Incr(background, "secKillGood:"+strconv.Itoa(secKillId))
		controllers.Response(ctx, 500501, "已下单", nil)
		return
	}

	services.SecKill(&user, secKillId)
	//发送消息给消息队列，进行异步下单
	var massage common.SecKillMassage
	massage.User = &user
	massage.GoodId = secKillId
	marshal, err := jsoniter.Marshal(massage)
	if err != nil {
		panic(err)
	}
	utils.SendSecKillMessage(marshal)

	controllers.Response(ctx, http.StatusOK, "排队中", nil)
}

/**
获取秒杀地址（动态获取地址）
*/
func GetSecKillPath(ctx *gin.Context) {
	//
	//userInterFace, exists := ctx.Get("user")
	//if !exists {
	//	//如果用户未登陆或登陆状态过期，则重定向到登陆界面
	//	ctx.Redirect(http.StatusTemporaryRedirect, "/home/tologin")
	//	return
	//}
	//user, ok := userInterFace.(models.User)
	//if !ok {
	//	fmt.Println("类型不匹配")
	//}
	var md5Utile utils.MD5Util
	md5Utile.Seed = md5.New()
	//获取uuid
	u := uuid.New()
	_, err := md5Utile.Seed.Write([]byte(u.String() + "12345"))
	if err != nil {
		panic(err)
	}
	//存入redis
	//redisClinet := dbs.GetRedisClinet()
	//redisClinet.Set(ctx, "secKillPath_" + strconv.Itoa(user.UserId), fmt.Sprintf("%x",md5Utile.Seed.Sum(nil)), 0)
	//装填数据
	dataMap := make(map[string]interface{})
	dataMap["path"] = fmt.Sprintf("%x", md5Utile.Seed.Sum(nil))
	controllers.Response(ctx, 200, "生成成功", dataMap)
	fmt.Println("处理程序")
	return
}
