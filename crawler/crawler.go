package crawler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/piccuss/gonovel/trace"
	"github.com/axgle/mahonia"
	"github.com/PuerkitoBio/goquery"
)

const (
	novelAPI       string = "https://www.37zw.net"
	novelAPI2      string = "http://www.biquge.info"
	novelSearchAPI string = "https://www.37zw.net/s/so.php?type=articlename&s="
)

//SearchNovel return search result by novel name
func SearchNovel(name string) []*Novel {
	doc, err := getDocument(novelSearchAPI+name, "gbk")
	trace.CheckErr(err)
	novels := []*Novel{}
	doc.Find(".novellist").First().Find("li").Each(func(i int, s *goquery.Selection) {
		novel := &Novel{}
		//parse novel name, author and URI
		novelInfo := strings.Split(s.Text(), "/")
		novel.Name = novelInfo[0]
		novel.Author = novelInfo[1]
		novel.URI, _ = s.Find("a").First().Attr("href")
		novel.Type = type37zw
		novels = append(novels, novel)
	})
	return novels
}

//getChapters return novel chapters by novel uri
func getChapters(novel *Novel) []Chapter {
	charsetType := "gbk"
	api := novelAPI
	if novel.Type == typeBiquege {
		charsetType = "utf-8"
		api = novelAPI2
	}
	doc, err := getDocument(api+novel.URI, charsetType)
	trace.CheckErr(err)

	chapters := []Chapter{}
	doc.Find("#list").First().Find("dd").Each(func(i int, s *goquery.Selection) {
		//parse chapter index, name and URI
		chapter := Chapter{}
		chapter.Index = i
		chapter.Name = s.Text()
		chapter.URI, _ = s.Find("a").First().Attr("href")
		chapters = append(chapters, chapter)
	})
	novel.Chapters = chapters
	return chapters
}

//getContents return chapter content
func getContents(novel Novel, index int) []Content {
	charsetType := "gbk"
	api := novelAPI
	if novel.Type == typeBiquege {
		charsetType = "utf-8"
		api = novelAPI2
	}
	doc, err := getDocument(api+novel.URI+novel.Chapters[index].URI, charsetType)
	trace.CheckErr(err)

	contents := []Content{}
	doc.Find("#content").First().Each(func(i int, s *goquery.Selection) {
		//parse content
		html, _ := s.Html()
		for _, value := range strings.Split(html, "<br/><br/>") {
			content := Content(value)
			contents = append(contents, content)
		}
	})
	return contents
}

//getDocument parse html to document
func getDocument(url string, charset string) (*goquery.Document, error) {
	res, err := http.Get(url)
	trace.CheckErr(err)

	defer res.Body.Close()
	if res.StatusCode != 200 {
		trace.Error.Fatalln("Get chapters error.Status code :%d", res.StatusCode)

	}

	//decode gbk html to utf-8
	respByte, err := ioutil.ReadAll(res.Body)
	trace.CheckErr(err)
	if charset == "gbk" {
		respByte = DecodeGBKBytes(respByte)
	}
	decodeReader := bytes.NewReader(respByte)
	//parse html body
	doc, err := goquery.NewDocumentFromReader(decodeReader)
	return doc, err
}

//DecodeGBKBytes decode gbk bytes to utf-8 bytes
func DecodeGBKBytes(source []byte) []byte {
	decoder := mahonia.NewDecoder("GB18030")
	decodedString := decoder.ConvertString(string(source))
	return []byte(decodedString)
}
