import type { DataTableFeatures } from '#/components/data-table/table-features.ts';
import type { JobApplication } from '#/api/types.ts';

import { createFileRoute, Link, useNavigate } from '@tanstack/react-router';
import { apiClient } from '#/api/client.ts';
import { DataTable } from '#/components/data-table';
import { createColumnHelper } from '@tanstack/react-table';
import { Button } from '#/components/ui/button.tsx';
import { ArrowUpDown, ArrowUpRightIcon, Eye, Plus } from 'lucide-react';
import { getColorFromKind } from '#/lib/utils.ts';
import { Badge } from '#/components/ui/badge.tsx';
import UpdateJobStatusForm from '#/components/forms/update-job-status-form.tsx';
import { Tooltip, TooltipContent, TooltipTrigger } from '#/components/ui/tooltip.tsx';
import JobNoteForm from '#/components/forms/job-note-form.tsx';
import DeleteJobApplication from '#/components/actions/delete-job-application.tsx';

export const Route = createFileRoute('/_app/jobs/')({
	component: RouteComponent,
});

function RouteComponent() {
	const navigate = useNavigate();
	const { data, isLoading: isDataLoading, isError } = apiClient.useQuery('get', '/jobs');

	const columnHelper = createColumnHelper<DataTableFeatures, JobApplication>();
	const columns = columnHelper.columns([
		columnHelper.accessor('id', {
			header: ({ column }) => (
				<div className="flex justify-center">
					<Button
						variant="ghost"
						onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
						className="flex cursor-pointer items-center justify-center text-center"
					>
						<ArrowUpDown />
						ID
					</Button>
				</div>
			),

			cell: (cell) => <div className="text-center">{cell.getValue()}</div>,
		}),
		columnHelper.accessor('title', {
			header: () => <div className="text-center">Title</div>,
			cell: (cell) => <div className="text-center">{cell.getValue()}</div>,
		}),
		columnHelper.accessor('company.name', {
			header: () => <div className="text-center">Company</div>,
			cell: (cell) => <div className="text-center">{cell.getValue()}</div>,
		}),
		columnHelper.accessor('status', {
			sortFn: (rowA, rowB) => {
				const a = rowA.original.status.kind;
				const b = rowB.original.status.kind;
				return a.localeCompare(b);
			},
			header: ({ column }) => (
				<div className="flex justify-center">
					<Button
						variant="ghost"
						onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
						className="cursor-pointer text-center"
					>
						<ArrowUpDown />
						Status
					</Button>
				</div>
			),
			cell: (cell) => {
				const cellData = cell.getValue();
				return (
					<div className="flex items-center justify-center">
						<Badge className={'w-full py-2.5 ' + getColorFromKind(cellData.kind)}>
							{cellData.archived ? cellData.status + ' (Archived)' : cellData.status}
						</Badge>
					</div>
				);
			},
		}),
		columnHelper.accessor('url', {
			header: () => <div className="text-center">URL</div>,
			cell: (cell) => {
				const value = cell.getValue();
				return (
					<div className="flex items-center justify-center">
						<Badge
							className="w-fit py-2.5"
							render={
								<a href={value}>
									Link <ArrowUpRightIcon data-icon="inline-end" />
								</a>
							}
						/>
					</div>
				);
			},
		}),

		columnHelper.accessor('created_at', {
			header: ({ column }) => (
				<div className="flex justify-center">
					<Button
						variant="ghost"
						onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
						className="flex cursor-pointer items-center justify-between text-center"
					>
						<ArrowUpDown />
						Created At
					</Button>
				</div>
			),
			cell: (cell) => {
				const value = cell.getValue();
				const converted = new Date(value);
				return <div className="text-center">{converted.toDateString()}</div>;
			},
		}),
		columnHelper.accessor('updated_at', {
			header: ({ column }) => (
				<div className="flex justify-center">
					<Button
						variant="ghost"
						onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
						className="flex cursor-pointer items-center justify-between text-center"
					>
						<ArrowUpDown />
						Last Update
					</Button>
				</div>
			),
			cell: (cell) => {
				const value = cell.getValue();
				const converted = new Date(value);
				return <div className="text-center">{converted.toDateString()}</div>;
			},
		}),

		columnHelper.display({
			id: 'actions',
			header: () => <div className="text-center">Actions</div>,
			cell: ({ row }) => {
				const cellData = row.original;

				return (
					<div className="flex justify-center gap-x-2">
						<Tooltip>
							<TooltipTrigger
								render={
									<Button
										size="icon-sm"
										onClick={(e) => {
											e.preventDefault();
											e.stopPropagation();
											navigate({
												to: '/jobs/$jobID',
												params: { jobID: cellData.id.toString() },
											});
										}}
									>
										<Eye />
										<span className="sr-only">View Job Application</span>
									</Button>
								}
							/>
							<TooltipContent>
								<p>View Job Application</p>
							</TooltipContent>
						</Tooltip>

						<UpdateJobStatusForm
							currentStatus={cellData.status.id}
							jobID={cellData.id}
						/>
						<JobNoteForm jobID={cellData.id} />
						<DeleteJobApplication jobID={cellData.id} />
					</div>
				);
			},
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
