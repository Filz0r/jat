import type { APIResponse, JobApplicationStatus } from '#/api/types.ts';

import { toast } from '#components/ui/toast';
import { api } from '#/api/client.ts';
import { useQueryClient } from '@tanstack/react-query';
import { Button } from '#/components/ui/button.tsx';
import { useState } from 'react';
import { Tooltip, TooltipContent, TooltipTrigger } from '#/components/ui/tooltip.tsx';
import ModularInformationDialog from '#/components/dialogs/modular-action-dialog.tsx';
import { IconRestore } from '@tabler/icons-react';

interface UnarchiveApplicationStatusProps {
	statusID: number;
	data: JobApplicationStatus;
	inView?: boolean;
}

export default function UnarchiveApplicationStatus({
	data,
	statusID,
	inView = false,
}: UnarchiveApplicationStatusProps) {
	const queryClient = useQueryClient();
	const [open, setOpen] = useState(false);

	return (
		<ModularInformationDialog
			open={open}
			// onOpenChange={setOpen}
			message="You are un-archiving this Application Status. You will now be able to use this Application Status again, are you sure?"
			title={`Un-archive note with id of ${data.id}`}
			onConfirm={async () => {
				await toast.promise(
					Promise.all([
						api.PUT('/application_statuses/{statusID}/unarchive', {
							params: {
								path: {
									statusID,
								},
							},
						}),
					]),
					{
						loading: `Un-archiving Application Status with id ${data.id}...`,
						success: `Un-archiving Application Status id ${data.id}`,
						error: (response: APIResponse) =>
							`Error un-archiving Application Status: ${response.message}`,
					},
				);

				await queryClient.refetchQueries({
					queryKey: [
						'get',
						'/application_statuses',
						// { params: { query: { include_archived: inArchivedList } } },
					],
				});
				await queryClient.refetchQueries({
					queryKey: ['get', '/jobs'],
				});
				if (inView) {
					await queryClient.refetchQueries({
						queryKey: [
							'get',
							'/application_statuses/{statusID}',
							{ params: { path: { statusID } } },
						],
					});
				}
				setOpen(false);
			}}
			onClose={() => setOpen(false)}
			trigger={
				<Tooltip>
					<TooltipTrigger
						render={
							<Button
								onClick={() => setOpen(true)}
								className="bg-purple-500 hover:bg-purple-600"
								size="icon-sm"
							>
								<IconRestore />
							</Button>
						}
					/>
					<TooltipContent>
						<p>Un-archive this Application Status.</p>
					</TooltipContent>
				</Tooltip>
			}
		/>
	);
}
