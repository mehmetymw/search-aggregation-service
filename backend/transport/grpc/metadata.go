package grpc

import contentpb "github.com/mehmetymw/search-aggregation-service/backend/proto/gen"

var (
	sortOptionMetadata = []*contentpb.SortOptionMetadata{
		{Id: "score_desc", DisplayName: "Highest Score"},
		{Id: "score_asc", DisplayName: "Lowest Score"},
		{Id: "date_desc", DisplayName: "Newest First"},
		{Id: "date_asc", DisplayName: "Oldest First"},
	}
)
