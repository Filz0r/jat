import { Outlet, createRootRouteWithContext, useRouter } from '@tanstack/react-router';

import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools';
import { TanStackDevtools } from '@tanstack/react-devtools';

import '../styles.css';
import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { UserProvider } from '#/contexts/user-context.tsx';
import { ThemeProvider } from '#/contexts/theme-context.tsx';
import type { RouterContext } from '#/router';

export const Route = createRootRouteWithContext<RouterContext>()({
	component: RootComponent,
});

function RootComponent() {
	const router = useRouter();
	const { queryClient } = router.options.context;

	return (
		<ThemeProvider>
			<QueryClientProvider client={queryClient}>
				<UserProvider>
					<Outlet />
					<TanStackDevtools
						config={{
							position: 'bottom-right',
						}}
						plugins={[
							{
								name: 'TanStack Router',
								render: <TanStackRouterDevtoolsPanel />,
							},
							{
								name: 'TanStack Query Client',
								render: <ReactQueryDevtools />,
							},
						]}
					/>
				</UserProvider>
			</QueryClientProvider>
		</ThemeProvider>
	);
}
