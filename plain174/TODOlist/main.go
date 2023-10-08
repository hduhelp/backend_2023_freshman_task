package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
	"strings"
)

type todo struct {
	Name   string `json:"name"`
	Season int    `json:"season"`
	Done   bool   `json:"done"`
}

var todos []todo
var m = "H:\\program\\goworkfile\\TODOlist.json"
var m1 = "H:\\program\\goworkfile\\"

func init() {
	Save, _ := os.OpenFile(m, os.O_CREATE|os.O_RDWR, os.ModePerm)
	var content []byte
	buf := make([]byte, 1024*8, 1024*8)
	for {
		n, err := Save.Read(buf[:])
		if err == io.EOF {
			// 读取结束
			break
		}
		if err != nil {
			fmt.Println("read file err ", err)
			return
		}
		content = append(content, buf[:n]...)
		json.Unmarshal(content, &todos)
	}
	Save.Close()
	fmt.Println(todos)
}

func main() {

	r := gin.Default()
	//添加 TODO
	r.POST("/todo", func(c *gin.Context) {
		var todo1 todo
		c.BindJSON(&todo1)
		todos = append(todos, todo1)
		c.JSON(200, gin.H{"condition": "ok"})
	})

	//上传 excel
	r.GET("/todoex/:file", func(c *gin.Context) {
		file := c.Param("file")
		Save, _ := os.OpenFile(m1+file, os.O_CREATE|os.O_RDWR, os.ModePerm)
		reader := csv.NewReader(bufio.NewReader(Save))
		for {
			line, err := reader.Read()
			if err != nil {
				c.String(400, err.Error())
				return
			}
			x := strings.Split(line[0], " ")
			y := strings.Split(line[1], " ")
			var todo1 todo
			for x1 := 0; x1 < len(x); x1++ {
				todo1.Name = x[x1]
				todo1.Season, _ = strconv.Atoi(y[x1])
				todo1.Done = true
				todos = append(todos, todo1)
			}

		}

		Save.Close()
		c.JSON(200, gin.H{"condition": "ok"})
	})

	//删除 TODO
	r.DELETE("/todo/:index", func(c *gin.Context) {
		//index:=
		index, _ := strconv.Atoi(c.Param("index"))

		if index == -1 {
			todos = nil
		} else if index == -2 {
			todos = todos[:len(todos)-1]
		} else {
			todos = append(todos[:index], todos[index+1:]...)
		}
		c.JSON(200, gin.H{"condition": "ok"})
	})

	//修改 TODO
	r.PUT("/todo/:index/*season", func(c *gin.Context) {
		//index:=
		index, _ := strconv.Atoi(c.Param("index"))
		season, _ := strconv.Atoi(c.Param("season")[1:])
		var todo1 todo
		todo1 = todos[index]
		println(season)
		println(todo1.Season)
		if season >= todo1.Season {
			todo1.Season = season
		}
		todo1.Done = true
		todos[index] = todo1
		c.JSON(200, gin.H{"condition": "ok"})
	})

	//获取 TODO
	r.GET("/todo", func(c *gin.Context) {
		c.JSON(200, todos)
	})
	//查询 TODO
	r.GET("/todo/:status", func(c *gin.Context) {
		status := c.Param("status")
		if status == "1" {
			var todos1 []todo
			for x := 0; x < len(todos); x++ {
				if todos[x].Done == true {
					todos1 = append(todos1, todos[x])
				}
			}
			c.JSON(200, todos1)
		} else if status == "2" {
			var todos1 []todo
			for x := 0; x < len(todos); x++ {
				if todos[x].Done != true {
					todos1 = append(todos1, todos[x])
				}
			}
			c.JSON(200, todos1)
		}
		if status == "0" {
			c.JSON(200, todos)
		}

	})
	//保存 TODO
	r.GET("/todo/save", func(c *gin.Context) {
		Save, _ := os.OpenFile(m, os.O_RDWR|os.O_CREATE, os.ModePerm)
		jsonBytes, _ := json.Marshal(todos)
		Save.Write(jsonBytes)
		Save.Close()
		c.JSON(200, gin.H{"condition": "ok"})
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
