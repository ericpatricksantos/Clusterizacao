package Controllers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"main.go/Function/Auxiliares"

	"gopkg.in/mgo.v2/bson"
	"main.go/Function/API"
	"main.go/Function/Repository"
	Model "main.go/Models"
)

/*
	Recupera todos os Multi Endereços no MonhoDB e
	retorna em uma lista de strings
*/
func RecuperarAllMultiEnderecos(ConnectionMongoDB string, DataBase string, Collection string) []string {
	var addresses []string

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

	// le os documentos em partes, testei com 1000 documentos e deu certo
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var enderecos Model.MultiEndereco

		if err := cursor.Decode(&enderecos); err != nil {
			log.Fatal(err)
		}

		for _, item := range enderecos.Addresses {
			addresses = append(addresses, item.Address)
		}

	}

	return addresses
}

/*
Busca lista de endereços recuperado do mongoDB
Faz uma requisiçao com esses Endereços
Salva no MongoDB esses dados
Escreve em dois arquivos de log
	* LogIndiceMultiEndereco, ultimo valor q foi salvo no mongoDB
	* LogMultiEnderecosSemDados, Enderecos que voltaram sem dado e nao foram salvo no MongoDb
*/
func SalvaListaMultiEnderecos(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string,
	UrlAPI string, MultiAddr string, CollectionSalvaDados string,
	nomeArquivoSemApagar string, nomeArquivoIndice string, limiteEnderecos int) {

	listaEnderecos := Auxiliares.RecuperaEnderecosUnionArrayRemoveDuplicados(ConnectionMongoDB, DataBase, CollectionRecuperaDados)

	/* incrementa 1 no indice para remover o ultimo elemento salvo*/
	indiceInicial := GetIndiceLogIndice(nomeArquivoIndice) + 1
	listaMultiEnderecos, qtd := Auxiliares.TransformaArrayEmMatriz(listaEnderecos, indiceInicial, limiteEnderecos)
	indiceAtual := 0
	for contador := 0; contador < qtd; contador++ {
		if len(listaMultiEnderecos[contador]) > 0 {

			indiceAtual = indiceAtual + (len(listaMultiEnderecos[contador]))

			salvou := SalvarMultiEndereco(listaMultiEnderecos[contador], indiceAtual, UrlAPI, MultiAddr, ConnectionMongoDB, DataBase, CollectionSalvaDados,
				nomeArquivoSemApagar, nomeArquivoIndice)

			if !salvou {
				break
			}
		}

	}

	temp := []string{strconv.Itoa(indiceInicial + indiceAtual)}

	EscreverTexto(temp, nomeArquivoIndice)

	fmt.Println("Quantidade de enderecos salvos ", indiceInicial+indiceAtual)
}

/*Salvar multiEnderecos no MongoDb*/
func SalvarMultiEndereco(endereco []string, indice int, UrlAPI string, MultiAddr string,
	ConnectionMongoDB string, DataBase string, Collection string,
	nomeArquivoSemApagar string, nomeArquivoIndice string) bool {
	if len(endereco) > 0 && len(endereco) <= 100 {
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

		temp := []string{strconv.Itoa(indice)}

		EscreverTexto(temp, nomeArquivoIndice)

		fmt.Println("Quantidade de enderecos salvos ", indice)

		return true

	} else {
		valorTemp := " Indice: " + strconv.Itoa(indice) + " Endereco: " + strings.Join(endereco, "|")
		temp := []string{valorTemp}

		EscreverTextoSemApagar(temp, nomeArquivoSemApagar)
		fmt.Println("Nao tem Dados o Indice: " + strconv.Itoa(indice) + " Endereco: " + strings.Join(endereco, "|"))

		return false
	}
}
