package main

import (
	"bufio"
	"filestore/server/cache"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var (
	port string
	ver  bool
)

func init() {
	flag.StringVar(&port, "port", ":5000", "The port to listen on.")
	flag.BoolVar(&ver, "version", false, "Print server version.")
}

const (
	// base HTTP paths.
	apiVersion  = "v1"
	apiBasePath = "/files/" + apiVersion + "/"
	addfile     = apiBasePath + "store"
	hashexist   = apiBasePath + "hexist"
	copyfile    = apiBasePath + "copy"
	wc          = apiBasePath + "wc"
	rm          = apiBasePath + "rm"
	list        = apiBasePath + "ls"

	// server version.
	version = "1.0.0"
)

func main() {

	fmt.Println("Started File Store on port ", port, " Version : ", version)
	r := gin.New()
	r.Use(gin.Logger())
	r.POST(addfile, AddFile)
	r.POST(copyfile, CopyFile)
	r.GET(hashexist, HashExist)
	r.DELETE(rm, RemoveFile)
	r.GET(list, ListFiles)
	r.GET(wc, WordCount)
	r.Run(port)

}

type Response struct {
	Error string
}

type CopyFileData struct {
	FileName string
	Hash     string
}

func CopyFile(c *gin.Context) {
	b := &CopyFileData{}
	if err := c.ShouldBindWith(b, binding.JSON); err == nil {
		rdb := cache.GetConnection()
		hashCmd := rdb.HGet(rdb.Context(), "FileHash", b.Hash)
		if hashCmd.Err() != nil {
			c.JSON(http.StatusBadRequest, Response{Error: hashCmd.Err().Error()})
		} else {
			res, err := hashCmd.Result()
			if err != nil {
				c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
			} else {
				files := strings.Split(res, ",")
				if len(files) > 0 {
					//filenm := "filestore/" + b.FileName
					_, err := copy(files[0], "filestore/"+b.FileName)
					fmt.Println(err)
					if err != nil {
						c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
					} else {
						fmt.Println("Set Hash", b.Hash, "nm ", b.FileName)
						rdb.HSet(rdb.Context(), "FileHash", b.Hash, res+","+b.FileName)
						rdb.HSet(rdb.Context(), "HashFromFile", map[string]interface{}{b.FileName: b.Hash})
						c.JSON(http.StatusOK, Response{Error: ""})
					}
				} else {
					c.JSON(http.StatusInternalServerError, Response{Error: "Copy File not found"})
				}
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
	}
}

func AddFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	hash := c.Request.Header.Get("hash")
	fmt.Println("Hash is ", hash)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	rdb := cache.GetConnection()
	filename := header.Filename
	filename = filepath.Base(filename)
	fmt.Println(filename)
	out, err := os.Create("filestore/" + filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
	}
	hashCmd := rdb.HGet(rdb.Context(), "FileHash", hash)
	res, err := hashCmd.Result()
	fmt.Println(res, err)
	if res == "" {
		fmt.Println("Res empty Hash ", hash, filename)
		rdb.HSet(rdb.Context(), "FileHash", map[string]interface{}{hash: filename})
	} else {
		fmt.Println("Setting Hash ", hash, res+","+filename)
		rdb.HSet(rdb.Context(), "FileHash", map[string]interface{}{hash: res + "," + filename})
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}
	rdb.HSet(rdb.Context(), "HashFromFile", map[string]interface{}{filename: hash})
	c.JSON(http.StatusOK, Response{Error: ""})
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func RemoveFile(c *gin.Context) {
	b := &CopyFileData{}
	if err := c.ShouldBindWith(b, binding.JSON); err == nil {
		filename := b.FileName
		fmt.Println(filename)
		rdb := cache.GetConnection()
		hashfromfile := rdb.HGet(rdb.Context(), "HashFromFile", filename)
		hash, _ := hashfromfile.Result()
		res, err := (rdb.HGet(rdb.Context(), "FileHash", hash)).Result()
		if res != "" && err == nil {
			splitstr := strings.Split(res, ",")
			if len(splitstr) > 1 {
				newres := strings.Replace(res, filename, "", -1)
				rdb.HSet(rdb.Context(), "FileHash", hash, newres)
			} else {
				rdb.HDel(rdb.Context(), "FileHash", hash, filename)
			}
		}
		err = os.Remove("filestore/" + filename)
		fmt.Println(err)
		if err != nil {
			c.JSON(http.StatusNotFound, Response{Error: err.Error()})
		} else {
			rdb.HDel(rdb.Context(), "HashFromFile", b.FileName)
			c.JSON(http.StatusOK, Response{Error: ""})
		}
	} else {
		c.JSON(http.StatusBadRequest, Response{Error: "Failed to Marshal json"})
	}
}

func ListFiles(c *gin.Context) {
	dataarr, err := listFiles()
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"files": dataarr})
	}
}

func listFiles() ([]string, error) {
	strarr := make([]string, 0)
	err := filepath.Walk("filestore/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		strarr = append(strarr, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return strarr, nil
}

func WordCount(c *gin.Context) {
	strarr, err := listFiles()
	if err != nil {
		c.JSON(http.StatusNotFound, Response{Error: err.Error()})
	} else {
		count := 0
		wg := &sync.WaitGroup{}
		ch := make(chan int, len(strarr))
		for _, v := range strarr {
			wg.Add(1)
			go wordCountFromFile(v, ch, wg)
		}

		wg.Wait()
		for {
			count += <-ch
			if len(ch) == 0 {
				break
			}
		}
		c.JSON(http.StatusAccepted, gin.H{"count": count})
	}
}

func wordCountFromFile(v string, ch chan int, wg *sync.WaitGroup) {
	fileHandle, err := os.Open(v)
	if err != nil {
		wg.Done()
	}
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)
	fileScanner.Split(bufio.ScanWords)
	count := 0
	for fileScanner.Scan() {
		count++
	}
	wg.Done()
	ch <- count
}

func HashExist(c *gin.Context) {
	hash := c.Request.Header.Get("hash")
	filename := c.Request.Header.Get("file")
	rdb := cache.GetConnection()
	res := rdb.HGet(rdb.Context(), "FileHash", hash)
	result, _ := res.Result()
	if fs, err := os.Stat("filestore/" + filename); err == nil {
		fmt.Println(fs.Name())
		c.String(http.StatusConflict, "")
		return
	}
	if result != "" {
		fmt.Println("Found")
		c.String(http.StatusOK, "")
	} else {
		fmt.Println("Not Found")
		c.String(http.StatusNotFound, "")
	}
}
