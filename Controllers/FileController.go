package Controllers

import (
	"fmt"
	"strconv"

	"main.go/Function/File"
)

func LerTexto(caminhoDoArquivo string) ([]string, error) {

	return File.LerTexto(caminhoDoArquivo)
}

// Funcao que escreve um texto no arquivo e retorna um erro caso tenha algum problema
func EscreverTexto(linhas []string, caminhoDoArquivo string) error {

	return File.EscreverTexto(linhas, caminhoDoArquivo)
}

func EscreverTextoSemApagar(linhas []string, caminhoDoArquivo string) error {

	return File.EscreverTextoSemApagar(linhas, caminhoDoArquivo)
}

func GetIndiceLogIndice(nomeArquivoIndice string) int {
	valorLogIndice, err := LerTexto(nomeArquivoIndice)
	indiceInicial := 0
	if len(valorLogIndice) > 0 {
		if err != nil {
			fmt.Print(err.Error())
			fmt.Println("Erro na função GetIndiceLogIndice na leitura do arquivo")
		}
		var er error
		indiceInicial, er = strconv.Atoi(valorLogIndice[0])
		if er != nil {
			fmt.Print(er.Error())
			fmt.Println("Erro na função GetIndiceLogIndice na conversão de string para int")
		}

	}

	return indiceInicial
}
