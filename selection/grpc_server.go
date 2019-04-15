package selection

import (
	"context"

	"github.com/jukeizu/selection/api/protobuf-spec/selectionpb"
)

type GrpcServer struct {
	service Service
}

func NewGrpcServer(service Service) GrpcServer {
	return GrpcServer{service}
}

func (s GrpcServer) CreateSelection(ctx context.Context, req *selectionpb.CreateSelectionRequest) (*selectionpb.CreateSelectionReply, error) {
	selection, err := s.service.Create(createSelectionRequestToDto(req))
	if err != nil {
		return nil, err
	}

	return dtoToCreateSelectionReply(selection), nil
}

func (s GrpcServer) ParseSelection(ctx context.Context, req *selectionpb.ParseSelectionRequest) (*selectionpb.ParseSelectionReply, error) {
	rankedOptions, err := s.service.Parse(parseSelectionRequestToDto(req))
	if err != nil {
		return nil, err
	}

	return dtoToParseSelectionReply(rankedOptions), nil
}

func createSelectionRequestToDto(req *selectionpb.CreateSelectionRequest) CreateSelectionRequest {
	c := CreateSelectionRequest{
		AppId:     req.AppId,
		UserId:    req.UserId,
		ServerId:  req.ServerId,
		Randomize: req.Randomize,
	}

	for _, reqOption := range req.Options {
		option := Option{
			OptionId: reqOption.OptionId,
			Content:  reqOption.Content,
			Metadata: reqOption.Metadata,
		}

		c.Options = append(c.Options, option)
	}

	return c
}

func dtoToCreateSelectionReply(selection Selection) *selectionpb.CreateSelectionReply {
	reply := &selectionpb.CreateSelectionReply{}

	for _, dtoSelectionOption := range selection.Options {
		selectionOption := &selectionpb.SelectionOption{
			SelectionOptionIndex: int32(dtoSelectionOption.SelectionOptionIndex),
			Option:               dtoToOption(dtoSelectionOption.Option),
		}

		reply.SelectionOptions = append(reply.SelectionOptions, selectionOption)
	}

	return reply
}

func dtoToOption(dtoOption Option) *selectionpb.Option {
	option := &selectionpb.Option{
		OptionId: dtoOption.OptionId,
		Content:  dtoOption.Content,
		Metadata: dtoOption.Metadata,
	}

	return option
}

func parseSelectionRequestToDto(req *selectionpb.ParseSelectionRequest) ParseSelectionRequest {
	p := ParseSelectionRequest{
		AppId:    req.AppId,
		UserId:   req.UserId,
		ServerId: req.ServerId,
		Content:  req.Content,
	}

	return p
}

func dtoToParseSelectionReply(dtoRankedOptions []RankedOption) *selectionpb.ParseSelectionReply {
	reply := &selectionpb.ParseSelectionReply{}

	for _, dtoRankedOption := range dtoRankedOptions {
		rankedOption := &selectionpb.RankedOption{
			Rank:   int32(dtoRankedOption.Rank),
			Option: dtoToOption(dtoRankedOption.Option),
		}

		reply.RankedOptions = append(reply.RankedOptions, rankedOption)
	}

	return reply
}
