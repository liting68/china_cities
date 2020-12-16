package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
)

//GBK2UTF8 将GBK转换为UTF-8
func GBK2UTF8(s string) string {
	utf8 := mahonia.NewDecoder("gbk").ConvertString(s)
	return utf8
}

//UTF82GBK 将UTF-8转换为 GBK
func UTF82GBK(s string) string {
	gbk := mahonia.NewEncoder("gbk").ConvertString(s)
	return gbk
}

//CurlGET 向服务端发送get请求
func CurlGET(requestURL string) (bodystr string) {
	request, err := http.NewRequest("GET", strings.Trim(requestURL, " "), nil)
	if err != nil {
		fmt.Println("CurlGET==NewRequest ERROR", err.Error())
		return err.Error()
	}
	// 接收服务端返回给客户端的信息
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("CurlGET==client.Do ERROR", err.Error())
		return err.Error()
	}
	if res.StatusCode == 200 {
		str, _ := ioutil.ReadAll(res.Body)
		bodystr = string(str)
		// fmt.Println("CurlGET==", bodystr)
	}
	return bodystr

}

type Area struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}
type City struct {
	Code     string `json:"code"`
	Value    string `json:"value"`
	Children []Area `json:"children"`
}
type Province struct {
	Code     string `json:"code"`
	Value    string `json:"value"`
	Children []City `json:"children"`
}

func main() {
	indexURL := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2020/"
	resIndex := CurlGET(indexURL + "index.html")
	re := regexp.MustCompile(`<a href=\'(\d{2,4}).html\'>(.{3,20})<br\/><\/a>`)
	match := re.FindAllStringSubmatch(resIndex, -1)
	province := []Province{}
	for _, arr := range match {
		cities := []City{}
		resCity := CurlGET(indexURL + arr[1] + ".html")
		reCity := regexp.MustCompile(`<a href=\'\d{2}\/(.{1,30}).html\'>(.{1,30})<\/a><\/td><\/tr>`)
		matchCities := reCity.FindAllStringSubmatch(resCity, -1)
		for _, arrCitiy := range matchCities {
			areas := []Area{}
			resArea := CurlGET(indexURL + arr[1] + "/" + arrCitiy[1] + ".html")
			reArea := regexp.MustCompile(`<a href=\'\d{2}\/(.{1,30}).html\'>(.{1,30})<\/a><\/td><\/tr>`)
			matchAreas := reArea.FindAllStringSubmatch(resArea, -1)
			for _, arrArea := range matchAreas {
				fmt.Println(arrArea[1], ":", GBK2UTF8(arrArea[2]))
				areas = append(areas, Area{Code: arrArea[1], Value: GBK2UTF8(arrArea[2])})
			}
			fmt.Println(arrCitiy[1], ":", GBK2UTF8(arrCitiy[2]))
			cities = append(cities, City{Code: arrCitiy[1], Value: GBK2UTF8(arrCitiy[2]), Children: areas})
		}
		// if len(matchCities) == 1 {
		// 	fmt.Println(matchCities[0][1])
		// 	resArea := CurlGET(indexURL + arr[1] + "/" + matchCities[0][1] + ".html")
		// 	reArea := regexp.MustCompile(`<a href=\'\d{2}\/(.{1,30}).html\'>(.{1,30})<\/a><\/td><\/tr>`)
		// 	matchAreas := reArea.FindAllStringSubmatch(resArea, -1)
		// 	for _, arrArea := range matchAreas {
		// 		fmt.Println(arrArea[1], ":", GBK2UTF8(arrArea[2]))
		// 		cities = append(cities, Area{Code: arrArea[1], Value: GBK2UTF8(arrArea[2])})
		// 	}
		// }
		fmt.Println(arr[1], ":", GBK2UTF8(arr[2]))
		province = append(province, Province{Code: arr[1], Value: GBK2UTF8(arr[2]), Children: cities})
	}
	basedir, _ := os.Getwd()
	fileName := strconv.FormatInt(time.Now().Unix(), 10) + ".json"
	fmt.Println(basedir + "/" + fileName)
	if data, err := json.Marshal(province); err == nil {
		ioutil.WriteFile(basedir+"/"+fileName, data, 0644)
	}
}
