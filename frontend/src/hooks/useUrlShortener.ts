import { useState } from "react";
import * as api from "../services/api";
import type { ShortenResponse, UrlStats } from "../types";

export function useUrlShortener() {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    async function create(original: string): Promise<ShortenResponse | null> {
        setLoading(true);
        setError(null);
        try {
            const resp = await api.shorten(original);
            return resp.data as ShortenResponse;
        } catch (e: any) {
            setError(e?.response?.data?.error || e.message);
            return null;
        } finally {
            setLoading(false);
        }
    }

    async function stats(code: string): Promise<UrlStats | null> {
        setLoading(true);
        setError(null);
        try {
            const resp = await api.getStats(code);
            return resp.data as UrlStats;
        } catch (e: any) {
            setError(e?.response?.data?.error || e.message);
            return null;
        } finally {
            setLoading(false);
        }
    }

    return { create, stats, loading, error };
}