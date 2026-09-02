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

// ---------------------------------------------//
//				Company Queries					//
// ---------------------------------------------//

type companyListQuery struct {
	UserCount    bool `query:"user_count"`
	TotalCount   bool `query:"total_count"`
	PreloadUsers bool `query:"preload_users"`
}

type countCompanyQuery struct {
	TotalCount bool `query:"total_count"`
}
