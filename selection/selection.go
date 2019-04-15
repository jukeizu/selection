package selection

type CreateSelectionRequest struct {
	AppId     string
	UserId    string
	ServerId  string
	Randomize bool
	Options   []Option
}

type Option struct {
	OptionId string
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
	Id       string
	AppId    string
	UserId   string
	ServerId string
	Options  map[int]Option
}

type RankedOption struct {
	Rank   int
	Option Option
}
