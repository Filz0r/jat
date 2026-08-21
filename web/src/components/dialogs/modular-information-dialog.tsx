import type { ReactElement } from 'react';

import {
	Dialog,
	DialogContent,
	DialogTitle,
	DialogHeader,
	DialogDescription,
	DialogFooter,
	DialogClose,
} from '#/components/ui/dialog.tsx';
import { Button } from '#/components/ui/button.tsx';

interface ModularInformationDialogProps {
	title: string;
	message: string;
	CloseButton?: ReactElement;
	description?: string;
	open?: boolean;
	onClose?: () => void;
}

export default function ModularInformationDialog({
	CloseButton = <Button variant="outline">Close</Button>,
	title,
	description,
	message,
	open = false,
	onClose,
}: ModularInformationDialogProps) {
	return (
		<Dialog open={open}>
			<DialogContent className="sm:max-w-sm">
				<DialogHeader>
					<DialogTitle>{title}</DialogTitle>
					{description && <DialogDescription>{description}</DialogDescription>}
				</DialogHeader>
				{message}
				<DialogFooter>
					<DialogClose render={CloseButton} onClick={() => onClose?.()} />
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
