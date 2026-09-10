import { createFileRoute } from '@tanstack/react-router';
import { apiClient } from '#/api/client.ts';
import { useAuth } from '#/hooks/use-auth.ts';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card.tsx';
import { Skeleton } from '#/components/ui/skeleton.tsx';
import ApplicationStatusPieChart from '#/components/charts/application-status-pie-chart.tsx';
export const Route = createFileRoute('/_app/')({
	component: Home,
});

function Home() {
	const { user } = useAuth();
	const userName = user?.username || 'Unknown User';
	const { data, isLoading, error } = apiClient.useQuery('get', '/users/stats');

	return (
		<div className="p-8">
			<Card>
				<CardHeader>
					<CardTitle>
						<h1 className="text-center text-2xl font-bold">Welcome {userName}</h1>
					</CardTitle>
				</CardHeader>
				<CardContent>
					{isLoading && <Skeleton className="size-full" />}
					{error && (
						<p className="text-destructive text-xl">
							Error loading statistics: {error.message}
						</p>
					)}
					{data && data.data && (
						<div className="mx-8 flex flex-col">
							<div className="flex items-center justify-between space-x-4">
								<Card className="flex-1 text-center">
									<CardHeader className="border-b-border border-b">
										<CardTitle>Total Applications</CardTitle>
									</CardHeader>
									<CardContent>{data.data.total_applications}</CardContent>
								</Card>

								<Card className="flex-1 text-center">
									<CardHeader className="border-b-border border-b">
										<CardTitle>Rejection Percentage</CardTitle>
									</CardHeader>
									<CardContent>
										{data.data.rejection_percentage.toFixed(2)}%
									</CardContent>
								</Card>
								<Card className="flex-1 text-center">
									<CardHeader className="border-b-border border-b">
										<CardTitle>Total Companies Created</CardTitle>
									</CardHeader>
									<CardContent>{data.data.total_companies_created}</CardContent>
								</Card>
							</div>
							<div className="mt-2 flex w-full items-center space-x-5">
								<ApplicationStatusPieChart
									title="Job Applications"
									data={data.data.job_applications_by_status_kind}
								/>
								<ApplicationStatusPieChart
									title="Application Status Counts"
									data={data.data.application_status}
								/>
							</div>
						</div>
					)}
				</CardContent>
			</Card>
		</div>
	);
}
