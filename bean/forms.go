package bean

// ReCharge ReCharge
type ReCharge struct {
	Nums       string `form:"nums" json:"nums"`               // nums
	Address    string `form:"address" json:"address"`         // address
	CreateTime int64  `form:"create_time" json:"create_time"` // create_time
	IsAuth     string `form:"isauth" json:"isauth"`           // isauth
}
