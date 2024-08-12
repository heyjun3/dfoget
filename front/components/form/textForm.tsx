"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/components/ui/form";
import { Textarea } from "@/components/ui/textarea";
import { MemoService } from "@/gen/api/memo/v1/memo_connect";
import { useClient } from "@/hooks/client";
import React from "react";
import { Memo } from "@/gen/api/memo/v1/memo_pb";

const FormSchema = z.object({
	id: z.string().optional(),
	title: z.string(),
	text: z
		.string()
		.min(0, {
			message: "memo must be at least 0 characters.",
		})
		.max(1000, {
			message: "memo must not be longer than 1000 characters.",
		}),
});

type TextareaForm = {
	memo?: Memo;
	deleteMemo?: (id?: string) => void;
	mergeMemo: (memo?: Memo) => void;
};

export function TextareaForm({ memo, deleteMemo, mergeMemo }: TextareaForm) {
	const client = useClient(MemoService);
	const form = useForm<z.infer<typeof FormSchema>>({
		resolver: zodResolver(FormSchema),
		defaultValues: {
			id: memo?.id,
			title: memo?.title,
			text: memo?.text,
		},
	});

	function onSubmit(data: z.infer<typeof FormSchema>) {
		(async () => {
			const res = await client.registerMemo({
				memo: {
					id: data.id,
					title: data.title,
					text: data.text,
				},
			});
			mergeMemo(res.memo);
		})();
	}
	function onClick() {
		console.warn("id", memo?.id);
		if (deleteMemo) {
			deleteMemo(memo?.id);
		}
	}

	return (
		<Form {...form}>
			<form
				onSubmit={form.handleSubmit(onSubmit)}
				className="w-2/3 space-y-6 mx-auto"
			>
				<FormField
					control={form.control}
					name="id"
					render={() => (
						<FormItem>
							<FormControl></FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name="title"
					render={({ field }) => (
						<FormItem>
							<FormLabel>Memo</FormLabel>
							<FormControl>
								<Textarea
									placeholder="Tell us a little bit about yourself"
									className="resize-none min-h-[30px]"
									{...field}
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name="text"
					render={({ field }) => (
						<FormItem>
							<FormControl>
								<Textarea
									placeholder="Tell us a little bit about yourself"
									className="resize-none min-h-[500px]"
									{...field}
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<Button type="submit">Submit</Button>
				<Button className="space-x-1" type="button" onClick={onClick}>
					Delete
				</Button>
			</form>
		</Form>
	);
}
