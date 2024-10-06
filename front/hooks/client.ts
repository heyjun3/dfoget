import type { ServiceType } from "@bufbuild/protobuf";
import {
	Code,
	ConnectError,
	Interceptor,
	type PromiseClient,
	createPromiseClient,
} from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { AppRouterInstance } from "next/dist/shared/lib/app-router-context.shared-runtime";
import { useRouter } from "next/navigation";
import { useMemo } from "react";

const url =
	"http://localhost:8888/realms/myrealm/protocol/openid-connect/auth?response_type=code&client_id=myclient&redirect_uri=http://localhost:8080/oidc&scope=openid";

function handleUnauthenticated(e: unknown, router: AppRouterInstance) {
	if (e instanceof ConnectError && e.code == Code.Unauthenticated) {
		router.push(url);
	}
}

async function* streamMessages(
	stream: AsyncIterable<any>,
	router: AppRouterInstance,
) {
	try {
		for await (const m of stream) {
			yield m;
		}
	} catch (e) {
		handleUnauthenticated(e, router);
	}
}

const logger = (router: AppRouterInstance): Interceptor => {
	const _logger: Interceptor = (next) => async (req) => {
		try {
			const res = await next(req);
			if (res.stream) {
				return {
					...res,
					message: streamMessages(res.message, router),
				};
			}
			return res;
		} catch (e) {
			handleUnauthenticated(e, router);
			throw e;
		}
	};
	return _logger;
};

export function useClient<T extends ServiceType>(service: T): PromiseClient<T> {
	const router = useRouter();
	const transport = createConnectTransport({
		baseUrl: "http://localhost:8080",
		useBinaryFormat: false,
		credentials: "include",
		interceptors: [logger(router)],
	});
	return useMemo(() => createPromiseClient(service, transport), [service]);
}
