package selection

type Repository interface {
	Create(Selection) error
}
