package selection

import (
	"context"

	"github.com/jukeizu/selection/api/protobuf-spec/selectionpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	service Service
}

func NewGrpcServer(service Service) GrpcServer {
	return GrpcServer{service}
}

func (s GrpcServer) CreateSelection(ctx context.Context, req *selectionpb.CreateSelectionRequest) (*selectionpb.CreateSelectionResponse, error) {
	selection, err := s.service.Create(createSelectionRequestToDto(req))
	if err != nil {
		return nil, toStatusErr(err)
	}

	return dtoToCreateSelectionReply(selection), nil
}

func (s GrpcServer) ParseSelection(ctx context.Context, req *selectionpb.ParseSelectionRequest) (*selectionpb.ParseSelectionResponse, error) {
	rankedOptions, err := s.service.Parse(ParseSelectionRequest{
		AppId:      req.AppId,
		InstanceId: req.InstanceId,
		UserId:     req.UserId,
		ServerId:   req.ServerId,
		Content:    req.Content,
	})
	if err != nil {
		return nil, toStatusErr(err)
	}

	return &selectionpb.ParseSelectionResponse{
		RankedOptions: dtoToRankedOption(rankedOptions),
	}, nil
}

func (s GrpcServer) QuerySelection(ctx context.Context, req *selectionpb.QuerySelectionRequest) (*selectionpb.QuerySelectionResponse, error) {
	queryReply, err := s.service.Query(QuerySelectionRequest{
		AppId:      req.AppId,
		InstanceId: req.InstanceId,
		UserId:     req.UserId,
		ServerId:   req.ServerId,
		Options:    req.Options,
	})
	if err != nil {
		return nil, toStatusErr(err)
	}

	return &selectionpb.QuerySelectionResponse{
		Options: dtoToRankedOption(queryReply.Options),
		Content: queryReply.Content,
	}, nil
}

func createSelectionRequestToDto(req *selectionpb.CreateSelectionRequest) CreateSelectionRequest {
	c := CreateSelectionRequest{
		AppId:      req.AppId,
		InstanceId: req.InstanceId,
		UserId:     req.UserId,
		ServerId:   req.ServerId,
		Randomize:  req.Randomize,
		BatchSize:  int(req.BatchSize),
		SortMethod: SortMethod(req.SortMethod),
		SortKey:    req.SortKey,
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

func dtoToCreateSelectionReply(selectionReply SelectionReply) *selectionpb.CreateSelectionResponse {
	reply := &selectionpb.CreateSelectionResponse{
		Batches: []*selectionpb.Batch{},
	}

	for _, dtoBatch := range selectionReply.Batches {
		replyBatch := &selectionpb.Batch{
			Options: []*selectionpb.BatchOption{},
		}

		for _, dtoBatchOption := range dtoBatch.Options {
			replyBatchOption := &selectionpb.BatchOption{
				Number: int32(dtoBatchOption.Number),
				Option: dtoToOption(dtoBatchOption.Option),
			}

			replyBatch.Options = append(replyBatch.Options, replyBatchOption)
		}

		reply.Batches = append(reply.Batches, replyBatch)
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

func dtoToRankedOption(dtoRankedOptions []RankedOption) []*selectionpb.RankedOption {
	rankedOptions := []*selectionpb.RankedOption{}

	for _, dtoRankedOption := range dtoRankedOptions {
		rankedOption := &selectionpb.RankedOption{
			Rank:   int32(dtoRankedOption.Rank),
			Number: int32(dtoRankedOption.Number),
			Option: dtoToOption(dtoRankedOption.Option),
		}

		rankedOptions = append(rankedOptions, rankedOption)
	}

	return rankedOptions
}

func toStatusErr(err error) error {
	switch err.(type) {
	case ValidationError:
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return err
}
