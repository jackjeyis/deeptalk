package main

import (
	"encoding/json"
	"log"
	"net/http"
    //"crypto/md5"
    "io"
    //"encoding/hex"
    "path"
    "os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	Id       int
	Phone    int    `json:"phone"`
	Pwd      string `json:"password"`
	Email    string `json:"email"`
	UserName string `json:"username" gorm:"column:name"`
}

var (
	db  *gorm.DB
	err error
)

func main() {
	db, err = gorm.Open("mysql", "soultalk:soultalk@tcp(127.0.0.1:3306)/soultalk?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Printf("gorm.Open err %v", err)
		return
	}
	defer db.Close()
	r := gin.Default()
	r.POST("soultalk/api/user/register", Register)
	r.POST("soultalk/api/user/login", Login)
    r.POST("soultalk/api/user/find",Find)
    r.POST("soultalk/api/user/upload",Upload)
    r.GET("soultalk/api/user/image/:hash",GetImage)
	r.Run(":9999")
}

func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err == nil {
		d := db.Table("user").Where(&user).Find(&user)
		if d.Error == nil {
			body, _ := json.Marshal(user)
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": string(body)})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "msg": d.Error.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad request", "err": err.Error()})
	}
}

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err == nil {
		d := db.Table("user").Create(&user)
		if d.Error == nil {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "user register success"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 2, "msg": d.Error.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad request", "err": err.Error()})
	}
}

func Find(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err == nil {
		d := db.Table("user").Where(&user).Find(&user)
		if d.Error == nil {
			body, _ := json.Marshal(user)
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": string(body)})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "msg": d.Error.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad request", "err": err.Error()})
	}
}

func Upload(c *gin.Context) {
    file, header , err := c.Request.FormFile("file")
    filename := header.Filename
    //fmt.Println(header.Filename)
    //ext := path.Ext(header.Filename)
    /*hash := md5.New()
    md5String := hex.EncodeToString(hash.Sum([]byte("")))
    filename = md5String + ext
    io.Copy(hash,file)
    */
    out, err := os.Create("/tmp/image/"+filename)
    if err != nil {
        log.Fatal(err)
    }
    defer out.Close()
    _, err = io.Copy(out, file)
    if err != nil {
        log.Fatal(err)
    }
    c.JSON(http.StatusOK,gin.H{"code":0,"msg":"upload image success","hash":filename})
}

func GetImage(c *gin.Context){
    filename := c.Param("hash")
    log.Printf("filename %v",filename)
    imagePath := path.Join("/tmp","/image/",filename)
    log.Printf("image path %v",imagePath)
    http.ServeFile(c.Writer,c.Request,imagePath)
}
