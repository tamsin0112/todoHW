package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type TODO struct {
	Index  int    `gorm:"primaryKey" json:"index"`
	Task   string `json:"task"`
	Done   int    `json:"done"`
	Type   string `json:"type"`
	DueDay string `json:"dueDay"`
}

func main() {
	// 初始化 SQLite 資料庫連接
	var err error
	dsn := "root:ppassword@tcp(db:3307)/todo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自動遷移資料庫結構
	db.AutoMigrate(&TODO{})

	// 初始化範例任務
	db.Create(&TODO{
		Index:  1,
		Task:   "範例任務1",
		Done:   0,
		Type:   "工作",
		DueDay: time.Now().AddDate(0, 0, 1).Format("2006-01-02"), // 明天
	})
	db.Create(&TODO{
		Index:  2,
		Task:   "範例任務2",
		Done:   1,
		Type:   "學習",
		DueDay: time.Now().AddDate(0, 0, 3).Format("2006-01-02"), // 三天後
	})

	r := gin.Default()

	// API 路由
	r.GET("/todo/show/:filter", getTodos)
	r.POST("/todo/new", createTodo)
	r.POST("/todo/delete", deleteTodo)
	r.POST("/todo/update", updateTodo)

	r.Run(":8008")
}

func getTodos(c *gin.Context) {
	var todos []TODO
	filter := c.Param("filter")
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is nil"})
		return
	}
	switch filter {
	case "done":
		db.Where("done = ?", 2).Find(&todos)
	case "notdone":
		db.Where("done = ?", 0).Find(&todos)
	case "on-going":
		db.Where("done = ?", 1).Find(&todos)
	default:
		db.Find(&todos)
	}

	c.JSON(http.StatusOK, todos)
}

func createTodo(c *gin.Context) {
	task := c.PostForm("Task")
	doneStr := c.PostForm("Done")
	taskType := c.PostForm("Type")
	dueDay := c.PostForm("DueDay")

	done, err := strconv.Atoi(doneStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Done value, must be 'true' or 'false'"})
		return
	}

	newTODO := TODO{
		Task:   task,
		Done:   done,
		Type:   taskType,
		DueDay: dueDay,
	}
	result := db.Create(&newTODO)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "New task added successfully",
		"task":    newTODO,
	})

}

func deleteTodo(c *gin.Context) {
	// 從請求中獲取 Index
	strIndex := c.PostForm("Index")
	index, err := strconv.Atoi(strIndex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Index, must be an integer"})
		return
	}

	// 從資料庫中查找對應的任務
	var todo TODO
	result := db.First(&todo, index)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// 從資料庫中刪除任務
	result = db.Delete(&todo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
		"task":    todo,
	})
}

func updateTodo(c *gin.Context) {
	// 從請求中獲取 Index 和 Done
	strIndex := c.PostForm("Index")
	strDone := c.PostForm("Done")

	index, err := strconv.Atoi(strIndex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Index, must be an integer"})
		return
	}

	done, err := strconv.Atoi(strDone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Done value, must be an integer"})
		return
	}

	// 從資料庫中查找對應的任務
	var todo TODO
	result := db.First(&todo, index)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// 更新任務的 Done 狀態
	todo.Done = done
	result = db.Save(&todo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}
