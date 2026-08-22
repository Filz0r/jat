import type { JobApplicationHistory } from '#/api/types.ts';
import { Item, ItemContent, ItemGroup, ItemHeader } from '#/components/ui/item.tsx';
import { Badge } from '#/components/ui/badge.tsx';
import { getColorFromKind } from '#/lib/utils.ts';

export default function JobApplicationHistoryRenderer({
	data,
	jobID,
}: {
	data: JobApplicationHistory[];
	jobID: number;
}) {
	return (
		<ItemGroup className="mt-3">
			{data.map((d, idx) => (
				<Item
					key={`job-${jobID}-history-${d.id}`}
					className="flex w-full flex-col flex-nowrap items-stretch gap-0 border border-accent p-0"
				>
					<ItemHeader className="flex w-full items-center justify-center border-b border-b-accent px-3 py-2.5 text-lg">{`Status History #${idx + 1}`}</ItemHeader>
					<ItemContent className="flex items-center justify-center space-y-2 px-3 py-2.5">
						<div className="flex items-center space-x-6">
							{d.old_status && (
								<div className="flex flex-col items-center justify-center space-x-2">
									<span className="pb-2 font-bold text-primary">Old Status</span>
									<Badge className={getColorFromKind(d.old_status.kind)}>
										{d.old_status.status}
									</Badge>
								</div>
							)}
							<div className="flex flex-col items-center justify-center space-x-2">
								<span className="pb-2 font-bold text-primary">
									{d.old_status ? 'New Status' : 'Initial Status'}
								</span>
								<Badge className={getColorFromKind(d.new_status.kind)}>
									{d.new_status.status}
								</Badge>
							</div>
						</div>

						<div className="flex items-center space-x-2">
							<span className="font-bold text-primary">Timestamp</span>
							<span>{new Date(d.created_at).toLocaleString()}</span>
						</div>
					</ItemContent>
				</Item>
			))}
		</ItemGroup>
	);
}
