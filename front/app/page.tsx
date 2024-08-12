"use client";
import { useEffect, useState } from "react";

import { TextareaForm } from "@/components/form/textForm";
import { MemoService } from "@/gen/api/memo/v1/memo_connect";
import { Memo } from "@/gen/api/memo/v1/memo_pb";
import { useClient } from "@/hooks/client";

export default function Page() {
	const [isLoading, setIsLoading] = useState(true);
	const [memos, setMemos] = useState<Memo[]>([]);
	const client = useClient(MemoService);

	const mergeMemo = (n?: Memo) => {
		if (!n) {
			return;
		}
		setMemos((memos) =>
			memos.find((memo) => memo.id === n.id) ? memos : [n, ...memos],
		);
	};
	const deleteMemo = (id?: string) => {
		if (!id) {
			return;
		}
		setMemos((memos) => memos.filter((memo) => memo.id !== id));
		client.deleteMemo({ id: [id] });
	};

	useEffect(() => {
		(async () => {
			const res = await client.getMemo({});
			console.warn(res.memo);
			setMemos(res.memo);
			setIsLoading(false);
		})();
	}, []);

	if (isLoading) {
		return <></>;
	}
	return (
		<div className="container mx-auto content-center">
			{memos?.length ? (
				memos.map((memo) => (
					<TextareaForm
						key={memo.id}
						memo={memo}
						deleteMemo={deleteMemo}
						mergeMemo={mergeMemo}
					/>
				))
			) : (
				<TextareaForm mergeMemo={mergeMemo} />
			)}
		</div>
	);
}
