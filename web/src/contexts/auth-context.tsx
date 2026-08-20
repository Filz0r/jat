import {
	createContext,
	type ReactNode,
	useCallback,
	useContext,
	useMemo,
	useSyncExternalStore,
} from 'react';
import type { components } from '#/api/gen-spec';
import { api } from '#/api/client';
import { authStore } from '#/lib/auth-store';
import { queryClient } from '#/router';

export type User = components['schemas']['api.userCreateResponse'];

interface AuthContextValue {
	initialized: boolean | null;
	user: User | null;
	isLoading: boolean;
	refreshUser: () => Promise<void>;
	clearSession: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

function useAuthStoreState() {
	return useSyncExternalStore(
		(callback) => authStore.subscribe(callback),
		() => authStore.getState(),
		() => authStore.getState(),
	);
}

export function AuthProvider({ children }: { children: ReactNode }) {
	const { initialized, user, isLoading } = useAuthStoreState();

	const refreshUser = useCallback(async () => {
		authStore.resetBoot();
		authStore.setState({ isLoading: true });

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
				initialized: true,
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
	}, []);

	const clearSession = useCallback(() => {
		authStore.setState({ user: null });
		queryClient.removeQueries({ queryKey: ['me'] });
	}, []);

	const value = useMemo(
		() => ({
			initialized,
			user,
			isLoading,
			refreshUser,
			clearSession,
		}),
		[initialized, user, isLoading, refreshUser, clearSession],
	);

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
	const value = useContext(AuthContext);
	if (!value) {
		throw new Error('useAuth must be used within an AuthProvider');
	}
	return value;
}
