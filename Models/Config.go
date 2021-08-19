package Model

type Configuration struct {
	ConnectionMongoDB string `json:"ConnectionMongoDB"`
	FileLog           string `json:"FileLog"`
	DataBase          string `json:"DataBase"`
	Collection        string `json:"Collection"`
	UrlAPI            string `json:"UrlAPI"`
	UrlAPIBlock       string `json:"UrlAPIBlock"`
	UrlAPITransaction string `json:"UrlAPITransaction"`
}
