import type { Summary } from "../App";

interface Props {
  summary: Summary;
  onNewRun: () => void;
}

export default function ResultCard({ summary, onNewRun }: Props) {
  return (
    <div className="card card--result">
      <p className="card-title">
        {summary.cancelled ? "Resumo (cancelado)" : "Resumo final"}
      </p>

      <div className="result__stats">
        <div className="result__stat">
          <span className="result__stat-value">{summary.total}</span>
          <span className="result__stat-label">Analisados</span>
        </div>
        <div className="result__stat">
          <span className="result__stat-value result__stat-value--success">
            {summary.converted}
          </span>
          <span className="result__stat-label">Transformados</span>
        </div>
        <div className="result__stat">
          <span className="result__stat-value">{summary.kept}</span>
          <span className="result__stat-label">Mantidos</span>
        </div>
        <div className="result__stat">
          <span className="result__stat-value result__stat-value--error">
            {summary.failed}
          </span>
          <span className="result__stat-label">Falhas</span>
        </div>
      </div>

      <div className="result__footer">
        <span className="result__time">Tempo total: {summary.elapsedText}</span>
        <button className="btn btn--accent" onClick={onNewRun}>
          Nova execução
        </button>
      </div>
    </div>
  );
}
