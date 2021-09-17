package Controllers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
	"main.go/Function/API"
	"main.go/Function/Repository"
	Model "main.go/Models"
)

/*
	Apos chamar a funçao MapeandoMultiEndereco que retorna []Model.MapeandoTransacao
	Use SalvarMapeamentoTransacaoMongoDB para salva um []Model.MapeandoTransacao no MongoDb
*/
func SalvarMapeamentoTransacaoMongoDB(obj []Model.MapeandoTransacao, ConnectionMongoDB string, DataBase string, Collection string) {
	if len(obj) > 0 {
		cliente, contexto, cancel, errou := Repository.Connect(ConnectionMongoDB)
		if errou != nil {
			log.Fatal(errou)
		}

		er := Repository.Ping(cliente, contexto)

		if er != nil {
			log.Fatal(er)
		}

		defer Repository.Close(cliente, contexto, cancel)

		for i := 0; i < len(obj); i++ {
			Repository.ToDoc(obj[i])

			insertOneResult, err := Repository.InsertOne(cliente, contexto, DataBase, Collection, obj[i])

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

/*
	Recupera todos os valores de Multi Enderecos salvados no MongoDb
	Filtrar em quais Transações(hashTransacao) e a quantidade de vezes que um
	endereco(addr) aparece no array de input e out
	Retorna uma lista de objetos(Mapeando Transação)
*/
func MapeandoMultiEndereco(ConnectionMongoDB string, DataBase string, Collection string) []Model.MapeandoTransacao {
	var addressesMaping []Model.MapeandoTransacao

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

	addresses := RecuperarAllMultiEnderecos(ConnectionMongoDB, DataBase, Collection)

	// le os documentos em partes, testei com 1000 documentos e deu certo
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var enderecos Model.MultiEndereco

		if err := cursor.Decode(&enderecos); err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(addresses); i++ {
			// Definindo variáveis temporárias para atribuir os valores do endereco analisad e seus hash de transação
			var temp Model.MapeandoTransacao
			var tempInput Model.InputHash
			var tempOutput Model.OutputHash

			// Definindo o endereco analisado
			temp.Adresses = addresses[i]

			// Inicializando a variavel que indica quantas vezes no array de input e output
			// da Transação o endereco analisado(addr) aparece
			tempInput.Qtd = 0
			tempOutput.Qtd = 0

			for _, item := range enderecos.Txs {
				// Definindo o hash da transação
				tempInput.HashTransacao = item.Hash
				tempOutput.HashTransacao = item.Hash

				for _, inp := range item.Inputs {

					// Se o endereço analisado aparecer no array de input incrementar a variavel Qtd
					if addresses[i] == inp.Prev_out.Addr {

						tempInput.Qtd = tempInput.Qtd + 1
					}
				}
				for _, j := range item.Out {
					// Se o endereço analisado aparecer no array de Out incrementar a variavel Qtd
					if addresses[i] == j.Addr {
						tempOutput.Qtd = tempInput.Qtd + 1
					}

				}
			}

			// Atribuindo os valores da variavel temporaria
			temp.EntradaHash = append(temp.EntradaHash, tempInput)
			temp.SaidaHash = append(temp.SaidaHash, tempOutput)

			// Atribuindo os objetos que será retornado nessa função
			addressesMaping = append(addressesMaping, temp)
		}
	}

	return addressesMaping
}

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

	listaEnderecos := RecuperaEnderecosUnionArrayRemoveDuplicados(ConnectionMongoDB, DataBase, CollectionRecuperaDados)

	/* incrementa 1 no indice para remover o ultimo elemento salvo*/
	indiceInicial := GetIndiceLogIndice(nomeArquivoIndice) + 1
	listaMultiEnderecos, qtd := TransformaArrayEmMatriz(listaEnderecos, indiceInicial, limiteEnderecos)
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
