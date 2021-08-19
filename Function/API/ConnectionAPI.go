package API

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"main.go/Function/Config"
	Model "main.go/Models"
)

func GetTransaction(hashTransacao string) Model.Transaction {
	fmt.Println("Calling API...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", Config.GetConfig().UrlAPITransaction+hashTransacao, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject Model.Transaction
	json.Unmarshal(bodyBytes, &responseObject)

	return responseObject
}

func GetBloco(hastBlock string) Model.Block {

	fmt.Println("Calling API...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", Config.GetConfig().UrlAPIBlock+hastBlock, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject Model.Block
	json.Unmarshal(bodyBytes, &responseObject)

	return responseObject
}
