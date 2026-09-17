'use client';

import { format } from 'date-fns';
import * as React from 'react';

import { Button } from '#/components/ui/button.tsx';
import { Calendar } from '#/components/ui/calendar.tsx';
import { Popover, PopoverContent, PopoverTrigger } from '#/components/ui/popover.tsx';
import { cn } from '#lib/utils';

interface DatePickerProps {
	value?: Date;
	onChange?: (date?: Date) => void;
	onBlur?: () => void;
	id?: string;
	name?: string;
	placeholder?: string;
	className?: string;
	disabled?: boolean;
}

export function DatePicker({
	value,
	onChange,
	onBlur,
	id,
	name,
	placeholder = 'Pick a date',
	className,
	disabled,
}: DatePickerProps) {
	const [open, setOpen] = React.useState(false);

	return (
		<Popover open={open} onOpenChange={setOpen}>
			<PopoverTrigger
				render={
					<Button
						variant="outline"
						id={id}
						name={name}
						disabled={disabled}
						onBlur={onBlur}
						className={cn(
							'h-[max(1.75rem,calc(100vw/1920*28))] w-full justify-start font-normal',
							className,
						)}
					>
						{value ? (
							format(value, 'PPP')
						) : (
							<span className="text-muted-foreground">{placeholder}</span>
						)}
					</Button>
				}
			/>
			<PopoverContent className="w-auto p-0" align="start">
				<Calendar
					mode="single"
					selected={value}
					defaultMonth={value}
					onSelect={(date) => {
						onChange?.(date);
						setOpen(false);
					}}
				/>
			</PopoverContent>
		</Popover>
	);
}
