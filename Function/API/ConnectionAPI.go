package API

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	Model "main.go/Models"
)

func GetMultiEnderecos(enderecos []string, urlAPI string, multiAddr string) Model.MultiEndereco {

	fmt.Println("Calling API...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlAPI+multiAddr+strings.Join(enderecos, "|"), nil)
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
	var responseObject Model.MultiEndereco
	json.Unmarshal(bodyBytes, &responseObject)

	return responseObject
}

func GetUnicoEndereco(endereco string, urlAPI string, RawAddr string) Model.UnicoEndereco {
	fmt.Println("Calling API...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlAPI+RawAddr+endereco, nil)
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
	var responseObject Model.UnicoEndereco
	json.Unmarshal(bodyBytes, &responseObject)

	return responseObject
}

func GetTransaction(hashTransacao string, urlAPI string, rawTx string) Model.Transaction {
	fmt.Println("Calling API...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlAPI+rawTx+hashTransacao, nil)
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

func GetBloco(hashBlock string, urlAPI string, rawBlock string) Model.Block {

	fmt.Println("Calling API...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlAPI+rawBlock+hashBlock, nil)
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
