package Controllers

import (
	"log"

	"main.go/Function/Auxiliares"

	"gopkg.in/mgo.v2/bson"
	"main.go/Function/Repository"
	Model "main.go/Models"
)

/*
Busca um único documento atráves de uma chave e um valor
	Exemplo:
			Key = _id , Code = "6153a58d3700e70e40f8177a"
			Key = adresses , Code = "13adwKvLLpHdcYDh21FguCdJgKhaYP3Dse"
*/
func GetMapAdressId(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string, Key string, Code string) Model.ReturnAddrMapTx {
	return Auxiliares.GetMapAdressId(ConnectionMongoDB, DataBase, CollectionRecuperaDados, Key, Code)
}

func GetAllMapAddress(ConnectionMongoDB string, DataBase string, CollectionRecuperaDados string) []Model.ReturnAddrMapTx {
	var addressesMap []Model.ReturnAddrMapTx
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
		CollectionRecuperaDados, filter, option)
	// handle the errors.
	if err != nil {
		panic(err)
	}

	// le os documentos em partes, testei com 1000 documentos e deu certo
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var addrMap Model.ReturnAddrMapTx

		if err := cursor.Decode(&addrMap); err != nil {
			log.Fatal(err)
		}

		addressesMap = append(addressesMap, addrMap)

	}

	return addressesMap
}
