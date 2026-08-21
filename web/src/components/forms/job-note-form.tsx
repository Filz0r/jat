import type { CreateJobApplicationNote } from '#/schemas/job-applications.ts';

import { createJobApplicationNoteSchema } from '#/schemas/job-applications.ts';
import { useQueryClient } from '@tanstack/react-query';
import { useForm } from '@tanstack/react-form';
import ModularFormDialog from '#/components/dialogs/modular-form-dialog.tsx';
import { useState } from 'react';
import FieldWrapper from '#/components/field-wrapper.tsx';
import { Textarea } from '#/components/ui/textarea.tsx';
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from '#/components/ui/tooltip.tsx';
import { Button } from '#/components/ui/button.tsx';
import { FilePlusCorner, FileCodeCorner } from 'lucide-react';
import { api } from '#/api/client.ts';
import type { JobApplicationNote } from '#/api/types.ts';

interface JobNoteFormPropsBase {
	jobID: number;
	refreshSelf?: boolean;
	largeTrigger?: boolean;
}

interface CreateJobNoteFormProps extends JobNoteFormPropsBase {
	edit?: false;
}

interface EditJobNoteFormProps extends JobNoteFormPropsBase {
	edit: true;
	existingData: JobApplicationNote;
}

type JobNoteFormProps = CreateJobNoteFormProps | EditJobNoteFormProps;

export default function JobNoteForm(props: JobNoteFormProps) {
	const { jobID, refreshSelf = false, largeTrigger = false, edit = false } = props;
	const existingData = props.edit ? props.existingData : undefined;

	const queryClient = useQueryClient();
	const [isOpen, setIsOpen] = useState(false);
	const [errorMessage, setErrorMessage] = useState<string | null>(null);

	const form = useForm({
		defaultValues: {
			body: edit && existingData ? existingData.body : '',
		} satisfies CreateJobApplicationNote,
		validators: {
			onChange: createJobApplicationNoteSchema,
		},
		onSubmit: async ({ value }) => {
			if (edit && existingData) {
				const { data, error } = await api.PUT('/jobs/{jobID}/notes/{noteID}', {
					params: {
						path: {
							jobID,
							noteID: existingData.id,
						},
					},
					body: {
						body: value.body,
					},
				});
				if (error || !data.ok) {
					setErrorMessage(
						error?.message || data?.message || 'Error updating note content',
					);
					return;
				}
			} else {
				const { data, error } = await api.POST('/jobs/{jobID}/notes', {
					params: {
						path: {
							jobID,
						},
					},
					body: {
						body: value.body,
					},
				});
				if (error || !data.ok) {
					setErrorMessage(
						error?.message ||
							data?.message ||
							'Error saving the note for this application',
					);
					return;
				}
			}

			setIsOpen(false);
			setErrorMessage(null);

			if (refreshSelf) {
				await queryClient.refetchQueries({
					queryKey: ['get', '/jobs/{jobID}', { params: { path: { jobID } } }],
				});
			}
			await queryClient.refetchQueries({
				queryKey: ['get', '/jobs/{jobID}/notes', { params: { path: { jobID } } }],
			});

			// Todo: uncomment bellow after adding note count to jobs fetching
			// await queryClient.refetchQueries({ queryKey: ['get', '/jobs'] });
			// TODO: add a toast on success
		},
	});

	return (
		<ModularFormDialog
			open={isOpen}
			onClose={() => {
				setIsOpen(false);
				setErrorMessage(null);
			}}
			TriggerButton={
				!largeTrigger ? (
					<TooltipProvider>
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
										{edit ? <FileCodeCorner /> : <FilePlusCorner />}
										<span className="sr-only">
											{edit
												? 'Update this notes content'
												: 'Add a Note to this Job Application'}
										</span>
									</Button>
								}
							/>
							<TooltipContent>
								<p>
									{edit
										? 'Update this notes content'
										: 'Add a Note to this Job Application'}
								</p>
							</TooltipContent>
						</Tooltip>
					</TooltipProvider>
				) : (
					<div>
						<Button
							className={
								edit
									? 'bg-amber-500 hover:bg-amber-700'
									: 'bg-green-500 hover:bg-green-700'
							}
							onClick={() => setIsOpen(true)}
						>
							{edit
								? 'Update this notes content'
								: 'Add a Note to this Job Application'}
							<FilePlusCorner />
							<span className="sr-only">
								{edit
									? 'Update this notes content'
									: 'Add a Note to this Job Application'}
							</span>
						</Button>
					</div>
				)
			}
			title={
				edit && existingData
					? `Updating note with id ${existingData.id} for job with ID ${jobID}`
					: `Insert note for job with ID ${jobID}`
			}
			description={
				edit
					? undefined
					: 'You can add notes to your job application using this form, this is useful for information such as recruiter names, detailed position information, etc.'
			}
			form={form}
			Content={
				<div>
					{errorMessage && (
						<p className="mb-4 text-sm text-destructive">{errorMessage}</p>
					)}
					<form.Field
						name="body"
						validators={{
							onChange: createJobApplicationNoteSchema.shape.body,
						}}
						children={(field) => (
							<FieldWrapper field={field} label="Body">
								<Textarea
									id={field.name}
									name={field.name}
									value={field.state.value}
									onBlur={field.handleBlur}
									onChange={(e) => field.handleChange(e.target.value)}
								/>
							</FieldWrapper>
						)}
					/>
				</div>
			}
		/>
	);
}
