package Model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
  Essa struct é usada para o retorno de MapeandoEnderecoTransacao
*/
type ReturnAddrMapTx struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Adresses    string             `json:"addresses" `  // Endereco Analisado
	EntradaAddr []InputAddr        `json:"EntradaAddr"` // Todos os address que aparecem no input do endereco analisado
	SaidaAddr   []OutputAddr       `json:"SaidaAddr"`   // Todos os address que aparecem no output do endereco analisado
}

type MapeandoEnderecoTransacao struct {
	Adresses    string       `json:"addresses"`   // Endereco Analisado
	EntradaAddr []InputAddr  `json:"EntradaAddr"` // Todos os address que aparecem no input do endereco analisado
	SaidaAddr   []OutputAddr `json:"SaidaAddr"`   // Todos os address que aparecem no output do endereco analisado
}

type InputAddr struct {
	Addr string `json:"Addr"`
	Qtd  int    `json:"Qtd"` // quantidade de vezes que esse Addr aparece no array de input do endereco analisado
}

type OutputAddr struct {
	Addr string `json:"Addr"`
	Qtd  int    `json:"Qtd"` // quantidade de vezes que esse Addr aparece no array de out do endereco analisado
}
