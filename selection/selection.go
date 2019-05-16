package selection

type CreateSelectionRequest struct {
	AppId      string
	InstanceId string
	UserId     string
	ServerId   string
	BatchSize  int
	SortMethod SortMethod
	Options    []Option
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

type SortMethod int

const (
	None SortMethod = iota
	Random
	Alphabetical
)

type Batch struct {
	Start   int
	End     int
	Options map[int]Option
}

type RankedOption struct {
	Rank   int
	Option Option
}

type ByAlphabetical []Option

func (s ByAlphabetical) Len() int {
	return len(s)
}

func (s ByAlphabetical) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByAlphabetical) Less(i, j int) bool {
	return s[i].Content < s[j].Content
}
