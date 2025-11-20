import React, { useState } from "react";

type Props = {
  onShorten: (url: string) => Promise<void>;
  loading: boolean;
};

export default function UrlShortenerForm({ onShorten, loading }: Props) {
  const [url, setUrl] = useState("");

  function validate(u: string) {
    return u.trim().length > 3;
  }

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate(url)) return;
    await onShorten(url.trim());
    setUrl("");
  };

  return (
    <form onSubmit={submit} className="w-full max-w-md mx-auto p-4 ">
      <label className="block text-sm font-medium mb-2">Enter URL</label>
      <div className="flex gap-2">
        <input
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://example.com"
          className="flex-1 p-2 border rounded form-input"
          aria-label="url"
        />
        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 rounded bg-blue-600 text-white disabled:opacity-50"
        >
          {loading ? "Loading..." : "Shorten"}
        </button>
      </div>
    </form>
  );
}