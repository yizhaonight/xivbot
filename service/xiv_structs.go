package service

type Data struct {
	Embedded []EmbeddedData `json:"_embedded"`
}

type EmbeddedData struct {
	DataType    string `json:"data_type"`
	Name        string `json:"中文名"`
	JPName      string `json:"日文名"`
	Description string `json:"描述"`
	// for items
	ILevel int      `json:"品级"`
	Source Source   `json:"来源"`
	Usage  []string `json:"用途"`
	// for instances
	Type     string `json:"类型"`
	Location string `json:"地点"`
}

type Source struct {
	Collect           []Collect           `json:"采集"`
	Drop              []Drop              `json:"掉落"`
	RetainerAdventure []RetainerAdventure `json:"雇员探险"`
	Mall              []Mall              `json:"道具商城"`
	Quest             Quest               `json:"任务"`
}

type RetainerAdventure struct {
	Name  string `json:"中文名"`
	Job   string `json:"职业组"`
	Level int    `json:"等级"`
	Coin  int    `json:"探险币"`
	Time  int    `json:"时间"`
	Exp   int32  `json:"经验值"`
}

type Collect struct {
	Job   string `json:"职业"`
	Level int    `json:"等级"`
	Star  int    `json:"星级"`
}

type Mall struct {
	ItemName string `json:"商品名称"`
	Price    string `json:"价格"`
}

type Quest struct {
	Name string `json:"中文名"`
}

type Drop struct {
	Name string `json:"中文名"`
}
