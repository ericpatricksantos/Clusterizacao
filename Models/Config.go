package Model

type Configuration struct {
	ConnectionMongoDB       []string `json:"ConnectionMongoDB"`
	FileLog                 []string `json:"FileLog"`
	DataBase                []string `json:"DataBase"`
	Collection              []string `json:"Collection"`
	UrlAPI                  []string `json:"UrlAPI"`
	RawBlock                string   `json:"RawBlock"`
	RawTx                   string   `json:"RawTx"`
	MultiAddr               string   `json:"MultiAddr"`
	RawAddr                 string   `json:"RawAddr"`
	BlockHeight             string   `json:"BlockHeight"`
	FormatJson              string   `json:"FormatJson"`
	Charts                  string   `json:"Charts"`
	Unspent                 string   `json:"Unspent"`
	Balance                 string   `json:"Balance"`
	LastBlock               string   `json:"LastBlock"`
	Unconfirmedtransactions string   `json:"Unconfirmedtransactions"`
	Blocks                  string   `json:"Blocks"`
}
