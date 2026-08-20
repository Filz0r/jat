import { IconMoon, IconSun } from '@tabler/icons-react';

import { SidebarMenuButton } from '#components/ui/sidebar';
import { useTheme } from '#/contexts/theme-context.tsx';

export function ThemeToggle() {
	const { theme, toggleTheme } = useTheme();

	return (
		<SidebarMenuButton
			onClick={toggleTheme}
			tooltip={theme === 'light' ? 'Switch to dark mode' : 'Switch to light mode'}
		>
			{theme === 'light' ? <IconMoon /> : <IconSun />}
			<span>{theme === 'light' ? 'Dark mode' : 'Light mode'}</span>
		</SidebarMenuButton>
	);
}
