import { Item, ItemActions, ItemContent, ItemHeader } from '#/components/ui/item.tsx';
import { Button } from '#/components/ui/button.tsx';
import { IconLink } from '@tabler/icons-react';
import { Tooltip, TooltipContent, TooltipTrigger } from '#/components/ui/tooltip.tsx';
import type { JobApplication } from '#/api/types.ts';
import { Link, linkOptions } from '@tanstack/react-router';

export default function ApplicationStatusJobApplicationItemRenderer({
	data,
	idx,
}: {
	data: JobApplication;
	idx: number;
}) {
	// TODO: add navigation to company page on actions
	const linkOpts = linkOptions({
		to: '/jobs/$jobID',
		params: { jobID: data.id.toString() },
	});
	return (
		<Item className="border-accent flex w-full flex-col flex-nowrap items-stretch gap-0 border p-0">
			<ItemHeader className="border-b-accent flex w-full items-center justify-center border-b px-3 py-2.5 text-lg">{`Job Application #${idx + 1}`}</ItemHeader>
			<div className="flex w-full items-center px-3 py-2.5">
				<ItemContent className="flex h-fit grow flex-row items-center justify-evenly px-3 py-2.5">
					<div className="grid w-full grid-cols-2 gap-x-4 gap-y-2">
						<div className="flex flex-col items-center">
							<div className="text-primary font-bold">Job Application Title</div>
							<span>{data.title}</span>
						</div>
						<div className="flex flex-col items-center">
							<div className="text-primary font-bold">Created At</div>
							<span>{new Date(data.created_at).toLocaleString()}</span>
						</div>
						<div className="flex flex-col items-center">
							<div className="text-primary font-bold">Company</div>
							<span>{data.company.name}</span>
						</div>
						<div className="flex flex-col items-center">
							<div className="text-primary font-bold">Last Update</div>
							<span>{new Date(data.updated_at).toLocaleString()}</span>
						</div>
					</div>
				</ItemContent>
				<ItemActions className="border-accent bg-background/20 h-full shrink-0 rounded-xl border p-3">
					<Link {...linkOpts}>
						<Tooltip>
							<TooltipTrigger
								render={
									<Button size="icon-sm">
										<IconLink />
									</Button>
								}
							/>
							<TooltipContent>
								<p>Navigate to this Job Application</p>
							</TooltipContent>
						</Tooltip>
					</Link>
				</ItemActions>
			</div>
		</Item>
	);
}
