package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/weilaihui/fdfs_client"
)

func main() {
	writeLog("info", "Web Server Start.")
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("github api gataway 1.0"))
	})
	http.HandleFunc("/upload", upload)
	http.ListenAndServe("0.0.0.0:8080", nil)
	writeLog("info", "Web Server End.")
}

func upload(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		rw.Write([]byte(buildResponsBody(-1, "Please use the Post request.", "")))
		return
	}

	isOK, fileExtName := getFileExtName(req.Header.Get("Authorization"))
	if !isOK {
		rw.Write([]byte(buildResponsBody(-2, fileExtName, "")))
		return
	}

	// 读取 body 内容
	data, _ := ioutil.ReadAll(req.Body)
	if len(data) < 4 {
		rw.Write([]byte(buildResponsBody(-4, "The length of the content is less than 4", "")))
		return
	}

	verifyOK, errorStr := authorization(req.Header.Get("Authorization"), []byte{data[0], data[1], data[len(data)-2], data[len(data)-1]})
	if !verifyOK {
		rw.Write([]byte(buildResponsBody(-3, errorStr, "")))
		return
	}

	path := "client.conf"
	url := "http://resource.github.com"
	fds, err := fdfs_client.NewFdfsClient(path)
	if fds == nil {
		writeLog("error", "fdfs conn error:"+err.Error())
		rw.Write([]byte(buildResponsBody(-101, "fdfs conn error", "")))
		return
	}

	uploadResponse, err := fds.UploadByBuffer(data, fileExtName)
	if uploadResponse == nil {
		writeLog("error", "upload error:"+err.Error())
		rw.Write([]byte(buildResponsBody(-102, "fdfs upload error", "")))
		return
	}

	url = url + "/" + uploadResponse.RemoteFileId
	url = strings.Replace(url, "\\", "/", -1)
	writeLog("info", "success:"+url)
	rw.Write([]byte(buildResponsBody(0, "", url)))
}

func getFileExtName(authValue string) (bool, string) {
	if len(authValue) < 1 {
		return false, "authorization, params error."
	}

	params := strings.Split(authValue, "&")
	if len(params) == 0 {
		return false, "authorization, params error."
	}

	fileExtName := strings.ToLower(authValue)
	if strings.Index(authValue, "extname=") == -1 {
		return false, "authorization, not find params 'extname'."
	}

	fileExtName = string([]byte(authValue)[strings.Index(authValue, "extname=")+8:])

	return true, fileExtName
}

func authorization(authValue string, bodyData []byte) (bool, string) {
	if len(authValue) < 1 {
		return false, "authorization, params error."
	}

	params := strings.Split(authValue, "&")
	if len(params) == 0 {
		return false, "authorization, params error."
	}

	authValue = strings.ToLower(authValue)
	if strings.Index(authValue, "timestamp=") == -1 {
		return false, "authorization, not find params 'timestamp'."
	}

	if strings.Index(authValue, "sign=") == -1 {
		return false, "authorization, not find params 'sign'."
	}

	timestamp := string([]byte(authValue)[strings.Index(authValue, "timestamp=")+10:])
	if strings.Index(timestamp, "&") > -1 {
		timestamp = string([]byte(timestamp)[:strings.Index(timestamp, "&")])
	}

	sign := string([]byte(authValue)[strings.Index(authValue, "sign=")+5:])
	if strings.Index(sign, "&") > -1 {
		sign = string([]byte(sign)[:strings.Index(sign, "&")])
	}

	// 1、检验时间戳是否正确
	timeUnix, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false, "timestamp, format error."
	}

	secTime := time.Unix(timeUnix, 0)
	if (time.Now().Unix() - secTime.Unix()) > 60*10 {
		return false, "timestamp, has gone out of time."
	}

	// 2、检验签名是否正常
	body := string(bodyData)
	key := "You Private Key"
	if buildMD5(timestamp+key+body) != sign {
		return false, "authorization, sign error."
	}

	return true, ""
}

func buildMD5(value string) string {
	data := []byte(value)
	has := md5.Sum(data)

	return fmt.Sprintf("%x", has)
}

func buildResponsBody(code int, msg string, url string) string {
	return "{\"code\":" + strconv.Itoa(code) + ", \"msg\":\"" + msg + "\", \"url\":\"" + url + "\"}"
}
