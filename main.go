package main

import (
	"main.go/Controllers"
	"main.go/Service"
)

var ConnectionMongoDB string = Controllers.GetConfig().ConnectionMongoDB[0] //"connection string into your application code"
var DataBaseBlockchain string = Controllers.GetConfig().DataBase[0]         //blockchain

var CollectionAdresses string = Controllers.GetConfig().Collection[0]           // "Adresses"
var CollectionBlockchain string = Controllers.GetConfig().Collection[1]         // "blockchain"
var CollectionEnderecos string = Controllers.GetConfig().Collection[2]          // "enderecos"
var CollectionTeste string = Controllers.GetConfig().Collection[3]              // "teste"
var CollectionTesteMultiAdress string = Controllers.GetConfig().Collection[4]   // testeMultiAdress
var CollectionMapeandoEnderecos string = Controllers.GetConfig().Collection[5]  //MapeandoEnderecos
var CollectionTesteMapeando string = Controllers.GetConfig().Collection[6]      // TesteMapeando
var CollectionGraph string = Controllers.GetConfig().Collection[7]              // Graph
var CollectionTesteGraph string = Controllers.GetConfig().Collection[8]         // TesteGraph
var CollectionEnderecosAgrupados string = Controllers.GetConfig().Collection[9] // EnderecosAgrupados

var UrlAPI string = Controllers.GetConfig().UrlAPI[0] // "https://blockchain.info"

var RawAddr string = Controllers.GetConfig().RawAddr
var MultiAddr string = Controllers.GetConfig().MultiAddr

var LogBlockchain string = Controllers.GetConfig().FileLog[0]             // "LogBlockchain.txt"
var LogIndiceEndereco string = Controllers.GetConfig().FileLog[1]         // "LogIndiceEndereco.txt"
var LogEndereco string = Controllers.GetConfig().FileLog[2]               // "LogEndereco.txt"
var LogEnderecosSemDados string = Controllers.GetConfig().FileLog[3]      // "LogEnderecosSemDados.txt"
var LogIndiceMultiEndereco string = Controllers.GetConfig().FileLog[4]    // LogIndiceMultiEndereco.txt
var LogMultiEnderecosSemDados string = Controllers.GetConfig().FileLog[5] //LogMultiEnderecosSemDados.txt

func main() {
	//fmt.Println("Nao tem nenhuma função sendo executada")

	//x := Service.AdicionaAddress("1","145YPBBWRj4aquewvx59SAWNrSZFT5rvxr",ConnectionMongoDB,DataBaseBlockchain,CollectionEnderecosAgrupados)
	//fmt.Println(x)
	//graph :=  Service.GetGraph(ConnectionMongoDB, DataBaseBlockchain,CollectionTesteGraph,false)
	//fmt.Println(graph.Visitado)
	//graph.Visitado = true;
	//fmt.Println(graph.Visitado)

	//fmt.Println(graph)
	//atualizado:= Service.AtualizaCampoVisitado("1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F",ConnectionMongoDB,DataBaseBlockchain,CollectionTesteGraph)
	//if atualizado == true{
	//	fmt.Println("Atualizado")
	//}else{
	//	fmt.Println("Deu ruim")
	//}
	Service.ClusteringAddressesV1(ConnectionMongoDB, DataBaseBlockchain, CollectionEnderecosAgrupados, CollectionTesteGraph, 10)

	//Service.AgrupamentoEnderecos(ConnectionMongoDB, DataBaseBlockchain,CollectionEnderecosAgrupados,CollectionTesteGraph, false)

	//graphs := Service.GetGraphs(ConnectionMongoDB, DataBaseBlockchain,CollectionTesteGraph,"nodos.nextaddresses","13adwKvLLpHdcYDh21FguCdJgKhaYP3Dse")
	//if len(graphs) > 0 {
	//	for _, item := range graphs {
	//		fmt.Println(item.Addresses)
	//	}
	//}else{
	//	fmt.Println("Elemento não encontrado")
	//}
	//e := Service.GetGraph(ConnectionMongoDB, DataBaseBlockchain, CollectionTesteGraph)
	//fmt.Println(e)

	//Controllers.RecuperaMultiAddrESalvaEmEnderecos(ConnectionMongoDB, DataBaseBlockchain,
	//	CollectionEnderecos, UrlAPI, MultiAddr,
	//	CollectionTesteMultiAdress,LogMultiEnderecosSemDados ,LogIndiceMultiEndereco)

	// Salva enderecos
	//Controllers.SalvaListaEnderecos(ConnectionMongoDB, DataBaseBlockchain, CollectionEnderecos, UrlAPI, RawAddr, CollectionAdresses, LogEnderecosSemDados,
	//	LogIndiceEndereco)

	// Pega os valores da Collection
	// retorna uma lista mapeada
	//x := Controllers.MapeandoEndereco(ConnectionMongoDB, DataBaseBlockchain, CollectionTeste)
	//// salva no mongoDB
	//Controllers.SalvarMapeamentoTransacaoMongoDB(x, ConnectionMongoDB, DataBaseBlockchain, CollectionTesteMapeando)

	// Busca um unico endereco
	//y := Controllers.GetMapAdressId(ConnectionMongoDB, DataBaseBlockchain, CollectionTesteMapeando, "_id", "61565cd7db7d9366829e0531")
	//fmt.Println(y)

	// Busca todos os enderecos mapeandos e seus objectId
	//d := Auxiliares.GetAllMapeandoAdress(ConnectionMongoDB, DataBaseBlockchain, CollectionTesteMapeando)
	//fmt.Println(d)

	// Verifica o se o Grafo esta vazio
	//v := Service.CheckEmptyGraph(ConnectionMongoDB,DataBaseBlockchain,CollectionTesteGraph)
	//fmt.Println(v)

	// Builder Graph
	// graph := Service.BuilderGraph(Auxiliares.GetAllMapeandoAdress,
	// 	ConnectionMongoDB, DataBaseBlockchain, CollectionTesteMapeando, CollectionTesteGraph, "addresses")
	// fmt.Println(graph)

	//// Retorna um endereco de TesteGraph
	//gra, _ :=Service.GetAddress(ConnectionMongoDB, DataBaseBlockchain, CollectionTesteGraph, "addresses", "13adwKvLLpHdcYDh21FguCdJgKhaYP3Dse")
	//fmt.Println(gra)

}
