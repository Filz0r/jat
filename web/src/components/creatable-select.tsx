'use client';

import { useEffect, useMemo, useState } from 'react';

import { Button } from '#/components/ui/button.tsx';
import { Label } from '#/components/ui/label.tsx';
import {
	Combobox,
	ComboboxContent,
	ComboboxEmpty,
	ComboboxInput,
	ComboboxItem,
	ComboboxList,
} from '#/components/ui/combobox.tsx';
import type { ComboboxRootChangeEventDetails } from '#/components/ui/combobox.tsx';

export interface CreatableSelectItem {
	id?: number;
	label: string;
}

interface CreatableSelectProps {
	value: number;
	onChange: (value: number) => void;
	onBlur?: () => void;
	id?: string;
	name?: string;
	label?: string;
	placeholder?: string;
	searchPlaceholder?: string;
	createButtonLabel?: string;
	emptyMessage: string;
	items: CreatableSelectItem[];
	isLoading?: boolean;
	isError?: boolean;
	onCreateItem: (label: string) => Promise<CreatableSelectItem | number | undefined>;
	disabled?: boolean;
}

export function CreatableSelect({
	value,
	onChange,
	onBlur,
	id,
	name,
	label = 'Select an option',
	placeholder = 'Select an option',
	searchPlaceholder = 'Search...',
	createButtonLabel = 'Create new',
	emptyMessage,
	items,
	isLoading = false,
	isError = false,
	onCreateItem,
	disabled,
}: CreatableSelectProps) {
	const [search, setSearch] = useState('');
	const [isCreating, setIsCreating] = useState(false);
	const [createError, setCreateError] = useState<string | null>(null);
	const [createdItems, setCreatedItems] = useState<CreatableSelectItem[]>([]);

	// Drop locally cached created items once the real list contains them.
	useEffect(() => {
		setCreatedItems((prev) =>
			prev.filter((created) => !items.some((item) => item.id === created.id)),
		);
	}, [items]);

	const allItems = useMemo(() => {
		const byId = new Map<number, CreatableSelectItem>();

		for (const item of createdItems) {
			if (item.id) {
				byId.set(item.id, item);
			}
		}

		for (const item of items) {
			if (item.id) {
				byId.set(item.id, item);
			}
		}

		return Array.from(byId.values());
	}, [items, createdItems]);

	const filteredItems = useMemo(() => {
		const term = search.trim().toLowerCase();
		if (!term) return allItems;
		return allItems.filter((item) => item.label.toLowerCase().includes(term));
	}, [allItems, search]);

	const exactMatch = useMemo(() => {
		const term = search.trim().toLowerCase();
		if (!term) return false;
		return allItems.some((item) => item.label.toLowerCase() === term);
	}, [allItems, search]);

	const handleValueChange = (newValue: number | null) => {
		setCreateError(null);
		setSearch('');
		onChange(newValue ?? 0);
	};

	const handleInputValueChange = (
		newValue: string,
		eventDetails: ComboboxRootChangeEventDetails,
	) => {
		// Only treat input changes as search when the user is actually typing.
		// Programmatic changes (selection, open/close) must not replace the search
		// or Base UI can end up filtering the list to a single item and lose
		// the selected value.
		if (eventDetails.reason === 'input-change') {
			setSearch(newValue);
		}
		setCreateError(null);
	};

	const handleCreate = async () => {
		setCreateError(null);

		const trimmed = search.trim();
		if (!trimmed) {
			setCreateError('Name is required');
			return;
		}

		setIsCreating(true);
		try {
			const result = await onCreateItem(trimmed);
			const created = typeof result === 'number' ? { id: result, label: trimmed } : result;

			if (created?.id) {
				setCreatedItems((prev) => [...prev, created]);
				onChange(created.id);
				setSearch('');
			} else {
				setCreateError('Failed to create. Please try again.');
			}
		} catch {
			setCreateError('Failed to create. Please try again.');
		} finally {
			setIsCreating(false);
		}
	};

	const showCreateButton = search.trim() && !isLoading && !isError && !exactMatch && !isCreating;

	return (
		<div className="flex flex-col gap-1.5">
			<Label htmlFor={id ?? name}>{label}</Label>
			<Combobox
				value={value === 0 ? null : value}
				onValueChange={handleValueChange}
				onInputValueChange={handleInputValueChange}
				items={allItems}
				filteredItems={filteredItems}
				filter={null}
				itemToStringLabel={(itemValue) =>
					allItems.find((item) => item.id === itemValue)?.label ?? ''
				}
				disabled={disabled || isCreating}
			>
				<ComboboxInput
					id={id}
					name={name}
					placeholder={value === 0 ? placeholder : searchPlaceholder}
					onBlur={onBlur}
					disabled={disabled || isCreating}
					showTrigger
					showClear={false}
					className="w-full"
				/>
				<ComboboxContent>
					{createError && (
						<p className="text-destructive px-2 pt-2 text-xs">{createError}</p>
					)}
					<ComboboxList>
						{isLoading && (
							<div className="text-muted-foreground px-2 py-1.5 text-xs">
								Loading...
							</div>
						)}
						{!isLoading && isError && (
							<div className="text-destructive px-2 py-1.5 text-xs">
								Failed to load
							</div>
						)}
						{!isLoading &&
							!isError &&
							filteredItems.map((item) => (
								<ComboboxItem key={item.id ?? item.label} value={item.id ?? 0}>
									{item.label}
								</ComboboxItem>
							))}
						{!isLoading && !isError && filteredItems.length === 0 && (
							<ComboboxEmpty>{emptyMessage}</ComboboxEmpty>
						)}
					</ComboboxList>
					{showCreateButton && (
						<div className="border-t p-2">
							<Button
								type="button"
								variant="ghost"
								className="w-full justify-start"
								disabled={isCreating}
								onClick={handleCreate}
							>
								{`${createButtonLabel} "${search.trim()}"`}
							</Button>
						</div>
					)}
				</ComboboxContent>
			</Combobox>
		</div>
	);
}
