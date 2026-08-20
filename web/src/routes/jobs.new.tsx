import { createFileRoute } from '@tanstack/react-router';
import { AppShell } from '#/components/app-shell.tsx';
import { useAuth } from '#/contexts/auth-context.tsx';
import { requireAuth } from '#/lib/route-guards.ts';
import CreateJobApplicationForm from '#/components/forms/create-job-application-form.tsx';

export const Route = createFileRoute('/jobs/new')({
	beforeLoad: requireAuth,
	component: RouteComponent,
});

function RouteComponent() {
	const { user } = useAuth();

	return (
		<AppShell>
			<main className="m-4">
				<CreateJobApplicationForm
					defaultStatus={user?.default_status_id ? user.default_status_id : 0}
				/>
			</main>
		</AppShell>
	);
}
