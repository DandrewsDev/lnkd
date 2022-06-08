package main

// LinkRecord Structure of link redirect records.
type LinkRecord struct {
	Id          string `form:"_id"          json:"_id"          binding:"-"`
	LnkdUrl     string `form:"lnkd_url"     json:"lnkd_url"     binding:"-"`
	RedirectUrl string `form:"redirect_url" json:"redirect_url" binding:"required"`
	HitCount    int    `form:"hits"         json:"hits"         binding:"-"`
	Username    string `form:"username"     json:"username"     binding:"-"`
}

type LinkListing struct {
	Links []LinkRecord `json:"links"`
}

type LoginRequest struct {
	Username string `form:"username" json:"username"  binding:"required"`
	Password string `form:"password" json:"password"  binding:"required"`
}

type User struct {
	Id       string `form:"_id"         json:"_id"         binding:"-"`
	Username string `form:"username"    json:"username"    binding:"required"`
	Email    string `form:"email"       json:"email"       binding:"required"`
	Password string `form:"password"    json:"password"    binding:"required"`
	Roles    string `form:"roles"       json:"roles"       binding:"-"`
	Status   int    `form:"status"      json:"status"      binding:"-"`
}

type UserListing struct {
	Users []User `json:"users"`
}

type CustomClaims struct {
	Username string `json:"username"`
	Roles    string `json:"roles"`
}

type DisallowedRoutes struct {
	Routes []string `json:"disallowed_routes"`
}
