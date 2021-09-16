package main

import (
	"main.go/Controllers"
)

var ConnectionMongoDB string = Controllers.GetConfig().ConnectionMongoDB[0] //"connection string into your application code"
var DataBaseBlockchain string = Controllers.GetConfig().DataBase[0]         //blockchain

var CollectionAdresses string = Controllers.GetConfig().Collection[0] // "Adresses"
//var CollectionBlockchain string = Controllers.GetConfig().Collection[1] // "blockchain"
var CollectionEnderecos string = Controllers.GetConfig().Collection[2]         // "enderecos"
var CollectionTeste string = Controllers.GetConfig().Collection[3]             // "teste"
var CollectionTesteMultiAdress string = Controllers.GetConfig().Collection[4]  // testeMultiAdress
var CollectionMapeandoEnderecos string = Controllers.GetConfig().Collection[5] //MapeandoEnderecos

var UrlAPI string = Controllers.GetConfig().UrlAPI[0] // "https://blockchain.info"

//var rota string = Controllers.GetConfig().RawAddr

var LogBlockchain string = Controllers.GetConfig().FileLog[0]             // "LogBlockchain.txt"
var LogIndiceEndereco string = Controllers.GetConfig().FileLog[1]         // "LogIndiceEndereco.txt"
var LogEndereco string = Controllers.GetConfig().FileLog[2]               // "LogEndereco.txt"
var LogEnderecosSemDados string = Controllers.GetConfig().FileLog[3]      // "LogEnderecosSemDados.txt"
var LogIndiceMultiEndereco string = Controllers.GetConfig().FileLog[4]    // LogIndiceMultiEndereco.txt
var LogMultiEnderecosSemDados string = Controllers.GetConfig().FileLog[5] //LogMultiEnderecosSemDados.txt

func main() {
	// Controllers.SalvaListaEnderecos(ConnectionMongoDB, DataBaseBlockchain, CollectionEnderecos, UrlAPI, rota, CollectionAdresses, LogEnderecosSemDados,
	// 	LogIndiceEndereco)

	x := Controllers.MapeandoEndereco(ConnectionMongoDB, DataBaseBlockchain, CollectionTeste)

	Controllers.SalvarMapeamentoTransacaoMongoDB(x, ConnectionMongoDB, DataBaseBlockchain, CollectionMapeandoEnderecos)

}
