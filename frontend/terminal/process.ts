import { join, sep } from "path";
import { appendFile } from "node:fs/promises";
import { logPath } from "../terminal/log";

export const execDir = process.execPath.split(sep).slice(0, -1).join(sep);

export const startBackend = () => {
  const proc = Bun.spawn([join(execDir, "server")], {
    onExit(subprocess, exitCode, signalCode, error) {
      console.log("Exited with code", exitCode);
    },
    stdout: "pipe",
    stderr: "pipe",
  });

  (async () => {
    // @ts-ignore
    for await (const chunk of proc.stderr) {
      //console.log("Server [STDERR]: " + new TextDecoder().decode(chunk));
      appendFile(logPath, new TextDecoder().decode(chunk));
    }
  })();

  (async () => {
    // @ts-ignore
    for await (const chunk of proc.stdout) {
      //console.log("Server [STDOUT]: " + new TextDecoder().decode(chunk));
      appendFile(logPath, new TextDecoder().decode(chunk));
    }
  })();

  return {
    stopBackend() {
      proc.kill();
    },
  };
};
