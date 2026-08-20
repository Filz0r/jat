import createClient from 'openapi-fetch';
import createQueryClient from 'openapi-react-query';
import type { paths } from './gen-spec';

const WEB_CLIENT_TYPE = 'web-client';
const SESSION_EXPIRED_ERROR = 'session expired';
const SESSION_EXPIRED_EVENT = 'jat:session-expired';

let refreshPromise: Promise<boolean> | null = null;

const bodyCache = new WeakMap<Request, string>();

function dispatchSessionExpired() {
	if (typeof window !== 'undefined') {
		window.dispatchEvent(new CustomEvent(SESSION_EXPIRED_EVENT));
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

		const cloned = response.clone();
		const data = await cloned.json().catch(() => null);
		if (typeof data !== 'object' || data === null || data.error !== SESSION_EXPIRED_ERROR) {
			return response;
		}

		const isRefreshRequest = new URL(request.url).pathname === '/api/auth/refresh_token';
		if (isRefreshRequest) {
			dispatchSessionExpired();
			return response;
		}

		const refreshed = await refreshAccessToken(api);
		if (!refreshed) {
			dispatchSessionExpired();
			return response;
		}

		const originalBody = bodyCache.get(request);
		const retry = cloneRequestWithBody(request, originalBody);

		return fetch(retry);
	},
});

export const apiClient = createQueryClient(api);

export type ApiClient = typeof api;
