import { Label } from '#/components/ui/label.tsx';
import type { $ZodIssueBase } from 'zod/v4/core';
import type { ReactNode } from 'react';
import { cn } from '#/lib/utils.ts';

interface FieldWrapperProps {
	field: {
		name: string;
		state: {
			meta: {
				isDirty: boolean;
				isTouched: boolean;
				errors: unknown[];
				errorMap: Record<string, unknown>;
			};
		};
	};
	label: string;
	children: ReactNode;
	className?: string;
}

export default function FieldWrapper({ field, label, children, className }: FieldWrapperProps) {
	const showErrors = field.state.meta.isDirty || field.state.meta.isTouched;
	const errors = field.state.meta.errors;

	const messages = errors
		.filter((err): err is $ZodIssueBase => typeof err === 'object' && err !== null)
		.map((err) => {
			if ('message' in err && typeof err.message === 'string') {
				return err.message;
			}
			return 'Invalid value';
		});

	return (
		<div className={cn('flex flex-col gap-1.5', className)}>
			<Label htmlFor={field.name}>{label}</Label>
			{children}
			{showErrors && messages.length > 0 && (
				<p className="text-xs text-destructive">{messages.join(', ')}</p>
			)}
		</div>
	);
}
