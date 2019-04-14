package selection

type CreateSelectionRequest struct {
	AppId    string
	UserId   string
	ServerId string
	Options  []Option
}

type Option struct {
	Id       string
	Content  string
	Metadata map[string]string
}

type ParseSelectionRequest struct {
	AppId    string
	UserId   string
	ServerId string
	Content  string
}

type Selection struct {
	Options []SelectionOption
}

type SelectionOption struct {
	Id       int32
	Content  string
	Metadata map[string]string
}

type RankedOption struct {
	Id   string
	Rank int32
}
