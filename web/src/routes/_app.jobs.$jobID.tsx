import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_app/jobs/$jobID')({
	component: RouteComponent,
});

function RouteComponent() {
	const { jobID } = Route.useParams();
	const parsedID = parseInt(jobID);

	console.log(parsedID);

	return <div>Hello "/jobs/{jobID}"!</div>;
}