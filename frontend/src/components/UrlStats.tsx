import React from "react";
import type { UrlStats as Stats } from "../types/index.ts";

export default function UrlStats({ stats }: { stats: Stats | null }) {
  if (!stats) return null;
  return (
    <div className="max-w-md mx-auto mt-4 p-4 border rounded">
      <div><strong>Original:</strong> <a href={stats.original} target="_blank" rel="noreferrer">{stats.original}</a></div>
      <div><strong>Short code:</strong> {stats.short_code}</div>
      <div><strong>Clicks:</strong> {stats.clicks}</div>
      <div><strong>Created:</strong> {new Date(stats.created_at).toLocaleString()}</div>
    </div>
  );
}