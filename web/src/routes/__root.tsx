'use client';

import { useEffect } from 'react';
import { Outlet, createRootRouteWithContext, useRouter } from '@tanstack/react-router';

import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools';
import { TanStackDevtools } from '@tanstack/react-devtools';

import '../styles.css';
import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtoolsPanel } from '@tanstack/react-query-devtools';
import { AuthProvider } from '#/contexts/auth-context.tsx';
import { ThemeProvider } from '#/contexts/theme-context.tsx';
import { api, registerSessionExpiredHandler } from '#/api/client';
import { authStore } from '#/lib/auth-store';
import type { User } from '#/lib/auth-store';
import type { RouterContext } from '#/router';
import { Toaster } from '#/components/ui/toast.tsx';
import { TooltipProvider } from '#/components/ui/tooltip.tsx';

export const Route = createRootRouteWithContext<RouterContext>()({
	component: RootComponent,
	beforeLoad: async ({ context }) => {
		if (!context.authStore.getState().isLoading) {
			return {};
		}

		try {
			const initResult = await context.queryClient.fetchQuery({
				queryKey: ['initialized'],
				queryFn: () => api.GET('/initialized'),
				staleTime: 0,
			});

			const initialized =
				initResult.response.ok && initResult.data?.ok && initResult.data.data?.initialized;

			if (!initialized) {
				context.authStore.setState({
					initialized: false,
					user: null,
					isLoading: false,
				});
				return {};
			}

			const pathname = typeof window !== 'undefined' ? window.location.pathname : '/login';
			const isPublic = pathname === '/login' || pathname === '/setup';
			let nextUser: User | null = null;

			if (!isPublic) {
				const meResult = await context.queryClient.fetchQuery({
					queryKey: ['me'],
					queryFn: () => api.GET('/users/me'),
					staleTime: 0,
				});

				if (meResult.response.ok && meResult.data?.ok && meResult.data.data) {
					nextUser = meResult.data.data;
				}
			}

			context.authStore.setState({
				initialized: true,
				user: nextUser,
				isLoading: false,
			});
		} catch (err) {
			console.error('Fatal error during auth boot', err);
			context.authStore.setState({
				initialized: false,
				user: null,
				isLoading: false,
			});
		}

		return {};
	},
});

function SessionExpiredHandler() {
	const router = useRouter();

	useEffect(() => {
		return registerSessionExpiredHandler((redirectHref) => {
			if (authStore.getState().isLoading) {
				return;
			}

			authStore.setState({ user: null });
			router.navigate({
				to: '/login',
				search: { redirect: redirectHref },
				replace: true,
			});
		});
	}, [router]);

	return null;
}

function RootComponent() {
	const router = useRouter();
	const { queryClient } = router.options.context;

	return (
		<ThemeProvider>
			<TooltipProvider>
				<QueryClientProvider client={queryClient}>
					<AuthProvider queryClient={queryClient}>
						<SessionExpiredHandler />
						<Outlet />
						<TanStackDevtools
							config={{
								position: 'middle-left',
							}}
							plugins={[
								{
									name: 'TanStack Router',
									render: <TanStackRouterDevtoolsPanel />,
								},
								{
									name: 'TanStack Query Client',
									render: <ReactQueryDevtoolsPanel />,
								},
							]}
						/>
					</AuthProvider>
				</QueryClientProvider>
				<Toaster />
			</TooltipProvider>
		</ThemeProvider>
	);
}
