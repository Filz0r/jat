import type { DataTableFeatures } from '#/components/data-table/table-features.ts';
import type { JobApplication } from '#/api/types.ts';

import { createFileRoute, Link, useNavigate } from '@tanstack/react-router';
import { apiClient } from '#/api/client.ts';
import { DataTable } from '#/components/data-table';
import { createColumnHelper } from '@tanstack/react-table';
import { Button } from '#/components/ui/button.tsx';
import { IconLink, IconArrowsUpDown, IconEye, IconPlus } from '@tabler/icons-react';
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
						<IconArrowsUpDown />
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
						<IconArrowsUpDown />
						Status
					</Button>
				</div>
			),
			cell: (cell) => {
				const cellData = cell.getValue();
				return (
					<div className="flex items-center justify-center">
						<Badge
							className={
								'text-md my-1 w-full py-3.5 ' + getColorFromKind(cellData.kind)
							}
						>
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
							className="text-md my-1 w-fit py-3.5"
							render={
								<a href={value}>
									Link
									<IconLink size={40} className="flex-1" data-icon="inline-end" />
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
						<IconArrowsUpDown />
						Created At
					</Button>
				</div>
			),
			cell: (cell) => {
				const value = cell.getValue();
				const converted = new Date(value);
				const dateString = converted.toLocaleDateString();
				const timeString = converted.toLocaleTimeString();
				return <div className="text-center">{dateString + ' ' + timeString}</div>;
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
						<IconArrowsUpDown />
						Last Update
					</Button>
				</div>
			),
			cell: (cell) => {
				const value = cell.getValue();
				const converted = new Date(value);
				const dateString = converted.toLocaleDateString();
				const timeString = converted.toLocaleTimeString();
				return <div className="text-center">{dateString + ' ' + timeString}</div>;
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
										<IconEye />
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
			<IconPlus size={32} />
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
					keyName="jobs-table"
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
