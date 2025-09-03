import { spawnSync } from "node:child_process"
import { platform } from "node:os"
import path from "node:path"

const root = process.cwd()
const isWin = platform() === "win32"

const scriptPath = path.join(root, "script", isWin ? "hooks.bat" : "hooks")

const result = spawnSync(isWin ? "cmd" : scriptPath, isWin ? ["/c", scriptPath] : [], {
  stdio: "inherit",
  cwd: root,
  env: process.env,
})

if (result.status !== 0) {
  process.exit(result.status ?? 1)
}

