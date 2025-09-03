#!/usr/bin/env bun

import { $ } from "bun"

console.log("=== publishing ===\n")

const snapshot = process.env["KEN8N_CODER_SNAPSHOT"] === "true"
const version = snapshot
  ? `0.0.0-${new Date().toISOString().slice(0, 16).replace(/[-:T]/g, "")}`
  : process.env["KEN8N_CODER_VERSION"]
if (!version) {
  throw new Error("KEN8N_CODER_VERSION is required")
}
process.env["KEN8N_CODER_VERSION"] = version
console.log("version:", version)

const pkgjsons = await Array.fromAsync(
  new Bun.Glob("**/package.json").scan({
    absolute: true,
  }),
).then((arr) => arr.filter((x) => !x.includes("node_modules") && !x.includes("dist")))

const tree = await $`git add . && git write-tree`.text().then((x) => x.trim())
for (const file of pkgjsons) {
  let pkg = await Bun.file(file).text()
  pkg = pkg.replaceAll(/"version": "[^"]+"/g, `"version": "${version}"`)
  console.log("updated:", file)
  await Bun.file(file).write(pkg)
}
await $`bun install`

console.log("\n=== ken8n-coder ===\n")
await import(`../packages/ken8n-coder/script/publish.ts`)

console.log("\n=== sdk ===\n")
await import(`../packages/sdk/js/script/publish.ts`)

console.log("\n=== plugin ===\n")
await import(`../packages/plugin/script/publish.ts`)

const dir = new URL("..", import.meta.url).pathname
process.chdir(dir)

if (!snapshot) {
  await $`git commit -am "release: v${version}"`
  await $`git tag v${version}`
  await $`git fetch origin`
  await $`git cherry-pick HEAD..origin/dev`.nothrow()
  await $`git push origin HEAD --tags --no-verify --force`

  const previous = await fetch("https://api.github.com/repos/kenkaiii/ken8n-coder/releases/latest")
    .then((res) => {
      if (!res.ok) throw new Error(res.statusText)
      return res.json()
    })
    .then((data) => data.tag_name)

  console.log("finding commits between", previous, "and", "HEAD")
  const commits = await fetch(`https://api.github.com/repos/kenkaiii/ken8n-coder/compare/${previous}...HEAD`)
    .then((res) => res.json())
    .then((data) => data.commits || [])

  const raw = commits.map((commit: any) => `- ${commit.commit.message.split("\n").join(" ")}`)
  console.log(raw)

  const notes =
    raw
      .filter((x: string) => {
        const lower = x.toLowerCase()
        return (
          !lower.includes("release:") &&
          !lower.includes("ignore:") &&
          !lower.includes("chore:") &&
          !lower.includes("ci:") &&
          !lower.includes("wip:") &&
          !lower.includes("docs:") &&
          !lower.includes("doc:")
        )
      })
      .join("\n") || "No notable changes"

  await $`gh release create v${version} --title "v${version}" --notes ${notes} ./packages/ken8n-coder/dist/*.zip`
}
if (snapshot) {
  await $`git checkout -b snapshot-${version}`
  await $`git commit --allow-empty -m "Snapshot release v${version}"`
  await $`git tag v${version}`
  await $`git push origin v${version} --no-verify`
  await $`git checkout dev`
  await $`git branch -D snapshot-${version}`
  for (const file of pkgjsons) {
    await $`git checkout ${tree} ${file}`
  }
}
