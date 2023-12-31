package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gopkg.in/gomail.v2"
	"net/http"
	"strconv"
)

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type TODO struct {
	Id      int    `json:"Id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}
type UserInformation struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	PassWord       string `json:"password"`
	PhoneNumber    string `json:"phone_number"`
	Email          string `json:"email"`
	IdentityNumber string `json:"identity_number"`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 切片
var todos []TODO
var users []UserInformation
var todo TODO
var user UserInformation

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	r := gin.Default() //创建一个默认路由

	//------------------------------------------------------------------------------------
	//连接mysql用户信息数据库
	var conn *gorm.DB
	conn, err := gorm.Open("mysql", "root:xyf1029@tcp(127.0.0.1:3306)/userinformation")
	if err != nil { //连接失败
		fmt.Println("gorm.Open err", err)
		return
	}
	conn.DB().SetMaxIdleConns(20)  //初始连接数
	conn.DB().SetMaxIdleConns(200) //最大连接数

	defer conn.Close()                                        //关闭数据库
	conn.SingularTable(true)                                  //以单数表的形式
	fmt.Println(conn.AutoMigrate(new(UserInformation)).Error) //使用AutoMigrate ()方法来实现数据库表迁移，可以自动增加表中没有的字段和索引，在Gin main.go函数中使用非常方便，不用手动运行迁移了

	//连接todolist数据库
	var dconn *gorm.DB
	dconn, derr := gorm.Open("mysql", "root:xyf1029@tcp(127.0.0.1:3306)/todolist")
	if derr != nil { //连接失败
		fmt.Println("gorm.Open err", derr)
		return
	}

	dconn.DB().SetMaxIdleConns(20)  //初始连接数
	dconn.DB().SetMaxIdleConns(200) //最大连接数

	defer dconn.Close() //关闭数据库
	dconn.SingularTable(true)
	fmt.Println(dconn.AutoMigrate(new(TODO)).Error)

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//用户注册
	r.POST("/Register", func(c *gin.Context) { //绑定路由规则和函数，访问index的路由，将有对应的函数去处理
		var userinformation UserInformation
		c.BindJSON(&userinformation)
		c.JSON(200, gin.H{"状态": "ok", "已成功注册用户为": userinformation})

		//接下来是数据库的操作
		res := conn.Create(&userinformation) //保存数据到user_information
		fmt.Println(res.Error)               //查看是否有错误
		fmt.Println(res.RowsAffected)        //查看影响的行数

	})
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//	//////////////////////////////////////////////////////////////////////////////////////////////////////////

	//用户登录
	r.POST("/login", func(c *gin.Context) {

		username := c.PostForm("username")
		password := c.PostForm("password")
		conn.Select("pass_word").Where("name=?", username).First(&user)
		// 检查用户名和密码是否匹配
		if user.PassWord == password {
			// 生成访问令牌（可以使用JWT等方式）
			accessToken := "your_access_token"
			c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "登录失败，请检查用户名或者密码"})
		}
	})
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//新增TODO
	r.POST("/todo", func(c *gin.Context) {
		c.BindJSON(&todo) //添加TODO，接受前端传来的json数据
		c.JSON(200, gin.H{"状态": "ok", "已成功添加待办事项为": todo})
		//接下来是数据库的操作
		res := dconn.Create(&todo)    //保存数据到todolist
		fmt.Println(res.Error)        //查看是否有错误
		fmt.Println(res.RowsAffected) //查看影响的行数
	})
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//删除TODO
	r.DELETE("/todo/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(200, gin.H{"状态": "no.请检查输入"})
		} else {
			shanchude := dconn.Where("id=?", id).First(&todo)
			todo.Id = id
			dconn.Delete(&todo)
			c.JSON(200, gin.H{"状态": "ok", "已成功删除待办事项": shanchude})
		}

	})
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//修改TODO
	r.PUT("/todo/:id", func(c *gin.Context) { //地址+回调函数
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(200, gin.H{"状态": "no.请检查输入"})
		} else {
			c.BindJSON(&todo)
			todo.Id = id
			xiugaide := dconn.Where("id=?", id).First(&todos)
			dconn.Save(&todo)
			c.JSON(200, gin.H{"状态": "ok", "已成功去掉待办事项": xiugaide, "新的待办事项列表为": todo})
		}

	})
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//获取TODO
	r.GET("/todo", func(c *gin.Context) {

		c.JSON(200, gin.H{"状态": "ok", "当前待办事项列表为": dconn.Find(&todos)})

	})
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//查询todo
	r.GET("/todo/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"状态": "error", "消息": "无效的待办事项索引"})
			return
		}
		if dconn.Where("id=?", id).First(&todo).RecordNotFound() {
			c.JSON(404, gin.H{"状态": "error", "消息": "未找到待办事项"})
			return
		}

		c.JSON(200, gin.H{"状态": "ok", "查询待办事项为": todo})
	})
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//查询某个todo完成情况
	r.GET("/hello/:content", func(c *gin.Context) {
		content := c.Param("content")
		c.JSON(200, gin.H{"状态": "ok", "查询待办事项为": dconn.Where("content=?", content).First(&todos)})
	})
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//发送邮件
	r.GET("/hello", func(c *gin.Context) {
		//待办事项查询
		dconn.Where("done=?", 0).Find(&todo)

		// 创建一个新的消息
		m := gomail.NewMessage()

		// 设置发件人
		m.SetHeader("From", "xieyifeng@aliyun.com")

		// 设置收件人
		m.SetHeader("To", "2915901086@qq.com")

		// 设置主题
		m.SetHeader("Subject", "待办事项")

		// 设置邮件正文
		m.SetBody("text/html", todo.Content)

		// 设置 SMTP 服务器信息
		d := gomail.NewDialer("smtp.aliyun.com", 587, "xieyifeng@aliyun.com", "你猜")

		// 通过Dialer发送邮件
		if err := d.DialAndSend(m); err != nil {
			panic(err)
		} else {
			c.JSON(200, gin.H{"状态": "ok，发送成功"})
		}

	})
	r.Run(":8080") //运行

}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
