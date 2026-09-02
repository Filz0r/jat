import type { User } from '#/api/types.ts';

export type AuthState = {
	initialized: boolean | null;
	user: User | null;
	isLoading: boolean;
};

export class AuthStore {
	private state: AuthState = { initialized: null, user: null, isLoading: true };
	private listeners = new Set<() => void>();
	private bootPromise: Promise<void>;
	private resolveBoot: () => void;

	constructor() {
		let resolve: () => void;
		this.bootPromise = new Promise<void>((r) => {
			resolve = r;
		});
		this.resolveBoot = resolve!;
	}

	getState(): AuthState {
		return this.state;
	}

	whenReady(): Promise<void> {
		return this.state.isLoading ? this.bootPromise : Promise.resolve();
	}

	setState(updates: Partial<AuthState>) {
		this.state = { ...this.state, ...updates };
		if (updates.isLoading === false) {
			this.resolveBoot();
		}
		this.listeners.forEach((listener) => listener());
	}

	resetBoot() {
		let resolve: () => void;
		this.bootPromise = new Promise<void>((r) => {
			resolve = r;
		});
		this.resolveBoot = resolve!;
	}

	subscribe(listener: () => void) {
		this.listeners.add(listener);
		return () => this.listeners.delete(listener);
	}
}

export const authStore = new AuthStore();
