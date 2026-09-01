import type { FocusEventHandler } from 'react';
import { useState } from 'react';

import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuRadioGroup,
	DropdownMenuRadioItem,
	DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu.tsx';
import {
	InputGroup,
	InputGroupAddon,
	InputGroupButton,
	InputGroupInput,
} from '#/components/ui/input-group.tsx';
import { IconChevronDown } from '@tabler/icons-react';

const URL_SCHEMES = ['https', 'http'] as const;

interface HttpLinkInputProps {
	id?: string;
	name?: string;
	value?: string;
	onBlur?: FocusEventHandler<HTMLInputElement>;
	onChange?: (value: string) => void;
	placeholder?: string;
	disabled?: boolean;
}

export default function HttpLinkInput({
	id,
	name,
	value,
	onBlur,
	onChange,
	placeholder,
	disabled = false,
}: HttpLinkInputProps) {
	const fullURL = value ?? '';
	const rest = fullURL.replace(/^https?:\/\//, '');
	const [emptyScheme, setEmptyScheme] = useState<'http' | 'https'>('https');
	const scheme = URL_SCHEMES.find((s) => fullURL.startsWith(`${s}://`)) ?? emptyScheme;

	const compose = (nextScheme: string, raw: string) => {
		onChange?.(raw === '' ? '' : `${nextScheme}://${raw}`);
	};

	return (
		<InputGroup>
			<InputGroupAddon>
				<DropdownMenu>
					<DropdownMenuTrigger
						render={
							<InputGroupButton
								variant="ghost"
								aria-label="URL scheme"
								disabled={disabled}
							>
								{scheme}://
								<IconChevronDown className="size-3" />
							</InputGroupButton>
						}
					/>
					<DropdownMenuContent align="start" sideOffset={8} alignOffset={-4}>
						<DropdownMenuGroup>
							<DropdownMenuRadioGroup
								value={scheme}
								onValueChange={(nextScheme) => {
									const next = String(nextScheme) as 'http' | 'https';
									setEmptyScheme(next);
									compose(next, rest);
								}}
							>
								{URL_SCHEMES.map((s) => (
									<DropdownMenuRadioItem key={s} value={s}>
										{s}://
									</DropdownMenuRadioItem>
								))}
							</DropdownMenuRadioGroup>
						</DropdownMenuGroup>
					</DropdownMenuContent>
				</DropdownMenu>
			</InputGroupAddon>
			<InputGroupInput
				id={id}
				name={name}
				placeholder={placeholder}
				value={rest}
				disabled={disabled}
				onBlur={onBlur}
				onChange={(e) => {
					const raw = e.target.value.replace(/^https?:\/\//, '');
					compose(scheme, raw);
				}}
			/>
		</InputGroup>
	);
}
