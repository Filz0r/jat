import z from 'zod';

export const setupStepSchema = z.object({
	step: z.enum(['create-account', 'finish']).default('create-account').catch('create-account'),
});

export type SetupStep = z.infer<typeof setupStepSchema>['step'];
