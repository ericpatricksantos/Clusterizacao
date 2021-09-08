package Controllers

import (
	"fmt"
	"log"
	"strconv"

	"gopkg.in/mgo.v2/bson"
	"main.go/Function/API"
	"main.go/Function/Repository"
	Model "main.go/Models"
)

/*

	Recuperar Todos os enderecos que esta armazenado no mongoDB

	Retorno todas os input e out de todos os documentos

*/
func RecuperarEnderecos(ConnectionMongoDB string, DataBase string, Collection string) ([]string, []string) {
	var listaout []string
	var listainput []string

	// Get Client, Context, CalcelFunc and err from connect method.
	client, ctx, cancel, err := Repository.Connect(ConnectionMongoDB)
	if err != nil {
		panic(err)
	}

	// Free the resource when mainn dunction is  returned
	defer Repository.Close(client, ctx, cancel)

	// create a filter an option of type interface,
	// that stores bjson objects.
	var filter, option interface{}

	// filter  gets all document,
	// with maths field greater that 70
	filter = bson.M{}

	//  option remove id field from all documents
	option = bson.M{}

	// call the query method with client, context,
	// database name, collection  name, filter and option
	// This method returns momngo.cursor and error if any.
	cursor, err := Repository.Query(client, ctx, DataBase,
		Collection, filter, option)
	// handle the errors.
	if err != nil {
		panic(err)
	}

	// le os documentos em partes, testei com 1000 documentos e deu bom
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var enderecos Model.UnicoEndereco

		if err := cursor.Decode(&enderecos); err != nil {
			log.Fatal(err)
		}

		for _, doc := range enderecos.Txs {
			for _, x := range doc.Inputs {

				listainput = append(listainput, x.Prev_out.Addr)
			}
			for _, h := range doc.Out {
				listaout = append(listaout, h.Addr)
			}
		}

	}

	return listainput, listaout
}

// Salva somente os enderecos no MongoDB
func SalvarEnderecoMongoDB(enderecos []string, ConnectionMongoDB string, DataBase string, Collection string) {
	if len(enderecos) > 0 {
		cliente, contexto, cancel, errou := Repository.Connect(ConnectionMongoDB)
		if errou != nil {
			log.Fatal(errou)
		}

		er := Repository.Ping(cliente, contexto)

		if er != nil {
			log.Fatal(er)
		}

		defer Repository.Close(cliente, contexto, cancel)
		fmt.Println(len(enderecos))

		for i, valor := range enderecos {
			if len(valor) > 0 {
				insertOneResult, err := Repository.InsertOne(cliente, contexto, DataBase, Collection, valor[i])

				// handle the error
				if err != nil {
					panic(err)
				}

				// print the insertion id of the document,
				// if it is inserted.
				fmt.Println("Result of InsertOne")
				fmt.Println(insertOneResult.InsertedID)

			}
		}

	}
}

/*Salvar multiEnderecos*/
func SalvarMultiEndereco(endereco []string, UrlAPI string, MultiAddr string, ConnectionMongoDB string, DataBase string, Collection string) {
	if len(endereco) > 0 {
		cliente, contexto, cancel, errou := Repository.Connect(ConnectionMongoDB)
		if errou != nil {
			log.Fatal(errou)
		}

		Repository.Ping(cliente, contexto)

		defer Repository.Close(cliente, contexto, cancel)

		valor := API.GetMultiEnderecos(endereco, UrlAPI, MultiAddr)

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

	}
}

/*Salvar Unicoendereco pegando o endereco da API*/
func SalvarUnicoEndereco(endereco string, indice int, UrlAPI string, RawAddr string, ConnectionMongoDB string, DataBase string, Collection string, nomeArquivo string) {
	if len(endereco) > 0 {
		cliente, contexto, cancel, errou := Repository.Connect(ConnectionMongoDB)
		if errou != nil {
			log.Fatal(errou)
		}

		Repository.Ping(cliente, contexto)

		defer Repository.Close(cliente, contexto, cancel)

		valor := API.GetUnicoEndereco(endereco, UrlAPI, RawAddr)

		if len(valor.Address) > 0 {
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
		} else {
			valorTemp := " Indice: " + strconv.Itoa(indice) + " Endereco: " + endereco
			temp := []string{valorTemp}

			EscreverTextoSemApagar(temp, nomeArquivo)
			fmt.Println("Nao tem Dados o Indice: " + strconv.Itoa(indice) + " Endereco: " + endereco)
		}

	}
}
