import { createFileRoute } from '@tanstack/react-router';

import { useAuth } from '#/hooks/use-auth.ts';
import CreateJobApplicationForm from '#/components/forms/create-job-application-form.tsx';

export const Route = createFileRoute('/_app/jobs/new')({
	component: RouteComponent,
});

function RouteComponent() {
	const { user } = useAuth();

	return (
		<main className="m-4">
			<CreateJobApplicationForm
				defaultStatus={user?.default_status_id ? user.default_status_id : 0}
			/>
		</main>
	);
}