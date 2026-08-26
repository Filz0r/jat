'use client';

import type { ColumnDef, RowData, SortingState } from '@tanstack/react-table';
import type { DataTableFeatures } from './table-features';
import type { ReactNode } from 'react';

import { useTable } from '@tanstack/react-table';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import { features } from './table-features';
import { Button } from '#/components/ui/button.tsx';
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '#/components/ui/select.tsx';
import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from 'lucide-react';
import { useState } from 'react';

interface DataTableProps<TData extends RowData> {
	columns: ColumnDef<DataTableFeatures, TData>[];
	data: TData[];
	noResultsMessage?: string;
	CreateButton?: ReactNode | null;
}

export function DataTable<TData extends RowData>({
	columns,
	data,
	noResultsMessage = 'No results.',
	CreateButton = null,
}: DataTableProps<TData>) {
	const [sorting, setSorting] = useState<SortingState>([]);

	const table = useTable({
		features,
		data,
		columns,
		onSortingChange: setSorting,
		state: {
			sorting,
		},
	});

	return (
		<div className="overflow-hidden rounded-md border">
			<Table>
				<TableHeader>
					{table.getHeaderGroups().map((headerGroup) => (
						<TableRow key={headerGroup.id} className="bg-muted hover:bg-muted/80">
							{headerGroup.headers.map((header) => {
								return (
									<TableHead key={header.id}>
										{header.isPlaceholder ? null : (
											<table.FlexRender header={header} />
										)}
									</TableHead>
								);
							})}
						</TableRow>
					))}
				</TableHeader>
				<TableBody>
					{table.getRowModel().rows.length ? (
						table.getRowModel().rows.map((row) => (
							<TableRow key={row.id} data-state={row.getIsSelected() && 'selected'}>
								{row.getVisibleCells().map((cell) => (
									<TableCell key={cell.id}>
										<table.FlexRender cell={cell} />
									</TableCell>
								))}
							</TableRow>
						))
					) : (
						<TableRow>
							<TableCell
								colSpan={columns.length}
								className="h-24 text-center text-sm"
							>
								<div className="flex flex-col items-center justify-evenly gap-y-2.5">
									{noResultsMessage}
									{CreateButton !== null ? CreateButton : null}
								</div>
							</TableCell>
						</TableRow>
					)}
				</TableBody>
			</Table>

			<div className="border-t-accent mt-2 flex items-center justify-between border-t p-2">
				<div className="flex items-center space-x-2">
					<p className="text-sm font-medium">Rows per page</p>
					<Select
						value={`${table.state.pagination.pageSize}`}
						onValueChange={(value) => {
							table.setPageSize(Number(value));
						}}
					>
						<SelectTrigger className="h-8 w-17.5">
							<SelectValue placeholder={table.state.pagination.pageSize} />
						</SelectTrigger>
						<SelectContent side="top">
							{[10, 20, 25, 30, 40, 50].map((pageSize) => (
								<SelectItem key={pageSize} value={`${pageSize}`}>
									{pageSize}
								</SelectItem>
							))}
						</SelectContent>
					</Select>
				</div>
				<div className="flex w-25 items-center justify-center text-sm font-medium">
					Page {table.state.pagination.pageIndex + 1} of {table.getPageCount()}
				</div>
				<div className="flex items-center space-x-2">
					<Button
						variant="outline"
						size="icon"
						className="hidden size-8 lg:flex"
						onClick={() => table.setPageIndex(0)}
						disabled={!table.getCanPreviousPage()}
					>
						<span className="sr-only">Go to first page</span>
						<ChevronsLeft />
					</Button>
					<Button
						variant="outline"
						size="icon"
						className="size-8"
						onClick={() => table.previousPage()}
						disabled={!table.getCanPreviousPage()}
					>
						<span className="sr-only">Go to previous page</span>
						<ChevronLeft />
					</Button>
					<Button
						variant="outline"
						size="icon"
						className="size-8"
						onClick={() => table.nextPage()}
						disabled={!table.getCanNextPage()}
					>
						<span className="sr-only">Go to next page</span>
						<ChevronRight />
					</Button>
					<Button
						variant="outline"
						size="icon"
						className="hidden size-8 lg:flex"
						onClick={() => table.setPageIndex(table.getPageCount() - 1)}
						disabled={!table.getCanNextPage()}
					>
						<span className="sr-only">Go to last page</span>
						<ChevronsRight />
					</Button>
				</div>
			</div>
		</div>
	);
}
