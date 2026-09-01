import { createFileRoute, Link, linkOptions, useNavigate } from '@tanstack/react-router';
import { apiClient } from '#/api/client.ts';
import { Button } from '#/components/ui/button.tsx';
import { Skeleton } from '#/components/ui/skeleton.tsx';
import { DataTable } from '#/components/data-table';
import { createColumnHelper } from '@tanstack/react-table';
import type { DataTableFeatures } from '#/components/data-table/table-features.ts';
import type { CompanyData } from '#/api/types.ts';
import { IconArrowsUpDown, IconExternalLink, IconLink } from '@tabler/icons-react';
import { useAuth } from '#/hooks/use-auth.ts';
import { ArrowUpDown, Eye } from 'lucide-react';
import { Tooltip, TooltipContent, TooltipTrigger } from '#/components/ui/tooltip.tsx';
import CompanyForm from '#/components/forms/company-form.tsx';
import DeleteCompany from '#/components/actions/delete-company.tsx';

export const Route = createFileRoute('/_app/companies/')({
	component: RouteComponent,
});

function RouteComponent() {
	const { user } = useAuth();
	const isAdmin = user?.is_admin || false;
	const navigate = useNavigate();

	const { data, isLoading, error } = apiClient.useQuery('get', '/company', {
		params: {
			query: {
				user_count: true,
				preload_users: isAdmin ? true : undefined,
			},
		},
	});

	const columnHelper = createColumnHelper<DataTableFeatures, CompanyData>();
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
		columnHelper.accessor('name', {
			header: () => <div className="text-center">Name</div>,
			cell: (cell) => <div className="text-center">{cell.getValue()}</div>,
		}),
		columnHelper.accessor('website', {
			header: () => <div className="text-center">Website</div>,
			cell: (cell) => {
				const value = cell.getValue();
				if (value) {
					return (
						<div className="hover:outline-primary text-primary hover:bg-primary flex items-center justify-center gap-x-1 rounded-lg py-1.5 hover:text-white hover:outline-1 hover:outline-current">
							<IconExternalLink size={14} />
							<a href={value}>Website</a>
						</div>
					);
				}
				return <div className="text-center">None</div>;
			},
		}),
		columnHelper.accessor('user_count', {
			header: ({ column }) => (
				<div className="flex justify-center">
					<Button
						variant="ghost"
						onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
						className="flex cursor-pointer items-center justify-center text-center"
					>
						<IconArrowsUpDown />
						Nº Applications
					</Button>
				</div>
			),
			cell: (cell) => {
				const cellData = cell.getValue() || 0;
				return <div className="text-center">{cellData}</div>;
			},
		}),
		...(isAdmin
			? [
					columnHelper.accessor('created_by_user', {
						header: () => <div className="text-center">Created By</div>,
						cell: (cell) => {
							const value = cell.getValue()!;
							const linkOpts = linkOptions({
								to: '/admin/users/$userID',
								params: { userID: value.user_id },
							});
							return (
								<div className="my-1 flex justify-center">
									<Link
										{...linkOpts}
										className="hover:outline-primary text-primary hover:bg-primary flex items-center gap-x-0.5 rounded-lg px-1.5 py-1.5 hover:text-white hover:outline-1 hover:outline-current"
									>
										<IconLink size={14} className="mb-0.5" />
										{value.username}
									</Link>
								</div>
							);
						},
					}),
				]
			: []),
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
				return <div className="text-center">{converted.toLocaleString()}</div>;
			},
		}),

		...(isAdmin
			? [
					columnHelper.accessor('updated_by_user', {
						header: () => <div className="text-center">Edited By</div>,
						cell: (cell) => {
							const value = cell.getValue()!;
							const linkOpts = linkOptions({
								to: '/admin/users/$userID',
								params: { userID: value.user_id },
							});
							return (
								<div className="my-1 flex justify-center">
									<Link
										{...linkOpts}
										className="hover:outline-primary text-primary hover:bg-primary flex items-center gap-x-0.5 rounded-lg px-1.5 py-1.5 hover:text-white hover:outline-1 hover:outline-current"
									>
										<IconLink size={14} className="mb-0.5" />
										{value.username}
									</Link>
								</div>
							);
						},
					}),
				]
			: []),
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
				return <div className="text-center">{converted.toLocaleString()}</div>;
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
											void navigate({
												to: '/companies/$companyID',
												params: { companyID: cellData.id.toString() },
											});
										}}
									>
										<Eye />
										<span className="sr-only">View Company Page</span>
									</Button>
								}
							/>
							<TooltipContent>
								<p>View Company Page</p>
							</TooltipContent>
						</Tooltip>
						<CompanyForm edit refreshSelf existingData={cellData} />
						{isAdmin && <DeleteCompany companyID={cellData.id} />}
					</div>
				);
			},
		}),
	]);

	const CreateButton = <CompanyForm largeTrigger />;
	return (
		<main className="mx-4">
			<div className="py-4 text-center text-2xl">
				<h1>Job Application Status</h1>
			</div>
			<div className="flex justify-end px-2 py-2">{CreateButton}</div>
			{isLoading ? (
				<Skeleton className="size-fit" />
			) : !error && data && data.data ? (
				<DataTable
					keyName="company-table"
					data={data.data}
					columns={columns}
					noResultsMessage="No Jobs were found! You can create a new one using the button bellow!"
					CreateButton={CreateButton}
				/>
			) : (
				<div className="text-center text-xl text-red-500">
					Error Loading data: {error ? error.message : 'Could not reach API server'}
				</div>
			)}
		</main>
	);
}
