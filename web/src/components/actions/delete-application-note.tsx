import type { APIResponse, JobApplicationNote } from '#/api/types.ts';

import ModularDeleteEntityDialog from '#/components/dialogs/modular-delete-entity-dialog.tsx';
import { toast } from '#components/ui/toast';
import { api } from '#/api/client.ts';
import { useQueryClient } from '@tanstack/react-query';
import { Button } from '#/components/ui/button.tsx';
import { Trash } from 'lucide-react';

interface DeleteApplicationNoteProps {
	jobID: number;
	data: JobApplicationNote;
}

export default function DeleteApplicationNote({ data, jobID }: DeleteApplicationNoteProps) {
	const queryClient = useQueryClient();
	return (
		<ModularDeleteEntityDialog
			message="You are deleting this note, this is only reversible by administrative users, are you sure?"
			title={`Delete note with id of ${data.id}`}
			onClose={async () => {
				await toast.promise(
					Promise.all([
						api.DELETE('/jobs/{jobID}/notes/{noteID}', {
							params: {
								path: {
									jobID: data.job_id,
									noteID: data.id,
								},
							},
						}),
					]),
					{
						loading: `Deleting note with id ${data.id}...`,
						success: `Deleted note with id ${data.id}`,
						error: (response: APIResponse) =>
							`Error deleting note: ${response.message}`,
					},
				);
				await queryClient.refetchQueries({
					queryKey: ['get', '/jobs/{jobID}/notes', { params: { path: { jobID } } }],
				});
			}}
			triggerButton={
				<Button className="rounded-3xl p-3" size="icon-sm" variant="destructive">
					<Trash />
					<span className="sr-only">Delete this data</span>
				</Button>
			}
			deleteButtonMessage="Delete Note"
		/>
	);
}
