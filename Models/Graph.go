package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Graph struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Addresses string             `json:"Addresses"`
	Visitado  bool               `json:"Visitado"` // Verifica se todos os elementos do Nodo foram visitados
	Nodos     []Nodo             `json:"Nodos"`
}

type Nodo struct {
	NextAddresses string             `json:"NextAddresses"`
	Qtd           int                `json:"Qtd"`
	MongoId       primitive.ObjectID `json:"MongoId"`
	Visitado      bool               `json:"Visitado"`
}
