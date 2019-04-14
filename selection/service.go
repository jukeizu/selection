package selection

type Service interface {
	Create(CreateSelectionRequest) (Selection, error)
	Parse(ParseSelectionRequest) ([]RankedOption, error)
}

type DefaultService struct {
}
