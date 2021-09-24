package Controllers

import (
	"main.go/Function/Arrays"
	Model "main.go/Models"
)

func RecuperaEnderecosUnionArrayRemoveDuplicados(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string) []string {
	input, out := RecuperarEnderecosInputOut(ConnectionMongoDB, DataBase, CollectionRecuperaDados)
	return UnionArrayRemoveDuplicados(input, out)
}

func UnionArrayRemoveDuplicados(input []string, out []string) []string {
	return Arrays.RemoveDuplicados(Arrays.UnionArray(input, out))
}

func RemoveDuplicados(lista []string) []string {
	return Arrays.RemoveDuplicados(lista)
}

func Inicializa(Txs []Model.Transaction) ([]Model.InputAddr, []Model.OutputAddr) {
	var TempInput []string
	var TempOut []string

	// Atribuindo todos os inputs e output em variaveis temporarias
	for _, elem := range Txs {

		for _, input := range elem.Inputs {
			TempInput = append(TempInput, input.Prev_out.Addr)
		}

		for _, out := range elem.Out {
			TempOut = append(TempOut, out.Addr)
		}
	}
	// Tirando os valores repetidos dessa variaveis temporarias
	In, tamIn := Arrays.RemoveDuplicadosStringVazia(TempInput)
	Out, tamOut := Arrays.RemoveDuplicadosStringVazia(TempOut)

	if tamIn < 0 || tamOut < 0 {
		return nil, nil
	}

	var Input []Model.InputAddr
	var Output []Model.OutputAddr
	// inicializando as listas de Input e Output com todos os input e output
	// com o addr e a sua quantidade
	for i := 0; i < tamIn; i++ {
		var temp Model.InputAddr
		temp.Addr = In[i]
		temp.Qtd = 0

		Input = append(Input, temp)
	}

	for j := 0; j < tamOut; j++ {
		var temp Model.OutputAddr
		temp.Addr = Out[j]
		temp.Qtd = 0

		Output = append(Output, temp)
	}

	return Input, Output
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
		//Alocar memoria para o array de string(Quantidade de colunas)
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
