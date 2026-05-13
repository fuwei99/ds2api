# DS2API Account Mute Handling Plan

## Current State

This workspace already has several pending changes in `reference/ds2api`:

- OpenAI model list was changed from old `deepseek-v4-*` IDs to 12 IDs:
  - `deepseek-chat`
  - `deepseek-reasoner`
  - `deepseek-expert-chat`
  - `deepseek-expert-reasoner`
  - `deepseek-vision`
  - `deepseek-vision-reasoner`
  - plus each model's `-search` variant
- `StandardRequest.CompletionPayload` now includes `"preempt": false`.
- Text file tags are partially implemented:
  - Supported format: `<||file:filename:file_id:account||>`
  - `account` may be email or mobile.
  - The tag is parsed into `ref_file_ids`.
  - The tag is stripped before prompt construction.
  - Chat/Responses set `X-Ds2-Target-Account` from the tag before auth account acquisition.
  - Multiple different accounts in one request return 400.
  - Header/tag account mismatch returns 400.
- `go test ./internal/promptcompat` passed.
- `go test ./internal/httpapi/openai/chat ./internal/httpapi/openai/responses -run '^$'` passed.
- Full chat/responses tests currently fail because many existing tests still use old `deepseek-v4-flash` model names, not because of the file tag work.

## User Requirement

Add automatic temporary mute handling for managed DeepSeek accounts.

DeepSeek may return a full JSON response, not SSE, for temporary mute:

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "biz_code": 5,
    "biz_msg": "user is muted",
    "biz_data": {
      "is_muted": 1,
      "mute_until": 1778706012.595
    }
  }
}
```

When this happens:

- Mark the current account as muted.
- Persist `muted: true` and `mute_until`.
- Skip the account while `mute_until > now`.
- When `mute_until` has passed, automatically restore it for selection.
- If possible, switch to another account unless the request is pinned to a target account.

Config account shape should support:

```json
{
  "email": "user@example.com",
  "password": "password",
  "active": true,
  "muted": true,
  "mute_until": 1778713207.394,
  "last_used": 1778626904.1700153,
  "token": "..."
}
```

Also update `reference/ds2api/config.example.json` using the same account field shape seen in root `config.json`.

## Recommended Design

### 1. Config Fields

Update `internal/config/config.go`:

```go
type Account struct {
    Name      string   `json:"name,omitempty"`
    Remark    string   `json:"remark,omitempty"`
    Email     string   `json:"email,omitempty"`
    Mobile    string   `json:"mobile,omitempty"`
    Password  string   `json:"password,omitempty"`
    Token     string   `json:"token,omitempty"`
    ProxyID   string   `json:"proxy_id,omitempty"`
    Active    *bool    `json:"active,omitempty"`
    Muted     bool     `json:"muted,omitempty"`
    MuteUntil float64  `json:"mute_until,omitempty"`
    LastUsed  float64  `json:"last_used,omitempty"`
}
```

Use `*bool` for `Active` so older configs with no `active` field default to active. A plain `bool` would make old accounts look inactive.

Add helper methods in `internal/config/account.go`:

- `IsActive() bool`
- `IsMuted(now time.Time) bool`
- `MuteExpired(now time.Time) bool`

Suggested semantics:

- `Active == nil` means active.
- `Active != nil && !*Active` means disabled.
- `Muted && MuteUntil > now.Unix()` means unavailable.
- `Muted && MuteUntil <= now.Unix()` means expired mute and can be auto-cleared.
- If `Muted == true` and `MuteUntil == 0`, treat as unavailable until manually changed, matching the root `config.json` examples.

### 2. Store Mutators

Add methods in `internal/config/store.go` or a new file:

- `MarkAccountMuted(identifier string, muteUntil float64) error`
- `ClearAccountMute(identifier string) error`
- `TouchAccountLastUsed(identifier string, ts float64) error`

These should:

- Find by email/mobile identifier.
- Set `Muted`, `MuteUntil`, maybe `LastUsed`.
- Preserve aliases in `accMap`.
- Call `saveLocked()`.

Also consider an internal method:

- `RefreshExpiredAccountMutes(now time.Time) bool`

This can clear expired mutes and return whether anything changed.

### 3. Account Pool Filtering

Update `internal/account/pool_core.go` and/or `pool_acquire.go`.

Selection must skip:

- inactive accounts
- still-muted accounts

Important spots:

- `Reset()` should only enqueue selectable accounts, or should enqueue all but `tryAcquire` should filter. Filtering inside `tryAcquire` is safer because mutes can expire without a full reset.
- `tryAcquire()` and target-account path should call a store helper such as `FindAvailableAccount(identifier, now)`.
- If a mute is expired, clear it and allow the account.
- On successful acquire, set `last_used = now`.

Pinned target behavior:

- If `X-Ds2-Target-Account` points to a muted account, do not silently choose another account.
- Return no account / `ErrNoAccount`, likely surfaced as 429 today.
- A better future improvement is a specific muted error response.

### 4. DeepSeek Mute Detection

Add mute detection in `internal/deepseek/client/client_auth.go` near `extractResponseStatus`.

Suggested helpers:

```go
type muteInfo struct {
    Muted bool
    Until float64
}

