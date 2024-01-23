package main

import (
	"epub-reader-web-server/account"
	"epub-reader-web-server/setting"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	// account.CreateUser("example", "example", "")

	exampleUser, err := account.GetUserById("example")
	if err != nil {
		fmt.Println(err)
		return
	}
	setting.Init()

	// books := []account.Book{}
	// filepath.Walk(setting.EpubsAbsPath, func(path string, info fs.FileInfo, err error) error {
	// 	if err != nil {
	// 		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
	// 		return err
	// 	}
	// 	if info.IsDir() && path != setting.EpubsAbsPath {
	// 		return filepath.SkipDir
	// 	}
	// 	fmt.Printf("visited file or dir: %q\n", path)
	// 	if filepath.Ext(path) == ".epub" {
	// 		book, err := exampleUser.UnzipAndGenerateEpubWebInfo(path)
	// 		if err == nil {
	// 			books = append(books, *book)
	// 		}
	// 	}

	// 	return nil
	// })
	// exampleUser.Books = books
	// exampleUser.Save()

	router := gin.Default()
	router.Static(setting.ConfigYaml.GinEpubsStaticPath, setting.UnzipAbsPath)
	router.GET("/api/bookshelf", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, exampleUser.Books)
		// c.String(http.StatusOK, fmt.Sprintf())
	})
	router.POST("/api/login", func(ctx *gin.Context) {
		id := ctx.PostForm("username")
		pw := ctx.PostForm("password")
		fmt.Println("login", id, pw)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}

		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(pw))
		fmt.Println("密码", err)
		if err != nil {
			data := gin.H{
				"message": "密码错误",
			}

			ctx.JSON(http.StatusOK, data)
		}
		fmt.Println("密码", string(hashedPassword))
		data := gin.H{
			"message": "",
		}
		ctx.JSON(http.StatusOK, data)
	})

	router.POST("/api/register", func(ctx *gin.Context) {
		// id := ctx.PostForm("username")~
		// pw := ctx.PostForm("password")
		// invitecode := ctx.PostForm("invitecode")

	})

	router.Run()

}
