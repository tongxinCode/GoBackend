package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type (
	// 定义原始的数据库字段
	userModel struct {
		gorm.Model
		Username  string `json:"username"`
		Password  string `json:"password"`
		Nickname  string `json:"nickname"`
		Apartment string `json:"apartment"`
	}
	// 处理返回的字段
	transformedUser struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		Nickname  string `json:"nickname"`
		Apartment string `json:"apartment"`
	}
)

// 初始化数据库
func init() {
	//open a db connection
	var err error
	db, err = gorm.Open("mysql", "root:veausmysql@/gobackend?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	// 自动创建表
	db.AutoMigrate(&userModel{})
}

// 设置路由
func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := r.Group("/api/v1/users")
	{
		v1.POST("/", createUser)
		v1.GET("/", fetchAllUser)
		v1.GET("/:id", fetchSingleUser)
		v1.PUT("/:id", updateUser)
		v1.DELETE("/:id", deleteUser)
	}

	return r
}

// 创建User
// v1.POST("/", createUser)
func createUser(c *gin.Context) {
	user := userModel{
		Username:  c.PostForm("username"),
		Password:  c.PostForm("password"),
		Nickname:  c.PostForm("nickname"),
		Apartment: c.PostForm("apartment"),
	}
	db.Save(&user)
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"message":    "User created successfully!",
		"resourceId": user.ID,
	})
}

// 获取User
// v1.GET("/", fetchAllUser)
func fetchAllUser(c *gin.Context) {
	var users []userModel
	var _users []transformedUser
	db.Find(&users)
	if len(users) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Not found!",
		})
		return
	}
	for _, item := range users {
		_users = append(_users, transformedUser{
			ID:        item.ID,
			Username:  item.Username,
			Nickname:  item.Nickname,
			Apartment: item.Apartment,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _users})
}

// 获取单个User
// v1.GET("/:id", fetchSingleUser)
func fetchSingleUser(c *gin.Context) {
	var user userModel
	userID := c.Param("id")
	db.First(&user, userID)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not found!"})
		return
	}
	_user := transformedUser{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Apartment: user.Apartment,
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _user})
}

// 更新User
// v1.PUT("/:id", updateUser)
func updateUser(c *gin.Context) {
	var user userModel
	userID := c.Param("id")
	db.First(&user, userID)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not found!"})
		return
	}
	//定义map类型，key为字符串，value为interface{}类型，方便保存任意值
	data := make(map[string]interface{})
	data["Username"] = c.PostForm("username")
	data["Nickname"] = c.PostForm("nickname")
	data["Password"] = c.PostForm("password")
	data["Apartment"] = c.PostForm("apartment")
	db.Model(&user).Updates(data)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "User updated successfully!"})
}

// 删除User
// v1.DELETE("/:id", deleteUser)
func deleteUser(c *gin.Context) {
	var user userModel
	userID := c.Param("id")
	db.First(&user, userID)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not found!"})
		return
	}
	db.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "User deleted successfully!"})
}

func main() {
	r := setupRouter()
	r.Run(":3000")
}