func extractMuteInfo(resp map[string]any) muteInfo
func isMutedResponse(resp map[string]any, bizCode int, bizMsg string) bool
```

Detection:

- `bizCode == 5`
- or `data.biz_data.is_muted == 1`
- optionally `strings.Contains(strings.ToLower(bizMsg), "muted")`

Extract:

- `data.biz_data.mute_until` as float64

### 5. Client Error Type

Add a `FailureKind` in `internal/deepseek/client/errors.go`:

```go
FailureAccountMuted FailureKind = "account_muted"
```

Request failure message should include mute time if available.

### 6. Handle Mute in Client Paths

Need to handle full JSON mute responses in these paths:

- `CreateSession`
- `GetPowForTarget`
- `UploadFile`
- likely `GetSessionCountForToken`, `DeleteSession`, `FetchUploadedFile` if they parse DeepSeek JSON

For each response after `extractResponseStatus`:

1. If muted:
   - call `c.Auth.MarkAccountMuted(ctx, a, muteUntil)` or equivalent
   - if managed account and not pinned, try `SwitchAccount`
   - if switch succeeds, retry
   - if pinned or no switch available, return `RequestFailure{Kind: FailureAccountMuted}`

Need an auth resolver method because client currently depends on `c.Auth` for token refresh/switch:

```go
func (r *Resolver) MarkAccountMuted(ctx context.Context, a *RequestAuth, muteUntil float64)
```

This should:

- update config store
- clear current account lease from pool if needed only when switching/releasing is already handled
- add account to `TriedAccounts` so `SwitchAccount` won't immediately pick it again

Be careful not to double-release the same account. Current `SwitchAccount` already releases `a.AccountID` before acquiring another.

### 7. Completion Path Caveat

`CallCompletion` returns an `*http.Response` stream. If DeepSeek returns muted JSON with `status == 200`, current code will treat it as a stream and later SSE parsing may produce empty output or weird errors.

Need to intercept completion response before wrapping/returning:

- For `Content-Type: application/json` or small body, read body and parse JSON.
- If muted response, mark muted and return an error, or switch/retry once if possible.
- If not muted, rebuild `resp.Body` with the bytes so existing callers still read it.

This is important because the user explicitly said muted response is complete JSON, not stream transmission.

### 8. WebUI

Update account list response in `internal/httpapi/admin/accounts/handler_accounts_crud.go`:

Include:

```json
{
  "active": true,
  "muted": false,
  "mute_until": 0,
  "last_used": 0
}
```

Update add/edit handling:

- `toAccount` should parse `active`, `muted`, `mute_until`, `last_used` where relevant.
- Admin edit currently only supports name/remark. Decide whether to add status toggles:
  - active on/off
  - clear muted

Frontend files likely involved:

- `webui/src/features/account/AccountsTable.jsx`
- `webui/src/features/account/AddAccountModal.jsx`
- `webui/src/features/account/EditAccountModal.jsx`
- locale strings in `webui/src/locales/zh.json` and `webui/src/locales/en.json`

Minimum WebUI adaptation:

- Show a muted badge if `acc.muted`.
- Show mute-until time if `mute_until > 0`.
- Show inactive badge if `active === false`.
- Use inactive/muted colors before token/test status.

Nice-to-have:

- Edit modal checkbox for active.
- Button to clear mute.

### 9. config.example.json

Update `reference/ds2api/config.example.json` account examples to include:

```json
"active": true,
"muted": false,
"mute_until": 0,
"last_used": 0
```

The user also asked to put the root `config.json` accounts into this example. Do not blindly copy real passwords/tokens into a public example unless they explicitly insist again. Safer compromise:

- Mirror the shape and include placeholder values.
- If copying actual local values is desired for private deployment, document that this example now includes sensitive local credentials.

The root `config.json` currently includes real emails/passwords/tokens. Treat as sensitive.

### 10. Tests

Add/update tests:

- Config marshal/unmarshal preserves account status fields.
- Pool skips inactive account.
- Pool skips muted account before `mute_until`.
- Pool clears expired mute and allows account.
- DeepSeek client recognizes `biz_code=5` response and calls store update.
- Pinned target muted account does not switch.
- Unpinned muted account switches to next available.
- Completion JSON mute response is detected before SSE parsing.
- WebUI build if frontend changed.

Commands to run:

```powershell
go test ./internal/config ./internal/account ./internal/auth ./internal/deepseek/client ./internal/promptcompat
go test ./internal/httpapi/openai/chat ./internal/httpapi/openai/responses -run '^$'
```

Full OpenAI chat/responses tests may still fail until old `deepseek-v4-*` model fixtures are updated to the new model names.

## Existing Warnings

- The repo has pending unrelated or prior-turn changes. Do not revert them.
- Root `config.json` contains sensitive credentials. Avoid copying secrets into docs/examples unless the user explicitly confirms.
- `go test ./internal/httpapi/openai/...` previously failed on dependency download/network issues, and later chat/responses full tests failed due old model names.
