package selection

import "github.com/rs/zerolog"

// Batcher handles splitting a single slice of objects into batches of objects with limited length.
type Batcher struct {
	logger zerolog.Logger
}

// NewBatcher constructs a new Batcher.
func NewBatcher(logger zerolog.Logger) Batcher {
	return Batcher{logger}
}

// CreateBatches distributes batch options into batches of size batchSize.
func (b Batcher) CreateBatches(batchOptions []BatchOption, batchSize int) []Batch {
	if batchSize == 0 {
		b.logger.Info().Msg("batchSize is zero. No batches will be created")
		return []Batch{}
	}

	numBatchOptions := len(batchOptions)

	if numBatchOptions == 0 {
		b.logger.Info().Msg("no batch options were provided. No batches will be created")
		return []Batch{}
	}

	b.logger.Info().
		Int("batchSize", batchSize).
		Int("numBatchOptions", numBatchOptions).
		Msg("creating batches")

	batches := []Batch{}

	for i := 0; i < numBatchOptions; i += batchSize {
		nextBound := i + batchSize

		if nextBound > numBatchOptions {
			nextBound = numBatchOptions
		}

		batch := Batch{
			Options: batchOptions[i:nextBound],
		}

		batches = append(batches, batch)
	}

	b.logger.Info().
		Int("batchSize", batchSize).
		Int("numCreatedBatches", len(batches)).
		Msg("finished creating batches")

	return batches
}
