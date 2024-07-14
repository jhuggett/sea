import { appendFile } from "node:fs/promises";
import { join } from "path";
import { execDir } from "../backend/process";

export const logPath = join(execDir, "frontend.log");

const formatLog = (level: string, contents: string) => {
  return `Frontend ${new Date().toLocaleTimeString()} ${level} ${contents} \n`;
};

export const overrideConsole = () => {
  console = {
    ...console,
    log: (...args: any[]) => {
      console.info(...args);
    },
    error: (...args: any[]) => {
      appendFile(logPath, formatLog(`Error`, Bun.inspect(args)));
    },
    warn: (...args: any[]) => {
      appendFile(logPath, formatLog(`Warning`, Bun.inspect(args)));
    },
    info: (...args: any[]) => {
      appendFile(logPath, formatLog(`Info`, Bun.inspect(args)));
    },
    debug: (...args: any[]) => {
      appendFile(logPath, formatLog(`Debug`, Bun.inspect(args)));
    },
  };
};
