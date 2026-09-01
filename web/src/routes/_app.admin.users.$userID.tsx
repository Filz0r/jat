import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { requireAdmin } from '#/lib/route-guards.ts';
import { validate } from 'uuid';

export const Route = createFileRoute('/_app/admin/users/$userID')({
	beforeLoad: requireAdmin,
	component: RouteComponent,
});

function RouteComponent() {
	const navigate = useNavigate();
	const { userID } = Route.useParams();
	const validID = validate(userID);
	if (!validID) {
		void navigate({ to: '/admin/users' });
	}
	return <div>Hello "/_app/admin/users/$userID"!</div>;
}
