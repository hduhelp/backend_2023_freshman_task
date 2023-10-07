package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type TODO struct {
	Content  string `json:"content"`
	Done     bool   `json:"done"`
	Deadline string `json:"deadline"`
}

var todos []TODO
var blink []int
var bl TODO
var ptodos []TODO

// 计算时间
func time(a string) int {
	year, _ := strconv.Atoi(a[0:4])
	mouth, _ := strconv.Atoi(a[5:7])
	day, _ := strconv.Atoi(a[8:10])
	hour, _ := strconv.Atoi(a[11:13])
	minute, _ := strconv.Atoi(a[14:16])
	second, _ := strconv.Atoi(a[17:19])
	sum := year*31536000 + mouth*2592000 + day*86400 + hour*3600 + minute*60 + second
	return sum
}
func main() {
	r := gin.Default()
	bl.Deadline = "Nil"
	bl.Done = false

	// 添加 TODO
	r.POST("/todo", func(c *gin.Context) {
		var todo TODO
		c.BindJSON(&todo)
		if len(blink) > 0 {
			todos[blink[0]] = todo
			blink = blink[1:]
		} else {
			todos = append(todos, todo)
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 删除 TODO
	r.DELETE("/todo/:index", func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		if index <= len(todos)-1 && todos[index] != bl {
			todos[index] = bl
			blink = append(blink, index)
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 修改 TODO
	r.PUT("/todo/:index", func(c *gin.Context) {
		var todo TODO
		c.BindJSON(&todo)
		index, _ := strconv.Atoi(c.Param("index"))
		if index+1 > len(todos) {
			n := len(todos)
			for i := 0; i < index-n; i++ {
				todos = append(todos, bl)
				blink = append(blink, n+i)
			}
			todos = append(todos, todo)
		} else {
			if todos[index].Deadline == "Nil" {
				for j, v := range blink {
					if v == index {
						blink = append(blink[:j], blink[j+1:]...)
					}
				}
			}
		}
		todos[index] = todo
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 获取 TODO
	r.GET("/todo", func(c *gin.Context) {
		var ptodos []TODO
		for i, _ := range todos {
			if todos[i].Deadline != "Nil" {
				ptodos = append(ptodos, todos[i])
			}
		}

		c.JSON(200, ptodos)
	})

	// 查询TODO
	r.GET("/todo/:index", func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		c.JSON(200, todos[index])
	})

	//删除已完成TODO
	r.DELETE("/todo", func(c *gin.Context) {
		ptodos = nil
		blink = nil
		for i, _ := range todos {
			if todos[i].Done == false {
				ptodos = append(ptodos, todos[i])
			}
		}
		todos = ptodos
		c.JSON(200, gin.H{"status": "ok"})
	})
	//按照截止时间升序排序TODO
	r.PUT("/todo", func(c *gin.Context) {
		ptodos = nil
		blink = nil
		for i, _ := range todos {
			if todos[i].Deadline != "Nil" {
				ptodos = append(ptodos, todos[i])
			}
		}
		//冒泡排序
		for i := 0; i < len(ptodos)-1; i++ {
			for j := 0; j < len(ptodos)-i-1; j++ {
				if time(ptodos[j].Deadline) < time(ptodos[j+1].Deadline) {
					temp := ptodos[j]
					ptodos[j] = ptodos[j+1]
					ptodos[j+1] = temp
				}
			}
		}
		todos = ptodos
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Run(":8080")
}
