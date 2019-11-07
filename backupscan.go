package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"net/url"
	"io/ioutil"
	"time"
)


func file_read(filem string) []string {

	file, err := os.Open(filem)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, "http://www."+strings.TrimSpace(scanner.Text()))
	}

	return lines

}


func file_write(data string) {

	f, err := os.OpenFile("results.txt",
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
	start := time.Now()

	urls := create_backup(target)

	c := make(chan urlStatus)
	for _, path := range urls {
		go checkUrl(target+path, c)

	}
	result := make([]urlStatus, len(urls))
	for i, _ := range result {
		result[i] = <-c
		if result[i].status {
				if(result[i].response_type!="none") && !strings.Contains(result[i].response_type,"text") {
					fmt.Println(result[i].url, " backup file ")
					file_write(result[i].url + "\n")
				}

		}
	}
	p("Request Count:",len(urls))
	p("Scan Time :",time.Since(start).Seconds())

}
func main() {

	fmt.Println("Backup Scanner v0.1")
	fmt.Println("----- twitter.com/0x94 ----- ")
	var target string
	flag.StringVar(&target, "w", "", "go run backupscan.go -w url_list.txt")
	flag.Parse()

	if target == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	urls := file_read(target)

	for _, u := range urls {
		fmt.Println(u)
		scan_start(u+"/")
	}
	//scan_start("")

}


func checkUrl(path string, c chan urlStatus) {
	//fmt.Println(path)
	client := &http.Client{}
	req, err := http.NewRequest("GET",strings.TrimSpace(path), nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.70 Safari/537.36")

	resp, err := client.Do(req)
	if err==nil {

		body, err := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(body))

		if (err==nil && resp.StatusCode == 200) {
			c <- urlStatus{path, string(body),true,resp.Header.Get("Content-Type")}
		} else {
			c <- urlStatus{path, "none",false,"none"}
		}
	}else {
		c <- urlStatus{path, "none",false,"none"}

	}

}

type urlStatus struct {
	url    string
	response string
	status bool
	response_type string
}
