import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { apiClient } from '#/api/client.ts';

import { Skeleton } from '#/components/ui/skeleton.tsx';

import CompanyRenderer from '#/components/entity-renderers/company/company-renderer.tsx';

export const Route = createFileRoute('/_app/companies/$companyID')({
	component: RouteComponent,
});

function RouteComponent() {
	const navigate = useNavigate();
	const { companyID } = Route.useParams();
	const converted = parseInt(companyID);

	if (isNaN(converted)) {
		void navigate({ to: '/companies' });
	}

	const { data, isLoading, error } = apiClient.useQuery('get', '/company/{companyID}', {
		params: {
			path: {
				companyID: converted,
			},
			query: {
				total_count: true,
				user_count: true,
			},
		},
	});

	return (
		<div className="mx-4">
			{isLoading && <Skeleton className="size-full" />}
			{error && (
				<div className="text-xl text-red-500">Error loading data {error.message}</div>
			)}
			{data && data.data && <CompanyRenderer data={data.data} />}
		</div>
	);
}
