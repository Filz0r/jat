import {
	IconBriefcase,
	IconBuilding,
	IconHome,
	IconSettings,
	IconLayoutSidebar,
} from '@tabler/icons-react';

import {
	Sidebar,
	SidebarContent,
	SidebarFooter,
	SidebarHeader,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarRail,
	useSidebar,
} from '#components/ui/sidebar';
import { NavGroup } from '#/components/nav-group.tsx';
import type { NavItem } from '#/components/nav-group.tsx';
import { ThemeToggle } from '#/components/theme-toggle.tsx';
import { LogoutButton } from '#/components/logout-button.tsx';
import { Link } from '@tanstack/react-router';

const mainNavItems: NavItem[] = [
	{ title: 'Home', to: '/', icon: IconHome },
	{ title: 'Jobs', to: '/jobs', icon: IconBriefcase },
	{ title: 'Companies', to: '/companies', icon: IconBuilding },
	{ title: 'Account Settings', to: '/settings', icon: IconSettings },
];

function SidebarToggle() {
	const { toggleSidebar } = useSidebar();
	return (
		<SidebarMenuButton onClick={toggleSidebar}>
			<IconLayoutSidebar />
			<span>Toggle Sidebar</span>
		</SidebarMenuButton>
	);
}

export function AppSidebar() {
	const { open } = useSidebar();
	return (
		<Sidebar collapsible="icon">
			<SidebarHeader>
				<SidebarMenu>
					<SidebarMenuItem>
						<SidebarMenuButton size="lg" render={<Link to="/" />} tooltip="JAT">
							{!open && (
								<span className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
									J
								</span>
							)}
							{open && (
								<div className="grid flex-1 text-left text-sm leading-tight">
									<span className="text-center text-lg">
										Job Application Tracker
									</span>
								</div>
							)}
						</SidebarMenuButton>
					</SidebarMenuItem>
				</SidebarMenu>
			</SidebarHeader>

			<SidebarContent>
				<NavGroup items={mainNavItems} />
			</SidebarContent>

			<SidebarFooter>
				<NavGroup items={footerNavItems} />
				<SidebarMenu>
					<SidebarMenuItem>
						<ThemeToggle />
					</SidebarMenuItem>
					<SidebarMenuItem>
						<SidebarToggle />
					</SidebarMenuItem>
					<SidebarMenuItem>
						<LogoutButton />
					</SidebarMenuItem>
				</SidebarMenu>
			</SidebarFooter>

			<SidebarRail />
		</Sidebar>
	);
}
