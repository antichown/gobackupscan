package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"net/url"
	"io/ioutil"

)

func file_write(data string) {

	f, err := os.OpenFile("sonuc.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(data); err != nil {
		log.Println(err)
	}

}
func backup_list(tt string) []string {

	lines := []string {"db", "x", "ftp","hdocs",
		"1",  "root", "bak",  "admin", "www", "2014",
		"2015", "2016", "2014", "2015", "2016",
		"2017", "2018","2019","2020", "123", "a",
		"back", "backup","data", "root", "release","sql",
		"test", "template", "upfile", "vip",
		"web", "website", "wwwroot","wz","portal","blog","main","file"}

	//lines = append(lines, strings.TrimSpace())

	u, err := url.Parse(tt)
	if err != nil {
		panic(err)
	}

	resulthost := strings.Replace(u.Host, ".", "", -1)
	resulthost2 := strings.Replace(u.Host, "www.", "", -1)

	anahost := strings.Split(u.Host, ".")

	lines = append(lines, strings.TrimSpace(resulthost))
	lines = append(lines, strings.TrimSpace(u.Host))
	lines = append(lines, strings.TrimSpace(resulthost2))
	if strings.Contains(u.Host,"www") {
		lines = append(lines, strings.TrimSpace(anahost[1]))
	}else {
		lines = append(lines, strings.TrimSpace(anahost[0]))
	}
	return lines

}

func create_backup(tt string) (lines []string) {


	var cr []string=backup_list(tt)

	var ext = []string {".zip",".rar",".tar",".tar.gz",".tgz",".tar.bz2",".7z"}

	for i := 0; i < len(ext); i++ {
		for y := 0; y < len(cr); y++ {
			//fmt.Println(tt+cr[y]+ext[i])
			lines = append(lines, strings.TrimSpace(cr[y]+ext[i]))

		}
	}

	//lines = append(cr, strings.TrimSpace("xxx"))

	return lines
}

func scan_start(target string) {
	p := fmt.Println

	urls := create_backup(target)
	start := time.Now()
	c := make(chan urlStatus)
	for _, path := range urls {
		go checkUrl(target+path, c)

	}
	result := make([]urlStatus, len(urls))
	for i, _ := range result {
		result[i] = <-c
		if result[i].status {
			if !strings.Contains(result[i].response,"<head>") {
				fmt.Println(result[i].url, " backup file ")
				file_write(result[i].url + "\n")

			}

		}
	}
	p(time.Since(start).Seconds())

}
func main() {

	var target string
	flag.StringVar(&target, "t", "", "Usage")
	flag.Parse()

	if !strings.HasSuffix(target,"/") {
		target=target+"/"
	}
	scan_start(target)

}

func checkUrl(path string, c chan urlStatus) {
	//fmt.Println(path)
	client := &http.Client{}
	req, err := http.NewRequest("GET",strings.TrimSpace(path), nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.70 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if (err==nil && resp.StatusCode == 200) {
		c <- urlStatus{path, string(body),true}
	} else {
		c <- urlStatus{path, string(body),false}
	}
}

type urlStatus struct {
	url    string
	response string
	status bool
}
