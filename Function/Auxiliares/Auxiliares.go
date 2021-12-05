package Auxiliares

/*
	Esse arquivo foi criado para armazenar todas as funções que são utilizadas frequentemente.

*/
import (
	"fmt"
	"log"
	"strconv"

	"gopkg.in/mgo.v2/bson"
	"main.go/Function/API"
	"main.go/Function/Arrays"
	"main.go/Function/File"
	"main.go/Function/Repository"
	Model "main.go/Models"
)

/*
	Buscar todos os enderecos que estão no mongoDb
	na collection de MapeandoEndereco
*/
type IGetAllMapeandoAdress func(string, string, string) (addrMap []Model.ReturnAddrMapTx)

func GetAllMapeandoAdress(ConnectionMongoDB string, DataBase string, Collection string) (addrMap []Model.ReturnAddrMapTx) {

	// Get Client, Context, CalcelFunc and err from connect method.
	client, ctx, cancel, err := Repository.Connect(ConnectionMongoDB)
	if err != nil {
		panic(err)
	}

	// Free the resource when mainn dunction is  returned
	defer Repository.Close(client, ctx, cancel)

	// create a filter an option of type interface,
	// that stores bjson objects.
	var filter, option interface{}

	// filter  gets all document,
	// with maths field greater that 70
	filter = bson.M{}

	//  option remove id field from all documents
	option = bson.M{}

	// call the query method with client, context,
	// database name, collection  name, filter and option
	// This method returns momngo.cursor and error if any.
	cursor, err := Repository.Query(client, ctx, DataBase,
		Collection, filter, option)
	// handle the errors.
	if err != nil {
		panic(err)
	}

	// le os documentos em partes, testei com 1000 documentos e deu bom
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var enderecos Model.ReturnAddrMapTx

		if err := cursor.Decode(&enderecos); err != nil {
			log.Fatal(err)
		}

		addrMap = append(addrMap, enderecos)
	}

	return addrMap
}

/*
Busca um único documento atráves de uma chave e um valor
	Exemplo:
			Key = _id , Code = "6153a58d3700e70e40f8177a"
			Key = adresses , Code = "13adwKvLLpHdcYDh21FguCdJgKhaYP3Dse"
*/
func GetMapAdressId(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string, Key string, Code string) Model.ReturnAddrMapTx {
	// Get Client, Context, CalcelFunc and err from connect method.
	client, ctx, cancel, err := Repository.Connect(ConnectionMongoDB)
	if err != nil {
		panic(err)
	}

	// Free the resource when mainn dunction is  returned
	defer Repository.Close(client, ctx, cancel)

	result, err := Repository.QueryOne(client, ctx, DataBase,
		CollectionRecuperaDados, Key, Code)
	// handle the errors.
	if err != nil {
		panic(err)
	}
	return result
}

/*
Recuperar Todos os enderecos que esta armazenado no mongoDB
Retorno todas os input e out de todos os documentos
*/
func RecuperarEnderecosInputOut(ConnectionMongoDB string, DataBase string, Collection string) ([]string, []string) {
	var listaout []string
	var listainput []string

	// Get Client, Context, CalcelFunc and err from connect method.
	client, ctx, cancel, err := Repository.Connect(ConnectionMongoDB)
	if err != nil {
		panic(err)
	}

	// Free the resource when mainn dunction is  returned
	defer Repository.Close(client, ctx, cancel)

	// create a filter an option of type interface,
	// that stores bjson objects.
	var filter, option interface{}

	// filter  gets all document,
	// with maths field greater that 70
	filter = bson.M{}

	//  option remove id field from all documents
	option = bson.M{}

	// call the query method with client, context,
	// database name, collection  name, filter and option
	// This method returns momngo.cursor and error if any.
	cursor, err := Repository.Query(client, ctx, DataBase,
		Collection, filter, option)
	// handle the errors.
	if err != nil {
		panic(err)
	}

	// le os documentos em partes, testei com 1000 documentos e deu bom
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var enderecos Model.UnicoEndereco

		if err := cursor.Decode(&enderecos); err != nil {
			log.Fatal(err)
		}

		for _, doc := range enderecos.Txs {
			for _, x := range doc.Inputs {

				listainput = append(listainput, x.Prev_out.Addr)
			}
			for _, h := range doc.Out {
				listaout = append(listaout, h.Addr)
			}
		}

	}

	return listainput, listaout
}

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

func Inicializa(enderecoAnalisado string, Txs []Model.Transaction) ([]Model.InputAddr, []Model.OutputAddr) {
	var TempInput []string
	var TempOut []string

	// Atribuindo todos os inputs e output em variaveis temporarias
	for _, elem := range Txs {

		for _, input := range elem.Inputs {
			if enderecoAnalisado != input.Prev_out.Addr {
				TempInput = append(TempInput, input.Prev_out.Addr)
			}
		}

		for _, out := range elem.Out {
			if enderecoAnalisado != out.Addr {
				TempOut = append(TempOut, out.Addr)
			}
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

/*
	Busca um Multi Endereco
*/
func GetMultiEndereco(endereco []string, UrlAPI string, MultiAddr string) Model.MultiEndereco {
	if len(endereco) == 1 {
		valor := API.GetMultiEnderecos(endereco, UrlAPI, MultiAddr)

		Repository.ToDoc(valor)

		return valor
	}
	return Model.MultiEndereco{}
}

/*
	Converte um Multi Endereco em um Endereco Unico
*/
func ConverteMultiAddrEMAddr(MultiEndereco Model.MultiEndereco) (endereco Model.UnicoEndereco) {
	if len(MultiEndereco.Addresses) == 1 {
		endereco.Address = MultiEndereco.Addresses[0].Address
		endereco.N_tx = MultiEndereco.Addresses[0].N_tx
		endereco.Txs = MultiEndereco.Txs
	}
	return endereco
}

/*
  Salva um Endereco no MongoDb
*/
func SalvaEnderecoMongoDb(endereco Model.UnicoEndereco, indice int,
	ConnectionMongoDB string, DataBase string, Collection string,
	nomeArquivoSemApagar string, nomeArquivoIndice string) bool {
	if len(endereco.Address) > 0 {
		cliente, contexto, cancel, errou := Repository.Connect(ConnectionMongoDB)
		if errou != nil {
			log.Fatal(errou)
		}
		Repository.Ping(cliente, contexto)

		defer Repository.Close(cliente, contexto, cancel)

		insertOneResult, err := Repository.InsertOne(cliente, contexto, DataBase, Collection, endereco)

		// handle the error
		if err != nil {
			panic(err)
		}

		// print the insertion id of the document,
		// if it is inserted.
		fmt.Println("Result of InsertOne")
		fmt.Println(insertOneResult.InsertedID)

		temp := []string{strconv.Itoa(indice)}

		File.EscreverTexto(temp, nomeArquivoIndice)

		fmt.Println("Indice atualizado para ", indice)

		fmt.Printf("Salvamento concluido do %dº endereco %s \n\n", indice, endereco.Address)

		fmt.Printf("-------------------------------------------------------------------------------")
		fmt.Printf("\n")

		return true
	} else {
		valorTemp := " Indice: " + strconv.Itoa(indice) + " Endereco: " + endereco.Address
		temp := []string{valorTemp}

		File.EscreverTextoSemApagar(temp, nomeArquivoSemApagar)
		fmt.Println("Nao tem Dados o Indice: " + strconv.Itoa(indice) + " Endereco: " + endereco.Address)

		return false
	}
}
