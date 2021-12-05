package Service

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"log"
	"main.go/Function/Repository"
	Model "main.go/Models"
)

func ClusteringAddressesV2(ConnectionMongoDB string, DataBaseBlockchain string, CollectionEnderecoAgrup string, CollectionGraph string,
	limit int) bool {

	var enderecosAgrupados Model.EnderecosAgrupados

	elementoAnalisado := GetGraph(ConnectionMongoDB, DataBaseBlockchain, CollectionGraph, false)

	nodosAnalisados := OrdenaCrescenteNodo(elementoAnalisado.Nodos)

	// Retorna os objetos que possuem o addresses do enderecoAnalisado nos nodos
	graphElementoAnalisado := GetGraphs(ConnectionMongoDB, DataBaseBlockchain, CollectionGraph, "nodos.nextaddresses", elementoAnalisado.Addresses, false)
	listaEnderecos := ConvertListGraphEmListString(graphElementoAnalisado)

	if len(listaEnderecos) > 0 {
		for i := 0; i < limit; i++ {
			for _, item := range elementoAnalisado.Nodos {
				if listaEnderecos[i] == item.NextAddresses {
					enderecosAgrupados.Addresses = append(enderecosAgrupados.Addresses, item.NextAddresses)
				} else {
					graphElementoAnalisado2 := GetGraphUnico(ConnectionMongoDB, DataBaseBlockchain, CollectionGraph, "addresses", item.NextAddresses)
					if len(graphElementoAnalisado2.Addresses) > 0 {
						// se dentro desse graphElementoAnalisado2 tiver uns dos elementos do listaEnderecos
						// foi formado a relacao
					}
				}
			}
		}
	}
	for index := 0; index < limit; index++ {

		NextAddress := nodosAnalisados[index].NextAddresses
		elementoApontado := GetGraphUnico(ConnectionMongoDB, DataBaseBlockchain, CollectionGraph, "addresses", NextAddress)

		// Buscando o endereco do elementoAnalisado dentro do array de nodos do elementoApontado
		for _, itemElementoApontado := range elementoApontado.Nodos {

			if elementoAnalisado.Addresses == itemElementoApontado.NextAddresses {
				enderecosAgrupados.Addresses = append(enderecosAgrupados.Addresses, elementoApontado.Addresses)

				// Atualiza o campo visitado(dentro do array) do ElementoApontado para true
				confirmElementoApontado := AtualizaCampoVisitado(elementoApontado.Addresses, itemElementoApontado.NextAddresses, ConnectionMongoDB, DataBaseBlockchain, CollectionGraph)
				if confirmElementoApontado {
					fmt.Println("Atualizado Address Apontado: ", elementoApontado.Addresses, " NextAdresses Apontado: ", itemElementoApontado.NextAddresses)
				} else {
					fmt.Println("Não foi Atualizado Address Apontado: ", elementoApontado.Addresses, " NextAdresses Apontado: ", itemElementoApontado.NextAddresses)
				}

				// atualiza o campo visitado(dentro do array) do ElementoAnalisado para true
				confirmElementoAnalisado := AtualizaCampoVisitado(elementoAnalisado.Addresses, NextAddress, ConnectionMongoDB, DataBaseBlockchain, CollectionGraph)
				if confirmElementoAnalisado {
					fmt.Println("Atualizado Address Analisado: ", elementoAnalisado.Addresses, " NextAdresses Analisado: ", NextAddress)
				} else {
					fmt.Println("Não foi Atualizado Address Analisado: ", elementoAnalisado.Addresses, " NextAdresses Analisado: ", NextAddress)
				}

			} else {
				// Se o itemElementoApontado.NextAddresses != EnderecoAnalisado.Addresses,
				// procurar no Nodo do   itemElementoApontado.NextAddresses o EnderecoAnalisado.Addresses
				NextAddressSegundoNivel := itemElementoApontado.NextAddresses
				elementoApontadoSegundoNivel := GetGraphUnico(ConnectionMongoDB, DataBaseBlockchain, CollectionGraph, "addresses", NextAddressSegundoNivel)

				for _, itemElementoApontadoSegundoNivel := range elementoApontadoSegundoNivel.Nodos {
					if elementoAnalisado.Addresses == itemElementoApontadoSegundoNivel.NextAddresses {
						enderecosAgrupados.Addresses = append(enderecosAgrupados.Addresses, itemElementoApontadoSegundoNivel.NextAddresses)

						// Atualiza o campo visitado(dentro do array) do ElementoApontadoSegundoNivel para true
						confirmElementoApontadoSegundoNivel := AtualizaCampoVisitado(elementoApontadoSegundoNivel.Addresses, itemElementoApontadoSegundoNivel.NextAddresses, ConnectionMongoDB, DataBaseBlockchain, CollectionGraph)
						if confirmElementoApontadoSegundoNivel {
							fmt.Println("Atualizado Address Apontado: ", elementoApontadoSegundoNivel.Addresses, " NextAdresses Apontado: ", itemElementoApontadoSegundoNivel.NextAddresses)
						} else {
							fmt.Println("Não foi Atualizado Address Apontado: ", elementoApontadoSegundoNivel.Addresses, " NextAdresses Apontado: ", itemElementoApontadoSegundoNivel.NextAddresses)
						}

						// Atualiza o campo visitado(dentro do array) do ElementoApontado para true
						confirmElementoApontado := AtualizaCampoVisitado(elementoApontado.Addresses, itemElementoApontado.NextAddresses, ConnectionMongoDB, DataBaseBlockchain, CollectionGraph)
						if confirmElementoApontado {
							fmt.Println("Atualizado Address Apontado: ", elementoApontado.Addresses, " NextAdresses Apontado: ", NextAddressSegundoNivel)
						} else {
							fmt.Println("Não foi Atualizado Address Apontado: ", elementoApontado.Addresses, " NextAdresses Apontado: ", NextAddressSegundoNivel)
						}
					}
				}
			}
		}
	}

	if len(enderecosAgrupados.Addresses) > 0 {
		enderecosAgrupados.Addresses = append(enderecosAgrupados.Addresses, elementoAnalisado.Addresses)
	}

	salvou := false
	if len(enderecosAgrupados.Addresses) > 0 {
		salvou = SalvaEnderecosAgrupados(enderecosAgrupados, ConnectionMongoDB, DataBaseBlockchain, CollectionEnderecoAgrup)
	}
	if salvou {
		fmt.Println("Salvo endereco agrupado com sucesso")
		return true
	} else {
		fmt.Println("Nao foi salvo")
		return false
	}
}

