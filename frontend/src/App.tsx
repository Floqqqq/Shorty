import React, { useState } from "react";
import UrlShortenerForm from "./components/UrlShortenerForm";
import ShortUrlResult from "./components/ShortUrlResult";
import UrlStats from "./components/UrlStats";
import HistoryList from "./components/HistoryList";
import { useUrlShortener } from "./hooks/useUrlShortener";
import { useLocalStorage } from "./hooks/useLocalStorage";

export default function App() {
  const { create, stats, loading, error } = useUrlShortener();
  const [last, setLast] = useState<string | null>(null);
  const [history, setHistory] = useLocalStorage<string[]>("shorty.history", []);
  const [statsData, setStatsData] = useState<any>(null);

  async function handleShorten(url: string) {
    const res = await create(url);
    if (res) {
      setLast(res.short_url);
      setHistory((prev: any) => [res.short_url, ...prev].slice(0, 20));
    }
  }

  return (
    <main className="min-h-screen p-6">
      <div className="square-pic rotation"></div>
      <h1 className="text-2xl font-bold text-center">Shorty</h1>
      <UrlShortenerForm onShorten={handleShorten} loading={loading} />
      {error && <div className="text-red-600 max-w-md mx-auto">{error}</div>}
      {last && <ShortUrlResult shortUrl={last} onCopy={()=>{}} />}
      <div className="max-w-md mx-auto mt-4">
        <button className="px-3 py-1 border rounded" onClick={async ()=> {
          if (!last) return;
          const code = last.split("/").pop() || "";
          const s = await stats(code);
          setStatsData(s);
        }}>Load Stats for Last</button>
      </div>
      <UrlStats stats={statsData} />
      <HistoryList items={history} onSelect={(it: React.SetStateAction<string | null>)=>setLast(it)} />
    </main>
  );
}