package account

import (
	"epub-reader-web-server/setting"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pirmd/epub"
	"gopkg.in/yaml.v3"
)

type ReadingPos struct {
	Link       string `yaml:"link"`
	Percentage string `yaml:"percentage"`
}

type Book struct {
	Id         string `yaml:"id"`
	Name       string `yaml:"name"`
	Path       string `yaml:"path"`
	CoverPath  string `yaml:"cover_path"`
	Opf        string `yaml:"opf"`
	ReadingPos ReadingPos
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

var LastTimestamp int64 = 0
var randomNum int = 0

func Init() {
	rand.NewSource(time.Now().UnixNano())
}

func CreateUser(id string, name string, pw string) {
	u := User{id, name, pw, []Book{}}
	u.Save()
}

func IsUserExist(id string) bool {
	accFilePath := filepath.Join(AccFolder, id+YamlExt)
	_, err := os.ReadFile(accFilePath)
	return err == nil
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
	return &user, nil
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

func (u *User) GetBookInfo(id string) (*Book, error) {
	books := u.Books
	for index, bookinfo := range books {
		if bookinfo.Id == id {
			return &books[index], nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no book found with id %s", id))
}

func (u *User) Save() {
	bytes, err := yaml.Marshal(&u)
	if err != nil {
		fmt.Println("save user data err", err)
	}
	accPath := filepath.Join(AccFolder, u.Id+YamlExt)
	err = os.WriteFile(accPath, bytes, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}

func (u *User) UnzipAndGenerateEpubWebInfo(epubAbsPath string) (*Book, error) {
	var book Book
	e, err := epub.Open(epubAbsPath)
	if err != nil {
		return &book, err
	}
	defer e.Close()
	opf, err := e.Package()
	if err != nil {
		return &book, err
	}

	if len(opf.Metadata.Title) > 0 {
		book.Name = opf.Metadata.Title[0].Value
	}

	currentTimestamp := time.Now().Unix()
	if currentTimestamp != LastTimestamp {
		LastTimestamp = currentTimestamp
		randomNum = rand.Intn(2 ^ 16)
	} else {
		randomNum++
	}
	timeStampHex := strconv.FormatInt(currentTimestamp, 16)
	indexHex := strconv.FormatInt(int64(randomNum), 16)
	if len(indexHex) > 4 {
		indexHex = indexHex[len(indexHex)-4:]
	}

	book.Id = timeStampHex + fmt.Sprintf("%04s", indexHex)
	// fmt.Println("bookid", book.Id, timeStampHex, randomNum, indexHex)
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

		_, err = os.Stat(fileAbsPath)
		if os.IsNotExist(err) {
			destFile, err := os.Create(fileAbsPath)
			if err != nil {
				fmt.Println("create file err", err)
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, fileReader)
			if err != nil {
				fmt.Println("copy failed", err)
			}

		}

		if filepath.Base(file.Name) == filepath.Base(coverRelPath) {
			book.CoverPath = filepath.Join(setting.ConfigYaml.GinEpubsStaticPath, bookFolderName, file.Name)
		}

		if filepath.Ext(filepath.Base(file.Name)) == OpfExt {
			book.Opf = filepath.Join(setting.ConfigYaml.GinEpubsStaticPath, bookFolderName, file.Name)
		}
	}
	return &book, nil
}
