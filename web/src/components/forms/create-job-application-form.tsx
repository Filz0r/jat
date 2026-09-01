import { useForm } from '@tanstack/react-form';
import { useQueryClient } from '@tanstack/react-query';
import { createJobApplicationSchema } from '#/schemas/job-applications.ts';
import type { CreateJobApplication } from '#/schemas/job-applications.ts';
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
import { Input } from '#/components/ui/input.tsx';
import { DatePicker } from '#/components/date-picker.tsx';
import { CreatableSelect } from '#/components/creatable-select.tsx';
import LogForm from '#/components/forms/log-form.tsx';
import { Card, CardContent, CardFooter, CardHeader } from '#/components/ui/card.tsx';
import { Button } from '#components/ui/button';
import { useState } from 'react';
import { useNavigate } from '@tanstack/react-router';

export default function CreateJobApplicationForm({
	defaultStatus,
	preSelectedCompany,
}: {
	defaultStatus: number;
	preSelectedCompany: number | undefined;
}) {
	const navigate = useNavigate();
	const queryClient = useQueryClient();
	const [serverError, setServerError] = useState<string | null>(null);

	const form = useForm({
		defaultValues: {
			title: '',
			url: '',
			createdAt: new Date(),
			company: preSelectedCompany ? preSelectedCompany : 0,
			status: defaultStatus,
		} satisfies CreateJobApplication,
		validators: {
			onChange: createJobApplicationSchema,
		},
		onSubmit: async ({ value }) => {
			setServerError(null);

			const { data: createData, error: createError } = await api.POST('/jobs', {
				body: {
					company_id: value.company,
					title: value.title,
					status_id: value.status,
					url: value.url,
					created_at: value.createdAt.toISOString(),
				},
			});

			if (createError || !createData.ok || !createData.data) {
				setServerError(createError ? createError.message : createData.message);
				return;
			}
			form.reset();
			await queryClient.refetchQueries({
				queryKey: ['get', '/jobs'],
			});
			void navigate({ to: '/jobs', replace: true });
		},
	});

	const {
		data: dataStatus,
		isError: isErrorStatus,
		isLoading: isLoadingStatus,
	} = apiClient.useQuery('get', '/application_statuses');

	const {
		data: dataCompanies,
		isError: isErrorCompanies,
		isLoading: isLoadingCompanies,
	} = apiClient.useQuery('get', '/company');

	const createCompany = apiClient.useMutation('post', '/company');

	if (
		(!isLoadingStatus && isErrorStatus) ||
		!dataStatus ||
		dataStatus.data === undefined ||
		(!isLoadingCompanies && isErrorCompanies) ||
		!dataCompanies ||
		dataCompanies.data === undefined
	) {
		return <h1>Error loading application data</h1>;
	}
	const applicationStatus = dataStatus.data.map((status) => ({
		label: status.status,
		value: status.id,
	}));

	const companies = dataCompanies.data.map((company) => ({
		id: company.id,
		label: company.name,
	}));

	return (
		// eslint-disable-next-line tailwindcss/no-arbitrary-value
		<div className="flex min-h-[calc(100vh-2rem)] w-full p-4 sm:min-h-0 sm:justify-center">
			<form
				className="flex w-full flex-col justify-start"
				noValidate
				onSubmit={(e) => {
					e.preventDefault();
					e.stopPropagation();
					form.handleSubmit();
				}}
			>
				<Card className="flex w-full flex-col sm:h-fit">
					<CardHeader>
						<h1 className="text-center text-2xl font-bold sm:text-3xl">
							Create Job Application
						</h1>
					</CardHeader>
					<CardContent className="flex-1 space-y-4 sm:flex-initial">
						<form.Field
							name="title"
							validators={{ onChange: createJobApplicationSchema.shape.title }}
							children={(field) => (
								<FieldWrapper field={field} label={'Title'}>
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
							name="url"
							validators={{ onChange: createJobApplicationSchema.shape.url }}
							children={(field) => (
								<FieldWrapper field={field} label="Posting URL">
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
							name="company"
							validators={{
								onChange: createJobApplicationSchema.shape.company,
							}}
							children={(field) => (
								<CreatableSelect
									id={field.name}
									name={field.name}
									value={field.state.value}
									onChange={(value) => field.handleChange(value)}
									onBlur={field.handleBlur}
									label="Company"
									placeholder="Select a company"
									searchPlaceholder="Search companies..."
									createButtonLabel="Create new company"
									emptyMessage="No companies found. Type a name to create one."
									items={companies}
									isLoading={isLoadingCompanies}
									isError={isErrorCompanies}
									onCreateItem={async (name) => {
										const company = await createCompany.mutateAsync({
											body: { name },
										});

										// Wait for the refetch so the newly created company is
										// present in the items list before selecting it.
										await queryClient.refetchQueries({
											queryKey: ['get', '/company'],
										});

										return company.data?.id;
									}}
								/>
							)}
						/>
						<div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:gap-3">
							<form.Field
								name="status"
								validators={{
									onChange: createJobApplicationSchema.shape.status,
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
														applicationStatus.find(
															(s) => s.value === field.state.value,
														)?.label
													}
												</SelectValue>
											</SelectTrigger>
											<SelectContent>
												<SelectGroup>
													{applicationStatus.map((status) => (
														<SelectItem
															key={status.value}
															value={status.value}
														>
															{status.label}
														</SelectItem>
													))}
												</SelectGroup>
											</SelectContent>
										</Select>
									</FieldWrapper>
								)}
							/>
							<form.Field
								name="createdAt"
								validators={{
									onChange: createJobApplicationSchema.shape.createdAt,
								}}
								children={(field) => (
									<FieldWrapper
										field={field}
										label="Applied on"
										className="flex-1"
									>
										<DatePicker
											id={field.name}
											name={field.name}
											value={field.state.value}
											onBlur={field.handleBlur}
											onChange={(date) => {
												if (date) {
													field.handleChange(date);
												}
											}}
											className="w-full"
										/>
									</FieldWrapper>
								)}
							/>
						</div>
						{serverError && <p className="text-destructive text-xs">{serverError}</p>}
					</CardContent>
					<CardFooter className="flex flex-col gap-3 sm:flex-row sm:justify-stretch">
						<form.Subscribe
							selector={(state) => [state.canSubmit, state.isSubmitting]}
							children={([canSubmit, isSubmitting]) => (
								<Button
									type="submit"
									className="w-full sm:min-w-32 sm:flex-1"
									disabled={!canSubmit || isSubmitting}
								>
									{isSubmitting ? 'Creating job posting...' : 'Save'}
								</Button>
							)}
						/>

						<Button
							type="reset"
							className="w-full bg-red-500 hover:bg-red-700 sm:min-w-32 sm:flex-1"
							onClick={(e) => {
								e.preventDefault();
								e.stopPropagation();
								form.reset();
							}}
						>
							Reset
						</Button>
						<div className="w-full sm:flex-1">
							<LogForm form={form} />
						</div>
					</CardFooter>
				</Card>
			</form>
		</div>
	);
}
