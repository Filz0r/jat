import type { APIResponse, JobApplicationStatus } from '#/api/types.ts';
import type { ApplicationStatusSchema } from '#/schemas/job-applications.ts';

import { useForm } from '@tanstack/react-form';
import {
	applicationStatusKindSchema,
	jobApplicationStatusSchema,
} from '#/schemas/job-applications.ts';
import ModularFormDialog from '#/components/dialogs/modular-form-dialog.tsx';
import { Tooltip, TooltipContent, TooltipTrigger } from '#/components/ui/tooltip.tsx';
import { Button } from '#/components/ui/button.tsx';
import { IconFilter2Plus, IconFilter2Edit } from '@tabler/icons-react';
import { useState } from 'react';
import FieldWrapper from '#/components/field-wrapper.tsx';
import { Input } from '#/components/ui/input.tsx';
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '#/components/ui/select.tsx';
import { toast } from '#/components/ui/toast.tsx';
import { api } from '#/api/client.ts';
import { useQueryClient } from '@tanstack/react-query';

interface ApplicationStatusFormProps {
	largeTrigger?: boolean;
	edit?: boolean;
	data?: JobApplicationStatus;
	inView?: boolean;
}

export default function ApplicationStatusForm({
	largeTrigger = false,
	edit = false,
	inView = false,
	data,
}: ApplicationStatusFormProps) {
	const queryClient = useQueryClient();
	const [isOpen, setIsOpen] = useState(false);
	const form = useForm({
		defaultValues: {
			kind: edit && data ? data.kind : applicationStatusKindSchema.enum.applied,
			status: edit && data ? data.status : '',
		} satisfies ApplicationStatusSchema,
		validators: {
			onChange: jobApplicationStatusSchema,
		},
		onSubmit: async ({ value }) => {
			await toast.promise(
				Promise.all([
					edit && data
						? api.PUT('/application_statuses/{statusID}', {
								params: {
									path: {
										statusID: data.id,
									},
								},
								body: {
									status: value.status,
									kind: value.kind,
								},
							})
						: api.POST('/application_statuses', {
								body: {
									status: value.status,
									kind: value.kind,
								},
							}),
				]),
				{
					loading:
						edit && data
							? `Saving changes to Application Status with id: ${data.id}...`
							: 'Creating new Application status...',
					success:
						edit && data
							? `Saved changes to Application Status with id: ${data.id}`
							: 'Created new Application Status!',
					error: (response: APIResponse) => {
						if (edit && data) {
							return `Error updating Application Status with ID ${data.id}: ${response.message}`;
						} else {
							return `Error creating new application Status: ${response.message}`;
						}
					},
				},
			);
			if (inView && edit && data) {
				await queryClient.refetchQueries({
					queryKey: [
						'get',
						'/application_statuses/$statusID',
						{ params: { path: { statusID: data.id } } },
					],
				});
			}
			await queryClient.refetchQueries({
				queryKey: ['get', '/application_statuses'],
			});
			form.reset();
			setIsOpen(false);
		},
	});

	const kindValues = Object.values(applicationStatusKindSchema.enum).map((v) => ({
		label: v.toUpperCase(),
		value: v,
	}));

	return (
		<ModularFormDialog
			open={isOpen}
			onClose={() => {
				setIsOpen(false);
				form.reset();
			}}
			title={
				edit && data
					? `Updating Application Status with id of ${data.id}`
					: 'Create a new Application Status'
			}
			description={
				edit
					? undefined
					: 'You can create application status to attach to your job applications'
			}
			form={form}
			TriggerButton={
				!largeTrigger ? (
					<Tooltip>
						<TooltipTrigger
							render={
								<Button
									size="icon-sm"
									className={
										edit
											? 'bg-amber-500 hover:bg-amber-700'
											: 'bg-green-500 hover:bg-green-700'
									}
									onClick={() => setIsOpen(true)}
								>
									{edit ? <IconFilter2Edit /> : <IconFilter2Plus />}
									<span className="sr-only">
										{edit
											? 'Update this Application Status'
											: 'Add a new Application Status'}
									</span>
								</Button>
							}
						/>
						<TooltipContent>
							<p>
								{edit
									? 'Update this Application Status'
									: 'Add a new Application Status'}
							</p>
						</TooltipContent>
					</Tooltip>
				) : (
					<Button
						className={
							edit
								? 'bg-amber-500 hover:bg-amber-700'
								: 'bg-green-500 hover:bg-green-700'
						}
						onClick={() => setIsOpen(true)}
					>
						{edit ? 'Update this Application Status' : 'Add a new Application Status'}
						{edit ? <IconFilter2Edit /> : <IconFilter2Plus />}
						<span className="sr-only">
							{edit
								? 'Update this Application Status'
								: 'Add a new Application Status'}
						</span>
					</Button>
				)
			}
			Content={
				<div className="space-y-3">
					<form.Field
						name="status"
						validators={{
							onChange: jobApplicationStatusSchema.shape.status,
						}}
						children={(field) => (
							<FieldWrapper field={field} label="Status Name">
								<Input
									id={field.name}
									name={field.name}
									value={field.state.value}
									onBlur={field.handleBlur}
									onChange={(e) => field.handleChange(e.target.value)}
								/>
							</FieldWrapper>
						)}
					/>
					<form.Field
						name="kind"
						validators={{
							onChange: jobApplicationStatusSchema.shape.kind,
						}}
						children={(field) => (
							<FieldWrapper field={field} label="Kind" className="flex-1">
								<Select
									value={field.state.value}
									onValueChange={(value) => {
										if (value) {
											form.setFieldValue('kind', value);
										}
									}}
								>
									<SelectTrigger className="w-full">
										<SelectValue>
											{
												kindValues.find(
													(s) => s.value === field.state.value,
												)?.label
											}
										</SelectValue>
									</SelectTrigger>
									<SelectContent>
										<SelectGroup>
											{kindValues.map((kind) => (
												<SelectItem key={kind.value} value={kind.value}>
													{kind.label}
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
