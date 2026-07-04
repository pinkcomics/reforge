import { useState, useEffect, useCallback } from "react";
import { EventsOn } from "../wailsjs/runtime/runtime";
import {
  ScanFolder,
  StartConversion,
  CancelConversion,
  SelectFolder,
  CheckExtractor,
  GetSettings,
  SaveSettings,
} from "../wailsjs/go/main/App";

import FileSelector from "./components/FileSelector";
import PreviewCard from "./components/PreviewCard";
import OptionsCard from "./components/OptionsCard";
import ProgressCard from "./components/ProgressCard";
import LogCard from "./components/LogCard";
import ResultCard from "./components/ResultCard";

import "./App.css";

export interface LogEntry {
  id: number;
  status: "converting" | "done" | "kept" | "error" | "info";
  message: string;
}

export interface Summary {
  total: number;
  converted: number;
  kept: number;
  failed: number;
  cancelled: boolean;
  elapsedMs: number;
  elapsedText: string;
}

interface ScanData {
  cbrCount: number;
  cbzCount: number;
  cbrFiles: string[];
  cbzFiles: string[];
  error?: string;
}

let logIdSeq = 0;

export default function App() {
  const [folder, setFolder] = useState("");
  const [scan, setScan] = useState<ScanData | null>(null);
  const [skip, setSkip] = useState(false);
  const [replace, setReplace] = useState(false);
  const [running, setRunning] = useState(false);
  const [done, setDone] = useState(false);
  const [progress, setProgress] = useState({ current: 0, total: 0 });
  const [log, setLog] = useState<LogEntry[]>([]);
  const [extractorOk, setExtractorOk] = useState<boolean | null>(null);
  const [summary, setSummary] = useState<Summary | null>(null);
  const [settingsLoaded, setSettingsLoaded] = useState(false);

  useEffect(() => {
    CheckExtractor().then(setExtractorOk);
  }, []);

  useEffect(() => {
    GetSettings()
      .then((s) => {
        setReplace(!!s?.replace);
      })
      .finally(() => setSettingsLoaded(true));
  }, []);

  useEffect(() => {
    if (!settingsLoaded) return;
    SaveSettings({ replace });
  }, [replace, settingsLoaded]);

  const addLog = useCallback((status: LogEntry["status"], message: string) => {
    setLog((prev) => [...prev, { id: ++logIdSeq, status, message }]);
  }, []);

  const applyScan = useCallback(
    (dir: string, result: ScanData) => {
      setFolder(dir);
      setScan(result);
      setDone(false);
      setSummary(null);
      setLog([]);
      if (result.error) {
        addLog("error", `Erro ao escanear: ${result.error}`);
      }
    },
    [addLog],
  );

  useEffect(() => {
    const offs = [
      EventsOn("progress:converting", (data: { file: string }) => {
        addLog("converting", `⏳ Processando: ${data.file}`);
      }),
      EventsOn("progress:done", (data: { file: string; dest: string }) => {
        addLog("done", `✓ ${data.file} → ${data.dest}`);
        setProgress((p) => ({ ...p, current: p.current + 1 }));
      }),
      EventsOn("progress:kept", (data: { file: string }) => {
        addLog("kept", `✓ ${data.file} (mantido)`);
        setProgress((p) => ({ ...p, current: p.current + 1 }));
      }),
      EventsOn("progress:error", (data: { file: string; error: string }) => {
        addLog("error", `✗ ${data.file || "erro"}: ${data.error}`);
      }),
      EventsOn("progress:finished", (data: Summary) => {
        setRunning(false);
        setDone(true);
        setSummary(data);
        addLog(
          "info",
          data?.cancelled
            ? "⏹ Conversão cancelada."
            : "✔ Conversão finalizada!",
        );
      }),
    ];

    return () => offs.forEach((off) => off && off());
  }, [addLog, applyScan]);

  const handleSelectFolder = useCallback(async () => {
    const dir = await SelectFolder();
    if (!dir) return;

    const result = await ScanFolder(dir);
    applyScan(dir, result);
  }, [applyScan]);

  const handleFolderChange = useCallback(
    async (value: string) => {
      setFolder(value);
      setDone(false);
      setScan(null);
      setSummary(null);
      setLog([]);

      if (!value.trim()) return;

      const result = await ScanFolder(value.trim());
      setScan(result);

      if (result.error) {
        addLog("error", `Erro ao escanear: ${result.error}`);
      }
    },
    [addLog],
  );

  const handleClear = () => {
    setFolder("");
    setScan(null);
    setLog([]);
    setDone(false);
    setSummary(null);
    setProgress({ current: 0, total: 0 });
  };

  const handleStart = async () => {
    if (!folder || !scan || running) return;

    const total = scan.cbrCount + scan.cbzCount;
    setProgress({ current: 0, total });
    setLog([]);
    setDone(false);
    setSummary(null);
    setRunning(true);

    addLog("info", `Iniciando conversão de ${total} arquivo(s)…`);

    const err = await StartConversion(folder, skip, replace);
    if (err) {
      addLog("error", err);
      setRunning(false);
    }
  };

  const handleCancel = async () => {
    addLog("info", "Cancelando conversão…");
    await CancelConversion();
  };

  const handleNewRun = () => {
    setDone(false);
    setSummary(null);
    setLog([]);
    setProgress({ current: 0, total: 0 });
  };

  const hasFiles = scan && (scan.cbrCount > 0 || scan.cbzCount > 0);
  const canStart = !!folder && hasFiles && !running && extractorOk !== false;
  const showWorkflow = !!folder && scan !== null;

  return (
    <div className="app">
      <header className="app-header">
        <h1 className="app-title">pInk - reforge</h1>
        <p className="app-subtitle">Conversor de arquivos .cbr para .cbz.</p>
      </header>

      {extractorOk === false && (
        <div className="banner banner--error">
          ⚠ Nenhum extrator RAR encontrado. Instale <strong>unrar</strong> ou{" "}
          <strong>7z</strong>.
        </div>
      )}

      <main className="app-main">
        <FileSelector
          folder={folder}
          onSelect={handleSelectFolder}
          onChange={handleFolderChange}
          onClear={handleClear}
        />

        {showWorkflow && (
          <>
            <PreviewCard scan={scan} />

            <OptionsCard
              skip={skip}
              replace={replace}
              onSkipChange={setSkip}
              onReplaceChange={setReplace}
              disabled={running}
            />

            {!running && (
              <button
                className="btn-start"
                onClick={handleStart}
                disabled={!canStart}
              >
                Iniciar conversão
              </button>
            )}

            {running && (
              <button
                className="btn-start btn-start--cancel"
                onClick={handleCancel}
              >
                Cancelar conversão
              </button>
            )}
          </>
        )}

        {(running || done || log.length > 0) && (
          <>
            <ProgressCard
              current={progress.current}
              total={progress.total}
              done={done}
            />
            <LogCard entries={log} />
          </>
        )}

        {done && summary && (
          <ResultCard summary={summary} onNewRun={handleNewRun} />
        )}
      </main>
    </div>
  );
}
