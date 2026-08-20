import { createFileRoute, Link } from '@tanstack/react-router';

import { apiClient } from '#/api/client.ts';
import { DataTable } from '#/components/data-table';
import { createColumnHelper } from '@tanstack/react-table';
import type { DataTableFeatures } from '#/components/data-table/table-features.ts';
import type { JobApplication } from '#/api/types.ts';
import { Button } from '#/components/ui/button.tsx';
import { Plus } from 'lucide-react';

export const Route = createFileRoute('/_app/jobs/')({
	component: RouteComponent,
});

function RouteComponent() {
	const { data, isLoading: isDataLoading, isError } = apiClient.useQuery('get', '/jobs');

	const columnHelper = createColumnHelper<DataTableFeatures, JobApplication>();

	const columns = columnHelper.columns([
		columnHelper.accessor('id', {
			header: () => <div className="text-center">ID</div>,
		}),
		columnHelper.accessor('title', {
			header: 'Title',
		}),
		columnHelper.accessor('url', {
			header: 'URL',
		}),
		columnHelper.accessor('company.name', {
			header: 'Company',
		}),
		columnHelper.accessor('created_at', {
			header: 'Created At',
		}),
	]);

	const CreateButton = (
		<Button className="space-x-2 p-2" size="lg">
			<Link to="/jobs/new">New Job Application</Link>
			<Plus size={32} />
		</Button>
	);

	return (
		<main className="mx-4">
			<div className="py-4 text-center text-2xl">
				<h1>Job Applications</h1>
			</div>
			<div className="flex justify-end py-2">{CreateButton}</div>
			{isDataLoading ? (
				<div>Loading...</div>
			) : !isError && data && data.data ? (
				<DataTable
					data={data.data}
					columns={columns}
					noResultsMessage="No Jobs were found! You can create a new one using the button bellow!"
					CreateButton={CreateButton}
				/>
			) : (
				'fack'
			)}
		</main>
	);
}