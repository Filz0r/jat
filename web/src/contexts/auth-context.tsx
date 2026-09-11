import type { ReactNode, Context } from 'react';
import type { User } from '#/api/types.ts';
import type { QueryClient } from '@tanstack/react-query';

import { createContext, useCallback, useMemo, useSyncExternalStore } from 'react';
import { api } from '#/api/client';
import { authStore } from '#/lib/auth-store';

export interface AuthContextValue {
	initialized: boolean | null;
	user: User | null;
	isLoading: boolean;
	refreshUser: () => Promise<void>;
	clearSession: () => void;
	loadInit: () => Promise<void>;
}

const AUTH_CONTEXT_KEY = '__jat_auth_context__';

const existingContext =
	typeof globalThis !== 'undefined'
		? (globalThis as Record<string, unknown>)[AUTH_CONTEXT_KEY]
		: undefined;

export const AuthContext: Context<AuthContextValue | undefined> =
	(existingContext as Context<AuthContextValue | undefined> | undefined) ??
	createContext<AuthContextValue | undefined>(undefined);

if (typeof globalThis !== 'undefined') {
	(globalThis as Record<string, unknown>)[AUTH_CONTEXT_KEY] = AuthContext;
}

function useAuthStoreState() {
	return useSyncExternalStore(
		(callback) => authStore.subscribe(callback),
		() => authStore.getState(),
		() => authStore.getState(),
	);
}

export function AuthProvider({
	children,
	queryClient,
}: {
	children: ReactNode;
	queryClient: QueryClient;
}) {
	const { initialized, user, isLoading } = useAuthStoreState();

	const loadInit = useCallback(async () => {
		authStore.setState({ isLoading: true });
		authStore.resetBoot();
		if (initialized) {
			return;
		}
		try {
			const initResult = await queryClient.fetchQuery({
				queryKey: ['initialized'],
				queryFn: () => api.GET('/initialized'),
				staleTime: 0,
			});

			const initInitialized =
				initResult.response.ok && initResult.data?.ok && initResult.data.data?.initialized;

			if (!initInitialized) {
				authStore.setState({
					initialized: false,
					user: null,
					isLoading: false,
				});
				return;
			}
		} catch (e) {
			console.error('Failed to refresh auth state', e);
			authStore.setState({
				initialized: false,
				user: null,
				isLoading: false,
			});
			return;
		}
	}, [queryClient]);

	const refreshUser = useCallback(async () => {
		authStore.resetBoot();
		authStore.setState({ isLoading: true });

		try {
			let nextUser: User | null = null;

			const meResult = await queryClient.fetchQuery({
				queryKey: ['me'],
				queryFn: () => api.GET('/users/me'),
				staleTime: 0,
			});

			if (meResult.response.ok && meResult.data?.ok && meResult.data.data) {
				nextUser = meResult.data.data;
			}

			authStore.setState({
				user: nextUser,
				isLoading: false,
			});
		} catch (err) {
			console.error('Failed to refresh auth state', err);
			authStore.setState({
				initialized: false,
				user: null,
				isLoading: false,
			});
		}
	}, [queryClient]);

	const clearSession = useCallback(() => {
		authStore.setState({ user: null });
		queryClient.clear();
	}, [queryClient]);

	const value = useMemo(
		() => ({
			initialized,
			user,
			isLoading,
			refreshUser,
			clearSession,
			loadInit,
		}),
		[initialized, user, isLoading, loadInit, refreshUser, clearSession],
	);

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
