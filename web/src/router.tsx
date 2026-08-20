import { QueryClient } from '@tanstack/react-query';
import { createRouter } from '@tanstack/react-router';
import { routeTree } from './routeTree.gen';
import { authStore } from '#/lib/auth-store';
import type { AuthStore } from '#/lib/auth-store';

const FIVE_MINUTES = 1000 * 60 * 5;

export const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			staleTime: FIVE_MINUTES,
		},
	},
});

export interface RouterContext {
	queryClient: QueryClient;
	authStore: AuthStore;
}

export function getRouter(context?: Partial<RouterContext>) {
	return createRouter({
		routeTree,
		scrollRestoration: true,
		defaultPreload: 'intent',
		defaultPreloadStaleTime: FIVE_MINUTES,
		context: {
			queryClient,
			authStore,
			...context,
		},
	});
}
