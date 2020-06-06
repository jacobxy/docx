package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"

	// "time"
	"strings"

	"github.com/unidoc/unioffice/document"
)

var BeginWard = []string{
	"摘要",
	"提要",
}

var EndWard = []string{
	"参考文献",
	"致谢",
}

var HeadWard = []string{}
var HeadWardPre = []string{
	"第",
	"（",
	"(",
}
var HeadWard2 = []string{
	"一",
	"二",
	"三",
	"四",
	"五",
	"六",
	"七",
	"八",
	"九",
	"十",
	// "首先",
	// "其次",
}
var HeadWard1 = []string{
	// "1",
	// "2",
	// "3",
	// "4",
	// "5",
	// "6",
	// "7",
	// "8",
	// "9", // "最后",
}

func init() {
	for _, v1 := range HeadWardPre {
		for _, v2 := range HeadWard2 {
			HeadWard = append(HeadWard, v1+v2)
		}
	}
	for _, v1 := range HeadWard1 {
		HeadWard = append(HeadWard, v1)
	}

}

func IsPre(str string, Stands []string) bool {
	if len(str) < 15 {
		return false
	}
	for _, v := range Stands {
		if strings.HasPrefix(str, v) {
			return true
		}
	}
	return false
}

func IsIn(str string, Stands []string) bool {
	for _, v := range Stands {
		if strings.Index(str, v) != -1 {
			return true
		}
	}
	return false
}

func GetPreAndLast(str string) (string, string) {
	strs := strings.Split(str, "。")
	l := len(strs)
	if l == 1 {
		return str, str
	}
	pre := strs[0]
	last := strs[l-1]
	if pre != "" {
		pre = pre + "。"
	}
	if last != "" {
		last = last + "。"
	}
	return pre, last
}

const (
	From = "docxsrc"
	To   = "result"
)

var toWord *document.Document

func main() {
	toWord = document.New()
	Dir, _ := ioutil.ReadDir(From)
	for _, f := range Dir {
		if f.IsDir() {
			continue
		}
		HandDocx(f.Name())
	}

	toFile := path.Join(To, "result.docx")
	toWord.SaveToFile(toFile)
}

func HandDocx(FileName string) {
	// FileName := "2.docx"
	file := path.Join(From, FileName)
	doc, err := document.Open(file)
	if err != nil {
		log.Fatalf("error opening document: %s", err, file)
	}

	fs := strings.Split(FileName, ".")

	paragraphs := []document.Paragraph{}
	for _, p := range doc.Paragraphs() {
		paragraphs = append(paragraphs, p)
	}

	for _, sdt := range doc.StructuredDocumentTags() {
		for _, p := range sdt.Paragraphs() {
			paragraphs = append(paragraphs, p)
		}
	}
	// dd, _ := document.New()
	flag := false
	Pre := ""
	NextGood := false

	_para := toWord.AddParagraph()
	_para.SetStyle("Heading1")
	_para.AddRun().AddText(fs[0])

	paraLength := len(paragraphs)
	for k0, p := range paragraphs {
		temp := make([]string, 0, 10)
		for _, r := range p.Runs() {
			temp = append(temp, r.Text())
		}
		res := strings.Join(temp, "")
		if !flag {
			if IsIn(res, BeginWard) {
				flag = true
			}
		}

		res = strings.ReplaceAll(res, " ", "")

		if !flag {
			continue
		}

		if k0 > paraLength/2 && IsIn(res, EndWard) {
			break
		}

		if false || IsPre(res, HeadWard) {
			_, last := GetPreAndLast(Pre)
			if Pre != "" && last != "" {
				toWord.AddParagraph().AddRun().AddText(last)
				fmt.Println(last)
			}
			pre, _ := GetPreAndLast(res)
			toWord.AddParagraph().AddRun().AddText(pre)
			fmt.Println(res)
			Pre = ""
			NextGood = true
			continue
		} else {
			if NextGood {
				pre, _ := GetPreAndLast(res)
				if pre != "" {
					toWord.AddParagraph().AddRun().AddText(pre)
					fmt.Println(pre)
				}
				NextGood = false
				Pre = ""
				continue
			}
			Pre = res
		}
		NextGood = false
	}

}
