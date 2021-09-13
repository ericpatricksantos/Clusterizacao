package Controllers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
	"main.go/Function/API"
	"main.go/Function/Arrays"
	"main.go/Function/Repository"
	Model "main.go/Models"
)

/*
Busca lista de endereços recuperado do mongoDB
Faz uma requisiçao com esses Endereços
Salva no MongoDB esses dados
Escreve em dois arquivos de log
	* LogIndiceMultiEndereco, ultimo valor q foi salvo no mongoDB
	* LogMultiEnderecosSemDados, Enderecos que voltaram sem dado e nao foram salvo no MongoDb
*/
func SalvaListaMultiEnderecos(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string,
	UrlAPI string, rota string, CollectionSalvaDados string,
	nomeArquivoSemApagar string, nomeArquivoIndice string) {

	listaEnderecos := RecuperaEnderecosUnionArrayRemoveDuplicados(ConnectionMongoDB, DataBase, CollectionRecuperaDados)

	indiceInicial := GetIndiceLogIndice(nomeArquivoIndice)
	listaMultiEnderecos, qtd := TransformaArrayEmMatriz(listaEnderecos, indiceInicial)
	indiceAtual := 0
	for contador := 0; contador < qtd; contador++ {
		if len(listaMultiEnderecos[contador]) > 0 {

			indiceAtual = indiceAtual + (len(listaMultiEnderecos[contador]))

			SalvarMultiEndereco(listaMultiEnderecos[contador], indiceAtual, UrlAPI, rota, ConnectionMongoDB, DataBase, CollectionSalvaDados,
				nomeArquivoSemApagar, nomeArquivoIndice)

		}
	}

	temp := []string{strconv.Itoa(indiceInicial + indiceAtual)}

	EscreverTexto(temp, nomeArquivoIndice)

	fmt.Println("Quantidade de enderecos salvos ", indiceInicial+indiceAtual)
}

func TransformaArrayEmMatriz(listaEnderecos []string, indiceInicial int) ([][]string, int) {
	QtdLinhas := (len(listaEnderecos) / 100) + 1
	tamanhoMaximoArray := 100 //Qtd de Colunas na Matriz
	matrizEnderecos := make([][]string, QtdLinhas)
	if indiceInicial > 0 {
		//Remove os elementos que foram salvos na matriz
		listaEnderecos = append(listaEnderecos[:0], listaEnderecos[indiceInicial:]...)

	}

	for contador := 0; contador < QtdLinhas; contador++ {

		enderecosSeparados := listaEnderecos[:tamanhoMaximoArray]
		if len(enderecosSeparados) == 0 {
			break
		}
		//Alocar memoria para o array de string
		matrizEnderecos[contador] = make([]string, len(enderecosSeparados))
		for j := 0; j < len(enderecosSeparados); j++ {
			//Atribuir os valores do slice para a matriz
			matrizEnderecos[contador][j] = enderecosSeparados[j]
		}

		//Remove os elementos que foram salvos na matriz
		listaEnderecos = append(listaEnderecos[:0], listaEnderecos[tamanhoMaximoArray:]...)

		// Se o tamanho do slice for menor do que o tamanho maximo do Array,
		// setar tamanhoMaximo = tamanho Array
		// Evitar acesso invalido de memoria
		if tamanhoMaximoArray > (len(listaEnderecos)) {
			tamanhoMaximoArray = len(listaEnderecos)
		}

	}

	return matrizEnderecos, QtdLinhas
}

/*
Busca lista de endereços recuperado do mongoDB
Faz uma requisiçao com esses Endereços
Salva no MongoDB esses dados
Escreve em dois arquivos de log
	* LogIndiceEndereco, ultimo valor q foi salvo no mongoDB
	* logEnderecosSemDados, Enderecos que voltaram sem dado e nao foram salvo no MongoDb
*/
func SalvaListaEnderecos(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string,
	UrlAPI string, rota string, CollectionSalvaDados string,
	nomeArquivoSemApagar string, nomeArquivoIndice string) {

	listaEnderecos := RecuperaEnderecosUnionArrayRemoveDuplicados(ConnectionMongoDB, DataBase, CollectionRecuperaDados)

	indiceInicial := GetIndiceLogIndice(nomeArquivoIndice)

	for contador := indiceInicial; contador < len(listaEnderecos); contador++ {
		if len(listaEnderecos[contador]) > 0 {
			SalvarUnicoEndereco(listaEnderecos[contador], contador, UrlAPI, rota, ConnectionMongoDB, DataBase, CollectionSalvaDados,
				nomeArquivoSemApagar, nomeArquivoIndice)

		}
	}
}

func RecuperaEnderecosUnionArrayRemoveDuplicados(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string) []string {
	input, out := RecuperarEnderecos(ConnectionMongoDB, DataBase, CollectionRecuperaDados)
	return UnionArrayRemoveDuplicados(input, out)
}

func UnionArrayRemoveDuplicados(input []string, out []string) []string {
	return Arrays.RemoveDuplicados(Arrays.UnionArray(input, out))
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

/*Salva somente os enderecos no MongoDB*/
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

/*Salvar multiEnderecos no MongoDb*/
func SalvarMultiEndereco(endereco []string, indice int, UrlAPI string, MultiAddr string,
	ConnectionMongoDB string, DataBase string, Collection string,
	nomeArquivoSemApagar string, nomeArquivoIndice string) {
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

	} else {
		valorTemp := " Indice: " + strconv.Itoa(indice) + " Endereco: " + strings.Join(endereco, "|")
		temp := []string{valorTemp}

		EscreverTextoSemApagar(temp, nomeArquivoSemApagar)
		fmt.Println("Nao tem Dados o Indice: " + strconv.Itoa(indice) + " Endereco: " + strings.Join(endereco, "|"))

	}
}

/*Salvar Unicoendereco pegando o endereco da API*/
func SalvarUnicoEndereco(endereco string, indice int, UrlAPI string, RawAddr string,
	ConnectionMongoDB string, DataBase string, Collection string,
	nomeArquivoSemApagar string, nomeArquivoIndice string) {
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
		} else {
			valorTemp := " Indice: " + strconv.Itoa(indice) + " Endereco: " + endereco
			temp := []string{valorTemp}

			EscreverTextoSemApagar(temp, nomeArquivoSemApagar)
			fmt.Println("Nao tem Dados o Indice: " + strconv.Itoa(indice) + " Endereco: " + endereco)
		}

	}
}
