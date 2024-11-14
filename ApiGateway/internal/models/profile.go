package models

// ProfileData representa datos de perfil del usuario.
type ProfileData struct {
	ID                       string `json:"id"`
	URLPagina                string `json:"personal_url"`
	Apodo                    string `json:"nickname"`
	DireccionCorrespondencia string `json:"address"`
	Biografia                string `json:"biography"`
	Organizacion             string `json:"organization"`
	Pais                     string `json:"country"`
	Token                    string `json:"token,omitempty"`
	ContactPublic            bool   `json:"contact_public"`
}

type UpdateUserData struct {
	Token         string   `json:"token"`
	Nickname      string   `json:"nickname"`
	PersonalURL   string   `json:"personal_url"`
	ContactPublic bool     `json:"contact_public"`
	Address       string   `json:"address"`
	Biography     string   `json:"biography"`
	Organization  string   `json:"organization"`
	Country       string   `json:"country"`
	SocialLinks   []string `json:"social_links"`
}

// CombinedUserData representa la combinación de AuthData y ProfileData.
type CombinedUserData struct {
	AuthData    *AuthData    `json:"auth_data"`
	ProfileData *ProfileData `json:"profile_data"`
}

// CombinedUserUpdateData representa la combinación de datos para actualizar autenticación y perfil.
type CombinedUserUpdateData struct {
	AuthData    AuthData    `json:"auth_data"`
	ProfileData ProfileData `json:"profile_data"`
}
