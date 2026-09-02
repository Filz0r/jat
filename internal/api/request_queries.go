package api

// ---------------------------------------------//
//			Application Status Queries			//
// ---------------------------------------------//

type applicationStatusListQuery struct {
	IncludeArchived bool `query:"include_archived"`
}

type applicationStatusSoftDeleteQuery struct {
	SoftDelete bool `query:"soft_delete"`
}
