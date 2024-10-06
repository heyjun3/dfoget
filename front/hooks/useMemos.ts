import { useMemo, useState, useEffect } from "react";

import { useClient } from "./client";
import { MemoService } from "@/gen/api/memo/v1/memo_connect";
import { Memo } from "@/gen/api/memo/v1/memo_pb";

export type MemoType = Pick<Memo, "id" | "title" | "text">;

export function useMemos() {
	const [isLoading, setIsLoading] = useState(true);
	const [memos, setMemos] = useState<MemoType[]>([]);
	const client = useClient(MemoService);

	const mergeMemo = useMemo(() => {
		return (n?: MemoType) => {
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

	const fetchMemos = async () => {
		const stream = client.getMemoServerStream({});
		for await (const res of stream) {
			setMemos(res.memo);
			setIsLoading(false);
		}
	};
	useEffect(() => {
		fetchMemos();
	}, []);
	return {
		isLoading,
		memos,
		fetchMemos,
		mergeMemo,
		deleteMemo,
	};
}
