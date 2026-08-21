import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '#/components/ui/dialog.tsx';
import type { ReactElement } from 'react';
import { Button } from '#/components/ui/button.tsx';

interface ModularDeleteEntityDialogProps {
	message: string;
	description?: string;
	title: string;
	onClose: () => void | Promise<void>;
	triggerButton: ReactElement;
	deleteButtonMessage: string;
}

export default function ModularDeleteEntityDialog({
	message,
	description,
	title,
	onClose,
	triggerButton,
	deleteButtonMessage,
}: ModularDeleteEntityDialogProps) {
	return (
		<Dialog>
			<DialogTrigger render={triggerButton} />

			<DialogContent className="sm:max-w-sm" showCloseButton={false}>
				<DialogHeader>
					<DialogTitle>{title}</DialogTitle>
					{description && <DialogDescription>{description}</DialogDescription>}
				</DialogHeader>
				{message}
				<DialogFooter>
					<DialogClose render={<Button variant="secondary">Cancel</Button>} />
					<DialogClose
						render={
							<Button variant="destructive" onClick={onClose}>
								{deleteButtonMessage}
							</Button>
						}
					/>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
