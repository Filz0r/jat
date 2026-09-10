import type { ApplicationStatusChart } from '#/api/types.ts';
import {
	type ChartConfig,
	ChartContainer,
	ChartLegend,
	ChartLegendContent,
	ChartTooltip,
	ChartTooltipContent,
} from '#/components/ui/chart.tsx';
import { Pie, PieChart } from 'recharts';

interface ApplicationStatusPieChartProps {
	data: ApplicationStatusChart;
	title: string;
}
const STATUS_COLORS: Record<string, string> = {
	applied: 'var(--color-orange-500)',
	interviewed: 'var(--color-primary)',
	accepted: 'var(--color-green-500)',
	rejected: 'var(--color-red-500)',
	ghosted: 'var(--color-red-600)',
	irrelevant: 'var(--color-muted)',
};

export default function ApplicationStatusPieChart({ title, data }: ApplicationStatusPieChartProps) {
	const chartConfig = {
		application_status: {
			label: 'Number',
		},
		applied: {
			label: 'Applied',
			color: 'var(--chart-1)',
		},
		irrelevant: {
			label: 'Irrelevant',
			color: 'var(--chart-2)',
		},
		interviewed: {
			label: 'Interviewed',
			color: 'var(--chart-3)',
		},
		rejected: {
			label: 'Rejected',
			color: 'var(--chart-4)',
		},
		ghosted: {
			label: 'Ghosted',
			color: 'var(--chart-5)',
		},
		accepted: {
			label: 'Accepted',
			color: 'var(--chart-6)',
		},
	} satisfies ChartConfig;

	const chartData = Object.entries(data).map(([key, value]) => ({
		label: key,
		value,
		fill: STATUS_COLORS[key] ?? 'var(--color-cyan-500)',
	}));
	return (
		<div className="mt-2 flex-1 rounded-lg border">
			<div className="flex items-center justify-center border-b py-3 text-center text-xl">
				{title}
			</div>
			<ChartContainer config={chartConfig}>
				<PieChart margin={{ top: 0, right: 0, bottom: 0, left: 0 }}>
					<ChartTooltip content={<ChartTooltipContent hideLabel />} />
					<Pie data={chartData} dataKey="value" nameKey="label" />
					<ChartLegend
						content={<ChartLegendContent nameKey="label" />}
						className="my-2 -translate-y-2 flex-wrap gap-2 border-t *:basis-1/4 *:justify-center"
					/>
				</PieChart>
			</ChartContainer>
		</div>
	);
}
