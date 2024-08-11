import type { ServiceType } from "@bufbuild/protobuf";
import { type PromiseClient, createPromiseClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { useMemo } from "react";

const transport = createConnectTransport({
	baseUrl: "http://localhost:8080",
});

export function useClient<T extends ServiceType>(service: T): PromiseClient<T> {
	return useMemo(() => createPromiseClient(service, transport), [service]);
}
