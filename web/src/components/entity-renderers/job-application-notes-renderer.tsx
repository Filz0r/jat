import JobNoteForm from '#/components/forms/job-note-form.tsx';
import type { JobApplicationNote } from '#/api/types.ts';
import {
	Item,
	ItemActions,
	ItemContent,
	ItemDescription,
	ItemGroup,
	ItemHeader,
} from '#/components/ui/item.tsx';
import { Badge } from '#/components/ui/badge.tsx';
import { getColorFromKind } from '#/lib/utils.ts';
import DeleteApplicationNote from '#/components/actions/delete-application-note.tsx';

export default function JobApplicationNotesRenderer({
	data,
	jobID,
}: {
	data: JobApplicationNote[];
	jobID: number;
}) {
	if (data.length === 0) {
		return (
			<div className="my-4 flex flex-1 flex-col items-center justify-center gap-y-2">
				<span className="text-gray-500">
					No notes added to this job, you can add one using the button bellow.
				</span>
				<JobNoteForm jobID={jobID} largeTrigger />
			</div>
		);
	}
	return (
		<ItemGroup className="mt-3 min-w-0">
			{data.map((d, idx) => (
				<Item
					key={`job-${jobID}-note-${d.id}`}
					className="mt-1 flex w-full flex-col flex-nowrap items-stretch gap-0 border border-accent p-0"
				>
					<ItemHeader className="w-full border-b border-b-accent px-3 py-2.5">
						<div className="flex flex-col">
							<span className="text-lg">Note {idx + 1}</span>
						</div>
						<div className="flex min-w-0 shrink items-center justify-center">
							<span className="mr-2 flex flex-col border-r pr-2">
								<span className="pb-1 text-center text-xs text-gray-500">
									Status
								</span>
								<Badge className={getColorFromKind(d.status.kind)}>
									{d.status.status}
								</Badge>
							</span>
							<div className="flex flex-col items-start space-y-2 text-xs text-gray-500">
								<span>Created: {new Date(d.created_at).toDateString()}</span>
								<span>Last Update: {new Date(d.updated_at).toDateString()}</span>
							</div>
						</div>
					</ItemHeader>
					<div className="flex w-full items-center gap-2.5 px-3 py-2.5">
						<ItemContent className="min-w-0 flex-1">
							<ItemDescription className="flex flex-col rounded-lg p-1">
								<span className="pl-2 font-bold">Content</span>
								<span className="mt-2 rounded-lg border border-accent p-2 wrap-anywhere">
									{d.body}
								</span>
							</ItemDescription>
						</ItemContent>
						<ItemActions className="h-full shrink-0 rounded-xl border border-accent bg-background/20 p-3">
							<JobNoteForm jobID={jobID} edit existingData={d} />
							<DeleteApplicationNote jobID={jobID} data={d} />
						</ItemActions>
					</div>
				</Item>
			))}
		</ItemGroup>
	);
}