func ClusteringAddressesV1(ConnectionMongoDB string, DataBaseBlockchain string, CollectionEnderecoAgrup string, CollectionGraph string,
	limit int) bool {

	var enderecosAgrupados Model.EnderecosAgrupados

	elementoAnalisado := GetGraph(ConnectionMongoDB, DataBaseBlockchain, CollectionGraph, false)

	nodosAnalisados := OrdenaCrescenteNodo(elementoAnalisado.Nodos)

	for index := 0; index < limit; index++ {

		NextAddress := nodosAnalisados[index].NextAddresses
		elementoApontado := GetGraphUnico(ConnectionMongoDB, DataBaseBlockchain, CollectionGraph, "addresses", NextAddress)

		// Buscando o endereco do elementoAnalisado dentro do array de nodos do elementoApontado
		for _, itemElementoApontado := range elementoApontado.Nodos {

			if elementoAnalisado.Addresses == itemElementoApontado.NextAddresses {
				enderecosAgrupados.Addresses = append(enderecosAgrupados.Addresses, elementoApontado.Addresses)

				// Atualiza o campo visitado(dentro do array) do ElementoApontado para true
				confirmElementoApontado := AtualizaCampoVisitado(elementoApontado.Addresses, itemElementoApontado.NextAddresses, ConnectionMongoDB, DataBaseBlockchain, CollectionGraph)
				if confirmElementoApontado {
					fmt.Println("Atualizado Address Apontado: ", elementoApontado.Addresses, " NextAdresses Apontado: ", itemElementoApontado.NextAddresses)
				} else {
					fmt.Println("Não foi Atualizado Address Apontado: ", elementoApontado.Addresses, " NextAdresses Apontado: ", itemElementoApontado.NextAddresses)
				}

				// atualiza o campo visitado(dentro do array) do ElementoAnalisado para true
				confirmElementoAnalisado := AtualizaCampoVisitado(elementoAnalisado.Addresses, NextAddress, ConnectionMongoDB, DataBaseBlockchain, CollectionGraph)
				if confirmElementoAnalisado {
					fmt.Println("Atualizado Address Analisado: ", elementoAnalisado.Addresses, " NextAdresses Analisado: ", NextAddress)
				} else {
					fmt.Println("Não foi Atualizado Address Analisado: ", elementoAnalisado.Addresses, " NextAdresses Analisado: ", NextAddress)
				}

			} else {
				// Se o itemElementoApontado.NextAddresses != EnderecoAnalisado.Addresses,
				// procurar no Nodo do   itemElementoApontado.NextAddresses o EnderecoAnalisado.Addresses
				NextAddressSegundoNivel := itemElementoApontado.NextAddresses
				elementoApontadoSegundoNivel := GetGraphUnico(ConnectionMongoDB, DataBaseBlockchain, CollectionGraph, "addresses", NextAddressSegundoNivel)

				for _, itemElementoApontadoSegundoNivel := range elementoApontadoSegundoNivel.Nodos {
					if elementoAnalisado.Addresses == itemElementoApontadoSegundoNivel.NextAddresses {
						enderecosAgrupados.Addresses = append(enderecosAgrupados.Addresses, itemElementoApontadoSegundoNivel.NextAddresses)

						// Atualiza o campo visitado(dentro do array) do ElementoApontadoSegundoNivel para true
						confirmElementoApontadoSegundoNivel := AtualizaCampoVisitado(elementoApontadoSegundoNivel.Addresses, itemElementoApontadoSegundoNivel.NextAddresses, ConnectionMongoDB, DataBaseBlockchain, CollectionGraph)
						if confirmElementoApontadoSegundoNivel {
							fmt.Println("Atualizado Address Apontado: ", elementoApontadoSegundoNivel.Addresses, " NextAdresses Apontado: ", itemElementoApontadoSegundoNivel.NextAddresses)
						} else {
							fmt.Println("Não foi Atualizado Address Apontado: ", elementoApontadoSegundoNivel.Addresses, " NextAdresses Apontado: ", itemElementoApontadoSegundoNivel.NextAddresses)
						}

						// Atualiza o campo visitado(dentro do array) do ElementoApontado para true
						confirmElementoApontado := AtualizaCampoVisitado(elementoApontado.Addresses, itemElementoApontado.NextAddresses, ConnectionMongoDB, DataBaseBlockchain, CollectionGraph)
						if confirmElementoApontado {
							fmt.Println("Atualizado Address Apontado: ", elementoApontado.Addresses, " NextAdresses Apontado: ", NextAddressSegundoNivel)
						} else {
							fmt.Println("Não foi Atualizado Address Apontado: ", elementoApontado.Addresses, " NextAdresses Apontado: ", NextAddressSegundoNivel)
						}
					}
				}
			}
		}
	}

	if len(enderecosAgrupados.Addresses) > 0 {
		enderecosAgrupados.Addresses = append(enderecosAgrupados.Addresses, elementoAnalisado.Addresses)
	}

	salvou := false
	if len(enderecosAgrupados.Addresses) > 0 {
		salvou = SalvaEnderecosAgrupados(enderecosAgrupados, ConnectionMongoDB, DataBaseBlockchain, CollectionEnderecoAgrup)
	}
	if salvou {
		fmt.Println("Salvo endereco agrupado com sucesso")
		return true
	} else {
		fmt.Println("Nao foi salvo")
		return false
	}
}
func OrdenaCrescenteNodo(nodos []Model.Nodo) (retorno []Model.Nodo) {
	if len(nodos) > 0 {

	}
	return retorno
}

