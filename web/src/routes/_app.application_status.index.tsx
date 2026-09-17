import type { JobApplicationStatus } from '#/api/types.ts';
import type { DataTableFeatures } from '#/components/data-table/table-features.ts';

import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { apiClient } from '#/api/client.ts';
import { createColumnHelper } from '@tanstack/react-table';
import { Button } from '#/components/ui/button.tsx';
import { DataTable } from '#/components/data-table';
import { Skeleton } from '#/components/ui/skeleton.tsx';
import { IconArrowsUpDown, IconCheck, IconEye, IconX } from '@tabler/icons-react';
import { Badge } from '#/components/ui/badge.tsx';
import { getColorFromKind } from '#/lib/utils.ts';
import CreateApplicationStatusForm from '#/components/forms/application-status-form.tsx';
import { Tooltip, TooltipContent, TooltipTrigger } from '#/components/ui/tooltip.tsx';
import { useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import ArchiveOrDeleteApplicationStatus from '#/components/actions/archive-or-delete-application-status.tsx';
import UnarchiveApplicationStatus from '#/components/actions/unarchive-application-status.tsx';

export const Route = createFileRoute('/_app/application_status/')({
	component: RouteComponent,
});

function RouteComponent() {
	const queryClient = useQueryClient();
	const [showArchived, setShowArchived] = useState(false);
	const navigate = useNavigate();
	const { data, isLoading, isError, error } = apiClient.useQuery('get', '/application_statuses', {
		params: {
			query: {
				include_archived: showArchived,
			},
		},
	});
	const columnHelper = createColumnHelper<DataTableFeatures, JobApplicationStatus>();
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
		columnHelper.accessor('status', {
			header: () => <div className="text-center">Status</div>,
			cell: (cell) => <div className="text-center">{cell.getValue()}</div>,
		}),
		columnHelper.accessor('kind', {
			header: ({ column }) => (
				<div className="flex justify-center">
					<Button
						variant="ghost"
						onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
						className="cursor-pointer text-center"
					>
						<IconArrowsUpDown />
						Kind
					</Button>
				</div>
			),
			cell: (cell) => {
				const cellData = cell.getValue();
				return (
					<div className="flex items-center justify-center">
						<Badge className={'w-full py-3.5 ' + getColorFromKind(cellData)}>
							{cellData.toUpperCase()}
						</Badge>
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
						<IconArrowsUpDown />
						Last Update
					</Button>
				</div>
			),
			cell: (cell) => {
				const value = cell.getValue();
				const converted = new Date(value);
				return <div className="text-center">{converted.toLocaleString()}</div>;
			},
		}),
		columnHelper.accessor('archived', {
			header: () => <div className="flex justify-center">Archived</div>,
			cell: (cell) => {
				const value = cell.getValue();
				return (
					<div className="flex justify-center">
						{value ? (
							<IconCheck className="text-green-500" />
						) : (
							<IconX className="text-red-500" />
						)}
					</div>
				);
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
												to: '/application_status/$statusID',
												params: { statusID: cellData.id.toString() },
											});
										}}
									>
										<IconEye />
										<span className="sr-only">
											View Application Status Page
										</span>
									</Button>
								}
							/>
							<TooltipContent>
								<p>View Application Status Page</p>
							</TooltipContent>
						</Tooltip>
						{cellData.archived && (
							<UnarchiveApplicationStatus statusID={cellData.id} data={cellData} />
						)}
						<ArchiveOrDeleteApplicationStatus
							statusID={cellData.id}
							data={cellData}
							softDelete={cellData.archived}
						/>
						<CreateApplicationStatusForm edit data={cellData} />
					</div>
				);
			},
		}),
	]);

	const CreateButton = <CreateApplicationStatusForm largeTrigger />;
	return (
		<main className="mx-4">
			<div className="py-4 text-center text-2xl">
				<h1>Job Application Status</h1>
			</div>
			<div className="flex justify-end space-x-2 py-2">
				<Button
					onClick={async () => {
						const newState = !showArchived;
						setShowArchived(newState);
						await queryClient.refetchQueries({
							queryKey: [
								'get',
								'/application_status',
								{
									params: {
										query: {
											include_archived: newState,
										},
									},
								},
							],
						});
					}}
				>
					{showArchived ? 'Hide' : 'Show'} Archived
				</Button>
				{CreateButton}
			</div>
			{isLoading ? (
				<Skeleton className="size-fit" />
			) : !isError && data && data.data ? (
				<DataTable
					keyName="application-status-table"
					data={data.data}
					columns={columns}
					noResultsMessage="No Jobs were found! You can create a new one using the button bellow!"
					CreateButton={CreateButton}
				/>
			) : (
				<div className="text-center text-xl text-red-500">
					Error Loading data: {isError ? error.message : 'Could not reach API server'}
				</div>
			)}
		</main>
	);
}
