package Controllers

import (
	"fmt"
	"time"

	"main.go/Function/Config"
	"main.go/Function/File"
)

func TesteArquivo(nomeArquivo string) {
	for {
		fmt.Println("Sleep Start.....")
		// Calling Sleep method
		time.Sleep(1 * time.Second)

		// Printed after sleep is over
		fmt.Println("Sleep Over.....")

		texto, _ := File.LerTexto(nomeArquivo)

		fmt.Println(texto[0])
	}
}

func TesteConfig() {
	conf := Config.GetConfig()

	fmt.Println("ConnectionMongoDB: ", conf.ConnectionMongoDB)
	fmt.Println("FileLog: ", conf.FileLog)
	fmt.Println("DataBase: ", conf.DataBase, " Collection: ", conf.Collection)
	fmt.Println("BlocoUnico UrlApi: ", conf.UrlAPI, " RawBlock: ", conf.RawBlock)
	fmt.Println("TransacaoUnica UrlApi: ", conf.UrlAPI, " RawTx: ", conf.RawTx)
	fmt.Println("DadosGrafico UrlApi: ", conf.UrlAPI, " Charts: ", conf.Charts+" $chart-type "+conf.FormatJson)
	fmt.Println("AlturaBloco UrlApi: ", conf.UrlAPI, " BlockHeight: ", conf.BlockHeight+" $block_height "+conf.FormatJson)
	fmt.Println("EnderecoUnico UrlApi: ", conf.UrlAPI, " RawAddr: ", conf.RawAddr)
	fmt.Println("EnderecoMulti UrlApi: ", conf.UrlAPI, " MultiAddr: ", conf.MultiAddr)
	fmt.Println("Unspent UrlApi: ", conf.UrlAPI, " Unspent: ", conf.Unspent)
	fmt.Println("Balance UrlApi: ", conf.UrlAPI, " Balance: ", conf.Balance)
	fmt.Println("Latestblock UrlApi: ", conf.UrlAPI, " Latestblock: ", conf.LastBlock)
	fmt.Println("Unconfirmed-transactions UrlApi: ", conf.UrlAPI, " Unconfirmed-transactions: ", conf.Unconfirmedtransactions)
	fmt.Println("Blocks UrlApi: ", conf.UrlAPI, " Blocks: ", conf.Blocks+" "+conf.FormatJson)
}