func ConvertListGraphEmListString(graphs []Model.Graph) (listaStrings []string) {

	if len(graphs) > 0 {
		for _, item := range graphs {
			listaStrings = append(listaStrings, item.Addresses)
		}
	}

	return listaStrings
}

/*
	Realiza o Agrupamento dos Elementos que estão no Graph
*/
func AgrupamentoEnderecos(ConnectionMongoDB string, DataBaseBlockchain string, CollectionEnderecoAgrup string, CollectionGraph string, visitado bool) bool {
	fmt.Println("Consultando um elemento do Graph para iniciar o Agrupamento de Endereços")
	// Retorna um unico elemento para ser agrupado
	graph := GetGraph(ConnectionMongoDB, DataBaseBlockchain, CollectionGraph, visitado)

	var enderecosAgrupados Model.EnderecosAgrupados

	fmt.Println("Analisando o elemento do Graph")

	// Retorna varios elementos que tem o endereço do graph nodo.nextaddresses
	graphs := GetGraphs(ConnectionMongoDB, DataBaseBlockchain, CollectionGraph, "nodos.nextaddresses", graph.Addresses, false)
	if len(graphs) > 0 {
		for _, item := range graphs {

			for _, itemGraph := range graph.Nodos {
				if itemGraph.NextAddresses == item.Addresses && itemGraph.Visitado == false {
					enderecosAgrupados.Addresses = append(enderecosAgrupados.Addresses, item.Addresses)
					// atualiza o campo visitado do array
					atualizado1 := AtualizaCampoVisitado(graph.Addresses, itemGraph.NextAddresses, ConnectionMongoDB, DataBaseBlockchain, CollectionGraph)
					atualizado2 := AtualizaCampoVisitado(item.Addresses, graph.Addresses, ConnectionMongoDB, DataBaseBlockchain, CollectionGraph)
					if atualizado1 == true && atualizado2 == true {
						fmt.Println("Atualizado")
					} else {
						fmt.Println("Deu ruim")
					}
				} else {
					fmt.Println("procurando o objeto que esta sendo referenciado pelo elemento do array")
					graphsaux := GetGraphUnico(ConnectionMongoDB, DataBaseBlockchain, CollectionGraph, "addresses", itemGraph.NextAddresses)
					fmt.Println("verifica se esse elemento esta dentro do array de graphs")
					for _, item2 := range item.Nodos {
						if graphsaux.Addresses == item2.NextAddresses && item2.Visitado == false {
							enderecosAgrupados.Addresses = append(enderecosAgrupados.Addresses, graphsaux.Addresses)
							atualizado1 := AtualizaCampoVisitado(graph.Addresses, item2.NextAddresses, ConnectionMongoDB, DataBaseBlockchain, CollectionGraph)
							atualizado2 := AtualizaCampoVisitado(item.Addresses, item.Addresses, ConnectionMongoDB, DataBaseBlockchain, CollectionGraph)
							if atualizado1 == true && atualizado2 == true {
								fmt.Println("Atualizado")
							} else {
								fmt.Println("Deu ruim")
							}
						}
					}
				}
			}
		}

	} else {
		fmt.Println("Elemento não encontrado")
	}

	if len(enderecosAgrupados.Addresses) > 0 {
		enderecosAgrupados.Addresses = append(enderecosAgrupados.Addresses, graph.Addresses)
	}

	salvou := false
	if len(enderecosAgrupados.Addresses) > 0 {
		salvou = SalvaEnderecosAgrupados(enderecosAgrupados, ConnectionMongoDB, DataBaseBlockchain, CollectionEnderecoAgrup)
	}
	if salvou {
		fmt.Println("Salvo endereco agrupado com sucesso")
		return true
	} else {
		fmt.Println("Nao foi salvo")
		return false
	}
}

