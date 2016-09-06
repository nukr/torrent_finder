package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enought arguments")
		os.Exit(1)
	}
	ysURL := os.Args[1]
	doc, err := goquery.NewDocument(ysURL)
	if err != nil {
		log.Fatal(err)
	}
	r := strings.NewReplacer("/", "")
	dirname := path.Join("/Users/nukr/Downloads", r.Replace(doc.Find("h4 a").Eq(0).Text()))
	doc.Find("#zoomtext .quote-content a").
		Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			idx := strings.Index(href, "fid=")
			downloadFile("http://www.16ys.org/attachment.php?fid="+href[idx+4:], dirname)
		})
}

func downloadFile(u, dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	res, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	r := strings.NewReplacer(`"`, "")
	contentDisposition := res.Header.Get("Content-Disposition")
	filenamePosition := strings.Index(contentDisposition, "filename=")
	escapedFilename := r.Replace(strings.Split(contentDisposition[filenamePosition:], "=")[1])
	filename, _ := url.QueryUnescape(escapedFilename)
	fmt.Println(filename)
	out, err := os.Create(path.Join(dir, filename))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, errCopy := io.Copy(out, res.Body)
	if errCopy != nil {
		log.Fatal(errCopy)
	}
}
