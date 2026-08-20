import { redirect } from '@tanstack/react-router';
import type { ParsedLocation } from '@tanstack/react-router';
import type { RouterContext } from '#/router';

interface GuardContext {
	context: RouterContext;
	location: ParsedLocation;
}

export async function requireAuth(ctx: GuardContext) {
	await ctx.context.authStore.whenReady();
	const { initialized, user } = ctx.context.authStore.getState();

	if (!initialized) {
		throw redirect({ to: '/setup' });
	}

	if (!user) {
		const redirectHref = ctx.location.pathname + ctx.location.search + ctx.location.hash;
		throw redirect({
			to: '/login',
			search: { redirect: redirectHref },
		});
	}
}

export async function requireSetup(ctx: GuardContext) {
	await ctx.context.authStore.whenReady();
	const { initialized } = ctx.context.authStore.getState();

	if (initialized) {
		throw redirect({ to: '/' });
	}
}

export async function redirectIfAuthenticated(ctx: GuardContext) {
	await ctx.context.authStore.whenReady();
	const { initialized, user } = ctx.context.authStore.getState();

	if (!initialized) {
		throw redirect({ to: '/setup' });
	}

	if (user) {
		throw redirect({ to: '/' });
	}
}