// Atualiza o campo visitado fora do array
func PutVisitado(valorAtualizado Model.Graph, ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string) bool {
	client, ctx, cancel, err := Repository.Connect(ConnectionMongoDB)
	if err != nil {
		panic(err)
	}

	defer Repository.Close(client, ctx, cancel)

	var filter, update interface{}

	filter = bson.M{"_id": valorAtualizado.Id}

	update = bson.M{
		"$set": valorAtualizado,
	}

	sucesso, err := Repository.UpdateOne(client, ctx, DataBase,
		CollectionRecuperaDados, filter, update)
	if err != nil {
		panic(err)
	}

	if sucesso.ModifiedCount == 1 {
		return true
	} else {
		return false
	}
}

/*
	Atualiza o campo Visitado do array par true
*/
func AtualizaCampoVisitado(addresses string, nextaddresses string, ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string) bool {
	client, ctx, cancel, err := Repository.Connect(ConnectionMongoDB)
	if err != nil {
		panic(err)
	}

	defer Repository.Close(client, ctx, cancel)

	var filter, update interface{}

	filter = bson.M{
		"addresses":           addresses,
		"nodos.nextaddresses": nextaddresses,
	}
	update = bson.M{"$set": bson.M{"nodos.$.visitado": true}}

	sucesso, err := Repository.UpdateOne(client, ctx, DataBase,
		CollectionRecuperaDados, filter, update)
	if err != nil {
		panic(err)
	}

	if sucesso.ModifiedCount == 1 {
		return true
	} else {
		return false
	}
}

