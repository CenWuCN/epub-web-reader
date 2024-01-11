package main

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pirmd/epub"
	"gopkg.in/yaml.v3"
)

type Book struct {
	Name      string `yaml:"name"`
	Path      string `yaml:"path"`
	CoverPath string `yaml:"cover_path"`
}

type User struct {
	Name  string `yaml:"name"`
	Id    string `yaml:"id"`
	Pw    string `yaml:"pw"`
	Books []Book `yaml:"books"`
}

type ConfigYaml struct {
	EpubsPath string `yaml:"epubsPath"`
}

var epubsPath string = ""
var epubsImagesPath string = ""
var epubExt = ".epub"

func main() {

	bytes, err := os.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println(err)
	}
	configyaml := ConfigYaml{}
	yaml.Unmarshal(bytes, &configyaml)
	epubsPath = configyaml.EpubsPath

	epubsPath, err = ExpandTilde(epubsPath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(epubsPath)
	epubsImagesPath = filepath.Join(epubsPath, "images")
	fmt.Println(epubsImagesPath)

	var exampleConfigName = "example.yaml"
	var accFolder = "./acc"
	var examplePath = filepath.Join(accFolder, exampleConfigName)

	_, err = os.Stat(accFolder)
	if os.IsNotExist(err) {
		err = os.Mkdir(accFolder, os.ModePerm)
		if err != nil {
			fmt.Println("创建 acc 文件夹失败")
		}
	}
	if err != nil {
		fmt.Println("加载 example.yaml 失败", err)
	}

	bytes, err = os.ReadFile(examplePath)
	if err != nil {
		fmt.Println(err)
	}
	exampleUser := User{}
	yaml.Unmarshal(bytes, &exampleUser)

	if exampleUser.Name == "" {
		exampleUser.Name = "example"
		exampleUser.Id = "0"
		exampleUser.Pw = ""
	}

	exampleUser.Books = []Book{}
	filepath.Walk(epubsPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		fmt.Printf("visited file or dir: %q\n", path)
		if filepath.Ext(path) == epubExt {
			book := GenerateEpubWebInfo(path)
			if book.Path != "" {
				fmt.Println(book)
			}
			exampleUser.Books = append(exampleUser.Books, book)
		}
		return nil
	})

	bytes, err = yaml.Marshal(&exampleUser)
	os.WriteFile(examplePath, bytes, os.ModePerm)
	for index, book := range exampleUser.Books {
		fmt.Println(index, book)
	}

	router := gin.Default()
	router.GET("/bookshelf", func(c *gin.Context) {
		c.JSON(http.StatusOK, exampleUser.Books)
	})

	router.Run()

	return
}

func ExpandTilde(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, path[1:]), nil
}

func GenerateEpubWebInfo(path string) Book {
	metaName := "cover"
	// metaContent := "cover-image"
	var book Book
	e, err := epub.Open(path)
	if err != nil {
		fmt.Println(err)
		return book
	}
	opf, err := e.Package()
	if err != nil {
		fmt.Println(err)
		return book
	}
	metaTagList := opf.Metadata.Meta
	itemList := opf.Manifest.Items

	book.Path = path
	if len(opf.Metadata.Title) > 0 {
		book.Name = opf.Metadata.Title[0].Value
	}

	for _, metaTag := range metaTagList {
		if metaTag.Name == metaName {

			for _, item := range itemList {
				coverPath := item.Href
				if item.ID == metaTag.Content {

					for _, file := range e.File {
						// fmt.Println("文件名对比", file.Name, filepath.Base(file.Name), coverPath, filepath.Base(coverPath))
						if filepath.Base(file.Name) == filepath.Base(coverPath) {
							fmt.Println(file.Name)
							epubWithoutExt := strings.TrimSuffix(filepath.Base(path), filepath.Ext(filepath.Base(path)))
							coverPathInSys := filepath.Join(epubsImagesPath, (epubWithoutExt + filepath.Ext(file.Name)))
							// fmt.Println("图片路径 ", epubWithoutExt, coverPathInSys)
							_, err = os.Stat(coverPathInSys)
							if os.IsNotExist(err) {
								destFile, err := os.Create(coverPathInSys)
								if err != nil {
									fmt.Println(err)
									return book
								}
								reader, err := file.Open()
								if err != nil {
									fmt.Println(err)
									return book
								}
								_, err = io.Copy(destFile, reader)
								if err != nil {
									fmt.Println("复制失败", err)
									return book
								}
							}

							book.CoverPath = coverPathInSys
							goto end
						}
					}
				}
			}
		}
	}

end:
	return book
}
