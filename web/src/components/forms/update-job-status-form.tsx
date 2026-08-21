import type { UpdateJobApplicationStatusSchema } from '#/schemas/job-applications.ts';

import { useForm } from '@tanstack/react-form';
import { updateJobApplicationStatusSchema } from '#/schemas/job-applications.ts';
import { api, apiClient } from '#/api/client.ts';
import FieldWrapper from '#/components/field-wrapper.tsx';
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '#/components/ui/select.tsx';
import ModularInformationDialog from '#/components/dialogs/modular-information-dialog.tsx';
import ModularFormDialog from '#/components/dialogs/modular-form-dialog.tsx';
import { useMemo, useState } from 'react';
import { Button } from '#/components/ui/button.tsx';
import { ChartNoAxesColumn } from 'lucide-react';
import { useQueryClient } from '@tanstack/react-query';
import {
	Tooltip,
	TooltipProvider,
	TooltipTrigger,
	TooltipContent,
} from '#/components/ui/tooltip.tsx';

interface UpdateJobApplicationStatusFormProps {
	currentStatus: number;
	jobID: number;
	refreshSelf?: boolean;
}

export default function UpdateJobStatusForm({
	currentStatus,
	jobID,
	refreshSelf = false,
}: UpdateJobApplicationStatusFormProps) {
	const queryClient = useQueryClient();
	const [isOpen, setIsOpen] = useState(false);
	const [errorMessage, setErrorMessage] = useState<string | null>(null);

	const form = useForm({
		defaultValues: {
			status: currentStatus,
		} satisfies UpdateJobApplicationStatusSchema,
		validators: {
			onChange: updateJobApplicationStatusSchema,
		},
		onSubmit: async ({ value }) => {
			const { data, error } = await api.PUT('/jobs/{jobID}/status/{statusID}', {
				params: {
					path: {
						jobID,
						statusID: value.status,
					},
				},
			});

			if (error || !data.ok) {
				setErrorMessage(
					error?.message ||
						data?.message ||
						'Error saving the status for this application',
				);
				return;
			}

			setIsOpen(false);
			setErrorMessage(null);

			if (refreshSelf) {
				await queryClient.refetchQueries({ queryKey: ['get', '/jobs/{jobID}'] });
			}
			await queryClient.refetchQueries({ queryKey: ['get', '/jobs'] });
			// TODO: add a toast on success
		},
	});

	const { data, isError } = apiClient.useQuery('get', '/application_statuses');

	const convertedData = useMemo(() => {
		if (!data || !data.ok || !data.data) return [];
		return data.data.map((status) => ({
			label: status.status,
			value: status.id,
		}));
	}, [data]);

	const queryErrorMessage =
		isError || !data || !data.ok || !data.data
			? data?.message || 'Error loading application statuses'
			: null;

	if (queryErrorMessage && isOpen) {
		return (
			<ModularInformationDialog
				open
				title="Error Loading Data"
				message={queryErrorMessage}
				onClose={() => setIsOpen(false)}
			/>
		);
	}

	return (
		<ModularFormDialog
			open={isOpen}
			onClose={() => {
				setIsOpen(false);
				setErrorMessage(null);
			}}
			TriggerButton={
				<TooltipProvider>
					<Tooltip>
						<TooltipTrigger
							render={
								<Button
									size="icon-sm"
									className="bg-amber-500 hover:bg-amber-700"
									onClick={() => setIsOpen(true)}
								>
									<ChartNoAxesColumn />
									<span className="sr-only">Update Job Application Status</span>
								</Button>
							}
						/>
						<TooltipContent>
							<p>Update Job Application Status</p>
						</TooltipContent>
					</Tooltip>
				</TooltipProvider>
			}
			title={`Update status of job with id ${jobID}`}
			description="You can update the status of the current job application so that you have better traceability during your job search."
			form={form}
			Content={
				<div>
					{errorMessage && (
						<p className="text-destructive mb-4 text-sm">{errorMessage}</p>
					)}
					<form.Field
						name="status"
						validators={{
							onChange: updateJobApplicationStatusSchema.shape.status,
						}}
						children={(field) => (
							<FieldWrapper field={field} label="Status" className="flex-1">
								<Select
									value={field.state.value}
									onValueChange={(value) => {
										if (value) {
											form.setFieldValue('status', value);
										}
									}}
								>
									<SelectTrigger className="w-full">
										<SelectValue>
											{
												convertedData.find(
													(s) => s.value === field.state.value,
												)?.label
											}
										</SelectValue>
									</SelectTrigger>
									<SelectContent>
										<SelectGroup>
											{convertedData.map((status) => (
												<SelectItem key={status.value} value={status.value}>
													{status.label}
												</SelectItem>
											))}
										</SelectGroup>
									</SelectContent>
								</Select>
							</FieldWrapper>
						)}
					/>
				</div>
			}
		/>
	);
}
