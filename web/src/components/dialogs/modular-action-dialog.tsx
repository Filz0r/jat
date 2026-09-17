import type { ReactElement } from 'react';

import {
	AlertDialog,
	AlertDialogContent,
	AlertDialogTitle,
	AlertDialogHeader,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogCancel,
	AlertDialogAction,
	AlertDialogTrigger,
} from '#/components/ui/alert-dialog.tsx';

interface ModularActionDialogProps {
	title: string;
	message: string;
	open: boolean;
	onClose: () => void;
	onConfirm: () => void;
	trigger: ReactElement;
	destructive?: boolean;
	actionLabel?: string;
}

export default function ModularActionDialog({
	title,
	message,
	open,
	onClose,
	onConfirm,
	trigger,
	destructive = false,
	actionLabel = 'Confirm',
}: ModularActionDialogProps) {
	let actionStyles = 'cursor-pointer flex-1';
	if (destructive && actionLabel === 'Confirm') {
		actionLabel = 'Delete';
	}
	if (!destructive) {
		actionStyles += ' bg-green-500 hover:bg-green-600';
	}
	return (
		<AlertDialog open={open} onOpenChange={onClose}>
			<AlertDialogTrigger render={trigger} />
			<AlertDialogContent className="border px-0 sm:max-w-sm">
				<AlertDialogHeader>
					<AlertDialogTitle className="text-primary w-full border-b px-2 pb-4 text-center text-lg">
						{title}
					</AlertDialogTitle>
					<AlertDialogDescription className="mx-2.5 py-1 text-center">
						<p>{message}</p>
					</AlertDialogDescription>
				</AlertDialogHeader>
				<AlertDialogFooter className="justify-center border-t px-2 pt-3">
					<AlertDialogCancel
						variant={!destructive ? 'destructive' : 'secondary'}
						size="lg"
						className="flex-1 cursor-pointer"
					>
						Cancel
					</AlertDialogCancel>
					<AlertDialogAction
						className={actionStyles}
						size="lg"
						onClick={() => onConfirm()}
						variant={!destructive ? 'default' : 'destructive'}
					>
						{actionLabel}
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	);
}
