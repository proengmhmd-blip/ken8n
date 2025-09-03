import { App } from "../app/app"
import { Config } from "../config/config"
import z from "zod"
import { Provider } from "../provider/provider"
import { generateObject, type ModelMessage } from "ai"
import PROMPT_GENERATE from "./generate.txt"
import PROMPT_WORKFLOW from "../session/prompt/workflow.txt"
import PROMPT_SUPERCODE from "../session/prompt/supercode.txt"
import { SystemPrompt } from "../session/system"
import { mergeDeep } from "remeda"

export namespace Agent {
  export const Info = z
    .object({
      name: z.string(),
      description: z.string().optional(),
      mode: z.union([z.literal("subagent"), z.literal("primary"), z.literal("all")]),
      builtIn: z.boolean(),
      topP: z.number().optional(),
      temperature: z.number().optional(),
      permission: z.object({
        edit: Config.Permission,
        bash: z.record(z.string(), Config.Permission),
        webfetch: Config.Permission.optional(),
      }),
      model: z
        .object({
          modelID: z.string(),
          providerID: z.string(),
        })
        .optional(),
      prompt: z.string().optional(),
      tools: z.record(z.boolean()),
      options: z.record(z.string(), z.any()),
    })
    .openapi({
      ref: "Agent",
    })
  export type Info = z.infer<typeof Info>

  function getDefaultPermission(): Info["permission"] {
    return {
      edit: "allow",
      bash: {
        "*": "allow",
      },
      webfetch: "allow",
    }
  }

  function getBuiltInAgents(agentPermission: Info["permission"]): Record<string, Info> {
    return {
      general: {
        name: "general",
        description:
          "General-purpose agent for researching complex questions, searching for code, and executing multi-step tasks. When you are searching for a keyword or file and are not confident that you will find the right match in the first few tries use this agent to perform the search for you.",
        tools: {
          todoread: false,
          todowrite: false,
        },
        options: {},
        permission: agentPermission,
        mode: "subagent",
        builtIn: true,
      },
      build: {
        name: "build",
        prompt: PROMPT_WORKFLOW,
        tools: {},
        options: {},
        permission: agentPermission,
        mode: "primary",
        builtIn: true,
      },
      plan: {
        name: "plan",
        options: {},
        permission: agentPermission,
        tools: {
          write: false,
          edit: false,
          patch: false,
        },
        mode: "primary",
        builtIn: true,
      },
      "super-code": {
        name: "super-code",
        description:
          "JavaScript code snippet creator for n8n Super Code nodes - writes and validates code that runs inside nodes",
        prompt: PROMPT_SUPERCODE,
        tools: {},
        options: {},
        permission: agentPermission,
        mode: "primary",
        builtIn: true,
      },
    }
  }

  function processConfigAgent(
    result: Record<string, Info>,
    key: string,
    value: any,
    agentPermission: Info["permission"],
    cfg: any,
  ) {
    if (value.disable) {
      delete result[key]
      return
    }

    let item = result[key]
    if (!item) {
      item = result[key] = {
        name: key,
        mode: "all",
        permission: agentPermission,
        options: {},
        tools: {},
        builtIn: false,
      }
    }

    const { model, prompt, tools, description, temperature, top_p, mode, permission, ...extra } = value
    item.options = { ...item.options, ...extra }

    if (model) item.model = Provider.parseModel(model)
    if (prompt) item.prompt = prompt
    if (tools) item.tools = { ...item.tools, ...tools }
    if (description) item.description = description
    if (temperature != undefined) item.temperature = temperature
    if (top_p != undefined) item.topP = top_p
    if (mode) item.mode = mode

    if (permission ?? cfg.permission) {
      item.permission = mergeAgentPermissions(cfg.permission ?? {}, permission ?? {})
    }
  }

  const state = App.state("agent", async () => {
    const cfg = await Config.get()
    const defaultPermission = getDefaultPermission()
    const agentPermission = mergeAgentPermissions(defaultPermission, cfg.permission ?? {})
    const result = getBuiltInAgents(agentPermission)

    for (const [key, value] of Object.entries(cfg.agent ?? {})) {
      processConfigAgent(result, key, value, agentPermission, cfg)
    }

    return result
  })

  export async function get(agent: string) {
    return state().then((x) => x[agent])
  }

  export async function list() {
    return state().then((x) => Object.values(x))
  }

  export async function generate(input: { description: string }) {
    const defaultModel = await Provider.defaultModel()
    const model = await Provider.getModel(defaultModel.providerID, defaultModel.modelID)
    const system = SystemPrompt.header(defaultModel.providerID)
    system.push(PROMPT_GENERATE)
    const existing = await list()
    const result = await generateObject({
      temperature: 0.3,
      prompt: [
        ...system.map(
          (item): ModelMessage => ({
            role: "system",
            content: item,
          }),
        ),
        {
          role: "user",
          content: `Create an agent configuration based on this request: "${input.description}".\n\nIMPORTANT: The following identifiers already exist and must NOT be used: ${existing.map((i) => i.name).join(", ")}\n  Return ONLY the JSON object, no other text, do not wrap in backticks`,
        },
      ],
      model: model.language,
      schema: z.object({
        identifier: z.string(),
        whenToUse: z.string(),
        systemPrompt: z.string(),
      }),
    })
    return result.object
  }
}

function mergeAgentPermissions(basePermission: any, overridePermission: any): Agent.Info["permission"] {
  const merged = mergeDeep(basePermission ?? {}, overridePermission ?? {}) as any
  let mergedBash
  if (merged.bash) {
    if (typeof merged.bash === "string") {
      mergedBash = {
        "*": merged.bash,
      }
    }
    // if granular permissions are provided, default to "ask"
    if (typeof merged.bash === "object") {
      mergedBash = mergeDeep(
        {
          "*": "ask",
        },
        merged.bash,
      )
    }
  }

  const result: Agent.Info["permission"] = {
    edit: merged.edit ?? "allow",
    webfetch: merged.webfetch ?? "allow",
    bash: mergedBash ?? { "*": "allow" },
  }

  return result
}
