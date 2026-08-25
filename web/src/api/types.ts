import type { components } from './gen-spec.ts';

export type JobApplication = components['schemas']['api.applicationResponse'];
export type JobApplicationStatus = components['schemas']['api.applicationStatusResponse'];
export type JobApplicationNote = components['schemas']['api.noteData'];
export type JobApplicationHistory = components['schemas']['api.applicationStatusHistoryResponse'];

export type APIResponse = components['schemas']['api.apiResponse'];
