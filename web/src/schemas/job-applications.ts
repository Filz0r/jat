import z from 'zod';

export const preSelectCompanyQuerySchema = z.object({
	company_id: z.number().positive().optional(),
});

export const createJobApplicationSchema = z.object({
	title: z.string('You need to provide a title for this job').min(3).max(100),
	url: z.httpUrl('You need to provide an URL'),
	createdAt: z.date(),
	company: z.uint32('You need to select a company, you can also create a new one').min(1),
	status: z.uint32('You need to select a status').min(1),
});

export const updateJobApplicationStatusSchema = z.object({
	status: z.uint32('You need to select a status').min(1),
});

export const createJobApplicationNoteSchema = z.object({
	body: z.string('You need to provide content').min(3).max(2048),
});

const applicationStatusKinds = [
	'applied',
	'rejected',
	'ghosted',
	'interviewed',
	'irrelevant',
	'accepted',
];
export const applicationStatusKindSchema = z.enum(applicationStatusKinds);

export const jobApplicationStatusSchema = z.object({
	status: z.string().min(3).max(64),
	kind: applicationStatusKindSchema,
});

export type CreateJobApplication = z.infer<typeof createJobApplicationSchema>;
export type UpdateJobApplicationStatusSchema = z.infer<typeof updateJobApplicationStatusSchema>;
export type CreateJobApplicationNote = z.infer<typeof createJobApplicationNoteSchema>;
export type ApplicationStatusSchema = z.infer<typeof jobApplicationStatusSchema>;
