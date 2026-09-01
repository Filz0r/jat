import { AccordionContent, AccordionItem } from '#/components/ui/accordion.tsx';
import type { CompanyData } from '#/api/types.ts';
import CompanyForm from '#/components/forms/company-form.tsx';
import { IconExternalLink } from '@tabler/icons-react';

export default function CompanyInformationRenderer({ data }: { data: CompanyData }) {
	return (
		<AccordionItem key="info" value="info">
			<div className="relative border-b py-2">
				<p className="text-primary -ml-10 text-center text-2xl font-bold">Information</p>
				<span className="absolute top-3 right-2 space-x-1">
					<CompanyForm edit refreshSelf existingData={data} />
				</span>
			</div>
			<AccordionContent>
				<div className="flex items-center justify-evenly pt-4">
					<div className="space-y-2">
						<div className="flex items-center justify-center space-x-2 text-lg">
							<div className="text-primary text-center">Company Name: </div>
							<span className="text-foreground">{data.name}</span>
						</div>
						{data.website && (
							<div className="flex items-center justify-center space-x-2 text-lg">
								<div className="text-primary text-center">Company Website: </div>
								<div className="text-foreground flex cursor-pointer items-center gap-x-1.5 hover:text-blue-300!">
									<a
										href={data.website}
										className="no-underline! hover:text-blue-300!"
									>
										Website
									</a>
									<IconExternalLink size={14} className="mt-0.5" />
								</div>
							</div>
						)}
					</div>
					<div className="space-y-2">
						<div className="flex items-center justify-center space-x-2 text-lg">
							<div className="text-primary text-center">Created At: </div>
							<span className="text-foreground">
								{new Date(data.created_at).toLocaleDateString() +
									' ' +
									new Date(data.created_at).toLocaleTimeString()}
							</span>
						</div>
						<div className="flex items-center justify-center space-x-2 text-lg">
							<div className="text-primary text-center">Last Update: </div>
							<span className="text-foreground">
								{new Date(data.updated_at).toLocaleDateString() +
									' ' +
									new Date(data.updated_at).toLocaleTimeString()}
							</span>
						</div>
					</div>
				</div>
			</AccordionContent>
		</AccordionItem>
	);
}
