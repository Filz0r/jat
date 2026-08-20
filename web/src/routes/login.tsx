import { createFileRoute, useNavigate } from '@tanstack/react-router';
import * as React from 'react';
import { LoginForm } from '#/components/forms/login-form.tsx';
import { useUser } from '#/contexts/user-context.tsx';
import { api } from '#/api/client';

export const Route = createFileRoute('/login')({
	component: RouteComponent,
});

function RouteComponent() {
	const { initialized, user, setUser, isLoading } = useUser();
	const navigate = useNavigate();

	React.useEffect(() => {
		if (isLoading) return;

		if (!initialized) {
			navigate({ to: '/setup', replace: true });
			return;
		}

		if (user) {
			navigate({ to: '/', replace: true });
		}
	}, [isLoading, initialized, user, navigate]);

	const handleSuccess = async () => {
		const { data: userData, response: userResponse } = await api.GET('/users/me');

		if (userResponse.ok && userData.ok && userData.data) {
			setUser(userData.data);
		}

		await navigate({ to: '/' });
	};

	if (isLoading || !initialized || user) {
		return null;
	}

	return (
		<div className="flex min-h-screen items-center justify-center p-4">
			<LoginForm onSuccess={handleSuccess} />
		</div>
	);
}
