package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID       uint `gorm:"primarykey;AUTO_INCREMENT"`
	Name     string
	Password string
}

func main() {
	//数据库
	dsn := "root:12345@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	db.AutoMigrate(&User{})
	// db.Create(User{Name: "haha", Password: "123"})
	//服务器
	server := gin.Default()
	server.Use(cors.Default())
	// 注册
	server.POST("/user/register", func(C *gin.Context) {
		u := User{}
		C.BindJSON(&u)
		// 查询用户名
		res := db.Where("name=?", u.Name).First(&User{})
		if res.RowsAffected != 0 {
			C.JSON(http.StatusOK, gin.H{
				"message": "注册失败，用户名已存在",
			})
		} else {
			db.Create(&u)
			C.JSON(http.StatusOK, gin.H{
				"msg": "注册成功",
			})
		}
	})

	server.Run(":8080")
}
