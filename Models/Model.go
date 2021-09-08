package Model

type MultiEndereco struct {
	Addresses []ListAddress `json:"addresses"`
	Wallet    Carteira      `json:"wallet"`
	Txs       []Transaction `json:"txs"`
}

type Carteira struct {
	Final_balance  int64 `json:"final_balance"`
	N_tx           int   `json:"n_tx"`
	N_tx_filtered  int   `json:"n_tx_filtered"`
	Total_received int64 `json:"total_received"`
	Total_sent     int64 `json:"total_sent"`
}

type ListAddress struct {
	Address        string `json:"address"`
	Final_balance  int64  `json:"final_balance"`
	N_tx           int    `json:"n_tx"`
	Total_received int64  `json:"total_received"`
	Total_sent     int64  `json:"total_sent"`
}

type UnicoEndereco struct {
	Hash160        string        `json:"hash160"`
	Address        string        `json:"address"`
	N_tx           int           `json:"n_tx"`
	N_unredeemed   int64         `json:"n_unredeemed"`
	Total_received int64         `json:"total_received"`
	Total_sent     int64         `json:"total_sent"`
	Final_balance  int64         `json:"final_balance"`
	Txs            []Transaction `json:"txs"`
}

type Block struct {
	Hash        string        `json:"hash"`
	Ver         int           `json:"ver"`
	Prev_block  string        `json:"prev_block"`
	Mrkl_root   string        `json:"mrkl_root"`
	Time        int           `json:"time"`
	Bits        int           `json:"bits"`
	Next_Block  []string      `json:"next_block"`
	Fee         int           `json:"fee"`
	Nonce       int           `json:"nonce"`
	N_tx        int           `json:"n_tx"`
	Size        int           `json:"size"`
	Block_index int           `json:"block_index"`
	Main_chain  bool          `json:"main_chain"`
	Height      int           `json:"height"`
	Weight      int           `json:"weight"`
	Tx          []Transaction `json:"tx"`
}

type Transaction struct {
	Hash         string        `json:"hash"`
	Ver          int           `json:"ver"`
	Vin_sz       int           `json:"vin_sz"`
	Vout_sz      int           `json:"vout_sz"`
	Size         int           `json:"size"`
	Weight       int           `json:"weight"`
	Fee          int           `json:"fee"`
	Relayed_by   string        `json:"relayed_by"`
	Lock_time    int           `json:"lock_time"`
	Tx_index     int64         `json:"tx_index"`
	Block_height int           `json:"block_height"`
	Inputs       []ListaInputs `json:"inputs"`
	Out          []ListaOut    `json:"out"`
}

type ListaInputs struct {
	Sequence int64   `json:"sequence"`
	Witness  string  `json:"witness"`
	Script   string  `json:"script"`
	Index    int     `json:"index"`
	Prev_out PrevOut `json:"prev_out"`
}

type ListaOut struct {
	Type               int                 `json:"type"`
	Spent              bool                `json:"spent"`
	Value              int64               `json:"value"`
	Spending_outpoints []Spendingoutpoints `json:"spending_outpoints"`
	N                  int                 `json:"n"`
	Tx_index           int64               `json:"tx_index"`
	Script             string              `json:"script"`
	Addr               string              `json:"addr"`
}

type PrevOut struct {
	Spent              bool                `json:"spent"`
	Script             string              `json:"script"`
	Spending_outpoints []Spendingoutpoints `json:"spending_outpoints"`
	Tx_index           int64               `json:"tx_index"`
	Value              int64               `json:"value"`
	Addr               string              `json:"addr"`
	N                  int                 `json:"n"`
	Type               int                 `json:"type"`
}

type Spendingoutpoints struct {
	Tx_index int64 `json:"tx_index"`
	N        int   `json:"n"`
}
