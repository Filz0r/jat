import z from 'zod';

export const createAccountSchema = z
	.object({
		username: z.string().min(3, 'Username must be at least 3 characters'),
		email: z.email('Please enter a valid email address'),
		password: z
			.string()
			.min(6, 'Password must be at least 6 characters')
			.max(64, 'Password must be at most 64 characters'),
		confirmPassword: z.string().min(1, 'Please confirm your password'),
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: 'Passwords do not match',
		path: ['confirmPassword'],
	});

export const loginSchema = z.object({
	email: z.email('Please enter a valid email address'),
	password: z.string().min(1, 'Password is required'),
});

export const updateAccountSchema = z
	.object({
		username: z.string().min(3, 'Username must be at least 3 characters'),
		email: z.email('Please enter a valid email address'),
		password: z.optional(
			z
				.string()
				.transform((v) => (v === '' ? undefined : v))
				.pipe(
					z
						.string()
						.min(6, 'Password must be at least 6 characters')
						.max(64, 'Password must be at most 64 characters')
						.optional(),
				),
		),
		confirmPassword: z.optional(
			z
				.string()
				.transform((v) => (v === '' ? undefined : v))
				.pipe(z.string().min(1, 'Please confirm your password').optional()),
		),
	})
	.refine((data) => data.password === undefined || data.password === data.confirmPassword, {
		message: 'Passwords do not match',
		path: ['confirmPassword'],
	});

export const defaultUserTabSchema = z.object({
	tab: z.enum(['settings', 'account']).default('settings').catch('settings').optional(),
});

export const defaultApplicationStatusSelectSchema = z.object({
	status_id: z.number().positive('The server only accepts unsigned integers'),
});

export type CreateAccountInput = z.infer<typeof createAccountSchema>;
export type LoginInput = z.infer<typeof loginSchema>;
export type UpdateAccountSchema = z.infer<typeof updateAccountSchema>;
export type DefaultUserTabSchema = z.infer<typeof defaultUserTabSchema>;
export type DefaultApplicationStatusSelectSchema = z.infer<
	typeof defaultApplicationStatusSelectSchema
>;
