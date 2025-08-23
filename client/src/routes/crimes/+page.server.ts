import type { CrimesDumpResponse } from "$lib/types";
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
    const res = await fetch('/api/public/crimes/details?limit=1000');
    const data: CrimesDumpResponse = await res.json();
    return {
        crimes: data.crimes
    };
}