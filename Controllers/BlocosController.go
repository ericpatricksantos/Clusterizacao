package Controllers

import (
	"fmt"
	"log"

	"main.go/Function/API"
	"main.go/Function/File"
	"main.go/Function/Repository"
)

func SalvarBlocosMassa(nomeArquivo string, urlAPI string, rawBlock string, ConnectionMongoDB string, DataBase string, Collection string) {

	hash, _ := File.LerTexto(nomeArquivo)

	cliente, contexto, cancel, errou := Repository.Connect(ConnectionMongoDB)
	if errou != nil {
		log.Fatal(errou)
	}

	Repository.Ping(cliente, contexto)

	defer Repository.Close(cliente, contexto, cancel)
	for {

		valor := API.GetBloco(hash[0], urlAPI, rawBlock)

		Repository.ToDoc(valor)

		insertOneResult, err := Repository.InsertOne(cliente, contexto, DataBase, Collection, valor)

		// handle the error
		if err != nil {
			panic(err)
		}

		// print the insertion id of the document,
		// if it is inserted.
		fmt.Println("Result of InsertOne")
		fmt.Println(insertOneResult.InsertedID)

		hash = valor.Next_Block

		File.EscreverTexto(valor.Next_Block, nomeArquivo)

	}

}
