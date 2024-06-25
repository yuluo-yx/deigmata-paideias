package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"ospp_search/pkg/model"
)

var (
	dataMap = make(map[string]FilteredRow)
)

type FilteredRow struct {
	ProId       int    `json:"proId"`
	ProgramName string `json:"programName"`
}

func readFile(filePath string) (map[string]FilteredRow, error) {
	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 解析 JSON 内容
	var rows []struct {
		ProId       int    `json:"proId"`
		ProgramCode string `json:"programCode"`
		ProgramName string `json:"programName"`
	}
	err = json.Unmarshal(data, &rows)
	if err != nil {
		return nil, err
	}

	// 将数据转换为 map[string]FilteredRow
	rowMap := make(map[string]FilteredRow)
	for _, row := range rows {
		rowMap[row.ProgramCode] = FilteredRow{
			ProId:       row.ProId,
			ProgramName: row.ProgramName,
		}
	}

	return rowMap, nil
}

func loadData() {

	for i := 1; i <= 12; i++ {
		fileName := fmt.Sprintf("data/%d--50.json", i)
		rowMap, err := readFile(fileName)
		if err != nil {
			fmt.Printf("read files %s failed: %v\n", fileName, err)
			continue
		}

		for programCode, row := range rowMap {
			dataMap[programCode] = row
		}
	}
}

func search() int {
	var proId int
	var customId = model.GetProject().Id
	if row, exists := dataMap[customId]; exists {
		fmt.Printf("ProgramCode: %s, Pro %v\n", customId, row)
		proId = row.ProId
	} else {
		fmt.Printf("ProgramCode: %s 不存在\n", customId)
	}

	return proId
}

func getPdf(proId int) {

	var url = "https://summer-ospp.ac.cn/previewPdf/" + fmt.Sprintf("%d", proId)
	download(url, "data/preview.pdf")
}

func download(url, fileName string) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/pdf")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to perform request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to download file, status code:", resp.StatusCode)
		return
	}

	out, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Failed to save file:", err)
		return
	}

	fmt.Println("File downloaded successfully")
}

func Query() {

	loadData()
	getPdf(search())
}
