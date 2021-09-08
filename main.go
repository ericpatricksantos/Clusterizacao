package main

import (
	"strconv"

	"main.go/Controllers"
	"main.go/Function/Arrays"
)

var ConnectionMongoDB string = Controllers.GetConfig().ConnectionMongoDB[0]
var DataBase string = Controllers.GetConfig().DataBase[0]
var Collection0 string = Controllers.GetConfig().Collection[3]
var Collection2 string = Controllers.GetConfig().Collection[2]

var UrlAPI string = Controllers.GetConfig().UrlAPI[0]

var rota string = Controllers.GetConfig().RawAddr

var Filelog string = Controllers.GetConfig().FileLog[1]
var LogEndereco string = Controllers.GetConfig().FileLog[2]
var FilelogSemDados string = Controllers.GetConfig().FileLog[3]

func main() {

	input, out := Controllers.RecuperarEnderecos(ConnectionMongoDB, DataBase, Collection2)

	x := Arrays.RemoveDuplicados(Arrays.UnionArray(input, out))

	for index, elem := range x {

		if len(elem) > 0 {
			Controllers.SalvarUnicoEndereco(elem, index, UrlAPI, rota, ConnectionMongoDB, DataBase, Collection0, FilelogSemDados)
			temp := []string{strconv.Itoa(index)}
			Controllers.EscreverTexto(temp, Filelog)
		}
	}

}
