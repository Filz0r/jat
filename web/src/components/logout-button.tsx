import { useNavigate } from '@tanstack/react-router';
import { IconLogout } from '@tabler/icons-react';
import { SidebarMenuButton } from '#components/ui/sidebar';
import { api } from '#/api/client';
import { useAuth } from '#/contexts/auth-context.tsx';
import { useState } from 'react';

export function LogoutButton() {
	const navigate = useNavigate();
	const { clearSession } = useAuth();
	const [isLoading, setIsLoading] = useState(false);

	const handleLogout = async () => {
		setIsLoading(true);
		try {
			await api.POST('/auth/logout');
		} finally {
			setIsLoading(false);
			clearSession();
			await navigate({ to: '/login' });
		}
	};

	return (
		<SidebarMenuButton onClick={handleLogout} disabled={isLoading} tooltip="Log out">
			<IconLogout />
			<span>{isLoading ? 'Logging out...' : 'Log out'}</span>
		</SidebarMenuButton>
	);
}
