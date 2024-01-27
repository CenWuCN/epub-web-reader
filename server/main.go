package main

import (
	"epub-reader-web-server/account"
	"epub-reader-web-server/setting"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	router.POST("/api/login", func(ctx *gin.Context) {
		id := ctx.PostForm("username")
		pw := ctx.PostForm("password")
		fmt.Println("login", id, pw)
		// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		user, err := account.GetUserById(id)
		if err != nil {
			data := gin.H{
				"errorcode": 1,
			}

			ctx.JSON(http.StatusOK, data)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Pw), []byte(pw))
		fmt.Println("密码", err)
		if err != nil {
			data := gin.H{
				"errorcode": 2,
			}

			ctx.JSON(http.StatusOK, data)
			return
		}
		// fmt.Println("密码", string(hashedPassword))

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = id
		claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

		tokenString, err := token.SignedString([]byte(setting.ConfigYaml.Jwtkey))
		if err != nil {
			fmt.Println(err)
			return
		}

		data := gin.H{
			"errorcode": 0,
			"token":     tokenString,
		}
		ctx.JSON(http.StatusOK, data)
	})

	router.POST("/api/register", func(ctx *gin.Context) {
		id := ctx.PostForm("id")
		pw := ctx.PostForm("password")
		invitecode := ctx.PostForm("invitecode")

		if account.IsUserExist(id) {
			data := gin.H{
				"errorcode": 1,
			}
			ctx.JSON(http.StatusOK, data)
			return
		}

		found := false
		for _, code := range setting.ConfigYaml.Invitecodes {
			if invitecode == code {
				found = true
				break
			}
		}

		if !found {
			data := gin.H{
				"errorcode": 4,
			}
			ctx.JSON(http.StatusOK, data)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			data := gin.H{
				"errorcode": 3,
			}
			ctx.JSON(http.StatusOK, data)
			return
		}
		account.CreateUser(id, "", string(hashedPassword))
		data := gin.H{
			"errorcode": 0,
		}
		ctx.JSON(http.StatusOK, data)
	})

	authorized := router.Group("/")
	authorized.Use(AuthMiddleware())

	authorized.POST("/api/bookshelf", func(ctx *gin.Context) {
		user, ok := GetUser(ctx)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"errorcode": 8})
			return
		}
		ctx.JSON(http.StatusOK, user.Books)
	})

	authorized.POST("/api/readingposset", func(ctx *gin.Context) {
		bookid := ctx.PostForm("id")
		percentage := ctx.PostForm("percentage")
		link := ctx.PostForm("link")

		fmt.Println(bookid, percentage, link)
		// exampleUser, err := account.GetUserById("example")
		user, ok := GetUser(ctx)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"errorcode": 8})
			return
		}
		bookinfo, err := user.GetBookInfo(bookid)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"errorcode": 9})
			return
		}
		bookinfo.ReadingPos.Link = link
		bookinfo.ReadingPos.Percentage = percentage
		fmt.Println(bookinfo)
		user.Save()
	})

	authorized.POST("/api/bookinfo", func(ctx *gin.Context) {
		id := ctx.PostForm("bookid")

		user, ok := GetUser(ctx)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"errorcode": 8})
			return
		}
		books := user.Books
		for _, bookinfo := range books {
			if bookinfo.Id == id {
				ctx.JSON(http.StatusOK, bookinfo)
				return
			}
		}
		ctx.JSON(http.StatusOK, "{}")

	})

	authorized.POST("/api/uploadbooks", func(ctx *gin.Context) {
		// Multipart form
		form, _ := ctx.MultipartForm()
		files := form.File["file"]

		user, ok := GetUser(ctx)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"errorcode": 8})
			return
		}

		dupFiles := ""
		for _, file := range files {
			log.Println(file.Filename)

			dest := filepath.Join(setting.EpubsAbsPath, file.Filename)
			_, err := os.Stat(dest)
			if os.IsNotExist(err) {
				// Upload the file to specific dst.
				ctx.SaveUploadedFile(file, dest)
				user.AddBook(dest)
			} else {
				dupFiles = dupFiles + file.Filename
			}
		}
		if dupFiles != "" {
			data := gin.H{
				"errorcode": 6,
				"message":   dupFiles,
			}
			ctx.JSON(http.StatusOK, data)
		}
	})

	router.Run()

}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		fmt.Println("路由中间件")
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"errorcode": 7})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return setting.JwtkeyBytes, nil
		})

		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"errorcode": 5})
			ctx.Abort()
			return
		}

		if !token.Valid {
			data := gin.H{
				"errorcode": 5,
			}
			ctx.JSON(http.StatusOK, data)
			ctx.Abort()
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		ctx.Set("id", claims["id"])
		ctx.Next()
	}
}

func GetUser(ctx *gin.Context) (*account.User, bool) {
	id, exist := ctx.Get("id")
	if !exist {
		return nil, false
	}
	fmt.Println("获取书数据")
	stringId, exist := id.(string)
	if !exist {
		return nil, false
	}
	user, err := account.GetUserById(stringId)
	if err != nil {
		return nil, false
	}
	return user, true
}
