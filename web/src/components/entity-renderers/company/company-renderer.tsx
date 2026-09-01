import { Accordion } from '#/components/ui/accordion.tsx';
import CompanyInformationRenderer from '#/components/entity-renderers/company/company-information-renderer.tsx';
import { useAuth } from '#/hooks/use-auth.ts';
import type { CompanyData } from '#/api/types.ts';
import CompanyStatsRenderer from '#/components/entity-renderers/company/company-stats-renderer.tsx';
import CompanyApplicationRenderer from '#/components/entity-renderers/company/company-application-renderer.tsx';
import CompanyChangeHistoryRenderer from '#/components/entity-renderers/company/company-change-history-renderer.tsx';

export default function CompanyRenderer({ data }: { data: CompanyData }) {
	const { user } = useAuth();
	const isAdmin = user?.is_admin || false;

	return (
		<div className="my-3.5">
			<Accordion multiple defaultValue={['info', 'applications', 'stats']}>
				<CompanyInformationRenderer data={data} />
				<CompanyStatsRenderer
					companyID={data.id}
					totalCount={data.total_count || 0}
					userCount={data.user_count || 0}
				/>
				<CompanyApplicationRenderer companyID={data.id} />

				{isAdmin && <CompanyChangeHistoryRenderer companyID={data.id} />}
			</Accordion>
		</div>
	);
}
