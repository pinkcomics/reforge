interface Props {
  skip: boolean;
  replace: boolean;
  onSkipChange: (v: boolean) => void;
  onReplaceChange: (v: boolean) => void;
  disabled: boolean;
}

export default function OptionsCard({
  skip,
  replace,
  onSkipChange,
  onReplaceChange,
  disabled,
}: Props) {
  return (
    <div className="card">
      <p className="card-title">Opções</p>
      <div className="options__list">
        <label className="options__item">
          <input
            type="checkbox"
            checked={replace}
            disabled={disabled}
            onChange={(e) => onReplaceChange(e.target.checked)}
          />
          <div>
            <div className="options__label">Remover CBRs após conversão</div>
          </div>
        </label>
      </div>
    </div>
  );
}
