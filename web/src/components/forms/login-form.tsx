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
import { loginSchema } from '#/schemas/users.ts';
import type { LoginInput } from '#/schemas/users.ts';
import FieldWrapper from '#/components/field-wrapper.tsx';
import { useState } from 'react';

interface LoginFormProps {
	onSuccess?: () => void;
}

export function LoginForm({ onSuccess }: LoginFormProps) {
	const [serverError, setServerError] = useState<string | null>(null);

	const form = useForm({
		defaultValues: {
			email: '',
			password: '',
		} satisfies LoginInput,
		validators: {
			onChange: loginSchema,
		},
		onSubmit: async ({ value }) => {
			setServerError(null);

			const { data: loginData, error: loginError } = await api.POST('/auth/login', {
				body: value,
			});

			if (loginError || !loginData.ok) {
				setServerError(
					loginData?.message ?? 'Invalid email or password. Please try again.',
				);
				return;
			}

			onSuccess?.();
		},
	});

	return (
		<Card className="w-full max-w-md">
			<form
				onSubmit={(e) => {
					e.preventDefault();
					e.stopPropagation();
					form.handleSubmit();
				}}
			>
				<CardHeader>
					<CardTitle>Log in</CardTitle>
					<CardDescription>
						Enter your credentials to access your account.
					</CardDescription>
				</CardHeader>

				<CardContent className="flex flex-col gap-4">
					<form.Field
						name="email"
						validators={{
							onChange: loginSchema.shape.email,
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
							onChange: loginSchema.shape.password,
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
									autoComplete="current-password"
								/>
							</FieldWrapper>
						)}
					/>

					{serverError && <p className="text-destructive text-xs">{serverError}</p>}
				</CardContent>

				<CardFooter>
					<form.Subscribe
						selector={(state) => [state.canSubmit, state.isSubmitting]}
						children={([canSubmit, isSubmitting]) => (
							<Button type="submit" disabled={!canSubmit || isSubmitting}>
								{isSubmitting ? 'Logging in...' : 'Log in'}
							</Button>
						)}
					/>
				</CardFooter>
			</form>
		</Card>
	);
}
