package crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/piccuss/gonovel/trace"
)

type (
	//Novel identity
	Novel struct {
		Name      string    `json:"name"`
		Author    string    `json:"author"`
		URI       string    `json:"uri"`
		LastIndex int       `json:"last"`
		Type      int       `json:"type"`
		Chapters  []Chapter `json:"-"`
	}

	//Chapter identity
	Chapter struct {
		Index int
		Name  string
		URI   string
	}

	//Content identity
	Content string
)

const (
	lineLength  int    = 50
	typeBiquege int    = 2
	type37zw    int    = 1
	historyFile string = "history.json"
)

//Println show novel info
func (novel Novel) Println() {
	fmt.Println("Novelname:", novel.Name)
	fmt.Println("Author:   ", novel.Author)
	fmt.Println("URI:      ", novel.URI)
}

//ShowChapters display all chapters
func (novel Novel) ShowChapters() {
	chapters := novel.Chapters
	for index, value := range chapters {
		fmt.Println("Index:", index, "Chapter:", value.Name)
	}
}

//InitNovel initialize novel
func (novel *Novel) InitNovel() {
	getChapters(novel)
}

//Read get novel content by chapter index
func (novel *Novel) Read(index int) []Content {
	contents := getContents(*novel, index)
	return contents
}

//LoadHistory load read history
func LoadHistory() ([]*Novel, error) {
	historyData, err := ioutil.ReadFile(historyFile)
	if err != nil {
		return nil, err
	}
	var novels []*Novel
	err = json.Unmarshal(historyData, &novels)
	if err != nil {
		return nil, err
	}
	return novels, nil
}

//SaveHistory save read history
func SaveHistory(novels []*Novel) {
	saveData, _ := json.Marshal(novels)
	err := ioutil.WriteFile(historyFile, saveData, os.ModeAppend)
	trace.CheckErr(err)
}
