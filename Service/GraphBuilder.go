package Service

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2/bson"
	"main.go/Function/Auxiliares"
	"main.go/Function/Repository"
	Model "main.go/Models"
)

/*
	Busca todos os enderecos que estao na Collection Mapeando Endereco
	Procurar Se o Addr esta no Grafo
	Adicionar Addr ao grafo
	Salvar no MongoDb
*/
func BuilderGraph(EnderecoMapeados Auxiliares.IGetAllMapeandoAdress, ConnectionMongoDB string, DataBase string, CollectionMap string, CollectionGraph string, Key string) bool {
	// Busca todos os enderecos que estão na Collection Mapeando Enderecos
	for _, item := range EnderecoMapeados(ConnectionMongoDB, DataBase, CollectionMap) {
		var First Model.Graph
		var Node []Model.Nodo
		// Verifica se o grafo esta vazio
		if CheckEmptyGraph(ConnectionMongoDB, DataBase, CollectionGraph) {
			fmt.Println("-------------------------------------------------------------------------")
			fmt.Println("O grafo esta vazio, sendo assim será inserido o primeiro elemento do grafo")
			// Definindo o endereo analisado
			First.Addresses = item.Addresses
			First.Id = item.Id

			for _, itemGraph := range item.SaidaAddr {
				// Procura o endereco para adiciona o seu ObjectID e addr no Array de Nodo
				// Para sabermos para quais enderecos o endereco analisado esta enviando as moedas
				adressesMap := Auxiliares.GetMapAdressId(ConnectionMongoDB, DataBase, CollectionMap, "addresses", itemGraph.Addr)
				// cria variavel temporaria para preencher os valores necessarios
				var temp Model.Nodo
				temp.MongoId = adressesMap.Id
				temp.NextAddresses = itemGraph.Addr
				temp.Qtd = itemGraph.Qtd
				// adiciona no array de Node contem o proximo endereco e o seu MongoId
				Node = append(Node, temp)
			}
			First.Nodos = Node
			// Salva o graph
			confirm := SalvaGraph(First, ConnectionMongoDB, DataBase, CollectionGraph)
			if confirm {
				fmt.Println("Graph criado com Sucesso")
			}
			fmt.Println("--------------------------------------------------------------------")
		} else {
			fmt.Println("--------------------------------------------------------------------")
			fmt.Println("O grafo não está vazio, sendo assim será verificado se o possível candidato já foi inserifo no grafo.")
			//fmt.Println(" Se ele foi inserido não será feita a inserção, caso contrário será ")
			//fmt.Println("inserido no grafo")
			Count, err := GetCountElem(ConnectionMongoDB, DataBase, CollectionGraph, "addresses", item.Addresses)
			if err != nil {
				fmt.Println(err)
			}
			if Count > 0 {
				fmt.Println("Esse endereço já foi inserido no Grafo")
			} else {
				First.Addresses = item.Addresses
				First.Id = item.Id
				for _, itemGraph := range item.SaidaAddr {
					fmt.Println("----------------------------------------------------------------------------")
					fmt.Println("Verificando se o endereco tem na Collection Mapeando Endereco")
					Count2, _ := GetCountElem(ConnectionMongoDB, DataBase, CollectionMap, "addresses", itemGraph.Addr)
					/*
						Em vez de consultar um por um. Pegar todos os enderecos que tiver no MOngoDb e colocar
						em um array estatico para o program ser mais eficiente
					*/
					if Count2 > 0 {
						fmt.Println("Adicionando todos os valores na variavel temporaria")
						adressesMap := Auxiliares.GetMapAdressId(ConnectionMongoDB, DataBase, CollectionMap, "addresses", itemGraph.Addr)
						// cria variavel temporaria para preencher os valores necessarios
						var temp Model.Nodo
						temp.MongoId = adressesMap.Id
						temp.NextAddresses = itemGraph.Addr
						temp.Qtd = itemGraph.Qtd
						// adiciona no array de Node contem o proximo endereco e o seu MongoId
						Node = append(Node, temp)
					} /*else{
						fmt.Println("Caso contrário so adiciona somente o addr")
						// cria variavel temporaria para preencher os valores necessarios
						var temp Model.Nodo
						temp.NextAddresses = itemGraph.Addr
						// adiciona no array de Node contem o proximo endereco e o seu MongoId
						Node = append(Node, temp)
					}*/
				}
				First.Nodos = Node
				// Salva o graph
				confirm := SalvaGraph(First, ConnectionMongoDB, DataBase, CollectionGraph)
				if confirm {
					fmt.Println("Elemento inserido com Sucesso")
				}
				fmt.Println("--------------------------------------------------------------------")
			}

		}

	}

	return true
}

// Salva Graph
func SalvaGraph(graph Model.Graph, ConnectionMongoDB string, DataBase string, Collection string) bool {
	cliente, contexto, cancel, errou := Repository.Connect(ConnectionMongoDB)
	if errou != nil {
		log.Fatal(errou)
	}

	er := Repository.Ping(cliente, contexto)

	if er != nil {
		log.Fatal(er)
	}

	defer Repository.Close(cliente, contexto, cancel)

	Repository.ToDoc(graph)

	insertOneResult, err := Repository.InsertOne(cliente, contexto, DataBase, Collection, graph)

	// handle the error
	if err != nil {
		fmt.Println("O graph nao foi salvo")
		panic(err)
	}

	// print the insertion id of the document,
	// if it is inserted.
	fmt.Println("O graph foi salvo com sucesso")
	fmt.Println(insertOneResult.InsertedID)

	return true
}

/*
	Função retorna a quantidade de vezes que um elemento esta em uma collection
*/
func GetCountElem(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string, Key string, Code string) (Count int64, err error) {
	// Get Client, Context, CalcelFunc and err from connect method.
	client, ctx, cancel, err := Repository.Connect(ConnectionMongoDB)
	if err != nil {
		panic(err)
	}

	// Free the resource when mainn dunction is  returned
	defer Repository.Close(client, ctx, cancel)

	Count, err = Repository.CountElemento(client, ctx, DataBase,
		CollectionRecuperaDados, Key, Code)
	// handle the errors.
	if err != nil {
		panic(err)
	}
	return Count, err
}

/*
Busca um único graph atráves de uma chave e um valor
	Exemplo:
			Key = _id , Code = "6153a58d3700e70e40f8177a"
			Key = addresses , Code = "13adwKvLLpHdcYDh21FguCdJgKhaYP3Dse"
*/
func GetAddress(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string, Key string, Code string) (graph Model.Graph, err error) {
	// Get Client, Context, CalcelFunc and err from connect method.
	client, ctx, cancel, err := Repository.Connect(ConnectionMongoDB)
	if err != nil {
		panic(err)
	}

	// Free the resource when mainn dunction is  returned
	defer Repository.Close(client, ctx, cancel)

	graph, err = Repository.QueryOneGraph(client, ctx, DataBase,
		CollectionRecuperaDados, Key, Code)
	// handle the errors.
	if err != nil {
		panic(err)
	}
	return graph, err
}

// Verifica se o grafo esta vazio
func CheckEmptyGraph(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string) bool {
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
	var graphs []Model.Graph
	// le os documentos em partes, testei com 1000 documentos e deu bom
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var graph Model.Graph

		if err := cursor.Decode(&graph); err != nil {
			log.Fatal(err)
		}

		graphs = append(graphs, graph)

		if len(graphs) > 0 {
			return false
		}

	}
	return true
}
