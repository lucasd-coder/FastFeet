package shared

const ADMIN = "ADMIN"

type Register struct {
	Name      string `json:"name,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Authority string `json:"authority,omitempty"`
}

type RegisterUserResponse struct {
	ID string `json:"id,omitempty"`
}

type GetUserResponse struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Enabled  bool   `json:"enabled,omitempty"`
}

type GetRolesResponse struct {
	Roles []string `json:"roles,omitempty"`
}

type GetToken struct {
	AccessToken      string `json:"access_token,omitempty"`
	ExpiresIn        string `json:"expires_in,omitempty"`
	RefreshExpiresIn string `json:"refresh_expires_in,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	NotBeforePolicy  int    `json:"not_before_policy,omitempty"`
	SessionState     string `json:"session_state,omitempty"`
	Scope            string `json:"scope,omitempty"`
}

type IsActiveUser struct {
	Active bool `json:"active,omitempty"`
}

type ViaCepAddressResponse struct {
	Address      string `json:"logradouro,omitempty"`
	PostalCode   string `json:"cep,omitempty"`
	Neighborhood string `json:"bairro,omitempty"`
	City         string `json:"localidade,omitempty"`
	State        string `json:"uf,omitempty"`
}

func (v *ViaCepAddressResponse) GetPostalCode() string {
	if v != nil {
		return v.PostalCode
	}

	return ""
}
