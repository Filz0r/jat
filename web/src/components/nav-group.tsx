import * as React from 'react';
import { useRouterState, Link } from '@tanstack/react-router';

import {
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	useSidebar,
} from '#components/ui/sidebar';

export interface NavItem {
	title: string;
	to: string;
	icon: React.ComponentType<{ className?: string }>;
}

interface NavGroupProps {
	items: NavItem[];
}

function isActiveRoute(currentPath: string, targetPath: string): boolean {
	return currentPath === targetPath || currentPath.startsWith(`${targetPath}/`);
}

export function NavGroup({ items }: NavGroupProps) {
	const currentPath = useRouterState({
		select: (state) => state.location.pathname,
	});
	const { open } = useSidebar();

	return (
		<SidebarMenu className="p-2">
			{items.map((item) => {
				const Icon = item.icon;
				const isActive = isActiveRoute(currentPath, item.to);

				return (
					<SidebarMenuItem key={item.to}>
						<SidebarMenuButton
							isActive={isActive}
							tooltip={item.title}
							render={<Link to={item.to} />}
						>
							<Icon />
							{open && <span>{item.title}</span>}
						</SidebarMenuButton>
					</SidebarMenuItem>
				);
			})}
		</SidebarMenu>
	);
}
