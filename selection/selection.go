package selection

type CreateSelectionRequest struct {
	AppId           string
	InstanceId      string
	UserId          string
	ServerId        string
	Randomize       bool
	BatchSize       int
	BatchSortMethod BatchSortMethod
	Options         []Option
}

type Option struct {
	OptionId string
	Content  string
	Metadata map[string]string
}

type ParseSelectionRequest struct {
	AppId      string
	InstanceId string
	UserId     string
	ServerId   string
	Content    string
}

type Selection struct {
	Id         string
	AppId      string
	InstanceId string
	UserId     string
	ServerId   string
	Batches    []Batch
}

type BatchSortMethod int

const (
	Number BatchSortMethod = iota
	Alphabetical
)

type Batch struct {
	Start   rune
	End     rune
	Options map[int]Option
}

type RankedOption struct {
	Rank   int
	Option Option
}
