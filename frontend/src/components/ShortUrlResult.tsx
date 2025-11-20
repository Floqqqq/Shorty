import React from "react";

type Props = {
  shortUrl: string;
  onCopy?: () => void;
};

export default function ShortUrlResult({ shortUrl, onCopy }: Props) {
  return (
    <div className="p-4 max-w-md mx-auto">
      <div className="p-3 border rounded flex items-center justify-between">
        <a href={shortUrl} target="_blank" rel="noreferrer" className="truncate">
          {shortUrl}
        </a>
        <button
          onClick={() => {
            navigator.clipboard.writeText(shortUrl).then(() => onCopy && onCopy());
          }}
          className="ml-2 px-3 py-1 border rounded"
        >
          Copy
        </button>
      </div>
    </div>
  );
}