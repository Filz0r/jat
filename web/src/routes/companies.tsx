import { createFileRoute } from '@tanstack/react-router';
import { AppShell } from '#/components/app-shell.tsx';
import { requireAuth } from '#/lib/route-guards.ts';

export const Route = createFileRoute('/companies')({
	beforeLoad: requireAuth,
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<AppShell>
			<div>Hello "/companies"!</div>
		</AppShell>
	);
}
