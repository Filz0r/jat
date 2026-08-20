import z from 'zod';

export const createJobApplicationSchema = z.object({
	title: z.string('You need to provide a title for this job').min(3).max(100),
	url: z.url('You need to provide an URL').min(3).max(100),
	createdAt: z.date(),
	company: z.uint32('You need to select a company, you can also create a new one').min(1),
	status: z.uint32('You need to select a status').min(1),
});

export type CreateJobApplication = z.infer<typeof createJobApplicationSchema>;
