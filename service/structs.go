package service

type Request struct {
	MessageID int32  `json:"message_id"`
	GroupID   int64  `json:"group_id"`
	Message   string `json:"message"`
	UserID    int64  `json:"user_id"`
}

type Message struct {
	GroupID int64       `json:"group_id"`
	Message interface{} `json:"message"`
}

type CQMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type CQImage struct {
	File string `json:"file"`
	Url  string `json:"url"`
}

type Data struct {
	Embedded []EmbeddedData `json:"_embedded"`
}

type EmbeddedData struct {
	DataType    string `json:"data_type"`
	Name        string `json:"中文名"`
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
	RetainerAdventure []RetainerAdventure `json:"雇员探险"`
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
