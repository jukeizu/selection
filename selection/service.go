package selection

import (
	"database/sql"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type DefaultService struct {
	logger          zerolog.Logger
	repository      Repository
	sorter          Sorter
	batcher         Batcher
	parseRegex      *regexp.Regexp
	validationRegex *regexp.Regexp
}

func NewDefaultService(logger zerolog.Logger, repository Repository, sorter Sorter, batcher Batcher) Service {
	parseRegex := regexp.MustCompile(`\b\d+\b`)
	validationRegex := regexp.MustCompile(`^[\d\s]+$`)
	return &DefaultService{logger, repository, sorter, batcher, parseRegex, validationRegex}
}

func (s DefaultService) Create(req CreateSelectionRequest) (SelectionReply, error) {
	selection, err := s.repository.Selection(req.AppId, req.InstanceId, req.UserId, req.ServerId)
	if err == nil {
		s.logger.Info().
			EmbedObject(selection).
			Msg("found existing selection")

		return s.createSelectionReply(req, selection), nil
	}
	if err != nil && err != sql.ErrNoRows {
		return SelectionReply{}, err
	}

	selection = Selection{
		AppId:      req.AppId,
		InstanceId: req.InstanceId,
		UserId:     req.UserId,
		ServerId:   req.ServerId,
		Options:    map[int]Option{},
	}

	if req.Randomize {
		req.Options = s.shuffleOptions(req.Options)
	}

	for i, option := range req.Options {
		selection.Options[i+1] = option
	}

	err = s.repository.CreateSelection(selection)
	if err != nil {
		return SelectionReply{}, err
	}

	s.logger.Info().
		EmbedObject(selection).
		Msg("created selection")

	return s.createSelectionReply(req, selection), nil
}

func (s DefaultService) Parse(req ParseSelectionRequest) ([]RankedOption, error) {
	if !s.validationRegex.MatchString(req.Content) {
		return nil, NewValidationError("Input may only contain numeric values.")
	}

	selection, err := s.repository.Selection(req.AppId, req.InstanceId, req.UserId, req.ServerId)
	if err != nil {
		return nil, err
	}

	choices := s.parseRegex.FindAllString(req.Content, -1)

	rankedOptions := []RankedOption{}

	for i, choice := range choices {
		c, err := strconv.Atoi(choice)
		if err != nil {
			return nil, NewValidationError("Input `%s` is not a valid selection.", choice)
		}

		option, ok := selection.Options[c]
		if !ok {
			return nil, NewValidationError("Input `%d` is not a valid selection.", c)
		}

		rankedOption := RankedOption{
			Rank:   i,
			Number: c,
			Option: option,
		}

		rankedOptions = append(rankedOptions, rankedOption)
	}

	return rankedOptions, nil
}

func (s DefaultService) Query(req QuerySelectionRequest) (QuerySelectionReply, error) {
	if req.Options == nil || len(req.Options) < 1 {
		return QuerySelectionReply{}, nil
	}

	selection, err := s.repository.Selection(req.AppId, req.InstanceId, req.UserId, req.ServerId)
	if err != nil {
		return QuerySelectionReply{}, err
	}

	if selection.Options == nil || len(selection.Options) < 1 {
		return QuerySelectionReply{}, err
	}

	rankedOptions := []RankedOption{}

	for number, option := range selection.Options {
		rank, ok := req.Options[option.OptionId]
		if ok {
			rankedOptions = append(rankedOptions, RankedOption{
				Rank:   int(rank),
				Number: number,
				Option: option,
			})
		}
	}

	sort.SliceStable(rankedOptions, func(i, j int) bool {
		return rankedOptions[i].Rank < rankedOptions[j].Rank
	})

	content := make([]string, len(rankedOptions))
	for i, rankedOption := range rankedOptions {
		content[i] = strconv.Itoa(int(rankedOption.Number))
	}

	return QuerySelectionReply{
		Options: rankedOptions,
		Content: strings.Join(content, " "),
	}, nil
}

func (s DefaultService) createSelectionReply(req CreateSelectionRequest, selection Selection) SelectionReply {
	batchOptions := s.createBatchOptions(req, selection)

	sortedBatchOptions := s.sorter.Sort(batchOptions, req.SortMethod, req.SortKey)

	selectionReply := SelectionReply{
		Selection: selection,
		Batches:   s.batcher.CreateBatches(sortedBatchOptions, req.BatchSize),
	}

	return selectionReply
}

func (s DefaultService) createBatchOptions(req CreateSelectionRequest, selection Selection) []BatchOption {
	batchOptions := BatchOptions{}

	for k, option := range selection.Options {
		batchOption := BatchOption{
			Number: k,
			Option: option,
		}

		batchOptions = append(batchOptions, batchOption)
	}

	return batchOptions
}

func (s DefaultService) shuffleOptions(options []Option) []Option {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range options {
		j := r.Intn(i + 1)
		options[i], options[j] = options[j], options[i]
	}

	return options
}
