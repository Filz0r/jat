import type { ReactFormExtendedApi } from '@tanstack/react-form';
import { Button } from '#/components/ui/button.tsx';

export default function LogForm({
	form,
}: {
	form: ReactFormExtendedApi<any, any, any, any, any, any, any, any, any, any, any, any>;
}) {
	return (
		<>
			{process.env.NODE_ENV === 'development' ? (
				<Button
					type="button"
					className="w-fit bg-orange-400 hover:bg-orange-500"
					onClick={(e) => {
						e.preventDefault();
						e.stopPropagation();
						console.log(form.state.values);
					}}
				>
					Log Form
				</Button>
			) : null}
		</>
	);
}
