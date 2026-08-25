import { apiClient } from '#/api/client.ts';
import { Accordion, AccordionContent, AccordionItem } from '#/components/ui/accordion.tsx';
import ApplicationStatusForm from '#/components/forms/application-status-form.tsx';
import { Badge } from '#/components/ui/badge.tsx';
import { getColorFromKind } from '#/lib/utils.ts';
import { Skeleton } from '#/components/ui/skeleton.tsx';
import { ItemGroup } from '#/components/ui/item.tsx';
import type { JobApplicationStatus } from '#/api/types.ts';
import ApplicationStatusJobApplicationItemRenderer from '#/components/entity-renderers/application-status/application-status-job-application-item-renderer.tsx';

export default function ApplicationStatusRenderer({ data }: { data: JobApplicationStatus }) {
	const {
		data: jobData,
		isLoading: isJobDataLoading,
		error: isJobDataError,
	} = apiClient.useQuery('get', '/jobs', {
		params: {
			query: {
				status_id: [data.id],
			},
		},
	});

	return (
		<div className="my-3.5">
			<Accordion multiple defaultValue={['info', 'notes']}>
				<AccordionItem key="info" value="info" disabled>
					<div className="text-primary grid grid-cols-[1fr_1fr_1fr] border-b py-2 text-2xl font-bold">
						<div />
						<div className="text-center">Information</div>
						<div className="flex justify-end pr-2">
							<ApplicationStatusForm data={data} edit inView />
						</div>
					</div>
					<AccordionContent>
						<div className="mt-2 flex flex-col justify-evenly space-y-2 md:flex-row">
							<div className="flex flex-col items-center">
								<p className="text-primary text-lg">Status</p>
								<Badge className={getColorFromKind(data.kind)}>
									{data.status ? data.status + ' (Archived)' : data.status}
								</Badge>
							</div>
							<div className="flex flex-col items-center">
								<p className="text-primary text-lg">Kind</p>
								<Badge className={getColorFromKind(data.kind)}>
									{data.kind.toUpperCase()}
								</Badge>
							</div>
							<div className="flex flex-col items-center">
								<p className="text-primary text-lg">Created At</p>
								<Badge>{new Date(data.created_at).toLocaleString()}</Badge>
							</div>
							<div className="flex flex-col items-center">
								<p className="text-primary text-lg">Last Update</p>
								<Badge>{new Date(data.updated_at).toLocaleString()}</Badge>
							</div>
						</div>
					</AccordionContent>
				</AccordionItem>

				<AccordionItem key="notes" value="notes" disabled>
					<div className="text-primary flex w-full items-center justify-center border-b py-2 text-2xl font-bold">
						Job Applications
					</div>
					<AccordionContent>
						{isJobDataLoading && <Skeleton className="size-full" />}
						{isJobDataError && (
							<div className="text-center text-2xl text-red-500">
								Error loading Job Applications that use this status{' '}
								{isJobDataError.message}
							</div>
						)}
						{jobData && jobData.data ? (
							jobData.data.length > 0 ? (
								<ItemGroup className="mt-3">
									{jobData.data.map((d, idx) => (
										<ApplicationStatusJobApplicationItemRenderer
											data={d}
											key={`application-status-${data.id}-job-${idx}`}
											idx={idx}
										/>
									))}
								</ItemGroup>
							) : (
								<div className="text-center text-lg">
									No Job Applications with this status were found
								</div>
							)
						) : (
							<div className="text-center text-2xl text-red-500">
								Error Rendering Job Data for this Application Status
							</div>
						)}
					</AccordionContent>
				</AccordionItem>
			</Accordion>
		</div>
	);
}
