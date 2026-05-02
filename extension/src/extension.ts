import * as vscode from "vscode";
import { execFile } from "child_process";
import * as path from "path";

interface SubtestEntry {
  func: string;
  name: string;
  file: string;
  line: number;
}

export function activate(context: vscode.ExtensionContext) {
  const provider = new SubtestCodeLensProvider();

  context.subscriptions.push(
    vscode.languages.registerCodeLensProvider(
      { language: "go", pattern: "**/*_test.go" },
      provider
    )
  );

  context.subscriptions.push(
    vscode.commands.registerCommand(
      "tdtls.runSubtest",
      (filePath: string, funcName: string, subtestName: string) => {
        runSubtest(filePath, funcName, subtestName, false);
      }
    )
  );

  context.subscriptions.push(
    vscode.commands.registerCommand(
      "tdtls.debugSubtest",
      (filePath: string, funcName: string, subtestName: string) => {
        runSubtest(filePath, funcName, subtestName, true);
      }
    )
  );
}

class SubtestCodeLensProvider implements vscode.CodeLensProvider {
  private _onDidChangeCodeLenses = new vscode.EventEmitter<void>();
  readonly onDidChangeCodeLenses = this._onDidChangeCodeLenses.event;

  constructor() {
    vscode.workspace.onDidSaveTextDocument(() => {
      this._onDidChangeCodeLenses.fire();
    });
  }

  async provideCodeLenses(
    document: vscode.TextDocument
  ): Promise<vscode.CodeLens[]> {
    const subtests = await getSubtests(document.uri.fsPath);
    if (!subtests) {
      return [];
    }

    const lenses: vscode.CodeLens[] = [];
    for (const sub of subtests) {
      // Line numbers from the CLI are 1-based; VS Code is 0-based.
      const line = sub.line - 1;
      const range = new vscode.Range(line, 0, line, 0);

      lenses.push(
        new vscode.CodeLens(range, {
          title: "run subtest",
          command: "tdtls.runSubtest",
          arguments: [sub.file, sub.func, sub.name],
        })
      );

      lenses.push(
        new vscode.CodeLens(range, {
          title: "debug subtest",
          command: "tdtls.debugSubtest",
          arguments: [sub.file, sub.func, sub.name],
        })
      );
    }
    return lenses;
  }
}

function getSubtests(filePath: string): Promise<SubtestEntry[] | null> {
  const config = vscode.workspace.getConfiguration("tdtls");
  const cliPath = config.get<string>("cliPath", "tdtls");

  return new Promise((resolve) => {
    execFile(cliPath, [filePath], { timeout: 10000 }, (err, stdout) => {
      if (err) {
        console.error("tdtls CLI error:", err.message);
        resolve(null);
        return;
      }
      try {
        const result = JSON.parse(stdout);
        resolve(result as SubtestEntry[]);
      } catch {
        console.error("tdtls CLI: invalid JSON output");
        resolve(null);
      }
    });
  });
}

function escapeRegExp(s: string): string {
  return s.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

function runSubtest(
  filePath: string,
  funcName: string,
  subtestName: string,
  debug: boolean
) {
  const dir = path.dirname(filePath);
  // Go test -run uses regex; spaces in subtest names become underscores at runtime.
  const escapedFunc = escapeRegExp(funcName);
  const escapedName = escapeRegExp(subtestName.replace(/ /g, "_"));
  const runFlag = `^${escapedFunc}$\/^${escapedName}$`;

  if (debug) {
    vscode.debug.startDebugging(vscode.workspace.workspaceFolders?.[0], {
      name: `Debug ${funcName}/${subtestName}`,
      type: "go",
      request: "launch",
      mode: "test",
      program: dir,
      args: ["-test.run", runFlag, "-test.v"],
    });
    return;
  }

  const terminal =
    vscode.window.terminals.find((t) => t.name === "tdtls") ??
    vscode.window.createTerminal("tdtls");
  terminal.show();
  terminal.sendText(`go test ${dir} -run '${runFlag}' -v`);
}

export function deactivate() {}
