import { useEffect, useState } from "react"

interface PingResponse {
    code: number;
    message: string;
}

export default function About() {
    const [ping, setPing] = useState<PingResponse | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        async function pinger() {
            try {
                console.log('Current URL:', window.location.href);

                const res = await fetch('/api/ping');

                console.log('Response received:', {
                    status: res.status,
                    statusText: res.statusText,
                    headers: Object.fromEntries(res.headers.entries()),
                    url: res.url
                });

                if (!res.ok) {
                    throw new Error(`HTTP ${res.status}: ${res.statusText}`);
                }

                const data = await res.json();
                setPing(data);

            } catch (err: any) {
                console.error('Fetch error occurred:', {
                    name: err.name,
                    message: err.message,
                    stack: err.stack
                });
                setError(err.message);
            } finally {
                setLoading(false);
                console.log('API call completed');
            }
        }

        pinger();
    }, []);

    if (loading) {
        return (
            <div className="container mx-auto p-8">
                <h1 className="text-4xl font-bold mb-6">About</h1>
                <p className="text-blue-600">Loading...</p>
            </div>
        );
    }

    if (error) {
        return (
            <div className="container mx-auto p-8">
                <h1 className="text-4xl font-bold mb-6">About</h1>
                <p className="text-red-600">Error loading ping: {error}</p>
                <div className="mt-4 p-4 bg-gray-100 rounded">
                    <details>
                        <summary className="cursor-pointer">Debug Info</summary>
                        <pre className="mt-2 text-sm">
                            Current URL: {window.location.href}
                            {'\n'}Attempted fetch: /api/ping
                            {'\n'}Error: {error}
                        </pre>
                    </details>
                </div>
            </div>
        );
    }

    return (
        <div className="container mx-auto p-8">
            <h1 className="text-4xl font-bold mb-6">About</h1>
            {ping ? (
                <div className="space-y-2">
                    <p className="text-green-600">Status Code: {ping.code}</p>
                    <p className="text-green-600">Message: {ping.message}</p>
                </div>
            ) : (
                <p className="text-gray-600">No data received</p>
            )}
        </div>
    );
}