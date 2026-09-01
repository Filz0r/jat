import { createFileRoute } from '@tanstack/react-router';
import { requireAdmin } from '#/lib/route-guards.ts';

export const Route = createFileRoute('/_app/admin/users/')({
	beforeLoad: requireAdmin,
	component: RouteComponent,
});

function RouteComponent() {
	return <div>Hello "/_app/admin/users/"!</div>;
}
