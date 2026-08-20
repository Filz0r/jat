import * as React from 'react';
import { createContext, useContext, useState, useCallback, useEffect } from 'react';
import type { components } from '#/api/gen-spec';
import { api } from '#/api/client';

export type User = components['schemas']['api.userCreateResponse'];

interface UserContextValue {
	user: User | null;
	setUser: (user: User | null) => void;
	initialized: boolean | null;
	setInitialized: (value: boolean) => void;
	isLoading: boolean;
	setIsLoading: (value: boolean) => void;
}

const UserContext = createContext<UserContextValue | undefined>(undefined);

export function UserProvider({
	children,
	initialUser,
	initialInitialized,
}: {
	children: React.ReactNode;
	initialUser?: User | null;
	initialInitialized?: boolean | null;
}) {
	const [user, setUserState] = useState<User | null>(initialUser ?? null);
	const [initialized, setInitializedState] = useState<boolean | null>(initialInitialized ?? null);
	const [isLoading, setIsLoading] = useState(true);

	useEffect(() => {
		if (initialUser !== undefined || initialInitialized !== undefined) {
			setIsLoading(false);
			return;
		}

		let cancelled = false;

		async function initialize() {
			try {
				const { data: initData, response: initResponse } = await api.GET('/initialized');

				if (!initResponse.ok || !initData?.ok || initData.data === undefined) {
					if (!cancelled) {
						setInitializedState(false);
						setUserState(null);
						setIsLoading(false);
					}
					return;
				}

				const nextInitialized = initData.data.initialized ?? false;
				let nextUser: User | null = null;

				if (nextInitialized) {
					const { data: userData, response: userResponse } = await api.GET('/users/me');
					if (userResponse.ok && userData?.ok && userData.data) {
						nextUser = userData.data;
					}
				}

				if (!cancelled) {
					setInitializedState(nextInitialized);
					setUserState(nextUser);
					setIsLoading(false);
				}
			} catch (err) {
				console.error('Fatal error communicating with api!', err);
				if (!cancelled) {
					setInitializedState(false);
					setUserState(null);
					setIsLoading(false);
				}
			}
		}

		initialize();

		return () => {
			cancelled = true;
		};
	}, [initialUser, initialInitialized]);

	const setUser = useCallback((next: User | null) => {
		setUserState(next);
	}, []);

	const setInitialized = useCallback((next: boolean) => {
		setInitializedState(next);
	}, []);

	return (
		<UserContext.Provider
			value={{
				user,
				setUser,
				initialized,
				setInitialized,
				isLoading,
				setIsLoading,
			}}
		>
			{children}
		</UserContext.Provider>
	);
}

export function useUser() {
	const value = useContext(UserContext);
	if (!value) {
		throw new Error('useUser must be used within a UserProvider');
	}
	return value;
}
