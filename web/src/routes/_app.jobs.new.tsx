import { createFileRoute, useSearch } from '@tanstack/react-router';

import { useAuth } from '#/hooks/use-auth.ts';
import CreateJobApplicationForm from '#/components/forms/create-job-application-form.tsx';
import { preSelectCompanyQuerySchema } from '#/schemas/job-applications.ts';

export const Route = createFileRoute('/_app/jobs/new')({
	component: RouteComponent,
	validateSearch: preSelectCompanyQuerySchema,
});

function RouteComponent() {
	const { user } = useAuth();
	const { company_id } = useSearch({ from: '/_app/jobs/new' });

	return (
		<main className="m-4">
			<CreateJobApplicationForm
				defaultStatus={user?.default_status_id ? user.default_status_id : 0}
				preSelectedCompany={company_id}
			/>
		</main>
	);
}
