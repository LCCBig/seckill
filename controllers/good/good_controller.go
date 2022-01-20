package good

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"seckill/controllers"
	"seckill/models"
	"seckill/services"
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

func DoSecKill(ctx *gin.Context) {
	//获取用户
	userInterFace, exists := ctx.Get("user")
	if !exists {
		//如果用户未登陆或登陆状态过期，则重定向到登陆界面
		ctx.Redirect(http.StatusTemporaryRedirect, "/home/tologin")
	}
	//goodsId
	idStirng := ctx.Param("secKillId")
	secKillId, err := strconv.Atoi(idStirng)
	if err != nil {
		//参数出错
		print(err)
	}

	//判断库存是否充足
	secKillGood := services.GetSecKillGoodByid(secKillId)
	if secKillGood.StockCount < 1 {
		controllers.Response(ctx, 500500, "库存不足", nil)
	}

	//判断是否重复抢购
	user, ok := userInterFace.(models.User)
	if !ok {
		fmt.Println("类型不匹配")
	}
	secKillOrderCache := services.GetSecKillOrderCache(user.UserId, secKillId)
	if secKillOrderCache != nil {
		controllers.Response(ctx, 500501, "该商品每人限购一件", nil)
	}

	//进行秒杀

}
