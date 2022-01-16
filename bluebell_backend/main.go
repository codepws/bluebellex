// 包注释
// @Title  请填写文件名称（需要改）
// @Description  请填写文件描述（需要改）
// @Author  请填写自己的真是姓名（需要改）  ${DATE} ${TIME}
// @Update  请填写自己的真是姓名（需要改）  ${DATE} ${TIME}
package main

import (
	"bluebell_backend/dao/redis"
	"bluebell_backend/dao/sqldb"
	"bluebell_backend/pkg/logger"
	"bluebell_backend/routers"
	"bluebell_backend/settings"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func reverseList(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}

	cur := head
	for head.Next != nil {
		t := head.Next.Next
		head.Next.Next = cur // 反转原指针方向
		cur = head.Next      // 将新头节点移到下一位
		head.Next = t        // 连接回断开的地方，继续重复上面操作
	}

	return cur
}

// 函数（方法）注释
// @title    函数名称
// @description   函数的详细描述
// @auth      作者             时间（2019/6/18   10:57 ）
// @param     输入参数名        参数类型         "解释"
// @return    返回参数名        参数类型         "解释"
func main() {
	fmt.Println("===================")
	var head *ListNode
	tmp := head
	for i := 0; i < 5; i++ {
		tmp = ListNode{i, }
	}

	return

	// 1.
	if err := SetupConfig(); err != nil {
		logger.MainLogger.Fatal(fmt.Sprintf("setup config failed, err:%v\n", err))
		return
	}
	defer RevokeConfig()

	logger.MainLogger.Info("Web初始化成功！")

	// 2. 注册路由
	route := routers.SetupRouter()

	// 3.注册自定义验证(自定义结构体校验)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("timing", timing)
		if err != nil {
			fmt.Println("success")
		}
	}
	route.GET("/time", getTime) //测试用例

	err := route.Run(fmt.Sprintf(":%d", settings.WebConf.Port))
	if err != nil {
		logger.MainLogger.Fatal(fmt.Sprintf("run server failed, err:%v\n", err))
		return
	}

	logger.MainLogger.Info(fmt.Sprintf("Web服务已启动，监听端口：%d", settings.WebConf.Port))

	logger.MainLogger.Error("Web服务退出！")
}

// 自定义验证规则断言
func timing(fl validator.FieldLevel) bool {

	fmt.Printf("FieldName=%v, StructFieldName=%v, Param=%v, GetTag=%v\n", fl.FieldName(), fl.StructFieldName(), fl.Param(), fl.GetTag())

	if date, ok := fl.Field().Interface().(time.Time); ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

//form:"create_time"
type Info struct {
	CreateTime time.Time `form:"create_time" binding:"required,timing" time_format:"2006-01-02"`
	UpdateTime time.Time `form:"update_time" binding:"required,timing" time_format:"2006-01-02"`
}

func getTime(c *gin.Context) {
	var b Info
	// 数据模型绑定查询字符串验证
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "time are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func RevokeConfig() {
	defer logger.Close()
	defer sqldb.Close() // 程序退出关闭数据库连接
	defer redis.Close()
}

func SetupConfig() error {
	//Step 1: 加载Web配置信息
	if err := settings.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return err
	}

	fmt.Println(settings.WebConf)
	fmt.Println(settings.WebConf.LogConfig)
	fmt.Println(settings.WebConf.DBs)
	fmt.Println(settings.WebConf.Caches)
	fmt.Println(*(settings.WebConf.Caches[0]))

	//Step 2: 加载日志配置信息
	if err := logger.Init(settings.WebConf.LogConfig, settings.WebConf.Level); err != nil {
		fmt.Printf("load log config failed, err:%v\n", err)
		return err
	}

	//Step 3: 加载数据库Mysql配置信息
	if err := sqldb.Init(settings.WebConf.DBs); err != nil {
		fmt.Printf("load db config failed, err:%v\n", err)
		return err
	}

	logger.MainLogger.Debug("Login服务数据库配置信息", zap.String("host", settings.WebConf.LoginDB.Host), zap.Int("port", settings.WebConf.LoginDB.Port), zap.String("database", settings.WebConf.LoginDB.Database))
	logger.MainLogger.Debug("Shop服务数据库配置信息", zap.String("host", settings.WebConf.ShopDB.Host), zap.Int("port", settings.WebConf.ShopDB.Port), zap.String("database", settings.WebConf.ShopDB.Database))

	//Step 4: 加载缓存Redis配置信息
	if err := redis.Init(settings.WebConf.Caches); err != nil {
		fmt.Printf("load redis config failed, err:%v\n", err)
		return err
	}

	logger.MainLogger.Debug("Redis服务缓存配置信息", zap.Int("Caches", len(settings.WebConf.Caches)))

	return nil
}

// 结构（接口）注释
// User   用户对象，定义了用户的基础信息
type User struct {
	Username string // 用户名
	Email    string // 邮箱
}
