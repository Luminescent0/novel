package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"novel/model"
	"novel/service"
	"novel/tool"
	"strconv"
)

const rtLen = 64

func login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	err := service.UsernameIsExist(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			tool.RespErrorWithData(ctx, "用户不存在")
		}
		tool.RespInternalError(ctx)
		return
	}
	flag, err := service.IsPasswordCorrect(username, password)
	if err != nil {
		fmt.Println("judge password correct err:", err)
		tool.RespInternalError(ctx)
		return
	}
	if !flag {
		tool.RespErrorWithData(ctx, "密码错误")
		return
	}
	token, retoken, err1 := CreateToken(username)
	if err1 != nil {
		tool.RespInternalError(ctx)
		return
	}
	tool.RespSuccessfulLogin(ctx, retoken, token)
}

func register(ctx *gin.Context) {
	username, password, err := verify(ctx)
	if err != nil {
		if username == "存在非法输入" {
			tool.RespErrorWithData(ctx, "用户名格式有误")
			return
		}
		if password == "" {
			tool.RespErrorWithData(ctx, "密码格式有误")
			return
		}
	}
	user := model.User{
		Username: username,
		Password: password,
	}
	flag, err := service.IsRepeatUsername(username)
	if err != nil {
		fmt.Println("judge repeat username err:", err)
		tool.RespInternalError(ctx)
		return
	}
	if flag {
		tool.RespErrorWithData(ctx, "用户名已存在")
		return
	}
	err = service.Register(user)
	if err != nil {
		fmt.Println("register err:", err)
		tool.RespInternalError(ctx)
		return
	}
	tool.RespSuccessfulWithData(ctx, "注册成功")
}

func verify(ctx *gin.Context) (string, string, error) {
	validate := validator.New() //创建验证器
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	u := model.User{Id: 0, Username: username, Password: password}

	err := validate.Struct(u)
	fmt.Println(err)
	if err != nil {
		return "存在非法输入", "", err
	}
	return username, password, nil
}

func RefreshToken(ctx *gin.Context) {
	rt, exist := ctx.GetPostForm("refresh_token")
	if !exist {
		tool.RespErrorWithData(ctx, "refresh_token不存在")
		return
	}
	if len(rt) != rtLen {
		tool.RespSuccessfulWithData(ctx, "refresh_token长度不符")
	}
	atoken, retoken, err := refreshToken(rt)
	if err != nil {
		log.Println("refresh token failed,err:", err)
		tool.RespInternalError(ctx)
		return
	}
	tool.RespSuccessfulLogin(ctx, retoken, atoken)
}

func changePassword(ctx *gin.Context) {
	oldPassword := ctx.PostForm("oldPassword")
	newPassword := ctx.PostForm("newPassword")
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string) //接口断言
	fmt.Println(username)
	//检验旧密码
	flag, err := service.IsPasswordCorrect(username, oldPassword)
	if err != nil {
		fmt.Println("judge password correct err:", err)
		tool.RespInternalError(ctx)
		return
	}
	if !flag {
		tool.RespErrorWithData(ctx, "旧密码有误")
		return
	}
	err = service.ChangePassword(username, newPassword)
	if err != nil {
		fmt.Println("change password err:", err)
		tool.RespErrorWithData(ctx, "修改失败")
		return
	}
	tool.RespSuccessfulWithData(ctx, "修改成功")
}

func liked(ctx *gin.Context) {
	sBookId := ctx.PostForm("bookId")
	bookId, _ := strconv.Atoi(sBookId)
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)
	user, err := service.SelectUserByUsername(username)
	if err != nil {
		tool.RespInternalError(ctx)
		return
	}
	userId := user.Id
	sUserId := strconv.Itoa(userId)

	service.AcquireLock(bookId, userId) //分布式锁
	defer service.ReleaseLock(bookId, userId)

	allow, err := service.CheckLikeRateLimit(bookId, userId)
	if err != nil {
		tool.RespInternalError(ctx)
		return
	}
	if !allow {
		tool.RespErrorWithData(ctx, "频繁操作")
	}

	flag := service.IsMemberInSet(sBookId, sUserId)
	if !flag { //Redis里没有 在MySQL里找
		flag, err = service.SelectLiked(bookId, userId)
		if err != nil {
			fmt.Println("select liked err:", err)
			tool.RespInternalError(ctx)
			return
		}
		if flag { //MySQL里有的话就把赞取消
			err = service.CancelLiked(bookId, userId)
			if err != nil {
				fmt.Println("MySQL delete liked failed:", err)
				tool.RespInternalError(ctx)
				return
			}
			tool.RespSuccessfulWithData(ctx, "取消点赞成功")
			return
		}
		service.SetAdd(sBookId, sUserId, 0)
		err = service.Liked(bookId, userId) //MySQL里也没有就点赞
		if err != nil {
			fmt.Println("MySQL liked failed:", err)
			tool.RespInternalError(ctx)
			return
		}
		tool.RespSuccessfulWithData(ctx, "点赞成功")
	}
	service.SetMemberDel(sBookId, sUserId) //redis里有的话取消点赞
	err = service.CancelLiked(bookId, userId)
	if err != nil {
		fmt.Println("MySQL delete liked failed:", err)
		tool.RespInternalError(ctx) //这里可以改成假返回
		return
	}
	tool.RespSuccessfulWithData(ctx, "取消点赞成功")
	return

}
