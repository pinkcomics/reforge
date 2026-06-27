interface Props {
  folder: string;
  onSelect: () => void;
  onChange: (value: string) => void;
  onClear: () => void;
}

export default function FileSelector({
  folder,
  onSelect,
  onChange,
  onClear,
}: Props) {
  return (
    <div className="card">
      <p className="card-title">Selecionar Pasta</p>
      <div className="file-selector__row">
        <input
          className="file-selector__input"
          type="text"
          placeholder="/caminho/para/pasta"
          value={folder}
          onChange={(e) => onChange(e.target.value)}
          spellCheck={false}
        />
        <button className="btn btn--accent" onClick={onSelect}>
          Selecionar
        </button>
        <button className="btn btn--ghost" onClick={onClear} disabled={!folder}>
          Limpar
        </button>
      </div>
    </div>
  );
}