/*
	Adiciona novo endereco no array de addresses
*/
func AdicionaAddress(etiqueta string, address string, ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string) bool {
	client, ctx, cancel, err := Repository.Connect(ConnectionMongoDB)
	if err != nil {
		panic(err)
	}

	defer Repository.Close(client, ctx, cancel)

	var filter, update interface{}

	filter = bson.M{
		"etiqueta": etiqueta,
	}
	// O operador $addToSet adiciona novo valor ao array se ele nao existir
	update = bson.M{"$addToSet": bson.M{"addresses": address}}

	opts := options.Update().SetUpsert(true)

	sucesso, err := Repository.AdicionaAddress(client, ctx, DataBase,
		CollectionRecuperaDados, filter, update, opts)
	if err != nil {
		panic(err)
	}

	if sucesso.ModifiedCount == 1 {
		return true
	} else {
		return false
	}
}

/*
	Retorna elementos que satisfaça uma condição
*/
func GetGraphs(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string, chave string, valor string, visitado bool) (graphs []Model.Graph) {
	client, ctx, cancel, err := Repository.Connect(ConnectionMongoDB)
	if err != nil {
		panic(err)
	}

	defer Repository.Close(client, ctx, cancel)

	var filter, option interface{}

	filter = bson.M{chave: valor, "visitado": visitado}

	option = bson.M{}

	cursor, err := Repository.Query(client, ctx, DataBase,
		CollectionRecuperaDados, filter, option)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var graph Model.Graph

		if err := cursor.Decode(&graph); err != nil {
			log.Fatal(err)
		}

		graphs = append(graphs, graph)

	}

	return graphs
}

/*
	Retorna um unico elemento do Graph atraves de key e value
*/
func GetGraphUnico(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string, key string, code string) Model.Graph {
	// Get Client, Context, CalcelFunc and err from connect method.
	client, ctx, cancel, err := Repository.Connect(ConnectionMongoDB)
	if err != nil {
		panic(err)
	}

	// Free the resource when mainn dunction is  returned
	defer Repository.Close(client, ctx, cancel)

	// call the query method with client, context,
	// database name, collection  name, filter and option
	// This method returns momngo.cursor and error if any.
	graph, err := Repository.QueryOneGraph(client, ctx, DataBase,
		CollectionRecuperaDados, key, code)
	// handle the errors.
	if err != nil {
		panic(err)
	}
	return graph
}

/*
	Retorna um unico elemento do Graph
*/
func GetGraph(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string, visitado bool) Model.Graph {
	// Get Client, Context, CalcelFunc and err from connect method.
	client, ctx, cancel, err := Repository.Connect(ConnectionMongoDB)
	if err != nil {
		panic(err)
	}

	// Free the resource when mainn dunction is  returned
	defer Repository.Close(client, ctx, cancel)

	// call the query method with client, context,
	// database name, collection  name, filter and option
	// This method returns momngo.cursor and error if any.
	graph, err := Repository.SearchOne(client, ctx, DataBase,
		CollectionRecuperaDados, visitado)
	// handle the errors.
	if err != nil {
		panic(err)
	}
	return graph
}

/*
	Salva os Endereços Agrupados
*/
func SalvaEnderecosAgrupados(enderecosAgrupados Model.EnderecosAgrupados, ConnectionMongoDB string, DataBase string, Collection string) bool {
	cliente, contexto, cancel, errou := Repository.Connect(ConnectionMongoDB)
	if errou != nil {
		log.Fatal(errou)
	}

	er := Repository.Ping(cliente, contexto)

	if er != nil {
		log.Fatal(er)
	}

	defer Repository.Close(cliente, contexto, cancel)

	Repository.ToDoc(enderecosAgrupados)

	insertOneResult, err := Repository.InsertOne(cliente, contexto, DataBase, Collection, enderecosAgrupados)

	// handle the error
	if err != nil {
		fmt.Println("O enderecoAgrupado nao foi salvo")
		panic(err)
	}

	// print the insertion id of the document,
	// if it is inserted.
	fmt.Println("O enderecoAgrupado foi salvo com sucesso")
	fmt.Println(insertOneResult.InsertedID)

	return true
}
