package selection

import (
	"github.com/rs/zerolog"
)

type CreateSelectionRequest struct {
	AppId      string
	InstanceId string
	UserId     string
	ServerId   string
	Randomize  bool
	BatchSize  int
	SortMethod SortMethod
	SortKey    string
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
	Options    map[int]Option
}

type SelectionReply struct {
	Selection Selection
	Batches   []Batch
}

type SortMethod string

const (
	Number       = SortMethod("number")
	Random       = SortMethod("random")
	Alphabetical = SortMethod("alphabetical")
	Metadata     = SortMethod("metadata")
)

type Batch struct {
	Options []BatchOption
}

type BatchOption struct {
	Number int
	Option Option
}

type BatchOptions []BatchOption

type RankedOption struct {
	Rank   int
	Number int
	Option Option
}

type Service interface {
	Create(CreateSelectionRequest) (SelectionReply, error)
	Parse(ParseSelectionRequest) ([]RankedOption, error)
}

func (selection Selection) MarshalZerologObject(e *zerolog.Event) {
	e.Str("selection.Id", selection.Id).
		Str("selection.AppId", selection.AppId).
		Str("selection.InstanceId", selection.InstanceId).
		Str("selection.UserId", selection.UserId).
		Str("selection.ServerId", selection.ServerId)
}
