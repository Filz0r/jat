import { createFileRoute } from '@tanstack/react-router';
import { AppShell } from '#/components/app-shell.tsx';
import { requireAuth } from '#/lib/route-guards.ts';

export const Route = createFileRoute('/jobs/$jobID')({
	beforeLoad: requireAuth,
	component: RouteComponent,
});

function RouteComponent() {
	const { jobID } = Route.useParams();
	const parsedID = parseInt(jobID);

	console.log(parsedID);

	return (
		<AppShell>
			<div>Hello "/jobs/{jobID}"!</div>
		</AppShell>
	);
}
