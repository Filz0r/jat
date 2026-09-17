import type { APIResponse } from '#/api/types.ts';
import { useQueryClient } from '@tanstack/react-query';
import { api } from '#/api/client.ts';
import { Button } from '#/components/ui/button.tsx';
import { Trash } from 'lucide-react';
import { toast } from '#/components/ui/toast.tsx';
import { useNavigate } from '@tanstack/react-router';
import { Tooltip, TooltipContent, TooltipTrigger } from '#/components/ui/tooltip.tsx';
import { useState } from 'react';
import ModularActionDialog from '#/components/dialogs/modular-action-dialog.tsx';

interface DeleteJobApplicationProps {
	jobID: number;
	navigateToJobs?: boolean;
}

export default function DeleteJobApplication({
	jobID,
	navigateToJobs = false,
}: DeleteJobApplicationProps) {
	const queryClient = useQueryClient();
	const navigate = useNavigate();
	const [open, setOpen] = useState(false);
	return (
		<ModularActionDialog
			destructive
			open={open}
			onClose={() => setOpen(false)}
			message="You are deleting this job application, this is only reversible by administrative users, are you sure?"
			title={`Delete Job Application with ID: ${jobID}`}
			onConfirm={async () => {
				await toast.promise(
					Promise.all([
						api.DELETE('/jobs/{jobID}', {
							params: {
								path: {
									jobID,
								},
							},
						}),
					]),
					{
						loading: `Deleting job with id ${jobID}...`,
						success: `Deleted job with id ${jobID}`,
						error: (response: APIResponse) =>
							`Error deleting note: ${response.message}`,
					},
				);
				await queryClient.refetchQueries({
					queryKey: ['get', '/jobs'],
				});
				if (navigateToJobs) {
					await navigate({ to: '/jobs', replace: true });
				}
				setOpen(false);
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
								<span className="sr-only">Delete this Job Application</span>
							</Button>
						}
					/>
					<TooltipContent>
						<p>Delete this Job Application</p>
					</TooltipContent>
				</Tooltip>
			}
			actionLabel="Delete Job Application"
		/>
	);
}
