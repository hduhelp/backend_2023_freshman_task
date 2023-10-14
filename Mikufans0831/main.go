package main

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var timeTemplate1 string = "2006-01-02 15:04:05"

type TODO struct {
	ID       int
	Content  string `json:"content"`
	Done     bool   `json:"done"`
	Deadline string `json:"deadline"`
	Priority int    `json:"priority"`
	Unix     int64
	Late     bool `json:"Late"`
}

// 返回一个支持至 秒 级别的 cron
func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

// 每日委托
func meiri(db *gorm.DB) {
	var todos []TODO
	db.Where("priority = ?", 3).Delete(&todos)
	db.Delete(&todos)

	tomUnix := time.Now().Unix() + 86400
	o := TODO{
		Content:  "原神委托。",
		Done:     false,
		Unix:     tomUnix,
		Priority: 3,
	}
	o.Deadline = time.Unix(tomUnix, 0).Format("2006-01-02 15:04:05")

	Bt := TODO{
		Content:  "崩铁活跃度。",
		Done:     false,
		Unix:     tomUnix,
		Priority: 3,
	}
	Bt.Deadline = time.Unix(tomUnix, 0).Format("2006-01-02 15:04:05")

	BA := TODO{
		Content:  "BA日常。",
		Done:     false,
		Unix:     tomUnix,
		Priority: 3,
	}
	BA.Deadline = time.Unix(tomUnix, 0).Format("2006-01-02 15:04:05")
	db.Create(&o)
	db.Create(&Bt)
	db.Create(&BA)
}

// 更新每日委托
func up(db *gorm.DB) {
	c := newWithSeconds()
	spec := "0 0 0 * * *"
	// spec := "5 */2 * * * ?"
	c.AddFunc(spec, func() {
		meiri(db)
	})
	c.Start()
}

// 逾期
func late(db *gorm.DB) {
	var todos []TODO
	c := newWithSeconds()
	spec := "0 0 * * * *"
	// spec := "5 */2 * * * ?"
	nowUnix := time.Now().Unix()
	db.Where("unix  <= ? AND late <> ? AND done = ?", nowUnix, true, false).Find(&todos)
	c.AddFunc(spec, func() {
		db.Model(TODO{}).Where("unix <= ?", nowUnix).Updates(TODO{Late: true})
	})
	c.Start()
	send1(db)
}

// 问候
func send1(db *gorm.DB) {
	var late []TODO
	var no []TODO
	var late1 string
	var no1 string
	db.Where("late = ?", true).Find(&late)
	db.Where("done = ?", false).Find(&no)
	for _, l := range late {
		late1 += l.Content
	}
	for _, n := range no {
		no1 += n.Content
	}
	if late1 != "" && no1 != "" {
		message := `
		    <p> Hello %s,</p>

			<p style="text-indent:2em">未完成事项：</p>
			<p style="text-indent:2em">` + no1 + `</p>

			<p style="text-indent:2em">逾期事项：</P>
			<p style="text-indent:2em">` + late1 + `</p>
			`
		host := "smtp.qq.com"
		port := 25
		userName := "1063243756@qq.com"
		password := "pdwbaqkhqshlbejg"

		m := gomail.NewMessage()
		m.SetHeader("From", userName) // 发件人
		m.SetHeader("To", "1063243756@qq.com")
		m.SetHeader("Subject", "TODOList") // 邮件主题
		m.SetBody("text/html", fmt.Sprintf(message, "testUser"))

		// text/plain的意思是将文件设置为纯文本的形式，浏览器在获取到这种文件时并不会对其进行处理

		d := gomail.NewDialer(
			host,
			port,
			userName,
			password,
		)
		// 关闭SSL协议认证
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
	}
}

func main() {
	r := gin.Default()
	dsn := "root:Qq563696767@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&TODO{})

	go up(db)
	go late(db)

	// 添加 TODO
	r.POST("/todo", func(c *gin.Context) {
		var todo TODO
		c.BindJSON(&todo)
		Un, _ := time.ParseInLocation(timeTemplate1, todo.Deadline, time.Local)
		todo.Unix = Un.Unix()
		db.Create(&todo)
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 删除 TODO
	r.DELETE("/todo/:index", func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		var todo TODO
		db.Where("ID = ?", index).Delete(&todo)
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 修改 TODO
	r.PUT("/todo/:index", func(c *gin.Context) {
		var todo TODO
		var new TODO
		c.BindJSON(&new)
		index, _ := strconv.Atoi(c.Param("index"))
		db.Where("ID = ?", index).First(&todo)
		Un, _ := time.ParseInLocation(timeTemplate1, new.Deadline, time.Local)
		new.Unix = Un.Unix()
		if !new.Done {
			db.Model(&todo).Select("done").Updates(map[string]interface{}{"done": false})
		}
		if !new.Late {
			db.Model(&todo).Select("late").Updates(map[string]interface{}{"late": false})
		}
		db.Model(&todo).Updates(new)

		// db.Where(xiShu{Name: "LiuBei"}).Assign(xiShu{Age: 35}).FirstOrInit(&todo)

		c.JSON(200, gin.H{"status": "ok"})

	})

	// 获取 TODO
	r.GET("/todo", func(c *gin.Context) {
		var todos []TODO
		db.Find(&todos)
		c.JSON(200, todos)
	})

	// 查询TODO
	r.GET("/todo/:index", func(c *gin.Context) {
		var todo TODO
		index, _ := strconv.Atoi(c.Param("index"))
		db.Where("id = ?", index).First(&todo)
		c.JSON(200, todo)

	})

	//删除已完成TODO
	r.DELETE("/todo", func(c *gin.Context) {
		var todos []TODO
		db.Where("done = ?", true).Delete(&todos)
		db.Delete(&todos)
		c.JSON(200, gin.H{"status": "ok"})
	})
	//按照优先级升序再按截止时间升序排序TODO
	r.PUT("/todo", func(c *gin.Context) {
		var todos []TODO
		db.Order("priority").Order("unix").Find(&todos)
		c.JSON(200, todos)
	})
	r.Run(":8080")
}
