package Model

type MapeandoMultiEnderecoTransacao struct {
	Adresses    string       `json:"addresses"`   // Endereco Analisado
	EntradaHash []InputHash  `json:"EntradaHash"` // Todos os hashTransacao que o endereco analisado aparece no input
	SaidaHash   []OutputHash `json:"SaidaHash"`   // Todos os hashTransacao que o endereco analisado aparece no output
}

type InputHash struct {
	HashTransacao string `json:"HashTransacao"`
	Qtd           int    `json:"Qtd"` // quantidade de vezes que esse endereco analisado aparece no array de input desse hashtransacao
}

type OutputHash struct {
	HashTransacao string `json:"HashTransacao"`
	Qtd           int    `json:"Qtd"` // quantidade de vezes que esse endereco analisado aparece no array de out desse hashtransacao
}
