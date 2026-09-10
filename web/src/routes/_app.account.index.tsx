import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router';
import { requireAuth } from '#/lib/route-guards.ts';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs.tsx';
import { IconSettings, IconUser } from '@tabler/icons-react';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card.tsx';
import { UpdateAccountForm } from '#/components/forms/update-account-form.tsx';
import { useAuth } from '#/hooks/use-auth.ts';
import { defaultUserTabSchema } from '#/schemas/users.ts';
import {
	Accordion,
	AccordionContent,
	AccordionItem,
	AccordionTrigger,
} from '#/components/ui/accordion.tsx';
import DefaultApplicationStatusForm from '#/components/forms/default-application-status-form.tsx';

export const Route = createFileRoute('/_app/account/')({
	component: RouteComponent,
	beforeLoad: requireAuth,
	validateSearch: defaultUserTabSchema,
});

function RouteComponent() {
	const { tab } = useSearch({ from: '/_app/account/' });
	const navigate = useNavigate();
	const { user, clearSession, refreshUser } = useAuth();
	if (!user) {
		clearSession();
		void navigate({ to: '/login' });
		return;
	}

	const currentTab = tab ? tab : 'settings';

	return (
		<div className="m-4 flex items-center justify-center">
			<Tabs
				value={currentTab}
				onValueChange={(value) =>
					navigate({
						to: '/account',
						search: (prev) => ({ ...prev, tab: value }),
						replace: true,
					})
				}
				className="flex-1"
			>
				<TabsList className="self-center">
					<TabsTrigger value="settings" onClick={() => {}}>
						<IconSettings />
						Settings
					</TabsTrigger>
					<TabsTrigger value="account">
						<IconUser />
						Account
					</TabsTrigger>
				</TabsList>
				<TabsContent value="account">
					<Card className="flex-1">
						<CardHeader className="p-0">
							<CardTitle className="text-primary border-b pb-4 text-center text-lg">
								User Profile
							</CardTitle>
						</CardHeader>
						<UpdateAccountForm
							username={user.username}
							email={user.email}
							onSuccess={async () => {
								await refreshUser();
							}}
						/>
					</Card>
				</TabsContent>
				<TabsContent value="settings">
					<Card className="flex-1">
						<CardHeader className="p-0">
							<CardTitle className="text-primary border-b pb-4 text-center text-lg">
								Default Settings
							</CardTitle>
						</CardHeader>
						<CardContent>
							<Accordion defaultValue={['defaults']} className="border-0">
								<AccordionItem value="defaults" className="bg-muted/50">
									<AccordionTrigger className="data-panel-open:border-b-border cursor-pointer text-lg no-underline!">
										<div className="text-foreground/80 flex-1 text-center">
											Default Application Status
										</div>
									</AccordionTrigger>
									<AccordionContent className="">
										<div className="pt-4">
											<DefaultApplicationStatusForm
												current={user.default_status_id}
											/>
										</div>
									</AccordionContent>
								</AccordionItem>
							</Accordion>
						</CardContent>
					</Card>
				</TabsContent>
			</Tabs>
		</div>
	);
}
