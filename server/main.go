package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pirmd/epub"
	"github.com/spf13/viper"
)

type Book struct {
	Name      string `yaml:"name"`
	Path      string `yaml:"path"`
	CoverPath string `yaml:"cover_path"`
}

type User struct {
	name  string `yaml:"name"`
	id    string `yaml:"id"`
	pw    string `yaml:"pw"`
	books []Book `yaml:"books"`
}

var epubsPath string = ""
var epubsImagesPath string = ""
var epubExt = ".epub"

func main() {

	viper.SetConfigType("yaml")
	viper.SetConfigFile("./config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	epubsPath = viper.GetString("epubsPath")
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
	exampleConfig := viper.New()
	exampleConfig.SetConfigType("yaml")
	exampleConfig.SetConfigFile(examplePath)
	_, err = os.Stat(accFolder)
	if err != nil {
		err = os.Mkdir(accFolder, os.ModePerm)
		if err != nil {
			fmt.Println("创建 acc 文件夹失败")
		}
	}
	exampleConfig.ReadInConfig()
	var exampleUser User
	err = exampleConfig.Unmarshal(&exampleUser)
	if err != nil {
		fmt.Println("加载 example.yaml 失败", err)
	}
	if exampleUser.name == "" {
		exampleUser.name = "example"
		exampleUser.id = "0"
		exampleUser.pw = ""
	}

	exampleUser.books = []Book{}
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
			exampleUser.books = append(exampleUser.books, book)
		}
		return nil
	})
	// exampleConfig.marsha
	exampleConfig.WriteConfig()
	for index, book := range exampleUser.books {
		fmt.Println(index, book)
	}

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
						fmt.Println("文件名对比", file.Name, filepath.Base(file.Name), coverPath, filepath.Base(coverPath))
						if filepath.Base(file.Name) == filepath.Base(coverPath) {
							fmt.Println(file.Name)
							epubWithoutExt := strings.TrimSuffix(filepath.Base(path), filepath.Ext(filepath.Base(path)))
							coverPathInSys := filepath.Join(epubsImagesPath, (epubWithoutExt + filepath.Ext(file.Name)))
							fmt.Println("图片路径 ", epubWithoutExt, coverPathInSys)
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
