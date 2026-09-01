import { Link, linkOptions } from '@tanstack/react-router';
import { apiClient } from '#/api/client.ts';
import { Skeleton } from '#/components/ui/skeleton.tsx';
import type { CompanyChangeData } from '#/api/types';
import {
	Item,
	ItemActions,
	ItemContent,
	ItemDescription,
	ItemGroup,
	ItemHeader,
} from '#/components/ui/item.tsx';
import { Tooltip, TooltipContent, TooltipTrigger } from '#/components/ui/tooltip.tsx';
import { IconUserShare } from '@tabler/icons-react';
import RestoreCompanyChange from '#/components/actions/restore-company-change.tsx';
import { AccordionContent, AccordionItem, AccordionTrigger } from '#/components/ui/accordion.tsx';

function CompanyChangeHistoryItem({ data, idx }: { data: CompanyChangeData; idx: number }) {
	const userLinkOpts = linkOptions({
		to: '/admin/users/$userID',
		params: { userID: data.changed_by.user_id },
	});
	const changedAt = new Date(data.created_at);
	return (
		<Item className="border-accent mt-1 flex w-full flex-col flex-nowrap items-stretch gap-0 border p-0">
			<ItemHeader className="border-b-accent flex w-full justify-between border-b px-3 py-2.5 text-center">
				<div className="text-lg">Change #{idx}</div>
				<div>
					Changed At:{' '}
					<span className="text-gray-500">
						{changedAt.toLocaleTimeString() + ' ' + changedAt.toLocaleDateString()}
					</span>
				</div>
			</ItemHeader>
			<div className="flex w-full items-center gap-2.5 px-3 py-2.5">
				<ItemContent className="min-w-0 flex-1">
					<ItemDescription className="flex flex-col rounded-lg p-1 text-lg">
						{data.old_name && (
							<div>
								<span className="text-primary pr-2">Old Name:</span>
								{data.old_name}
							</div>
						)}
						{data.new_name && (
							<div>
								<span className="text-primary pr-2">New Name:</span>
								{data.new_name}
							</div>
						)}
						{!data.old_website && data.new_website && (
							<div>
								<span className="text-primary pr-2">Old Website:</span>
								None
							</div>
						)}
						{data.old_website && (
							<div>
								<span className="text-primary pr-2">Old Website:</span>
								{data.old_website}
							</div>
						)}
						{data.old_website && !data.new_website && (
							<div>
								<span className="text-primary pr-2">New Website:</span>
								None
							</div>
						)}
						{data.new_website && (
							<div>
								<span className="text-primary pr-2">New Website:</span>
								{data.new_website}
							</div>
						)}
					</ItemDescription>
				</ItemContent>
				<ItemActions className="border-accent bg-background/20 h-full shrink-0 rounded-xl border p-3">
					<Tooltip>
						<TooltipTrigger
							render={
								<Link
									{...userLinkOpts}
									className="bg-primary hover:bg-primary/80 rounded-xl p-1"
								>
									<IconUserShare size={14} />
									<span className="sr-only">View User Profile</span>
								</Link>
							}
						/>
						<TooltipContent>
							<p>View User Profile</p>
						</TooltipContent>
					</Tooltip>
					{!data.reverted && (
						<RestoreCompanyChange companyID={data.company_id} changeID={data.id} />
					)}
				</ItemActions>
			</div>
		</Item>
	);
}

export default function CompanyChangeHistoryRenderer({ companyID }: { companyID: number }) {
	const { data, isLoading, error } = apiClient.useQuery('get', '/company/{companyID}/history', {
		params: {
			path: { companyID },
		},
	});

	return (
		<AccordionItem key="history" value="history" className="bg-muted/50">
			<AccordionTrigger className="data-panel-open:border-b-border hover:cursor-pointer hover:no-underline">
				<span className="text-primary flex-1 text-center text-2xl font-bold">
					History of Changes
				</span>
			</AccordionTrigger>
			<AccordionContent>
				{isLoading && <Skeleton className="size-fit" />}
				{error && <div className="text-lg text-red-500">{error.message}</div>}
				{data && data.data && data.data.length > 0 ? (
					<ItemGroup className="mt-5">
						{data.data.map((d, idx) => (
							<CompanyChangeHistoryItem
								data={d}
								// this is dumb because we already check if data.data exists above but TS was complaining
								idx={data.data!.length - idx}
								key={`company-${companyID}-change-${idx}`}
							/>
						))}
					</ItemGroup>
				) : (
					<div className="mt-3 text-center text-lg text-gray-500">
						No changes recorded on this company
					</div>
				)}
			</AccordionContent>
		</AccordionItem>
	);
}
