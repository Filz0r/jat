import createClient from 'openapi-fetch';
import createQueryClient from 'openapi-react-query';
import type { paths } from './gen-spec';

const WEB_CLIENT_TYPE = 'web-client';

let refreshPromise: Promise<boolean> | null = null;

const bodyCache = new WeakMap<Request, string>();

let sessionExpiredHandler: ((redirectHref: string) => void) | null = null;

export function registerSessionExpiredHandler(handler: (redirectHref: string) => void): () => void {
	sessionExpiredHandler = handler;
	return () => {
		if (sessionExpiredHandler === handler) {
			sessionExpiredHandler = null;
		}
	};
}

function notifySessionExpired(redirectHref: string) {
	if (sessionExpiredHandler) {
		sessionExpiredHandler(redirectHref);
	} else if (typeof window !== 'undefined') {
		window.location.replace(`/login?redirect=${encodeURIComponent(redirectHref)}`);
	}
}

async function refreshAccessToken(client: typeof api): Promise<boolean> {
	if (refreshPromise) {
		return refreshPromise;
	}

	refreshPromise = client
		.GET('/auth/refresh_token')
		.then(({ response }) => {
			refreshPromise = null;
			return response.status === 201;
		})
		.catch(() => {
			refreshPromise = null;
			return false;
		});

	return refreshPromise;
}

function isMutationMethod(method: string): boolean {
	return method === 'POST' || method === 'PUT' || method === 'PATCH' || method === 'DELETE';
}

function cloneRequestWithBody(request: Request, bodyText: string | undefined): Request {
	return new Request(request, {
		body: bodyText,
	});
}

function shouldSkipRefresh(pathname: string): boolean {
	return (
		pathname === '/api/auth/refresh_token' ||
		pathname === '/api/auth/login' ||
		pathname === '/api/initialized'
	);
}

export const api = createClient<paths>({
	baseUrl: '/api',
	credentials: 'include',
	headers: {
		'X-Jat-Client-Type': WEB_CLIENT_TYPE,
	},
});

api.use({
	async onRequest({ request }) {
		if (isMutationMethod(request.method)) {
			bodyCache.set(request, await request.clone().text());
		}
		return request;
	},

	async onResponse({ request, response }) {
		if (response.status !== 401) {
			return response;
		}

		const pathname = new URL(request.url).pathname;

		if (shouldSkipRefresh(pathname)) {
			return response;
		}

		const refreshed = await refreshAccessToken(api);
		if (!refreshed) {
			const redirectHref =
				typeof window !== 'undefined'
					? window.location.pathname + window.location.search + window.location.hash
					: '/';
			notifySessionExpired(redirectHref);
			return response;
		}

		const originalBody = bodyCache.get(request);
		const retry = cloneRequestWithBody(request, originalBody);

		return fetch(retry);
	},
});

export const apiClient = createQueryClient(api);

export type ApiClient = typeof api;
