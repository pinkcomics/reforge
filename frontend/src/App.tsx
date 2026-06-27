import { useState, useEffect, useCallback } from "react";
import { EventsOn } from "../wailsjs/runtime/runtime";
import {
  ScanFolder,
  StartConversion,
  SelectFolder,
  CheckExtractor,
} from "../wailsjs/go/main/App";

import FileSelector from "./components/FileSelector";
import PreviewCard from "./components/PreviewCard";
import OptionsCard from "./components/OptionsCard";
import ProgressCard from "./components/ProgressCard";
import LogCard from "./components/LogCard";

import "./App.css";

export interface LogEntry {
  id: number;
  status: "converting" | "done" | "kept" | "error" | "info";
  message: string;
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

  useEffect(() => {
    CheckExtractor().then(setExtractorOk);
  }, []);

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
      EventsOn("progress:finished", () => {
        setRunning(false);
        setDone(true);
        addLog("info", "✔ Conversão finalizada!");
      }),
    ];

    return () => offs.forEach((off) => off && off());
  }, []);

  const addLog = (status: LogEntry["status"], message: string) => {
    setLog((prev) => [...prev, { id: ++logIdSeq, status, message }]);
  };

  const handleSelectFolder = useCallback(async () => {
    const dir = await SelectFolder();
    if (!dir) return;

    setFolder(dir);
    setDone(false);
    setLog([]);

    const result = await ScanFolder(dir);
    setScan(result);

    if (result.error) {
      addLog("error", `Erro ao escanear: ${result.error}`);
    }
  }, []);

  const handleFolderChange = useCallback(async (value: string) => {
    setFolder(value);
    setDone(false);
    setScan(null);
    setLog([]);

    if (!value.trim()) return;

    const result = await ScanFolder(value.trim());
    setScan(result);

    if (result.error) {
      addLog("error", `Erro ao escanear: ${result.error}`);
    }
  }, []);

  const handleClear = () => {
    setFolder("");
    setScan(null);
    setLog([]);
    setDone(false);
    setProgress({ current: 0, total: 0 });
  };

  const handleStart = async () => {
    if (!folder || !scan || running) return;

    const total = scan.cbrCount + scan.cbzCount;
    setProgress({ current: 0, total });
    setLog([]);
    setDone(false);
    setRunning(true);

    addLog("info", `Iniciando conversão de ${total} arquivo(s)…`);

    const err = await StartConversion(folder, skip, replace);
    if (err) {
      addLog("error", err);
      setRunning(false);
    }
  };

  const hasFiles = scan && (scan.cbrCount > 0 || scan.cbzCount > 0);
  const canStart = !!folder && hasFiles && !running && extractorOk !== false;

  return (
    <div className="app">
      <header className="app-header">
        <h1 className="app-title">pInk | Conversor</h1>
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

        <PreviewCard scan={scan} />

        <OptionsCard
          skip={skip}
          replace={replace}
          onSkipChange={setSkip}
          onReplaceChange={setReplace}
          disabled={running}
        />

        <button
          className="btn-start"
          onClick={handleStart}
          disabled={!canStart}
        >
          {running ? "Convertendo…" : "Converter"}
        </button>

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
      </main>
    </div>
  );
}
