import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router';
import { LoginForm } from '#/components/forms/login-form.tsx';
import { useAuth } from '#/hooks/use-auth.ts';
import { redirectIfAuthenticated } from '#/lib/route-guards.ts';
import { loginSearchSchema } from '#/schemas/login.ts';

export const Route = createFileRoute('/login')({
	validateSearch: loginSearchSchema,
	beforeLoad: redirectIfAuthenticated,
	component: RouteComponent,
});

function RouteComponent() {
	const { refreshUser } = useAuth();
	const navigate = useNavigate();
	const { redirect } = useSearch({ from: '/login' });

	const handleSuccess = async () => {
		await refreshUser();

		const target = isValidRedirect(redirect) ? redirect : '/';
		await navigate({ to: target, replace: true });
	};

	return (
		<div className="flex min-h-screen items-center justify-center p-4">
			<LoginForm onSuccess={handleSuccess} />
		</div>
	);
}

function isValidRedirect(value: string | undefined): value is string {
	if (!value) {
		return false;
	}

	if (!value.startsWith('/')) {
		return false;
	}

	if (value.startsWith('//')) {
		return false;
	}

	try {
		const url = new URL(value, window.location.origin);
		return url.origin === window.location.origin && url.pathname !== '/login';
	} catch {
		return false;
	}
}
