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

export function TextareaForm() {
	const client = useClient(MemoService);
	const form = useForm<z.infer<typeof FormSchema>>({
		resolver: zodResolver(FormSchema),
		defaultValues: {
			id: "id",
		},
	});

	function onSubmit(data: z.infer<typeof FormSchema>) {
		console.warn("onSubmit run", data);
		client.registerMemo({
			memo: {
				title: data.title,
				text: data.text,
			},
		});
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
			</form>
		</Form>
	);
}
