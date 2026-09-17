import type { APIResponse, JobApplicationStatus } from '#/api/types.ts';

import { toast } from '#components/ui/toast';
import { api } from '#/api/client.ts';
import { useQueryClient } from '@tanstack/react-query';
import { Button } from '#/components/ui/button.tsx';
import { Trash } from 'lucide-react';
import { useNavigate } from '@tanstack/react-router';
import { useState } from 'react';
import { Tooltip, TooltipContent, TooltipTrigger } from '#/components/ui/tooltip.tsx';
import ModularActionDialog from '#/components/dialogs/modular-action-dialog.tsx';

interface ArchiveOrDeleteApplicationStatusProps {
	statusID: number;
	data: JobApplicationStatus;
	softDelete?: boolean;
	inView?: boolean;
}

export default function ArchiveOrDeleteApplicationStatus({
	data,
	statusID,
	softDelete = false,
	inView = false,
}: ArchiveOrDeleteApplicationStatusProps) {
	const queryClient = useQueryClient();
	const navigate = useNavigate();
	const [open, setOpen] = useState(false);
	const message = `You are ${!softDelete ? 'archiving' : 'deleting'} this Application Status. Archiving doesn't remove it from the visible data, while deleting removes it from the visible data and also removes all Job Applications and corresponding Notes and status history from the system, deleted data cannot be restored by you, only by Administrative users. You can always restore archived Application Status, and you don't loose visibility of Job Applications and corresponding Notes of archived Application Status. Are you sure?`;

	return (
		<ModularActionDialog
			open={open}
			// onOpenChange={setOpen}
			message={message}
			title={`${softDelete ? 'Delete' : 'Archive'} Application Status with id of ${data.id}`}
			onConfirm={async () => {
				await toast.promise(
					Promise.all([
						api.DELETE('/application_statuses/{statusID}', {
							params: {
								path: {
									statusID,
								},
								query: {
									soft_delete: softDelete,
								},
							},
						}),
					]),
					{
						loading: `${softDelete ? 'Deleting' : 'Archiving'} Application Status with id ${data.id}...`,
						success: `${softDelete ? 'Deleting' : 'Archiving'} Application Status id ${data.id}`,
						error: (response: APIResponse) =>
							`Error ${softDelete ? 'deleting' : 'archiving'} Application Status: ${response.message}`,
					},
				);

				await queryClient.refetchQueries({
					queryKey: ['get', '/application_statuses'],
				});
				await queryClient.refetchQueries({
					queryKey: ['get', '/jobs'],
				});
				if (inView) {
					await navigate({ to: '/application_status' });
				}
			}}
			trigger={
				<Tooltip>
					<TooltipTrigger
						render={
							<Button
								className="rounded-3xl p-3"
								size="icon-sm"
								variant="destructive"
								onClick={() => setOpen(true)}
							>
								<Trash />
								<span className="sr-only">
									{softDelete ? 'Delete' : 'Archive'} this Application Status
								</span>
							</Button>
						}
					/>
					<TooltipContent>
						<p>
							{softDelete
								? 'Delete this Application Status and all matching Job Applications'
								: 'Archive this Application Status'}
						</p>
					</TooltipContent>
				</Tooltip>
			}
			onClose={() => setOpen(false)}
			destructive
			actionLabel={`${softDelete ? 'Delete' : 'Archive'} Note`}
		/>
	);
}
