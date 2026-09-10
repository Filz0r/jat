import { useForm } from '@tanstack/react-form';
import { Button } from '#/components/ui/button.tsx';
import { Input } from '#/components/ui/input.tsx';
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from '#/components/ui/card.tsx';
import { api } from '#/api/client.ts';
import { createAccountSchema } from '#/schemas/users.ts';
import type { CreateAccountInput } from '#/schemas/users.ts';
import { useState } from 'react';
import FieldWrapper from '#/components/field-wrapper.tsx';

interface CreateAccountProps {
	onSuccess?: () => void;
}

export function CreateAccount({ onSuccess }: CreateAccountProps) {
	const [serverError, setServerError] = useState<string | null>(null);

	const form = useForm({
		defaultValues: {
			username: '',
			email: '',
			password: '',
			confirmPassword: '',
		} satisfies CreateAccountInput,
		validators: {
			onChange: createAccountSchema,
		},
		onSubmit: async ({ value }) => {
			setServerError(null);

			const { data: registerData, error: registerError } = await api.POST('/users', {
				body: {
					username: value.username,
					email: value.email,
					password: value.password,
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

			const { error: loginError } = await api.POST('/auth/login', {
				body: {
					email: value.email,
					password: value.password,
				},
			});

			if (loginError) {
				setServerError('Account created, but login failed. Please try again.');
				return;
			}

			onSuccess?.();
		},
	});

	return (
		<Card className="w-full max-w-md space-x-2">
			<form
				onSubmit={async (e) => {
					e.preventDefault();
					e.stopPropagation();
					await form.handleSubmit();
				}}
			>
				<CardHeader className="pb-2.5">
					<CardTitle className="text-center text-lg">Create your account</CardTitle>
					<CardDescription className="text-center">
						Set up the first admin account to start using JAT.
					</CardDescription>
				</CardHeader>

				<CardContent className="mx-2 mt-2 mb-4 flex flex-col gap-4">
					<form.Field
						name="username"
						validators={{
							onChange: createAccountSchema.shape.username,
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
							onChange: createAccountSchema.shape.email,
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
							onChange: createAccountSchema.shape.password,
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
								{isSubmitting ? 'Creating account...' : 'Create account'}
							</Button>
						)}
					/>
				</CardFooter>
			</form>
		</Card>
	);
}
