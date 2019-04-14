package selection

import "github.com/jukeizu/selection/api/protobuf-spec/selectionpb"

type GrpcServer struct {
	service Service
}

func NewGrpcServer(service Service) GrpcServer {
	return GrpcServer{service}
}

func (s GrpcServer) CreateSelection(req *selectionpb.CreateSelectionRequest) (*selectionpb.CreateSelectionReply, error) {
	selection, err := s.service.Create(CreateSelectionRequestToDto(req))
	if err != nil {
		return nil, err
	}

	return DtoToCreateSelectionReply(selection), nil
}

func (s GrpcServer) ParseSelection(req *selectionpb.ParseSelectionRequest) (*selectionpb.ParseSelectionReply, error) {
	rankedOptions, err := s.service.Parse(ParseSelectionRequestToDto(req))
	if err != nil {
		return nil, err
	}

	return DtoToParseSelectionReply(rankedOptions), nil
}

func CreateSelectionRequestToDto(req *selectionpb.CreateSelectionRequest) CreateSelectionRequest {
	c := CreateSelectionRequest{
		AppId:    req.AppId,
		UserId:   req.UserId,
		ServerId: req.ServerId,
	}

	for _, reqOption := range req.Options {
		option := Option{
			Id:       reqOption.Id,
			Content:  reqOption.Content,
			Metadata: reqOption.Metadata,
		}

		c.Options = append(c.Options, option)
	}

	return c
}

func DtoToCreateSelectionReply(selection Selection) *selectionpb.CreateSelectionReply {
	reply := &selectionpb.CreateSelectionReply{}

	for _, option := range selection.Options {
		selectionOption := &selectionpb.SelectionOption{
			Id:       option.Id,
			Content:  option.Content,
			Metadata: option.Metadata,
		}

		reply.SelectionOptions = append(reply.SelectionOptions, selectionOption)
	}

	return reply
}

func ParseSelectionRequestToDto(req *selectionpb.ParseSelectionRequest) ParseSelectionRequest {
	p := ParseSelectionRequest{
		AppId:    req.AppId,
		UserId:   req.UserId,
		ServerId: req.ServerId,
		Content:  req.Content,
	}

	return p
}

func DtoToParseSelectionReply(rankedOptions []RankedOption) *selectionpb.ParseSelectionReply {
	reply := &selectionpb.ParseSelectionReply{}

	for _, option := range rankedOptions {
		rankedOption := &selectionpb.RankedOption{
			Id:   option.Id,
			Rank: option.Rank,
		}

		reply.RankedOptions = append(reply.RankedOptions, rankedOption)
	}

	return reply
}
