import { createFileRoute, useNavigate } from '@tanstack/react-router';
import * as React from 'react';
import { AppShell } from '#/components/app-shell.tsx';
import { useUser } from '#/contexts/user-context.tsx';

export const Route = createFileRoute('/jobs/$jobID')({
	component: RouteComponent,
});

function RouteComponent() {
	const { jobID } = Route.useParams();
	const parsedID = parseInt(jobID);
	const { user, initialized, isLoading } = useUser();
	const navigate = useNavigate();

	React.useEffect(() => {
		if (isLoading) return;

		if (!initialized) {
			navigate({ to: '/setup', replace: true });
			return;
		}

		if (!user) {
			navigate({ to: '/login', replace: true });
			return;
		}
	}, [isLoading, initialized, user, navigate]);

	if (isLoading || !initialized || !user) {
		return null;
	}

	console.log(parsedID);

	return (
		<AppShell>
			<div>Hello "/jobs/{jobID}"!</div>
		</AppShell>
	);
}
