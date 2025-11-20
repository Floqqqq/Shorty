import { useState } from "react";

export function useLocalStorage<T>(key: string, initial: T) {
    const [state, setState] = useState<T>(() => {
        try {
            const raw = localStorage.getItem(key);
            return raw ? (JSON.parse(raw) as T) : initial;
        } catch {
            return initial;
        }
    });

    const set = (v: T | ((prev: T) => T)) => {
        const value = typeof v === "function" ? (v as (p: T) => T)(state) : v;
        setState(value);
        try {
            localStorage.setItem(key, JSON.stringify(value));
        } catch { }
    };

    return [state, set] as const;
}