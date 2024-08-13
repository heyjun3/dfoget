import { useMemo, useState } from "react";

import { useClient } from "./client";
import { MemoService } from "@/gen/api/memo/v1/memo_connect";
import { Memo } from "@/gen/api/memo/v1/memo_pb";

export function useMemos() {
	const [isLoading, setIsLoading] = useState(true);
	const [memos, setMemos] = useState<Memo[]>([]);
	const client = useClient(MemoService);

	const mergeMemo = useMemo(() => {
		return (n?: Memo) => {
			if (!n) {
				return;
			}
			setMemos((memos) =>
				memos.find((memo) => memo.id === n.id) ? memos : [n, ...memos],
			);
		};
	}, []);
	const deleteMemo = useMemo(() => {
		return (id?: string) => {
			if (!id) {
				return;
			}
			setMemos((memos) => memos.filter((memo) => memo.id !== id));
			client.deleteMemo({ id: [id] });
		};
	}, []);

	const fetchMemos = useMemo(() => {
		return async () => {
			const res = await client.getMemo({});
			console.warn(res.memo);
			setMemos(res.memo);
			setIsLoading(false);
		};
	}, []);
	return {
		isLoading,
		memos,
		fetchMemos,
		mergeMemo,
		deleteMemo,
	};
}
