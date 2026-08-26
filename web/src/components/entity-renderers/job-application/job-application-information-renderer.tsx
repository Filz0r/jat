import type { JobApplication } from '#/api/types.ts';

import { Badge } from '#/components/ui/badge.tsx';
import { ArrowUpRightIcon } from 'lucide-react';
import { getColorFromKind } from '#/lib/utils.ts';

interface JobApplicationInformationRendererProps {
	data: JobApplication;
}

export default function JobApplicationInformationRenderer({
	data,
}: JobApplicationInformationRendererProps) {
	return (
		<div className="flex flex-col items-center justify-center space-y-1.5">
			<div className="m-2 flex w-full flex-col items-center justify-center">
				<h1 className="flex w-full flex-col justify-evenly space-y-2 px-10 py-2 text-center text-lg">
					<div className="flex items-center justify-center">
						<div>
							<span className="text-primary font-bold">Title:</span>
							{data.title}
						</div>

						<Badge
							className="ml-10 decoration-0"
							render={
								<a href={data.url}>
									Link <ArrowUpRightIcon data-icon="inline-end" />
								</a>
							}
						/>
					</div>

					<span>
						<span>
							{/* TODO: change to platform link to website*/}
							<span className="text-primary font-bold">Company:</span>
							{data.company.name}
						</span>
						{data.company.website && (
							<a href={data.company.website} className="text-primary font-bold">
								Website <ArrowUpRightIcon data-icon="inline-end" />
							</a>
						)}
					</span>
				</h1>
				<div className="flex w-full items-center justify-evenly space-x-2 p-2">
					<div className="space-x-2">
						<span>Current Status:</span>
						<Badge className={'flex-1 ' + getColorFromKind(data.status.kind)}>
							{data.status.archived
								? data.status.status + ' (Archived)'
								: data.status.status}
						</Badge>
					</div>
					<div>
						<span className="mr-2 font-bold">Created At:</span>
						<Badge>{new Date(data.created_at).toLocaleString()}</Badge>
					</div>
					<div>
						<span className="mr-2 font-bold">Last Update:</span>
						<Badge>{new Date(data.updated_at).toLocaleString()}</Badge>
					</div>
				</div>
			</div>
		</div>
	);
}
