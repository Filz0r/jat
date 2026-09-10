import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { redirectIfAuthenticated } from '#/lib/route-guards.ts';
import { CreateAccount } from '#/components/forms/create-account-form.tsx';
import { useAuth } from '#/hooks/use-auth.ts';

export const Route = createFileRoute('/register')({
	component: RouteComponent,
	beforeLoad: redirectIfAuthenticated,
});

function RouteComponent() {
	const { refreshUser } = useAuth();
	const navigate = useNavigate();

	const afterCreated = async () => {
		await refreshUser();
		await navigate({ to: '/' });
	};

	return (
		<main className="flex min-h-screen items-center justify-center">
			<CreateAccount onSuccess={afterCreated} />
		</main>
	);
}
