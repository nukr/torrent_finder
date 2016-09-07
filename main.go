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

// Finder ...
type Finder struct {
	domain string
	dir    string
	urls   []string
}

func (f *Finder) download() {
	createDirIfNotExists(f.dir)
	for _, u := range f.urls {
		downloadFile(u, f.dir)
	}
}

func main() {
	siteURL := checkArgs()
	finder := newFinder(siteURL)
	finder.download()
}

func newExtractor(domain string) func(u string) (string, []string) {
	extractors := make(map[string]func(string) (string, []string))
	extractors["www.16ys.org"] = extractor16ys
	return extractors[domain]
}

func extractor16ys(u string) (string, []string) {
	doc, err := goquery.NewDocument(u)
	if err != nil {
		log.Fatal(err)
	}
	r := strings.NewReplacer("/", "")
	dir := path.Join("/Users/nukr/Downloads", r.Replace(doc.Find("h4 a").Eq(0).Text()))
	var urls []string
	doc.Find("#zoomtext .quote-content a").
		Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			idx := strings.Index(href, "fid=")
			urls = append(urls, "http://www.16ys.org/attachment.php?fid="+href[idx+4:])
		})
	return dir, urls
}

func newFinder(siteURL string) Finder {
	domain := getDomain(siteURL)
	dir, urls := newExtractor(domain)(siteURL)
	finder := Finder{
		domain: domain,
		dir:    dir,
		urls:   urls,
	}
	return finder
}

// TODO: implement real getDomain
func getDomain(siteURL string) string {
	return "www.16ys.org"
}

func checkArgs() string {
	if len(os.Args) < 2 {
		fmt.Println("not enought arguments")
		os.Exit(1)
	}
	return os.Args[1]
}

// TODO: return error
func createDirIfNotExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func downloadFile(u, dir string) {
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
