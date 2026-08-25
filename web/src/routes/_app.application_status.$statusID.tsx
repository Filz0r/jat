import { createFileRoute } from '@tanstack/react-router';
import { apiClient } from '#/api/client.ts';
import { Skeleton } from '#/components/ui/skeleton.tsx';
import ApplicationStatusRenderer from '#/components/entity-renderers/application-status/application-status-renderer.tsx';

export const Route = createFileRoute('/_app/application_status/$statusID')({
	component: RouteComponent,
});

function RouteComponent() {
	const { statusID: _statusID } = Route.useParams();
	const statusID = parseInt(_statusID);

	if (isNaN(statusID)) {
		return <div> Error fetching job data</div>;
	}

	const { data, isLoading, error } = apiClient.useQuery(
		'get',
		'/application_statuses/{statusID}',
		{ params: { path: { statusID } } },
	);

	return (
		<div className="mx-4">
			{isLoading && <Skeleton className="size-full" />}
			{error && (
				<div className="text-xl text-red-500">Error loading data {error.message}</div>
			)}
			{data && data.data && <ApplicationStatusRenderer data={data.data} />}
		</div>
	);
}
