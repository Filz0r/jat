import { useForm } from '@tanstack/react-form';
import { Button } from '#/components/ui/button.tsx';
import { Input } from '#/components/ui/input.tsx';
import { CardContent, CardFooter } from '#/components/ui/card.tsx';
import { api } from '#/api/client.ts';
import { updateAccountSchema, type UpdateAccountSchema } from '#/schemas/users.ts';
import { useState } from 'react';
import FieldWrapper from '#/components/field-wrapper.tsx';
import { toast } from '#/components/ui/toast.tsx';

interface UpdateAccountFormProps {
	onSuccess?: () => Promise<void> | void;
	username: string;
	email: string;
}

export function UpdateAccountForm({ onSuccess, username, email }: UpdateAccountFormProps) {
	const [serverError, setServerError] = useState<string | null>(null);

	const defaultValues: UpdateAccountSchema = {
		username,
		email,
		password: '',
		confirmPassword: '',
	};

	const form = useForm({
		defaultValues,
		validators: {
			onChange: updateAccountSchema,
		},
		onSubmit: async ({ value }) => {
			setServerError(null);

			const { data: registerData, error: registerError } = await api.PUT('/users', {
				body: {
					username: value.username !== username ? value.username : undefined,
					email: value.email !== email ? value.email : undefined,
					password: value.password ? value.password : undefined,
				},
			});

			if (
				(registerError && !registerError.ok) ||
				(registerData && !registerData.ok && !registerData.data)
			) {
				setServerError(
					registerData?.message ?? 'Failed to create account. Please try again.',
				);
				return;
			}

			// await queryClient.refetchQueries({ queryKey: ['me'] });
			if (onSuccess) {
				await onSuccess();
			}
			toast.add({ description: 'Updated user settings', timeout: 2500 });
			form.reset();
		},
	});

	return (
		<form
			onSubmit={async (e) => {
				e.preventDefault();
				e.stopPropagation();
				await form.handleSubmit();
			}}
		>
			<CardContent className="mx-2 mt-2 mb-4 flex flex-col gap-4">
				<form.Field
					name="username"
					validators={{
						onChange: updateAccountSchema.shape.username,
					}}
					children={(field) => (
						<FieldWrapper field={field} label="Username">
							<Input
								id={field.name}
								name={field.name}
								value={field.state.value}
								onBlur={field.handleBlur}
								onChange={(e) => field.handleChange(e.target.value)}
								placeholder="your_username"
								autoComplete="username"
							/>
						</FieldWrapper>
					)}
				/>

				<form.Field
					name="email"
					validators={{
						onChange: updateAccountSchema.shape.email,
					}}
					children={(field) => (
						<FieldWrapper field={field} label="Email">
							<Input
								id={field.name}
								name={field.name}
								type="email"
								value={field.state.value}
								onBlur={field.handleBlur}
								onChange={(e) => field.handleChange(e.target.value)}
								placeholder="you@example.com"
								autoComplete="email"
							/>
						</FieldWrapper>
					)}
				/>

				<form.Field
					name="password"
					validators={{
						onChange: updateAccountSchema.shape.password,
					}}
					children={(field) => (
						<FieldWrapper field={field} label="Password">
							<Input
								id={field.name}
								name={field.name}
								type="password"
								value={field.state.value}
								onBlur={field.handleBlur}
								onChange={(e) => field.handleChange(e.target.value)}
								placeholder="••••••"
								className={
									field.state.value ? 'tracking-[0.2rem]' : 'tracking-wide'
								}
								autoComplete="new-password"
							/>
						</FieldWrapper>
					)}
				/>

				<form.Field
					name="confirmPassword"
					children={(field) => (
						<FieldWrapper field={field} label="Confirm password">
							<Input
								id={field.name}
								name={field.name}
								type="password"
								value={field.state.value}
								onBlur={field.handleBlur}
								onChange={(e) => field.handleChange(e.target.value)}
								placeholder="••••••"
								className={
									field.state.value ? 'tracking-[0.2rem]' : 'tracking-wide'
								}
								autoComplete="new-password"
							/>
						</FieldWrapper>
					)}
				/>

				{serverError && <p className="text-destructive text-xs">{serverError}</p>}
			</CardContent>

			<CardFooter className="mt-4 flex justify-end">
				<form.Subscribe
					selector={(state) => [state.canSubmit, state.isSubmitting]}
					children={([canSubmit, isSubmitting]) => (
						<Button type="submit" disabled={!canSubmit || isSubmitting} size="lg">
							{isSubmitting ? 'Updating account...' : 'Update account'}
						</Button>
					)}
				/>
			</CardFooter>
		</form>
	);
}
