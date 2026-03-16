import React from "react";

export default function HistoryList({ items, onSelect }: { items: string[]; onSelect: (s: string)=>void }) {
  if (!items.length)  return null;
  return (
    <div className="max-w-md mx-auto mt-4">
      <h3 className="font-medium mb-2">History</h3>
      <ul className="space-y-2 text-decoration-none">
        {items.map((it, i) => (
          <li key={i} className="flex justify-between items-center border p-2 rounded">
            <a href={it} target="_blank" rel="noreferrer" className="truncate">{it}</a>
            <button onClick={()=> {
                onSelect(it)}}
                className="ml-2 px-2 py-1 border rounded hover:bg-gray-100">Select</button>
          </li>
        ))}
      </ul>
    </div>
  );
}