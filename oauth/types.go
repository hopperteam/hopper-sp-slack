package oauth

// OAuthV2Response ...
type OAuthV2Response struct {
	AccessToken string                    `json:"access_token"`
	TokenType   string                    `json:"token_type"`
	Scope       string                    `json:"scope"`
	BotUserID   string                    `json:"bot_user_id"`
	AppID       string                    `json:"app_id"`
	Team        OAuthV2ResponseTeam       `json:"team"`
	Enterprise  OAuthV2ResponseEnterprise `json:"enterprise"`
	AuthedUser  OAuthV2ResponseAuthedUser `json:"authed_user"`
}

// OAuthV2ResponseTeam ...
type OAuthV2ResponseTeam struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OAuthV2ResponseEnterprise ...
type OAuthV2ResponseEnterprise struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OAuthV2ResponseAuthedUser ...
type OAuthV2ResponseAuthedUser struct {
	ID          string `json:"id"`
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
