package Model

type Graph struct {
	Addresses string `json:"addresses"`
	Nodos     []Nodo `json:"Nodos"`
}

type Nodo struct {
	NextAddresses string `json:"NextAddresses"`
	MongoId       string `json:"MongoId"`
}
