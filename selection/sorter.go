package selection

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/rs/zerolog"
)

// SortStrategy defines an interface for sorting strategies.
type SortStrategy interface {
	Sort(batchOptions []BatchOption) []BatchOption
}

// Sorter sorts batch options.
type Sorter struct {
	logger zerolog.Logger
}

// NewSorter constructs a new Sorter.
func NewSorter(logger zerolog.Logger) Sorter {
	return Sorter{logger}
}

// Sort sorts batch options by method.
func (s Sorter) Sort(batchOptions []BatchOption, method SortMethod, sortKey string) []BatchOption {
	strategy := s.findSortStrategy(method, sortKey)

	s.logger.Info().
		Str("sortMethod", string(method)).
		Str("sortKey", sortKey).
		Str("strategy", fmt.Sprintf("%#v", strategy)).
		Msg("beginning batch options sort")

	sorted := strategy.Sort(batchOptions)

	s.logger.Info().
		Str("sortMethod", string(method)).
		Str("sortKey", sortKey).
		Str("strategy", fmt.Sprintf("%#v", strategy)).
		Msg("finished batch options sort")

	return sorted
}

func (s Sorter) findSortStrategy(method SortMethod, sortKey string) SortStrategy {
	s.logger.Info().
		Str("sortMethod", string(method)).
		Str("sortKey", sortKey).
		Msg("finding a strategy for the provided sort method")

	switch method {
	case Number:
		return SortByNumber{}
	case Random:
		return Shuffle{}
	case Alphabetical:
		return SortByContent{}
	case Metadata:
		return SortByMetadata{key: sortKey}
	}

	s.logger.Info().
		Str("sortMethod", string(method)).
		Str("sortKey", sortKey).
		Msg("couldn't find a sort strategy. Defaulting to SortByNumber strategy")

	return SortByNumber{}
}

// SortByNumber sorts by batch option number.
type SortByNumber struct{}

// Sort implements SortStrategy
func (s SortByNumber) Sort(batchOptions []BatchOption) []BatchOption {
	sort.Slice(batchOptions, func(i, j int) bool {
		return batchOptions[i].Number < batchOptions[j].Number
	})

	return batchOptions
}

// Shuffle randomly shuffles batch options.
type Shuffle struct{}

// Sort implements SortStrategy
func (s Shuffle) Sort(batchOptions []BatchOption) []BatchOption {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range batchOptions {
		j := r.Intn(i + 1)
		batchOptions[i], batchOptions[j] = batchOptions[j], batchOptions[i]
	}

	return batchOptions
}

// SortByContent sorts by batch option content.
type SortByContent struct{}

// Sort implements SortStrategy
func (s SortByContent) Sort(batchOptions []BatchOption) []BatchOption {
	sort.Slice(batchOptions, func(i, j int) bool {
		return batchOptions[i].Option.Content < batchOptions[j].Option.Content
	})

	return batchOptions
}

// SortByMetadata sorts by batch option metadata.
type SortByMetadata struct {
	key string
}

// Sort implements SortStrategy.
func (s SortByMetadata) Sort(batchOptions []BatchOption) []BatchOption {
	sort.Slice(batchOptions, func(i, j int) bool {
		if batchOptions[i].Option.Metadata == nil || batchOptions[j].Option.Metadata == nil {
			return true
		}

		iKey, ok := batchOptions[i].Option.Metadata[s.key]
		if !ok {
			return true
		}

		jKey, ok := batchOptions[j].Option.Metadata[s.key]
		if !ok {
			return true
		}

		return iKey < jKey
	})

	return batchOptions
}
