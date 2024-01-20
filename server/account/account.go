package account

import (
	"epub-reader-web-server/setting"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pirmd/epub"
	"gopkg.in/yaml.v3"
)

type Book struct {
	Name      string `yaml:"name"`
	Path      string `yaml:"path"`
	CoverPath string `yaml:"cover_path"`
	Opf       string `yaml:"opf"`
}

type User struct {
	Id    string `yaml:"id"`
	Name  string `yaml:"name"`
	Pw    string `yaml:"pw"`
	Books []Book `yaml:"books"`
}

var AccFolder string = "./acc"
var YamlExt string = ".yaml"
var CoverTagName string = "cover"
var EpubExt = ".epub"
var OpfExt = ".opf"

func CreateUser(id string, name string, pw string) {
	u := User{id, name, pw, []Book{}}
	u.Save()
}

func GetUserById(id string) (*User, error) {
	user := User{}
	accFilePath := filepath.Join(AccFolder, id+YamlExt)
	bytes, err := os.ReadFile(accFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("user is not exist")
		}
		return &user, err
	}

	yaml.Unmarshal(bytes, &user)
	return &user, err
}

func (u *User) AddBook(bookAbsName string) {
	if filepath.Ext(bookAbsName) == EpubExt {
		book, err := u.UnzipAndGenerateEpubWebInfo(bookAbsName)
		if err != nil {
			fmt.Println(err)
		}
		u.Books = append(u.Books, *book)
		u.Save()
	} else {
		fmt.Println("It is not a epub file")
	}
}

func (u *User) DelBook(epubsRelPath string) {
	books := u.Books
	found := -1

	for index, book := range books {
		if book.Path == epubsRelPath {
			found = index
			break
		}
	}
	if found != -1 {
		u.Books = append(books[:found], books[:found+1]...)
	}
}

func (u *User) Save() {
	bytes, err := yaml.Marshal(&u)
	if err != nil {
		fmt.Println("save user data err", err)
	}
	accPath := filepath.Join(AccFolder, u.Id+YamlExt)
	os.WriteFile(accPath, bytes, os.ModePerm)
}

func (u *User) UnzipAndGenerateEpubWebInfo(epubAbsPath string) (*Book, error) {
	var book Book
	e, err := epub.Open(epubAbsPath)
	if err != nil {
		return &book, err
	}
	opf, err := e.Package()
	if err != nil {
		return &book, err
	}

	if len(opf.Metadata.Title) > 0 {
		book.Name = opf.Metadata.Title[0].Value
	}

	book.Path = setting.ConfigYaml.GinEpubsStaticPath + epubAbsPath[len(setting.EpubsAbsPath):]

	metaTagList := opf.Metadata.Meta
	itemList := opf.Manifest.Items

	coverId := ""
	for _, metaTag := range metaTagList {
		if metaTag.Name == CoverTagName {
			coverId = metaTag.Content
			break
		}
	}

	coverRelPath := ""
	for _, item := range itemList {
		if coverId == item.ID {
			coverRelPath = item.Href
			break
		}
	}

	fileBaseName := filepath.Base(epubAbsPath)
	bookFolderName := strings.TrimSuffix(fileBaseName, filepath.Ext(fileBaseName))

	for _, file := range e.File {

		fileAbsPath := filepath.Join(setting.UnzipAbsPath, bookFolderName, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(fileAbsPath, os.ModePerm)
			continue
		}

		err = os.MkdirAll(filepath.Dir(fileAbsPath), os.ModePerm)
		if err != nil {
			fmt.Println("mkdir err", err)
		}

		fileReader, err := file.Open()
		if err != nil {
			fmt.Println("open zip file err", err)
		}
		defer fileReader.Close()

		destFile, err := os.Create(fileAbsPath)
		if err != nil {
			fmt.Println("create file err", err)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, fileReader)
		if err != nil {
			fmt.Println("copy failed", err)
		}

		if filepath.Base(file.Name) == coverRelPath {
			book.CoverPath = filepath.Join(setting.ConfigYaml.GinEpubsStaticPath, bookFolderName, file.Name)
		}

		if filepath.Ext(filepath.Base(file.Name)) == OpfExt {
			book.Opf = filepath.Join(setting.ConfigYaml.GinEpubsStaticPath, bookFolderName, file.Name)
		}
	}
	return &book, nil
}
