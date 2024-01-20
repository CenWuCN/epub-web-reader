package main

import (
	"epub-reader-web-server/account"
	"epub-reader-web-server/setting"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {

	account.CreateUser("example", "example", "")

	exampleUser, err := account.GetUserById("example")
	if err != nil {
		fmt.Println(err)
		return
	}
	setting.Init()
	books := []account.Book{}
	filepath.Walk(setting.EpubsAbsPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && path != setting.EpubsAbsPath {
			return filepath.SkipDir
		}
		fmt.Printf("visited file or dir: %q\n", path)
		book, err := exampleUser.UnzipAndGenerateEpubWebInfo(path)
		if err == nil {
			books := append(books, *book)
		}

		return nil
	})
	exampleUser.Books = books
	exampleUser.Save()

	router := gin.Default()
	router.Static(setting.ConfigYaml.GinEpubsStaticPath, setting.UnzipAbsPath)
	router.GET("/api/bookshelf", func(c *gin.Context) {
		c.JSON(http.StatusOK, exampleUser.Books)
		// c.String(http.StatusOK, fmt.Sprintf())
	})

	router.Run()

}
