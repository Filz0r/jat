import type { APIResponse } from '#/api/types.ts';

import { toast } from '#components/ui/toast';
import { api } from '#/api/client.ts';
import { useQueryClient } from '@tanstack/react-query';
import { Button } from '#/components/ui/button.tsx';
import { Trash } from 'lucide-react';
import { Tooltip, TooltipContent, TooltipTrigger } from '#/components/ui/tooltip.tsx';
import { useState } from 'react';
import ModularActionDialog from '#/components/dialogs/modular-action-dialog.tsx';

interface DeleteCompanyProps {
	companyID: number;
}

export default function DeleteCompany({ companyID }: DeleteCompanyProps) {
	const queryClient = useQueryClient();
	const [open, setOpen] = useState(false);
	return (
		<ModularActionDialog
			open={open}
			onClose={() => setOpen(false)}
			destructive
			actionLabel="Delete Company"
			message="You are deleting this company, this is will cause all job applications with this company to also be deleted, are you sure?"
			title={`Delete Company with id of ${companyID}`}
			onConfirm={async () => {
				await toast.promise(
					Promise.all([
						api.DELETE('/company/{companyID}', {
							params: {
								path: {
									companyID,
								},
							},
						}),
					]),
					{
						loading: `Deleting Company with id ${companyID}...`,
						success: `Deleted Company with id ${companyID}`,
						error: (response: APIResponse) =>
							`Error deleting note: ${response.message}`,
					},
				);
				await queryClient.refetchQueries({
					queryKey: ['get', '/company'],
				});
				await queryClient.refetchQueries({
					queryKey: ['get', '/jobs'],
				});
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
								<span className="sr-only">Delete this Company</span>
							</Button>
						}
					/>
					<TooltipContent>
						<p>Delete this Company</p>
					</TooltipContent>
				</Tooltip>
			}
		/>
	);
}
