import { clsx } from 'clsx';
import type { ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	// eslint-disable-next-line tailwindcss/no-custom-classname
	return twMerge(clsx(inputs));
}

export function getColorFromKind(kind: string): string {
	switch (kind) {
		case 'applied':
			return 'bg-orange-500 text-black';
		case 'rejected':
			return 'bg-red-500 text-black';
		case 'ghosted':
			return 'bg-red-600 text-black';
		case 'interviewed':
			return 'bg-primary text-white';
		case 'irrelevant':
			return 'bg-muted text-ghost';
		case 'accepted':
			return 'bg-green-500 text-green-100';
		default:
			return 'bg-background text-white';
	}
}
