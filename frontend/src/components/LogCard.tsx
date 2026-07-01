import { useEffect, useRef } from "react";
import type { LogEntry } from "../App";

interface Props {
  entries: LogEntry[];
}

export default function LogCard({ entries }: Props) {
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [entries]);

  return (
    <div className="card">
      <p className="card-title">Log de Progresso</p>
      <div className="log__list">
        {entries.map((e) => (
          <div key={e.id} className={`log__entry log__entry--${e.status}`}>
            {e.message}
          </div>
        ))}
        <div ref={bottomRef} />
      </div>
    </div>
  );
}
