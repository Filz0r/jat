import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router';

import { Button } from '#components/ui/button';
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from '#components/ui/card';

import { api } from '#/api/client';
import { CreateAccount } from '#/components/forms/create-account-form.tsx';
import { useAuth } from '#/hooks/use-auth.ts';
import { requireSetup } from '#/lib/route-guards.ts';
import { useCallback } from 'react';
import { setupStepSchema } from '#/schemas/setup.ts';
import type { SetupStep } from '#/schemas/setup.ts';

export const Route = createFileRoute('/setup')({
	validateSearch: setupStepSchema,
	beforeLoad: requireSetup,
	component: RouteComponent,
});

function RouteComponent() {
	const { step } = useSearch({ from: '/setup' });
	const navigate = useNavigate({ from: '/setup' });
	const { refreshUser, loadInit } = useAuth();

	const goToStep = useCallback(
		async (next: SetupStep) => {
			await navigate({
				to: '/setup',
				search: { step: next },
				replace: true,
			});
		},
		[navigate],
	);

	const handleAccountCreated = async () => {
		await refreshUser();
		await loadInit();
		await goToStep('finish');
	};

	const handleFinishSetup = async () => {
		const { data, error } = await api.GET('/initialized/set');

		if (error || !data.ok) {
			console.error('Failed to finish setup', error);
			return;
		}

		await refreshUser();
		await loadInit();
		await navigate({ to: '/' });
	};

	return (
		<div className="flex min-h-screen flex-col items-center justify-center gap-6 p-4">
			<StepIndicator currentStep={step} />

			{step === 'create-account' && <CreateAccount onSuccess={handleAccountCreated} />}

			{step === 'finish' && (
				<Card className="w-full max-w-md">
					<CardHeader>
						<CardTitle>Finalize setup</CardTitle>
						<CardDescription>
							Your admin account is ready. Click below to mark the server as
							initialized and start using JAT.
						</CardDescription>
					</CardHeader>
					<CardContent />
					<CardFooter className="flex justify-end">
						<Button onClick={handleFinishSetup}>Complete setup</Button>
					</CardFooter>
				</Card>
			)}
		</div>
	);
}

function StepIndicator({ currentStep }: { currentStep: SetupStep }) {
	const steps: { key: SetupStep; label: string }[] = [
		{ key: 'create-account', label: 'Create account' },
		{ key: 'finish', label: 'Finish' },
	];

	return (
		<ol className="flex items-center gap-4">
			{steps.map((step, index) => {
				const isActive = step.key === currentStep;
				const isPast = steps.findIndex((s) => s.key === currentStep) > index;

				return (
					<li key={step.key} className="flex items-center gap-2">
						<span
							className={`flex size-6 items-center justify-center rounded-full text-xs font-medium ${
								isActive || isPast
									? 'bg-primary text-primary-foreground'
									: 'bg-muted text-muted-foreground'
							}`}
						>
							{index + 1}
						</span>
						<span
							className={`text-xs ${
								isActive ? 'text-foreground font-medium' : 'text-muted-foreground'
							}`}
						>
							{step.label}
						</span>
					</li>
				);
			})}
		</ol>
	);
}
