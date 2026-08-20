import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { AppShell } from '#/components/app-shell.tsx';
import { useUser } from '#/contexts/user-context.tsx';
import { useEffect } from 'react';

export const Route = createFileRoute('/settings')({
	component: RouteComponent,
});

function RouteComponent() {
	const { user, initialized, isLoading } = useUser();
	const navigate = useNavigate();

	useEffect(() => {
		if (isLoading) return;

		if (!initialized) {
			navigate({ to: '/setup', replace: true });
			return;
		}

		if (!user) {
			navigate({ to: '/login', replace: true });
			return;
		}
	}, [isLoading, initialized, user, navigate]);

	if (isLoading || !initialized || !user) {
		return null;
	}

	return (
		<AppShell>
			<div>Hello "/settings"!</div>
		</AppShell>
	);
}
