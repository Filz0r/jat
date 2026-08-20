import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { useUser } from '#/contexts/user-context.tsx';
import { useEffect } from 'react';
import { AppShell } from '#/components/app-shell.tsx';
import CreateJobApplicationForm from '#/components/forms/create-job-application-form.tsx';

export const Route = createFileRoute('/jobs/new')({
	component: RouteComponent,
});

function RouteComponent() {
	const { user, initialized, isLoading } = useUser();
	const navigate = useNavigate();
	useEffect(() => {
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
	return (
		<AppShell>
			<main className="m-4">
				<CreateJobApplicationForm
					defaultStatus={user.default_status_id ? user.default_status_id : 0}
				/>
			</main>
		</AppShell>
	);
}
