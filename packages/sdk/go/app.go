// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package opencode

import (
	"context"
	"net/http"

	"github.com/kenkaiii/ken8n-coder-sdk-go/internal/apijson"
	"github.com/kenkaiii/ken8n-coder-sdk-go/internal/param"
	"github.com/kenkaiii/ken8n-coder-sdk-go/internal/requestconfig"
	"github.com/kenkaiii/ken8n-coder-sdk-go/option"
)

// AppService contains methods and other services that help with interacting with
// the opencode API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAppService] method instead.
type AppService struct {
	Options []option.RequestOption
}

// NewAppService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewAppService(opts ...option.RequestOption) (r *AppService) {
	r = &AppService{}
	r.Options = opts
	return
}

// List all agents
func (r *AppService) Agents(ctx context.Context, opts ...option.RequestOption) (res *[]Agent, err error) {
	opts = append(r.Options[:], opts...)
	path := "agent"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Get app info
func (r *AppService) Get(ctx context.Context, opts ...option.RequestOption) (res *App, err error) {
	opts = append(r.Options[:], opts...)
	path := "app"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Initialize the app
func (r *AppService) Init(ctx context.Context, opts ...option.RequestOption) (res *bool, err error) {
	opts = append(r.Options[:], opts...)
	path := "app/init"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// Write a log entry to the server logs
func (r *AppService) Log(ctx context.Context, body AppLogParams, opts ...option.RequestOption) (res *bool, err error) {
	opts = append(r.Options[:], opts...)
	path := "log"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// List all providers
func (r *AppService) Providers(ctx context.Context, opts ...option.RequestOption) (res *AppProvidersResponse, err error) {
	opts = append(r.Options[:], opts...)
	path := "config/providers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

type Agent struct {
	BuiltIn     bool                   `json:"builtIn,required"`
	Mode        AgentMode              `json:"mode,required"`
	Name        string                 `json:"name,required"`
	Options     map[string]interface{} `json:"options,required"`
	Permission  AgentPermission        `json:"permission,required"`
	Tools       map[string]bool        `json:"tools,required"`
	Description string                 `json:"description"`
	Model       AgentModel             `json:"model"`
	Prompt      string                 `json:"prompt"`
	Temperature float64                `json:"temperature"`
	TopP        float64                `json:"topP"`
	JSON        agentJSON              `json:"-"`
}

// agentJSON contains the JSON metadata for the struct [Agent]
type agentJSON struct {
	BuiltIn     apijson.Field
	Mode        apijson.Field
	Name        apijson.Field
	Options     apijson.Field
	Permission  apijson.Field
	Tools       apijson.Field
	Description apijson.Field
	Model       apijson.Field
	Prompt      apijson.Field
	Temperature apijson.Field
	TopP        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Agent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r agentJSON) RawJSON() string {
	return r.raw
}

type AgentMode string

const (
	AgentModeSubagent AgentMode = "subagent"
	AgentModePrimary  AgentMode = "primary"
	AgentModeAll      AgentMode = "all"
)

func (r AgentMode) IsKnown() bool {
	switch r {
	case AgentModeSubagent, AgentModePrimary, AgentModeAll:
		return true
	}
	return false
}

type AgentPermission struct {
	Bash     map[string]AgentPermissionBash `json:"bash,required"`
	Edit     AgentPermissionEdit            `json:"edit,required"`
	Webfetch AgentPermissionWebfetch        `json:"webfetch"`
	JSON     agentPermissionJSON            `json:"-"`
}

// agentPermissionJSON contains the JSON metadata for the struct [AgentPermission]
type agentPermissionJSON struct {
	Bash        apijson.Field
	Edit        apijson.Field
	Webfetch    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AgentPermission) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r agentPermissionJSON) RawJSON() string {
	return r.raw
}

type AgentPermissionBash string

const (
	AgentPermissionBashAsk   AgentPermissionBash = "ask"
	AgentPermissionBashAllow AgentPermissionBash = "allow"
	AgentPermissionBashDeny  AgentPermissionBash = "deny"
)

func (r AgentPermissionBash) IsKnown() bool {
	switch r {
	case AgentPermissionBashAsk, AgentPermissionBashAllow, AgentPermissionBashDeny:
		return true
	}
	return false
}

type AgentPermissionEdit string

const (
	AgentPermissionEditAsk   AgentPermissionEdit = "ask"
	AgentPermissionEditAllow AgentPermissionEdit = "allow"
	AgentPermissionEditDeny  AgentPermissionEdit = "deny"
)

func (r AgentPermissionEdit) IsKnown() bool {
	switch r {
	case AgentPermissionEditAsk, AgentPermissionEditAllow, AgentPermissionEditDeny:
		return true
	}
	return false
}

type AgentPermissionWebfetch string

const (
	AgentPermissionWebfetchAsk   AgentPermissionWebfetch = "ask"
	AgentPermissionWebfetchAllow AgentPermissionWebfetch = "allow"
	AgentPermissionWebfetchDeny  AgentPermissionWebfetch = "deny"
)

func (r AgentPermissionWebfetch) IsKnown() bool {
	switch r {
	case AgentPermissionWebfetchAsk, AgentPermissionWebfetchAllow, AgentPermissionWebfetchDeny:
		return true
	}
	return false
}

type AgentModel struct {
	ModelID    string         `json:"modelID,required"`
	ProviderID string         `json:"providerID,required"`
	JSON       agentModelJSON `json:"-"`
}

// agentModelJSON contains the JSON metadata for the struct [AgentModel]
type agentModelJSON struct {
	ModelID     apijson.Field
	ProviderID  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AgentModel) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r agentModelJSON) RawJSON() string {
	return r.raw
}

type App struct {
	Git      bool    `json:"git,required"`
	Hostname string  `json:"hostname,required"`
	Path     AppPath `json:"path,required"`
	Time     AppTime `json:"time,required"`
	JSON     appJSON `json:"-"`
}

// appJSON contains the JSON metadata for the struct [App]
type appJSON struct {
	Git         apijson.Field
	Hostname    apijson.Field
	Path        apijson.Field
	Time        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *App) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r appJSON) RawJSON() string {
	return r.raw
}

type AppPath struct {
	Config string      `json:"config,required"`
	Cwd    string      `json:"cwd,required"`
	Data   string      `json:"data,required"`
	Root   string      `json:"root,required"`
	State  string      `json:"state,required"`
	JSON   appPathJSON `json:"-"`
}

// appPathJSON contains the JSON metadata for the struct [AppPath]
type appPathJSON struct {
	Config      apijson.Field
	Cwd         apijson.Field
	Data        apijson.Field
	Root        apijson.Field
	State       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AppPath) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r appPathJSON) RawJSON() string {
	return r.raw
}

type AppTime struct {
	Initialized float64     `json:"initialized"`
	JSON        appTimeJSON `json:"-"`
}

// appTimeJSON contains the JSON metadata for the struct [AppTime]
type appTimeJSON struct {
	Initialized apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AppTime) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r appTimeJSON) RawJSON() string {
	return r.raw
}

