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
	setting.Init()

	// account.CreateUser("example", "example", "")

	// exampleUser, err := account.GetUserById("example")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
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
		exampleUser, err := account.GetUserById("example")
		if err != nil {
			return
		}
		ctx.JSON(http.StatusOK, exampleUser.Books)
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

	router.POST("/api/readingposset", func(ctx *gin.Context) {
		id := ctx.PostForm("id")
		percentage := ctx.PostForm("percentage")
		link := ctx.PostForm("link")
		fmt.Println(id, percentage, link)
		exampleUser, err := account.GetUserById("example")
		if err != nil {
			ctx.JSON(http.StatusOK, "")
		} else {
			fmt.Println("获取书数据")
			bookinfo, err := exampleUser.GetBookInfo(id)
			if err != nil {
				ctx.JSON(http.StatusOK, "")
				return
			}
			bookinfo.ReadingPos.Link = link
			bookinfo.ReadingPos.Percentage = percentage
			fmt.Println(bookinfo)
			exampleUser.Save()
		}
		ctx.JSON(http.StatusOK, "")
	})

	router.POST("/api/bookinfo", func(ctx *gin.Context) {
		id := ctx.PostForm("bookid")

		exampleUser, err := account.GetUserById("example")
		if err != nil {
			ctx.JSON(http.StatusOK, "{}")
		} else {
			books := exampleUser.Books
			for _, bookinfo := range books {
				if bookinfo.Id == id {
					ctx.JSON(http.StatusOK, bookinfo)
					return
				}
			}
			ctx.JSON(http.StatusOK, "{}")
		}

	})

	router.Run()

}
