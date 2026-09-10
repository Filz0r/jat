import { useForm } from '@tanstack/react-form';
import {
	defaultApplicationStatusSelectSchema,
	type DefaultApplicationStatusSelectSchema,
} from '#/schemas/users.ts';
import { api, apiClient } from '#/api/client.ts';
import { Skeleton } from '#/components/ui/skeleton.tsx';
import FieldWrapper from '#/components/field-wrapper.tsx';
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '#/components/ui/select.tsx';
import { Button } from '#/components/ui/button.tsx';
import { toast } from '#/components/ui/toast.tsx';
import type { APIResponse } from '#/api/types.ts';
import { useAuth } from '#/hooks/use-auth.ts';
import LogForm from '#/components/forms/log-form.tsx';

export default function DefaultApplicationStatusForm({ current }: { current: number }) {
	const { refreshUser } = useAuth();
	const form = useForm({
		defaultValues: {
			status_id: current,
		} satisfies DefaultApplicationStatusSelectSchema,
		validators: { onChange: defaultApplicationStatusSelectSchema },
		onSubmit: async ({ value }) => {
			await toast.promise(
				api.PUT('/users/default_status', {
					body: {
						status_id: value.status_id,
					},
				}),
				{
					loading: 'Updating default application status...',
					success: 'Updated default application status!',
					error: (response: APIResponse) =>
						`Error updating default status: ${response.message}`,
				},
			);
			await refreshUser();

			form.reset();
		},
	});

	const { data, error, isLoading } = apiClient.useQuery('get', '/application_statuses');

	if (isLoading) {
		return <Skeleton className="size-fit" />;
	}

	if (error || !data || !data.data || !data.ok) {
		return (
			<Skeleton
				className="size-fit"
				children={
					<div className="text-destructive">
						Error loading application data:{' '}
						{error ? error.message : data?.message || 'Unexpected error'}
					</div>
				}
			/>
		);
	}

	const applicationStatus = data.data.map((status) => ({
		label: status.status,
		value: status.id,
	}));
	return (
		<form
			onSubmit={async (e) => {
				e.preventDefault();
				e.stopPropagation();
				await form.handleSubmit();
			}}
			className="space-y-2"
		>
			<form.Field
				name="status_id"
				validators={{ onChange: defaultApplicationStatusSelectSchema.shape.status_id }}
				children={(field) => (
					<FieldWrapper field={field} label="Default Application Status">
						<Select
							value={field.state.value}
							onValueChange={(value) => {
								if (value) {
									form.setFieldValue('status_id', value);
								}
							}}
						>
							<SelectTrigger className="w-full">
								<SelectValue>
									{
										applicationStatus.find((s) => s.value === field.state.value)
											?.label
									}
								</SelectValue>
							</SelectTrigger>
							<SelectContent>
								<SelectGroup>
									{applicationStatus.map((s) => (
										<SelectItem key={s.value} value={s.value}>
											{s.label}
										</SelectItem>
									))}
								</SelectGroup>
							</SelectContent>
						</Select>
					</FieldWrapper>
				)}
			/>
			<div className="mx-2 mt-4 flex gap-x-2">
				<form.Subscribe
					selector={(s) => s.values.status_id}
					children={(status_id) => (
						<Button
							className="flex-1"
							size="lg"
							onClick={() => form.reset()}
							disabled={status_id === current}
						>
							Reset
						</Button>
					)}
				/>
				<form.Subscribe
					selector={(state) => [state.canSubmit, state.isSubmitting]}
					children={([canSubmit, isSubmitting]) => (
						<Button
							type="submit"
							disabled={
								!canSubmit ||
								form.getFieldValue('status_id') === current ||
								isSubmitting
							}
							size="lg"
							className="flex-1 cursor-pointer bg-green-500 hover:bg-green-600"
						>
							{isSubmitting ? 'Updating default Application Status...' : 'Save'}
						</Button>
					)}
				/>
				<LogForm form={form} />
			</div>
		</form>
	);
}
