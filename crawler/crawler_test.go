package crawler

import (
	"fmt"
	"gonovel/trace"
	"testing"
)

func TestDMZJ(t *testing.T) {
	doc, err := getDocument("http://q.dmzj.com/2458/index.shtml", "utf-8")
	trace.CheckErr(err)
	html, _ := doc.Html()
	fmt.Println(html)
}