type Model struct {
	ID          string                 `json:"id,required"`
	Attachment  bool                   `json:"attachment,required"`
	Cost        ModelCost              `json:"cost,required"`
	Limit       ModelLimit             `json:"limit,required"`
	Name        string                 `json:"name,required"`
	Options     map[string]interface{} `json:"options,required"`
	Reasoning   bool                   `json:"reasoning,required"`
	ReleaseDate string                 `json:"release_date,required"`
	Temperature bool                   `json:"temperature,required"`
	ToolCall    bool                   `json:"tool_call,required"`
	JSON        modelJSON              `json:"-"`
}

// modelJSON contains the JSON metadata for the struct [Model]
type modelJSON struct {
	ID          apijson.Field
	Attachment  apijson.Field
	Cost        apijson.Field
	Limit       apijson.Field
	Name        apijson.Field
	Options     apijson.Field
	Reasoning   apijson.Field
	ReleaseDate apijson.Field
	Temperature apijson.Field
	ToolCall    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Model) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelJSON) RawJSON() string {
	return r.raw
}

type ModelCost struct {
	Input      float64       `json:"input,required"`
	Output     float64       `json:"output,required"`
	CacheRead  float64       `json:"cache_read"`
	CacheWrite float64       `json:"cache_write"`
	JSON       modelCostJSON `json:"-"`
}

// modelCostJSON contains the JSON metadata for the struct [ModelCost]
type modelCostJSON struct {
	Input       apijson.Field
	Output      apijson.Field
	CacheRead   apijson.Field
	CacheWrite  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ModelCost) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelCostJSON) RawJSON() string {
	return r.raw
}

type ModelLimit struct {
	Context float64        `json:"context,required"`
	Output  float64        `json:"output,required"`
	JSON    modelLimitJSON `json:"-"`
}

// modelLimitJSON contains the JSON metadata for the struct [ModelLimit]
type modelLimitJSON struct {
	Context     apijson.Field
	Output      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ModelLimit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelLimitJSON) RawJSON() string {
	return r.raw
}

type Provider struct {
	ID     string           `json:"id,required"`
	Env    []string         `json:"env,required"`
	Models map[string]Model `json:"models,required"`
	Name   string           `json:"name,required"`
	API    string           `json:"api"`
	Npm    string           `json:"npm"`
	JSON   providerJSON     `json:"-"`
}

// providerJSON contains the JSON metadata for the struct [Provider]
type providerJSON struct {
	ID          apijson.Field
	Env         apijson.Field
	Models      apijson.Field
	Name        apijson.Field
	API         apijson.Field
	Npm         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Provider) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r providerJSON) RawJSON() string {
	return r.raw
}

type AppProvidersResponse struct {
	Default   map[string]string        `json:"default,required"`
	Providers []Provider               `json:"providers,required"`
	JSON      appProvidersResponseJSON `json:"-"`
}

// appProvidersResponseJSON contains the JSON metadata for the struct
// [AppProvidersResponse]
type appProvidersResponseJSON struct {
	Default     apijson.Field
	Providers   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AppProvidersResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r appProvidersResponseJSON) RawJSON() string {
	return r.raw
}

type AppLogParams struct {
	// Log level
	Level param.Field[AppLogParamsLevel] `json:"level,required"`
	// Log message
	Message param.Field[string] `json:"message,required"`
	// Service name for the log entry
	Service param.Field[string] `json:"service,required"`
	// Additional metadata for the log entry
	Extra param.Field[map[string]interface{}] `json:"extra"`
}

func (r AppLogParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Log level
type AppLogParamsLevel string

const (
	AppLogParamsLevelDebug AppLogParamsLevel = "debug"
	AppLogParamsLevelInfo  AppLogParamsLevel = "info"
	AppLogParamsLevelError AppLogParamsLevel = "error"
	AppLogParamsLevelWarn  AppLogParamsLevel = "warn"
)

func (r AppLogParamsLevel) IsKnown() bool {
	switch r {
	case AppLogParamsLevelDebug, AppLogParamsLevelInfo, AppLogParamsLevelError, AppLogParamsLevelWarn:
		return true
	}
	return false
}
