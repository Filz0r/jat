import { createFileRoute, Outlet } from '@tanstack/react-router';

import { AppShell } from '#/components/app-shell.tsx';
import { requireAuth } from '#/lib/route-guards.ts';

export const Route = createFileRoute('/_app')({
	beforeLoad: requireAuth,
	component: AppLayout,
});

function AppLayout() {
	return (
		<AppShell>
			<Outlet />
		</AppShell>
	);
}