import type { ReactElement } from 'react';
import type { ReactFormExtendedApi } from '@tanstack/react-form';

import {
	Dialog,
	DialogContent,
	DialogTrigger,
	DialogTitle,
	DialogHeader,
	DialogDescription,
	DialogFooter,
	DialogClose,
} from '#/components/ui/dialog.tsx';
import { Button } from '#/components/ui/button.tsx';
import LogForm from '#/components/forms/log-form.tsx';

interface ModularFormDialogProps {
	title: string;
	form: ReactFormExtendedApi<any, any, any, any, any, any, any, any, any, any, any, any>;
	Content: ReactElement;
	TriggerButton?: ReactElement;
	CloseButton?: ReactElement;
	description?: string;
	submitMessage?: string;
	open?: boolean;
	onClose?: () => void;
}

export default function ModularFormDialog({
	TriggerButton = <Button variant="outline">Open Dialog</Button>,
	CloseButton = <Button variant="outline">Cancel</Button>,
	title,
	description,
	form,
	Content,
	submitMessage = 'Submit',
	open = false,
	onClose,
}: ModularFormDialogProps) {
	return (
		<Dialog open={open}>
			<DialogTrigger render={TriggerButton} />
			<DialogContent className="border p-0 sm:max-w-sm" showCloseButton={false}>
				<form
					noValidate
					onSubmit={(e) => {
						e.preventDefault();
						e.stopPropagation();
						void form.handleSubmit();
					}}
					className="space-y-4"
				>
					<DialogHeader className="border-b px-2 pt-3">
						<DialogTitle className="text-primary pb-1 text-center text-lg">
							{title}
						</DialogTitle>
						{description && (
							<DialogDescription className="bg-muted mb-3 rounded-lg px-1 py-2 text-center text-[11px] font-extralight italic">
								{description}
							</DialogDescription>
						)}
					</DialogHeader>
					<div className="px-3">{Content}</div>
					<DialogFooter className="border-t px-2 pt-3 pb-3.5">
						<DialogClose render={CloseButton} onClick={onClose} />
						<form.Subscribe
							selector={(state) => [state.canSubmit, state.isSubmitting]}
							children={([canSubmit, isSubmitting]) => (
								<Button
									type="submit"
									className="w-full sm:min-w-32 sm:flex-1"
									disabled={!canSubmit || isSubmitting}
								>
									{isSubmitting ? 'Saving...' : submitMessage}
								</Button>
							)}
						/>
						<LogForm form={form} />
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	);
}
