import { SidebarInset, SidebarProvider, useSidebar } from '#components/ui/sidebar';
import AppSidebar from '#/components/sidebar';
import { Button } from '#components/ui/button';
import { IconLayoutSidebar } from '@tabler/icons-react';
import type { ReactNode } from 'react';

interface AppShellProps {
	children: ReactNode;
}

function MobileSidebarToggle() {
	const { toggleSidebar, isMobile } = useSidebar();
	if (!isMobile) return null;
	return (
		<Button
			onClick={toggleSidebar}
			size="icon-lg"
			className="fixed bottom-4 left-4 z-50 rounded-full shadow-lg"
		>
			<IconLayoutSidebar />
			<span className="sr-only">Toggle sidebar</span>
		</Button>
	);
}

export function AppShell({ children }: AppShellProps) {
	return (
		<SidebarProvider defaultOpen={false}>
			<AppSidebar />
			<SidebarInset className="flex min-w-0 flex-col">
				<main className="min-w-0 flex-1 overflow-auto">{children}</main>
			</SidebarInset>
			<MobileSidebarToggle />
		</SidebarProvider>
	);
}
