import { Resource } from "sst"
import { z } from "zod"
import { issuer } from "@openauthjs/openauth"
import { createSubjects } from "@openauthjs/openauth/subject"
import { CodeProvider } from "@openauthjs/openauth/provider/code"
import { GithubProvider } from "@openauthjs/openauth/provider/github"
import { GoogleOidcProvider } from "@openauthjs/openauth/provider/google"
import { CloudflareStorage } from "@openauthjs/openauth/storage/cloudflare"
import { Account } from ken8n-coder/cloud-core/account.js"

type Env = {
  AuthStorage: KVNamespace
}

export const subjects = createSubjects({
  account: z.object({
    accountID: z.string(),
    email: z.string(),
  }),
  user: z.object({
    userID: z.string(),
    workspaceID: z.string(),
  }),
})

export default {
  async fetch(request: Request, env: Env, ctx: ExecutionContext) {
    return issuer({
      providers: {
        github: GithubProvider({
          clientID: Resource.GITHUB_CLIENT_ID_CONSOLE.value,
          clientSecret: Resource.GITHUB_CLIENT_SECRET_CONSOLE.value,
          scopes: ["read:user", "user:email"],
        }),
        google: GoogleOidcProvider({
          clientID: Resource.GOOGLE_CLIENT_ID.value,
          scopes: ["openid", "email"],
        }),
        //        email: CodeProvider({
        //          async request(req, state, form, error) {
        //            console.log(state)
        //            const params = new URLSearchParams()
        //            if (error) {
        //              params.set("error", error.type)
        //            }
        //            if (state.type === "start") {
        //              return Response.redirect(process.env.AUTH_FRONTEND_URL + "/auth/email?" + params.toString(), 302)
        //            }
        //
        //            if (state.type === "code") {
        //              return Response.redirect(process.env.AUTH_FRONTEND_URL + "/auth/code?" + params.toString(), 302)
        //            }
        //
        //            return new Response("ok")
        //          },
        //          async sendCode(claims, code) {
        //            const email = z.string().email().parse(claims.email)
        //            const cmd = new SendEmailCommand({
        //              Destination: {
        //                ToAddresses: [email],
        //              },
        //              FromEmailAddress: `SST <auth@${Resource.Email.sender}>`,
        //              Content: {
        //                Simple: {
        //                  Body: {
        //                    Html: {
        //                      Data: `Your pin code is <strong>${code}</strong>`,
        //                    },
        //                    Text: {
        //                      Data: `Your pin code is ${code}`,
        //                    },
        //                  },
        //                  Subject: {
        //                    Data: "SST Console Pin Code: " + code,
        //                  },
        //                },
        //              },
        //            })
        //            await ses.send(cmd)
        //          },
        //        }),
      },
      storage: CloudflareStorage({
        namespace: env.AuthStorage,
      }),
      subjects,
      async success(ctx, response) {
        console.log(response)

        let email: string | undefined

        if (response.provider === "github") {
          const userResponse = await fetch("https://api.github.com/user", {
            headers: {
              Authorization: `Bearer ${response.tokenset.access}`,
              "User-Agent": "opencode",
              Accept: "application/vnd.github+json",
            },
          })
          const user = (await userResponse.json()) as { email: string }
          email = user.email
        } else if (response.provider === "google") {
          if (!response.id.email_verified) throw new Error("Google email not verified")
          email = response.id.email as string
        }
        //if (response.provider === "email") {
        //  email = response.claims.email
        //}
        else throw new Error("Unsupported provider")

        if (!email) throw new Error("No email found")

        let accountID = await Account.fromEmail(email).then((x) => x?.id)
        if (!accountID) {
          console.log("creating account for", email)
          accountID = await Account.create({
            email: email!,
          })
        }
        return ctx.subject("account", accountID, { accountID, email })
      },
    }).fetch(request, env, ctx)
  },
}
