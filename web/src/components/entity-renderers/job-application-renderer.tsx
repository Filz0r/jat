import type { JobApplication } from '#/api/types.ts';

import { apiClient } from '#/api/client.ts';
import { Skeleton } from '#/components/ui/skeleton.tsx';
import {
	Accordion,
	AccordionContent,
	AccordionItem,
	AccordionTrigger,
} from '#/components/ui/accordion.tsx';
import UpdateJobStatusForm from '#/components/forms/update-job-status-form.tsx';
import JobNoteForm from '#/components/forms/job-note-form.tsx';
import JobApplicationNotesRenderer from '#/components/entity-renderers/job-application-notes-renderer.tsx';
import JobApplicationHistoryRenderer from '#/components/entity-renderers/job-application-history-renderer.tsx';
import JobApplicationInformationRenderer from '#/components/entity-renderers/job-application-information-renderer.tsx';
import DeleteJobApplication from '#/components/actions/delete-job-application.tsx';

export default function JobApplicationRenderer({ data }: { data: JobApplication }) {
	const {
		data: noteData,
		isLoading: isNoteLoading,
		error: noteErrors,
		isError: isNoteError,
	} = apiClient.useQuery('get', '/jobs/{jobID}/notes', { params: { path: { jobID: data.id } } });

	const {
		data: historyData,
		isLoading: isHistoryLoading,
		error: historyErrors,
		isError: isHistoryError,
	} = apiClient.useQuery('get', '/jobs/{jobID}/history', {
		params: {
			path: { jobID: data.id },
		},
	});

	return (
		<div className="my-3.5">
			<Accordion multiple defaultValue={['info', 'notes']}>
				<AccordionItem key="info" value="info" disabled>
					<div className="relative flex items-center justify-center border-b py-2 text-2xl font-bold text-primary">
						Information
						<span className="absolute right-2 space-x-1">
							<UpdateJobStatusForm
								currentStatus={data.status.id}
								jobID={data.id}
								refreshSelf
							/>
							<DeleteJobApplication jobID={data.id} navigateToJobs />
						</span>
					</div>
					<AccordionContent>
						<JobApplicationInformationRenderer data={data} />
					</AccordionContent>
				</AccordionItem>

				<AccordionItem key="notes" value="notes">
					<div className="relative flex items-center justify-center border-b py-2 text-2xl font-bold text-primary">
						Notes
						{noteData && noteData.data && noteData.data.length > 0 && (
							<span className="absolute right-2 space-x-1">
								<JobNoteForm jobID={data.id} refreshSelf />
							</span>
						)}
					</div>
					<AccordionContent>
						{isNoteLoading && <Skeleton className="size-fit" />}
						{isNoteError && (
							<div className="text-lg text-red-500">{noteErrors.message}</div>
						)}
						{noteData && noteData.data && (
							<JobApplicationNotesRenderer data={noteData.data} jobID={data.id} />
						)}
					</AccordionContent>
				</AccordionItem>
				<AccordionItem key="history" value="history">
					<AccordionTrigger className="hover:no-underline">
						<span className="flex-1 text-center text-2xl font-bold text-primary">
							Status History
						</span>
					</AccordionTrigger>
					<AccordionContent>
						{isHistoryLoading && <Skeleton className="size-fit" />}
						{isHistoryError && (
							<div className="text-lg text-red-500">{historyErrors.message}</div>
						)}
						{historyData && historyData.data && (
							<JobApplicationHistoryRenderer
								data={historyData.data}
								jobID={data.id}
							/>
						)}
					</AccordionContent>
				</AccordionItem>
			</Accordion>
		</div>
	);
}
