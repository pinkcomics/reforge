interface Props {
  current: number;
  total: number;
  done: boolean;
}

export default function ProgressCard({ current, total, done }: Props) {
  const pct = total > 0 ? Math.round((current / total) * 100) : 0;

  return (
    <div className="card">
      <p className="card-title">Progresso</p>
      <div className="progress__track">
        <div className="progress__bar" style={{ width: `${pct}%` }} />
      </div>
      <div className="progress__label">
        <span>
          {current} / {total} arquivos
        </span>
        {done ? (
          <span className="progress__done">Concluído ✓</span>
        ) : (
          <span>{pct}%</span>
        )}
      </div>
    </div>
  );
}
