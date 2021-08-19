package main

import (
	"fmt"
	"log"
	"time"

	"main.go/Function/API"
	"main.go/Function/Config"
	"main.go/Function/File"
	"main.go/Function/Repository"
)

func main() {
	//SalvaDadosMongoDB()
	TesteConfig()
}

func SalvaDadosMongoDB() {

	hash, _ := File.LerTexto(Config.GetConfig().FileLog)

	cliente, contexto, cancel, errou := Repository.Connect(Config.GetConfig().ConnectionMongoDB)
	if errou != nil {
		log.Fatal(errou)
	}

	Repository.Ping(cliente, contexto)

	defer Repository.Close(cliente, contexto, cancel)
	for {

		valor := API.GetBloco(hash[0])

		Repository.ToDoc(valor)

		insertOneResult, err := Repository.InsertOne(cliente, contexto, Config.GetConfig().DataBase, Config.GetConfig().Collection, valor)

		// handle the error
		if err != nil {
			panic(err)
		}

		// print the insertion id of the document,
		// if it is inserted.
		fmt.Println("Result of InsertOne")
		fmt.Println(insertOneResult.InsertedID)

		hash = valor.Next_Block

		File.EscreverTexto(valor.Next_Block, Config.GetConfig().FileLog)

	}

}

func TesteArquivo() {
	for {
		fmt.Println("Sleep Start.....")
		// Calling Sleep method
		time.Sleep(1 * time.Second)

		// Printed after sleep is over
		fmt.Println("Sleep Over.....")

		texto, _ := File.LerTexto(Config.GetConfig().FileLog)

		fmt.Println(texto[0])
	}
}

func TesteConfig() {
	conf := Config.GetConfig()

	fmt.Println("ConnectionMongoDB: ", conf.ConnectionMongoDB)
	fmt.Println("FileLog: ", conf.FileLog)
	fmt.Println("DataBase: ", conf.DataBase, " Collection: ", conf.Collection)
	fmt.Println("UrlApi: ", conf.UrlAPI, " UrlApiBlock: ", conf.UrlAPIBlock, "UrlApiTransaction: ", conf.UrlAPITransaction)
}
