import z from 'zod';

export const loginSearchSchema = z.object({
	redirect: z.string().optional().catch(undefined),
});

export type LoginSearch = z.infer<typeof loginSearchSchema>;
