import type { NavItem } from '#/components/sidebar/nav-group.tsx';

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
import { Link } from '@tanstack/react-router';
import { NavGroup } from '#/components/sidebar/nav-group.tsx';
import { IconBriefcase, IconBuilding, IconChartColumn, IconHome } from '@tabler/icons-react';
import UserMenu from '#/components/sidebar/user-menu.tsx';

const mainNavItems: NavItem[] = [
	{ title: 'Home', to: '/', icon: IconHome },
	{ title: 'Jobs', to: '/jobs', icon: IconBriefcase },
	{ title: 'Application Status', to: '/application_status', icon: IconChartColumn },
	{ title: 'Companies', to: '/companies', icon: IconBuilding },
];

export default function AppSidebar() {
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
				<UserMenu />
			</SidebarFooter>

			<SidebarRail />
		</Sidebar>
	);
}
