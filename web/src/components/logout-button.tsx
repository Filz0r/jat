import { useNavigate } from '@tanstack/react-router';
import { IconLogout } from '@tabler/icons-react';
import { SidebarMenuButton } from '#components/ui/sidebar';
import { api } from '#/api/client';
import { useUser } from '#/contexts/user-context.tsx';
import { useState } from 'react';

export function LogoutButton() {
	const navigate = useNavigate();
	const { setUser } = useUser();
	const [isLoading, setIsLoading] = useState(false);

	const handleLogout = async () => {
		setIsLoading(true);
		try {
			await api.POST('/auth/logout');
		} finally {
			setIsLoading(false);
			setUser(null);
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
