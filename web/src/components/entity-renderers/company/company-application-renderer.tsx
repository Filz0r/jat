import { Link, linkOptions } from '@tanstack/react-router';
import {
	Item,
	ItemActions,
	ItemContent,
	ItemDescription,
	ItemGroup,
	ItemHeader,
} from '#/components/ui/item.tsx';
import { Tooltip, TooltipContent, TooltipTrigger } from '#/components/ui/tooltip.tsx';
import { IconLink, IconPlus } from '@tabler/icons-react';
import { apiClient } from '#/api/client.ts';
import UpdateJobStatusForm from '#/components/forms/update-job-status-form.tsx';
import JobNoteForm from '#/components/forms/job-note-form.tsx';
import DeleteJobApplication from '#/components/actions/delete-job-application.tsx';
import { AccordionContent, AccordionItem } from '#/components/ui/accordion.tsx';
import { Skeleton } from '#/components/ui/skeleton.tsx';
import type { JobApplication } from '#/api/types.ts';

function CompanyApplicationItemRenderer({ data, idx }: { data: JobApplication; idx: number }) {
	const linkOpts = linkOptions({ to: '/jobs/$jobID', params: { jobID: data.id.toString() } });
	const createdAt = new Date(data.created_at);
	const updatedAt = new Date(data.updated_at);
	return (
		<Item className="border-accent mt-1 flex w-full flex-col flex-nowrap items-stretch gap-0 border p-0">
			<ItemHeader className="border-b-accent w-full border-b px-3 py-2.5 text-center">
				<div className="text-lg">Application #{idx}</div>
			</ItemHeader>
			<div className="flex w-full items-center gap-2.5 px-3 py-2.5">
				<ItemContent className="min-w-0 flex-1">
					<ItemDescription className="flex flex-col rounded-lg p-1 text-lg">
						<div>
							<span className="text-primary pr-2">Position:</span>
							{data.title}
						</div>

						<div>
							<span className="text-primary pr-2">Created At:</span>
							{createdAt.toLocaleTimeString() + ' ' + createdAt.toLocaleDateString()}
						</div>
						<div>
							<span className="text-primary pr-2">Last Update:</span>
							{updatedAt.toLocaleTimeString() + ' ' + updatedAt.toLocaleDateString()}
						</div>
					</ItemDescription>
				</ItemContent>
				<ItemActions className="border-accent bg-background/20 h-full shrink-0 rounded-xl border p-3">
					<Tooltip>
						<TooltipTrigger
							render={
								<Link
									{...linkOpts}
									className="bg-primary hover:bg-primary/80 rounded-xl p-1"
								>
									<IconLink size={14} />
									<span className="sr-only">View this Job Application</span>
								</Link>
							}
						/>
						<TooltipContent>
							<p>View this Job Application</p>
						</TooltipContent>
					</Tooltip>

					<UpdateJobStatusForm currentStatus={data.status.id} jobID={data.id} />
					<JobNoteForm jobID={data.id} />
					<DeleteJobApplication jobID={data.id} />
				</ItemActions>
			</div>
		</Item>
	);
}

export default function CompanyApplicationRenderer({ companyID }: { companyID: number }) {
	const { data, isLoading, error } = apiClient.useQuery('get', '/jobs', {
		params: {
			query: {
				company_id: [companyID],
			},
		},
	});
	const linkOpts = linkOptions({
		to: '/jobs/new',
		search: { company_id: companyID },
	});
	const createButton = (
		<Link
			{...linkOpts}
			className="bg-primary hover:bg-primary/60 flex w-fit items-center justify-center space-x-1 rounded-lg p-1.5 text-xs text-white no-underline!"
		>
			<span className="ml-1">Create Job Application</span>
			<IconPlus size={14} />
		</Link>
	);
	return (
		<AccordionItem key="applications" value="applications">
			<div className="text-primary relative flex items-center justify-center border-b py-2 text-2xl font-bold">
				<p>Job Applications</p>
				<span className="absolute right-2">
					{data && data.data && data.data.length > 0 && createButton}
				</span>
			</div>
			<AccordionContent>
				{isLoading && <Skeleton className="size-fit" />}
				{error && <div className="text-lg text-red-500">{error.message}</div>}
				{data && data.data && data.data.length > 0 ? (
					<ItemGroup className="mt-2">
						{data.data.map((item, idx) => (
							<CompanyApplicationItemRenderer
								data={item}
								idx={idx + 1}
								key={`company-${companyID}-application-${idx}`}
							/>
						))}
					</ItemGroup>
				) : (
					<div className="mt-3 flex flex-col items-center justify-center">
						<p className="text-center text-lg text-gray-500">
							No Job Applications recorded for this company
						</p>
						{createButton}
					</div>
				)}
			</AccordionContent>
		</AccordionItem>
	);
}
