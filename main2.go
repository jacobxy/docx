package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	ClientId = "M5qShRNYdPw9pvF49Fjr3MG9"
	// ClientId    = "dc29769d0f9d4b41890b00cec22dd271"
	SecretKey   = "SVuTYUyxiTGawGiPcQmGCY9v0invT06Q"
	default_url = "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s"
	Token       = "24.bb4a90d85e311fa4fe518cab4705bf32.2592000.1593264932.282335-18015470"

	AskUrl = "https://aip.baidubce.com/rpc/2.0/nlp/v1/news_summary?charset=UTF-8&access_token=%s"
	// AskUrl = "https://aip.baidubce.com/rpc/2.0/nlp/v1/lexer?charset=UTF-8&access_token=%s"
)

func GetToken() {
	url := fmt.Sprintf(default_url, ClientId, SecretKey)
	res, _ := http.Post(url, "", nil)
	b, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))
}

func main() {
	GetToken()
}

type BaiduAsk struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	MaxSummaryLen int    `json:"max_summary_len"`
}

type BaiduResponse struct {
	LogId   int64  `json:"log_id"`
	Summary string `json:"summary"`
}

func HandleContentByBaidu(title, str string, max int) string {
	str = strings.TrimSpace(str)

	rSlice := strings.Split(str, "。")
	r1 := rSlice[0] + "。"
	max = 100
	str = strings.Join(rSlice[1:], "。")

	temp := BaiduAsk{
		Title:         rSlice[0],
		Content:       str,
		MaxSummaryLen: max,
	}
	b, _ := json.Marshal(temp)
	buffer := bytes.NewBuffer(b)
	url := fmt.Sprintf(AskUrl, Token)
	// fmt.Println(url)
	request, _ := http.NewRequest("POST", url, buffer)
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return r1
	}
	b, _ = ioutil.ReadAll(res.Body)
	result := BaiduResponse{}
	err = json.Unmarshal(b, &result)
	// fmt.Println(string(b))
	fmt.Println(title, len(str), max)
	fmt.Println(result, res.StatusCode)
	return r1 + result.Summary
}
