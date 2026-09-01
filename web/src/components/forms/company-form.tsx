import type { CompanySchema } from '#/schemas/job-applications.ts';
import type { CompanyData } from '#/api/types.ts';

import { companySchema } from '#/schemas/job-applications.ts';
import { useQueryClient } from '@tanstack/react-query';
import { useForm } from '@tanstack/react-form';
import ModularFormDialog from '#/components/dialogs/modular-form-dialog.tsx';
import { useState } from 'react';
import FieldWrapper from '#/components/field-wrapper.tsx';
import { Tooltip, TooltipContent, TooltipTrigger } from '#/components/ui/tooltip.tsx';
import { Button } from '#/components/ui/button.tsx';
import { api } from '#/api/client.ts';
import { toast } from '#/components/ui/toast.tsx';
import { Input } from '#/components/ui/input.tsx';
import HttpLinkInput from '#/components/http-link-input.tsx';
import { IconBuilding, IconBuildingPlus } from '@tabler/icons-react';
import { useAuth } from '#/hooks/use-auth.ts';

interface CompanyFormProps {
	edit?: boolean;
	existingData?: CompanyData;
	refreshSelf?: boolean;
	largeTrigger?: boolean;
}

export default function CompanyForm({
	existingData,
	edit = false,
	refreshSelf = false,
	largeTrigger = false,
}: CompanyFormProps) {
	const { user } = useAuth();
	const isAdmin = user?.is_admin || false;
	const queryClient = useQueryClient();
	const [isOpen, setIsOpen] = useState(false);
	const [errorMessage, setErrorMessage] = useState<string | null>(null);

	const form = useForm({
		defaultValues:
			edit && existingData
				? ({
						name: existingData.name,
						website: existingData.website ?? undefined,
					} satisfies CompanySchema)
				: ({ name: '' } satisfies CompanySchema),
		validators: {
			onChange: companySchema,
		},
		onSubmit: async ({ value }) => {
			if (edit && existingData) {
				const { data, error } = await api.PUT('/company/{companyID}', {
					params: {
						path: {
							companyID: existingData.id,
						},
					},
					body: {
						name: value.name,
						website: value.website,
					},
				});
				if (error || !data.ok) {
					setErrorMessage(error?.message || data?.message || 'Error updating company');
					return;
				}
			} else {
				const { data, error } = await api.POST('/company', {
					body: {
						name: value.name,
						website: value.website,
					},
				});
				if (error || !data.ok) {
					setErrorMessage(error?.message || data?.message || 'Error saving new company');
					return;
				}
			}

			setIsOpen(false);
			setErrorMessage(null);

			if (refreshSelf && existingData) {
				await queryClient.refetchQueries({
					queryKey: [
						'get',
						'/company/{companyID}',
						{ params: { path: { companyID: existingData.id } } },
					],
				});
				if (isAdmin) {
					await queryClient.refetchQueries({
						queryKey: [
							'get',
							'/company/{companyID}/history',
							{ params: { path: { companyID: existingData.id } } },
						],
					});
				}
			}
			await queryClient.refetchQueries({
				queryKey: ['get', '/company'],
			});

			form.reset();
			toast.add({
				title:
					edit && existingData
						? `Updated note with id of ${existingData.id}`
						: 'A new note was created',
			});
		},
	});

	return (
		<ModularFormDialog
			open={isOpen}
			onClose={() => {
				setIsOpen(false);
				setErrorMessage(null);
				form.reset();
			}}
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
									{edit ? <IconBuildingPlus /> : <IconBuilding />}
									<span className="sr-only">
										{edit ? 'Edit this company' : 'Add a new company'}
									</span>
								</Button>
							}
						/>
						<TooltipContent>
							<p>{edit ? 'Edit this company' : 'Add a new company'}</p>
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
						{edit ? 'Edit this company' : 'Add a new company'}
						{edit ? <IconBuildingPlus /> : <IconBuilding />}
						<span className="sr-only">
							{edit
								? 'Update this notes content'
								: 'Add a Note to this Job Application'}
						</span>
					</Button>
				)
			}
			title={edit ? 'Edit this company' : 'Add a new company'}
			description={
				edit ? undefined : 'You can add a new company to the system using this form.'
			}
			form={form}
			Content={
				<div className="space-y-3">
					{errorMessage && (
						<p className="text-destructive mb-4 text-sm">{errorMessage}</p>
					)}
					<form.Field
						name="name"
						validators={{
							onChange: companySchema.shape.name,
						}}
						children={(field) => (
							<FieldWrapper field={field} label="Name">
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
						name="website"
						validators={{
							onChange: companySchema.shape.website,
						}}
						children={(field) => (
							<FieldWrapper field={field} label="Website">
								<HttpLinkInput
									id={field.name}
									name={field.name}
									placeholder="example.com"
									value={field.state.value}
									onBlur={field.handleBlur}
									onChange={(value) => field.handleChange(value || undefined)}
								/>
							</FieldWrapper>
						)}
					/>
				</div>
			}
		/>
	);
}
