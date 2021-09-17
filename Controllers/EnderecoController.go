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
Busca lista de endereços recuperado do mongoDB
Faz uma requisiçao com esses Endereços
Salva no MongoDB esses dados
Escreve em dois arquivos de log
	* LogIndiceEndereco, ultimo valor q foi salvo no mongoDB
	* logEnderecosSemDados, Enderecos que voltaram sem dado e nao foram salvo no MongoDb
*/
func SalvaListaEnderecos(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string,
	UrlAPI string, RawAddr string, CollectionSalvaDados string,
	nomeArquivoSemApagar string, nomeArquivoIndice string) {

	listaEnderecos := RecuperaEnderecosUnionArrayRemoveDuplicados(ConnectionMongoDB, DataBase, CollectionRecuperaDados)

	indiceInicial := GetIndiceLogIndice(nomeArquivoIndice) + 1

	for contador := indiceInicial; contador < len(listaEnderecos); contador++ {
		if len(listaEnderecos[contador]) > 0 {
			fmt.Printf("\n")
			fmt.Printf("-------------------------------------------------------------------------------")

			fmt.Printf("\nSalvando %dº endereço %s \n", contador, listaEnderecos[contador])

			salvou := SalvarUnicoEndereco(listaEnderecos[contador], contador, UrlAPI, RawAddr, ConnectionMongoDB, DataBase, CollectionSalvaDados,
				nomeArquivoSemApagar, nomeArquivoIndice)

			if !salvou {
				break
			}
		}
	}
}

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

/*Salvar Unicoendereco pegando o endereco da API*/
func SalvarUnicoEndereco(endereco string, indice int, UrlAPI string, RawAddr string,
	ConnectionMongoDB string, DataBase string, Collection string,
	nomeArquivoSemApagar string, nomeArquivoIndice string) bool {
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

			temp := []string{strconv.Itoa(indice)}

			EscreverTexto(temp, nomeArquivoIndice)

			fmt.Println("Indice atualizado para ", indice)

			fmt.Printf("Salvamento concluido do %dº endereco %s \n\n", indice, endereco)

			fmt.Printf("-------------------------------------------------------------------------------")
			fmt.Printf("\n")

			return true
		} else {
			valorTemp := " Indice: " + strconv.Itoa(indice) + " Endereco: " + endereco
			temp := []string{valorTemp}

			EscreverTextoSemApagar(temp, nomeArquivoSemApagar)
			fmt.Println("Nao tem Dados o Indice: " + strconv.Itoa(indice) + " Endereco: " + endereco)

			return false
		}

	} else {
		return false
	}
}
