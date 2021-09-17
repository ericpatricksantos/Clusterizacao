package Controllers

import "main.go/Function/Arrays"

func RecuperaEnderecosUnionArrayRemoveDuplicados(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string) []string {
	input, out := RecuperarEnderecos(ConnectionMongoDB, DataBase, CollectionRecuperaDados)
	return UnionArrayRemoveDuplicados(input, out)
}

func UnionArrayRemoveDuplicados(input []string, out []string) []string {
	return Arrays.RemoveDuplicados(Arrays.UnionArray(input, out))
}

func RemoveDuplicados(lista []string) []string {
	return Arrays.RemoveDuplicados(lista)
}

func TransformaArrayEmMatriz(listaEnderecos []string, indiceInicial int, limiteEnderecos int) ([][]string, int) {
	QtdLinhas := (len(listaEnderecos) / 100) + 1
	limite := limiteEnderecos //Qtd de Colunas na Matriz
	matrizEnderecos := make([][]string, QtdLinhas)
	if indiceInicial > 0 {
		//Remove os elementos que foram salvos na matriz
		listaEnderecos = append(listaEnderecos[:0], listaEnderecos[indiceInicial:]...)

	}

	for contador := 0; contador < QtdLinhas; contador++ {

		// Recebe uma fatia do array baseado no LimiteEnderecos
		enderecosSeparados := listaEnderecos[:limite]
		if len(enderecosSeparados) == 0 {
			break
		}
		//Alocar memoria para o array de string
		matrizEnderecos[contador] = make([]string, len(enderecosSeparados))
		for j := 0; j < len(enderecosSeparados); j++ {
			//Atribuir os valores do slice para a matriz
			matrizEnderecos[contador][j] = enderecosSeparados[j]
		}

		//Remove os elementos que foram salvos na matriz
		listaEnderecos = append(listaEnderecos[:0], listaEnderecos[limite:]...)

		// Se o tamanho do slice for menor do que o tamanho maximo do Array,
		// setar tamanhoMaximo = tamanho Array
		// Evitar acesso invalido de memoria
		if limite > (len(listaEnderecos)) {
			limite = len(listaEnderecos)
		}

	}

	return matrizEnderecos, QtdLinhas
}
