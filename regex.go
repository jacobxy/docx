package main

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/unidoc/unioffice/document"
)

var reg1 = regexp.MustCompile(`[\p{Han}]+`)
var reg2 = regexp.MustCompile(`[1-9][0-9]{0,2}\.`)

var reg3 = regexp.MustCompile(`^[第(（]{0,1}([一二三四五六七八九])、{0,1}`)

// var reg3 = regexp.MustCompile(`[第（()]{0,1}([一二三四五六七八九123456789])[.|、]`)
var reg4 = regexp.MustCompile(`^[（(]{0,1}[0-9]+[)）]{0,1}`)

// var reg3 = regexp.MustCompile(`([0-9]*)$`)
var reg5 = regexp.MustCompile(`^([123]\.)+`)

var MustWrite []*regexp.Regexp

func init() {
	MustWrite = make([]*regexp.Regexp, 0, 10)
	MustWrite = append(MustWrite, reg3)
	MustWrite = append(MustWrite, reg4)
}

func IsMustWrite(res string) bool {
	for _, v := range MustWrite {
		r := v.FindStringSubmatch(res)
		if len(r) != 0 {
			return true
		}
	}
	return false
}

func main2() {
	str := `第一、1.2.3 .,，boy   20`
	// str := `1.2.3 .,，boy   20`
	r := reg3.FindStringSubmatch(str)
	fmt.Println(r)
}

type MyDocument struct {
	Doc     *document.Document
	LastStr string
	Write   bool
	ToDoc   *document.Document
	Path    string
}

func (this *MyDocument) Parser() bool {
	paragraphs := []document.Paragraph{}
	for _, p := range this.Doc.Paragraphs() {
		paragraphs = append(paragraphs, p)
	}
	for _, sdt := range this.Doc.StructuredDocumentTags() {
		for _, p := range sdt.Paragraphs() {
			paragraphs = append(paragraphs, p)
		}
	}

	for _, p := range paragraphs {
		temp := make([]string, 0, 10)
		for _, r := range p.Runs() {
			temp = append(temp, r.Text())
		}

		res := strings.Join(temp, "")
		res = Stand(res)
		if this.Write {
			strs := strings.Split(res, "。")
			this.ToDoc.AddParagraph().AddRun().AddText(strs[0] + "。")
			this.Write = false
			continue
		}
		if IsMustWrite(res) {
			strs := strings.Split(res, "。")
			if len(strs) > 1 {
				this.ToDoc.AddParagraph().AddRun().AddText(strings.Join(strs[:2], "。") + "。")
			} else {
				this.Write = true
				this.ToDoc.AddParagraph().AddRun().AddText(res)
			}
		}
	}
	this.ToDoc.SaveToFile(path.Join("toWord", this.Path))

	return true
}

func main() {
	filename := "13129435对外汉语教学法.docx"
	doc, _ := document.Open(filename)
	temp := &MyDocument{
		Path:  filename,
		Doc:   doc,
		ToDoc: document.New(),
	}
	temp.Parser()
}

var regDelete = regexp.MustCompile(`([(（][0-9a-zA-Z\p{Han}]{1,}[)）])`)

func Stand(str string) string {
	res := strings.ReplaceAll(str, " ", "")
	res = regDelete.ReplaceAllString(res, "")
	return res
}
