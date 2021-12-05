package Controllers

import (
	"fmt"
	"log"
	"strconv"

	"main.go/Function/Auxiliares"

	"gopkg.in/mgo.v2/bson"
	"main.go/Function/API"
	"main.go/Function/Repository"
	Model "main.go/Models"
)

/*
	Recupera Multi Endereços
	Converter esses Multi Endereços em Endereços Unicos
	Salva os Enderecos Unicos no MongoDbs
*/
func RecuperaMultiAddrESalvaEmEnderecos(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string,
	UrlAPI string, MultiAddr string, CollectionSalvaDados string,
	nomeArquivoSemApagar string, nomeArquivoIndice string) {

	listaEnderecos := Auxiliares.RecuperaEnderecosUnionArrayRemoveDuplicados(ConnectionMongoDB, DataBase, CollectionRecuperaDados)

	indiceInicial := GetIndiceLogIndice(nomeArquivoIndice) + 1

	for contador := indiceInicial; contador < len(listaEnderecos); contador++ {
		if len(listaEnderecos[contador]) > 0 {
			fmt.Printf("\n")
			fmt.Printf("-------------------------------------------------------------------------------")

			fmt.Printf("\nSalvando %dº endereço %s \n", contador, listaEnderecos[contador])
			addr := []string{listaEnderecos[contador]}
			multiAddr := Auxiliares.GetMultiEndereco(addr, UrlAPI, MultiAddr)
			address := Auxiliares.ConverteMultiAddrEMAddr(multiAddr)
			salvou := Auxiliares.SalvaEnderecoMongoDb(address, contador, ConnectionMongoDB, DataBase, CollectionSalvaDados,
				nomeArquivoSemApagar, nomeArquivoIndice)

			if !salvou {
				break
			}
		}
	}
}

/*
	Apos chamar a funçao MapeandoEndereco que retorna []Model.MapeandoEnderecoTransacao
	Use SalvarMapeamentoTransacaoMongoDB para salva um []Model.MapeandoEnderecoTransacao no MongoDb
*/
func SalvarMapeamentoTransacaoMongoDB(obj []Model.MapeandoEnderecoTransacao, ConnectionMongoDB string, DataBase string, Collection string) {
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
	Recupera todos os valores de Enderecos salvados no MongoDb
	Mapeia os enderecos que aparecem no input e output e quantidade de ocorrencia de um endereco analisado
	em um endereco Unico
	Retorna uma lista de objetos(MapeandoEnderecoTransação)
*/
func MapeandoEndereco(ConnectionMongoDB string, DataBase string, Collection string) []Model.MapeandoEnderecoTransacao {
	var addressesMaping []Model.MapeandoEnderecoTransacao

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

	contadorInput := 0
	contadorOut := 0
	for cursor.Next(ctx) {
		var enderecos Model.UnicoEndereco

		if err := cursor.Decode(&enderecos); err != nil {
			log.Fatal(err)
		}

		// Definindo variáveis temporárias para atribuir os valores do endereco analisado e seus addr
		var temp Model.MapeandoEnderecoTransacao
		// Inicializa os dois array com seus addr e seus qtd = 0
		tempInput, tempOutput := Auxiliares.Inicializa(enderecos.Address, enderecos.Txs)

		// Definindo o endereco analisado
		temp.Addresses = enderecos.Address
		for {

			if contadorInput == len(tempInput) {
				contadorInput = 0
				break
			} else {
				for _, item := range enderecos.Txs {
					if contadorInput < len(tempInput) {
						for _, inp := range item.Inputs {

							if tempInput[contadorInput].Addr == inp.Prev_out.Addr {
								tempInput[contadorInput].Qtd = tempInput[contadorInput].Qtd + 1
							}

						}
					}

				}
				contadorInput = contadorInput + 1
			}

		}
		for {
			if contadorOut == len(tempOutput) {
				contadorOut = 0
				break
			} else {
				for _, item := range enderecos.Txs {
					if contadorOut < len(tempOutput) {
						for _, out := range item.Out {
							if tempOutput[contadorOut].Addr == out.Addr {
								tempOutput[contadorOut].Qtd = tempOutput[contadorOut].Qtd + 1
							}
						}
					}
				}
				contadorOut = contadorOut + 1
			}

		}
		// Atribuindo os valores da variavel temporaria
		temp.EntradaAddr = tempInput
		temp.SaidaAddr = tempOutput

		// Atribuindo os objetos que será retornado nessa função
		addressesMaping = append(addressesMaping, temp)
	}

	return addressesMaping
}

func RecuperarAllEnderecos(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string) []string {
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
		CollectionRecuperaDados, filter, option)
	// handle the errors.
	if err != nil {
		panic(err)
	}

	// le os documentos em partes, testei com 1000 documentos e deu certo
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var enderecos Model.UnicoEndereco

		if err := cursor.Decode(&enderecos); err != nil {
			log.Fatal(err)
		}

		addresses = append(addresses, enderecos.Address)

	}

	return addresses
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
	UrlAPI string, RawAddr string, CollectionSalvaDados string,
	nomeArquivoSemApagar string, nomeArquivoIndice string) {

	listaEnderecos := Auxiliares.RecuperaEnderecosUnionArrayRemoveDuplicados(ConnectionMongoDB, DataBase, CollectionRecuperaDados)

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
func RecuperarEnderecosInputOut(ConnectionMongoDB string, DataBase string, Collection string) ([]string, []string) {
	return Auxiliares.RecuperarEnderecosInputOut(ConnectionMongoDB, DataBase, Collection)
}

/*Salva UnicoEndereco buscando o endereco da API*/
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
