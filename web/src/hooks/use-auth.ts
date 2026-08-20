import { useContext } from 'react';
import { AuthContext } from '#/contexts/auth-context.tsx';
import type { AuthContextValue } from '#/contexts/auth-context.tsx';

export function useAuth(): AuthContextValue {
	const value = useContext(AuthContext);
	if (!value) {
		throw new Error('useAuth must be used within an AuthProvider');
	}
	return value;
}
