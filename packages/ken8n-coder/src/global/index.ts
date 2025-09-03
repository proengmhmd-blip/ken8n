import fs from "fs/promises"
import { xdgData, xdgCache, xdgConfig, xdgState } from "xdg-basedir"
import path from "path"

const app = "ken8n-coder"

const data = path.join(xdgData || "~/.local/share", app)
const cache = path.join(xdgCache || "~/.cache", app)
const config = path.join(xdgConfig || "~/.config", app)
const state = path.join(xdgState || "~/.local/state", app)

export namespace Global {
  export const Path = {
    data,
    bin: path.join(data, "bin"),
    log: path.join(data, "log"),
    cache,
    config,
    state,
  } as const
}

await Promise.all([
  fs.mkdir(Global.Path.data, { recursive: true }),
  fs.mkdir(Global.Path.config, { recursive: true }),
  fs.mkdir(Global.Path.state, { recursive: true }),
  fs.mkdir(Global.Path.log, { recursive: true }),
  fs.mkdir(Global.Path.bin, { recursive: true }),
])

const CACHE_VERSION = "8"

const version = await Bun.file(path.join(Global.Path.cache, "version"))
  .text()
  .catch(() => "0")

if (version !== CACHE_VERSION) {
  await fs.rm(Global.Path.cache, { recursive: true, force: true })
  await Bun.file(path.join(Global.Path.cache, "version")).write(CACHE_VERSION)
}
