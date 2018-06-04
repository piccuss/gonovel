package main

import (
	"bufio"
	"fmt"
	"gonovel/crawler"
	"os"
	"strconv"
	"strings"
)

var historyNovels = []*crawler.Novel{}

func main() {
	fmt.Println("Novel online reader.Copyright:Piccus Peng 2018")
	fmt.Println("Loading history file...")
	historyNovelsData, err := crawler.LoadHistory()
	if err != nil {
		fmt.Println(err.Error())
	}
	if len(historyNovelsData) != 0 {
		historyNovels = historyNovelsData
	}
	for {
		showMenu()
	}
}

//get user input
func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.Trim(input, "\n")
	input = strings.Trim(input, "\r")
	return input
}

//show app menu
func showMenu() {
	fmt.Println("1. Continue reading by history")
	fmt.Println("2. Search book online by name")
	fmt.Println("3. Quit")
	fmt.Print("Input command index:")
	command := getInput()
	switch command {
	case "1":
		showHistory()
	case "2":
		search()
	case "3":
		fmt.Println("bye!")
		os.Exit(0)
	default:
		fmt.Println("Invalid command")
		fmt.Println()
	}
}

func showHistory() {
	for index, novel := range historyNovels {
		fmt.Println(index, ".")
		novel.Println()
		fmt.Println()
	}
	for {
		fmt.Println("Input book index to continue:")
		index := getInput()
		indexInt64, _ := strconv.ParseInt(index, 10, 0)
		indexInt := int(indexInt64)
		if indexInt >= 0 && indexInt < len(historyNovels) {
			fmt.Println("Continue reading book:", historyNovels[indexInt].Name)
			fmt.Println()
			read(historyNovels[indexInt], historyNovels[indexInt].LastIndex)
		} else {
			fmt.Println("Invalid book index.")
			fmt.Println()
		}
	}

}

//search book by name
func search() {
	for {
		fmt.Print("Input book name:")
		name := getInput()
		if name != "" {
			novels := crawler.SearchNovel(name)
			for index, novel := range novels {
				fmt.Println(index, ".")
				novel.Println()
				fmt.Println()
			}
			fmt.Println("Input book index to start read:")
			index := getInput()
			indexInt64, _ := strconv.ParseInt(index, 10, 0)
			indexInt := int(indexInt64)
			if indexInt >= 0 && indexInt < len(novels) {
				fmt.Println("Start reading book:", novels[indexInt].Name)
				fmt.Println()
				read(novels[indexInt], 0)
			} else {
				fmt.Println("Invalid book index.")
				fmt.Println()
			}

		}
	}
}

//read novel from chapter index
func read(novel *crawler.Novel, index int) {
	if len(novel.Chapters) == 0 {
		fmt.Println("Book is intializing...")
		novel.InitNovel()
	}
	if max := len(novel.Chapters) - 1; index > max {
		save(novel)
		fmt.Println("There is no content later.")
		os.Exit(0)
	}
	fmt.Println("Loading chapter...")
	contents := novel.Read(index)
	maxIndex := len(contents)
	fmt.Println(novel.Chapters[index].Name)
	fmt.Println()
	novel.LastIndex = index
	for nowIndex := 0; nowIndex < maxIndex; nowIndex++ {
		fmt.Println(contents[nowIndex])
		listenCommand(novel)
	}
	index++
	fmt.Println()
	read(novel, index)
}

//Jump to chapter index
func jump(novel *crawler.Novel) {
	fmt.Print("Input chapter index:")
	jumpIndex := getInput()
	jumpIndexInt64, _ := strconv.ParseInt(jumpIndex, 10, 0)
	jumpIndexInt := int(jumpIndexInt64)
	if jumpIndexInt >= 0 && jumpIndexInt < len(novel.Chapters) {
		read(novel, jumpIndexInt)
	} else {
		fmt.Println("Invalid chapter index!")
	}
}

//Listen user command during reading
func listenCommand(novel *crawler.Novel) {
	//Press to read
	command := getInput()
	switch command {
	case "q":
		save(novel)
		fmt.Println("bye!")
		os.Exit(0)
	case "jump":
		jump(novel)
	case "index":
		novel.ShowChapters()
	case "s":
		save(novel)
	case "fc":
		findChapter(novel)
	case "":
		fmt.Println()
	default:
		fmt.Println("Invalid command!")
	}
}

//save novel read history
func save(novel *crawler.Novel) {
	//if novel exist history
	tag := false
	for index, historyNovel := range historyNovels {
		if historyNovel.Name == novel.Name {
			historyNovels[index] = novel
			tag = true
		}
	}
	if !tag {
		historyNovels = append(historyNovels, novel)
	}
	crawler.SaveHistory(historyNovels)
	fmt.Println("Save complete...")
}

//find chapter return chapter index
func findChapter(novel *crawler.Novel) {
	fmt.Print("Input chaptername(contains):")
	chaptername := getInput()
	if chaptername != "" {
		for index, chapter := range novel.Chapters {
			if strings.Contains(chapter.Name, chaptername) {
				fmt.Println("Index:", index, "Chapter:", chapter.Name)
			}
		}
	}
}
