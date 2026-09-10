import {
	IconSettings,
	IconLayoutSidebar,
	IconLogout,
	IconUser,
	IconSun,
	IconMoon,
} from '@tabler/icons-react';

import { useSidebar } from '#components/ui/sidebar';
import { useNavigate } from '@tanstack/react-router';
import { Avatar, AvatarFallback } from '#/components/ui/avatar.tsx';
import { useAuth } from '#/hooks/use-auth.ts';
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu.tsx';
import { Button } from '#/components/ui/button.tsx';
import { api } from '#/api/client.ts';
import { useTheme } from '#/contexts/theme-context.tsx';

export default function UserMenu() {
	const { open, toggleSidebar } = useSidebar();
	const navigate = useNavigate();
	const { user, clearSession } = useAuth();
	const { theme, toggleTheme } = useTheme();

	const handleLogout = async () => {
		try {
			await api.POST('/auth/logout');
		} finally {
			clearSession();
			await navigate({ to: '/login' });
		}
	};

	const username = (user?.username || 'JAT').slice(0, 3).toUpperCase();
	const isDark = theme === 'dark';
	return (
		<DropdownMenu>
			<DropdownMenuTrigger
				openOnHover
				closeDelay={500}
				render={
					!open ? (
						<Button variant="ghost" size="icon" className="mb-4 rounded-full pl-1">
							<Avatar>
								<AvatarFallback>{username}</AvatarFallback>
							</Avatar>
						</Button>
					) : (
						<Button className="mb-4 text-center" variant="outline">
							{username}
						</Button>
					)
				}
			/>
			<DropdownMenuContent side="left" align="end" className="ml-3.5 w-full text-center">
				<DropdownMenuItem
					onClick={async () => {
						await navigate({
							to: '/account',
							search: {
								tab: 'settings',
							},
						});
					}}
				>
					<IconSettings />
					Settings
				</DropdownMenuItem>
				<DropdownMenuItem
					onClick={async () => {
						await navigate({
							to: '/account',
							search: {
								tab: 'account',
							},
						});
					}}
				>
					<IconUser /> Profile
				</DropdownMenuItem>
				<DropdownMenuSeparator />
				<DropdownMenuItem onClick={toggleTheme}>
					{isDark ? <IconSun /> : <IconMoon />}
					{isDark ? 'Light mode' : 'Dark mode'}
				</DropdownMenuItem>
				<DropdownMenuItem onClick={toggleSidebar}>
					<IconLayoutSidebar />
					{open ? 'Close' : 'Open'} Sidebar
				</DropdownMenuItem>
				<DropdownMenuSeparator />
				<DropdownMenuItem variant="destructive" onClick={handleLogout}>
					<IconLogout />
					Logout
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	);
}
