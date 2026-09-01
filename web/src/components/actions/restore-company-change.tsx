import type { APIResponse } from '#/api/types.ts';

import { toast } from '#components/ui/toast';
import { api } from '#/api/client.ts';
import { useQueryClient } from '@tanstack/react-query';
import { Button } from '#/components/ui/button.tsx';
import { useState } from 'react';
import { Tooltip, TooltipContent, TooltipTrigger } from '#/components/ui/tooltip.tsx';
import ModularInformationDialog from '#/components/dialogs/modular-action-dialog.tsx';
import { IconRestore } from '@tabler/icons-react';

interface RestoreCompanyChangeProps {
	companyID: number;
	changeID: number;
}

export default function RestoreCompanyChange({ changeID, companyID }: RestoreCompanyChangeProps) {
	const queryClient = useQueryClient();
	const [open, setOpen] = useState(false);

	return (
		<ModularInformationDialog
			open={open}
			message="You are Reverting this change to this company data, are you sure? You will also be able to revert this change again after the page is refreshed"
			title={`Revert change with ID ${changeID} on company with ID ${companyID}`}
			onConfirm={async () => {
				await toast.promise(
					Promise.all([
						api.PUT('/company/{companyID}/history/{changeID}', {
							params: {
								path: {
									companyID,
									changeID,
								},
							},
						}),
					]),
					{
						loading: `Reverting change with ID ${changeID}...`,
						success: `Reverted change with ID ${changeID}`,
						error: (response: APIResponse) =>
							`Error un-archiving Application Status: ${response.message}`,
					},
				);

				await queryClient.refetchQueries({
					queryKey: [
						'get',
						'/company/{companyID}/history',
						{ params: { path: { companyID } } },
					],
				});
				await queryClient.refetchQueries({
					queryKey: ['get', '/company/{companyID}', { params: { path: { companyID } } }],
				});
				setOpen(false);
			}}
			onClose={() => setOpen(false)}
			trigger={
				<Tooltip>
					<TooltipTrigger
						render={
							<Button
								onClick={() => setOpen(true)}
								className="bg-orange-500 hover:bg-orange-600"
								size="icon-sm"
							>
								<IconRestore />
								<span className="sr-only">Revert this change.</span>
							</Button>
						}
					/>
					<TooltipContent>
						<p>Revert this change.</p>
					</TooltipContent>
				</Tooltip>
			}
		/>
	);
}
