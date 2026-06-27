// frontend/src/components/PreviewCard.tsx
interface ScanData {
  cbrCount: number;
  cbzCount: number;
  error?: string;
}

interface Props {
  scan: ScanData | null;
}

export default function PreviewCard({ scan }: Props) {
  return (
    <div className="card">
      <p className="card-title">Prévia</p>

      {!scan && (
        <p className="preview__empty">
          Selecione uma pasta para ver os arquivos.
        </p>
      )}

      {scan?.error && (
        <p className="preview__empty" style={{ color: "var(--error)" }}>
          {scan.error}
        </p>
      )}

      {scan && !scan.error && (
        <div className="preview__stats">
          <div className="preview__stat">
            <span className="preview__stat-value">{scan.cbrCount}</span>
            <span className="preview__stat-label">Arquivos CBR</span>
          </div>
          <div className="preview__stat">
            <span className="preview__stat-value">{scan.cbzCount}</span>
            <span className="preview__stat-label">Arquivos CBZ (mantidos)</span>
          </div>
          <div className="preview__stat">
            <span className="preview__stat-value">
              {scan.cbrCount + scan.cbzCount}
            </span>
            <span className="preview__stat-label">Total</span>
          </div>
        </div>
      )}
    </div>
  );
}
