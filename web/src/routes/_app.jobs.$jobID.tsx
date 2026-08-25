import { createFileRoute } from '@tanstack/react-router';
import { apiClient } from '#/api/client.ts';
import { Skeleton } from '#/components/ui/skeleton.tsx';
import JobApplicationRenderer from '#/components/entity-renderers/job-application/job-application-renderer.tsx';

export const Route = createFileRoute('/_app/jobs/$jobID')({
	component: RouteComponent,
});

function RouteComponent() {
	const { jobID } = Route.useParams();
	const parsedID = parseInt(jobID);

	if (isNaN(parsedID)) {
		return <div> Error fetching job data</div>;
	}

	const { data, isLoading, isError, error } = apiClient.useQuery('get', '/jobs/{jobID}', {
		params: {
			path: {
				jobID: parsedID,
			},
		},
	});

	return (
		<div className="mx-4">
			{isLoading && <Skeleton className="size-full" />}
			{isError && (
				<div className="text-xl text-red-500">Error loading data {error.message}</div>
			)}
			{data && data.data && <JobApplicationRenderer data={data.data} />}
		</div>
	);
}
