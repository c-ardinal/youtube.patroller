package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

/* Input JSON format */
type InfoList struct {
	Infos []TeamInfo
}

type TeamInfo struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Entry       string `json:"entry"`
	Org         string `json:"org"`
	Description string `json:description`
	Title       string `json:"title"`
	Performance string `json:"performance"`
}

/* Output JSON format */
type ResultList struct {
	Date    string       `json:"date"`
	Results []ResultInfo `json:"result"`
}

type ResultInfo struct {
	Name        string `json:"name"`
	Entry       string `json:"entry"`
	Org         string `json:"org"`
	Description string `json:"description"`
	TVC         string `json:"title-view-count"`
	TGC         string `json:"title-good-count"`
	TBC         string `json:"title-bad-count"`
	PVC         string `json:"performance-view-count"`
	PGC         string `json:"performance-good-count"`
	PBC         string `json:"performance-bad-count"`
}

/* Main routine */
func main() {
	var infoList InfoList
	inputJson := json.NewDecoder(os.Stdin)
	inputJson.Decode(&infoList)

	var resultList ResultList
	resultList.Date = time.Now().Format("2006-01-02_15:04:05")
	fmt.Println("========== " + resultList.Date + " ==========")

	for _, info := range infoList.Infos {
		var resultInfo ResultInfo

		if info.Id < 100 {
			resultInfo.Name = "0" + strconv.Itoa(info.Id) + "_" + info.Name
		} else {
			resultInfo.Name = strconv.Itoa(info.Id) + "_" + info.Name
		}
		fmt.Print(resultInfo.Name + "...")
		resultInfo.Entry = info.Entry
		resultInfo.Org = info.Org
		resultInfo.Description = info.Description
		resultInfo.TVC, resultInfo.TGC, resultInfo.TBC = checkVideoInfo(info.Title)
		resultInfo.PVC, resultInfo.PGC, resultInfo.PBC = checkVideoInfo(info.Performance)
		resultList.Results = append(resultList.Results, resultInfo)
		fmt.Println("ok")
	}

	outputJson, _ := json.MarshalIndent(resultList, "", "\t")
	//fmt.Printf("%s", outputJson)
	writeFile("result/"+resultList.Date+".json", outputJson)
	fmt.Println("=======================================")
}

/* Get ShowCount & GoogEvalCount & BadEvalCount */
func checkVideoInfo(url string) (string, string, string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println("error")
	}
	html, _ := doc.Html()
	viewCountRegex := regexp.MustCompile("<div class=\"watch-view-count\">視聴回数 ([0-9,]+) 回</div>")
	goodEvalCountRegex := regexp.MustCompile("他 ([0-9,]+) 人もこの動画を高く評価しました")
	badEvalCountRegex := regexp.MustCompile("他 ([0-9,]+) 人もこの動画を低く評価しました")
	return viewCountRegex.FindStringSubmatch(html)[1],
		goodEvalCountRegex.FindStringSubmatch(html)[1],
		badEvalCountRegex.FindStringSubmatch(html)[1]
}

/* Create and write JSON */
func writeFile(path string, data []byte) {
	file, _ := os.Create(path)
	defer file.Close()
	file.Write(data)
	fmt.Println("JSON write to " + path + " complete.")
}
