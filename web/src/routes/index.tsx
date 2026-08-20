import { createFileRoute } from '@tanstack/react-router';
import { AppShell } from '#/components/app-shell.tsx';
import { requireAuth } from '#/lib/route-guards.ts';

export const Route = createFileRoute('/')({
	beforeLoad: requireAuth,
	component: Home,
});

function Home() {
	return (
		<AppShell>
			<div className="p-8">
				<h1 className="text-4xl font-bold">Welcome to TanStack Start</h1>
				<p className="mt-4 text-lg">
					Edit <code>src/routes/index.tsx</code> to get started.
				</p>
			</div>
		</AppShell>
	);
}
