"use client";
import { useEffect } from "react";

import { TextareaForm } from "@/components/form/textForm";
import { useMemos } from "@/hooks/useMemos";

export default function Page() {
	const { isLoading, memos, fetchMemos, mergeMemo, deleteMemo } = useMemos();

	useEffect(() => {
		fetchMemos();
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
				<TextareaForm mergeMemo={mergeMemo} deleteMemo={deleteMemo} />
			)}
		</div>
	);
}
