import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, fetch }) => {
	const lat = url.searchParams.get('lat');
	const lng = url.searchParams.get('lng');
	const radius = url.searchParams.get('radius');
	const year = url.searchParams.get('year');

	if (lat && lng && radius) {
		try {
			const response = await fetch(
				`/api/public/crimes/radius?lat=${lat}&lng=${lng}&radius=${radius}&year=${year}`
			);

			if (response.ok) {
				const data = await response.json();
				return {
					initialSearch: {
						latitude: parseFloat(lat),
						longitude: parseFloat(lng),
						radius: parseFloat(radius),
						result: data
					}
				};
			}
		} catch (error) {
			console.error('Error loading initial search:', error);
		}
	}

	return {
		initialSearch: null
	};
};