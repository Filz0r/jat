import { AccordionContent, AccordionItem, AccordionTrigger } from '#/components/ui/accordion.tsx';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card.tsx';

export default function CompanyStatsRenderer({
	totalCount,
	userCount,
}: {
	companyID: number;
	totalCount: number;
	userCount: number;
}) {
	return (
		<AccordionItem key="stats" value="stats" className="bg-muted/50">
			<AccordionTrigger className="data-panel-open:border-b-border hover:cursor-pointer hover:no-underline">
				<span className="text-primary flex-1 text-center text-2xl font-bold">
					Statistics
				</span>
			</AccordionTrigger>
			<AccordionContent>
				<div className="py-2">
					<div className="mt-2 flex w-full items-center justify-center space-x-2">
						<Card className="w-full">
							<CardHeader className="border-b">
								<CardTitle className="text-primary text-center text-lg">
									Your Applications
								</CardTitle>
							</CardHeader>
							<CardContent className="text-center text-lg">{userCount}</CardContent>
						</Card>
						<Card className="w-full">
							<CardHeader className="border-b">
								<CardTitle className="text-primary text-center text-lg">
									Total Applications
								</CardTitle>
							</CardHeader>
							<CardContent className="text-center text-lg">{totalCount}</CardContent>
						</Card>
					</div>
				</div>
			</AccordionContent>
		</AccordionItem>
	);
}
