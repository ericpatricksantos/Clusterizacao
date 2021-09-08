package Controllers

import (
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
