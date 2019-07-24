//go:generate go run ../cmd/rpcgen/subcmd/examples.go ../cmd/rpcgen/subcmd/doc.go examples
package rpc

import "github.com/october93/engine/model"

// Logging a user with username and password in.
var AuthParamsExample1 = AuthParams{
	Username: "chad",
	Password: "secret",
}

var AuthResponseExample1 = AuthResponse{
	User: &model.ExportedUser{
		ID:               "bb100fc4-29b6-4e2f-b19c-549687b27125",
		DisplayName:      "Chad Unicorn",
		FirstName:        "Chad",
		LastName:         "Unicorn",
		ProfileImagePath: "https://s3-us-west-2.amazonaws.com/assets.october.news/profiles/57f1c2fe-32a1-4b86-ab79-cc3c3e1425ef.jpeg",
		CoverImagePath:   "https://s3-us-west-2.amazonaws.com/assets.october.news/profiles/covers/4bc4db4e-6a0f-49b8-94c0-213c00ca7d47.jpeg",
		Email:            "chad@october.news",
		AllowEmail:       false,
		Username:         "konrad",
		Admin:            false,
	},
	Session: &model.Session{
		ID:     "aa5706e4-baf8-4011-aad0-6d8d55730b8b",
		UserID: "bb100fc4-29b6-4e2f-b19c-549687b27125",
	},
}

// Logging a user in with access token.
var AuthParamsExample2 = AuthParams{
	AccessToken: "EAAblUoKuqcIBAFIx4QNgx48hoUr5ONWdvv8Lr2GbgvBHzAVt9JJZBv1Oaos37koMDwZCHg9vJ2CGDXOz5Ukk9l0vCjs0uRXm3RQFlKcrNUkTi8HNaQvSie4Vdrwe3RZCHcaR2ET5CaMCjt3kiCiLOZB8EkhwrIe4wUfePaoi2ZBPculDSWWozKsV62MAuEjJBaVjlE7nofVC070fNpZBxrCuZBNpwZAzKaFsQZAZA6M7v328ggNDWpoCoj",
}

var AuthResponseExample2 = AuthResponse{
	User: &model.ExportedUser{
		ID:               "bb100fc4-29b6-4e2f-b19c-549687b27125",
		DisplayName:      "Chad Unicorn",
		FirstName:        "Chad",
		LastName:         "Unicorn",
		ProfileImagePath: "https://s3-us-west-2.amazonaws.com/assets.october.news/profiles/57f1c2fe-32a1-4b86-ab79-cc3c3e1425ef.jpeg",
		CoverImagePath:   "https://s3-us-west-2.amazonaws.com/assets.october.news/profiles/covers/4bc4db4e-6a0f-49b8-94c0-213c00ca7d47.jpeg",
		Email:            "chad@october.news",
		AllowEmail:       false,
		Username:         "konrad",
		Admin:            false,
	},
	Session: &model.Session{
		ID:     "aa5706e4-baf8-4011-aad0-6d8d55730b8b",
		UserID: "bb100fc4-29b6-4e2f-b19c-549687b27125",
	},
}

var ResetPasswordParamsExample1 = ResetPasswordParams{
	Email: "chad@october.news",
}
var ResetPasswordResponseExample1 = ResetPasswordResponse{}

var LogoutParamsExample1 = LogoutParams{}
var LogoutResponseExample1 = LogoutResponse{}

var DeleteCardParamsExample1 = DeleteCardParams{
	CardID: "cb9f18a2-2ee5-42b5-97e6-d7f0ec5c2726",
}
var DeleteCardResponseExample1 = DeleteCardResponse{}

var FollowUserParamsExample1 = FollowUserParams{
	UserID: "52ebb15b-ca7a-4682-b72c-4222aeae2e5c",
}
var FollowUserResponseExample1 = FollowUserResponse{}

var UnfollowUserParamsExample1 = UnfollowUserParams{
	UserID: "52ebb15b-ca7a-4682-b72c-4222aeae2e5c",
}
var UnfollowUserResponseExample1 = UnfollowUserResponse{}

var GetFollowingUsersParamsExample1 = GetFollowingUsersParams{}
var GetFollowingUsersResponseExample1 = GetFollowingUsersResponse([]*model.ExportedUser{
	&model.ExportedUser{
		Username:         "chad",
		DisplayName:      "Chad Unicorn",
		ProfileImagePath: "https://s3-us-west-2.amazonaws.com/assets.october.news/profiles/57f1c2fe-32a1-4b86-ab79-cc3c3e1425ef.jpeg",
	}, &model.ExportedUser{
		Username:         "richard",
		DisplayName:      "Richard Hendricks",
		ProfileImagePath: "https://s3-us-west-2.amazonaws.com/assets.october.news/profiles/db32869f-3ca5-4d97-b164-38375e8f83c8.png",
	},
})

var AddToWaitlistParamsExample1 = AddToWaitlistParams{
	AccessToken: "EAAblUoKuqcIBAFIx4QNgx48hoUr5ONWdvv8Lr2GbgvBHzAVt9JJZBv1Oaos37koMDwZCHg9vJ2CGDXOz5Ukk9l0vCjs0uRXm3RQFlKcrNUkTi8HNaQvSie4Vdrwe3RZCHcaR2ET5CaMCjt3kiCiLOZB8EkhwrIe4wUfePaoi2ZBPculDSWWozKsV62MAuEjJBaVjlE7nofVC070fNpZBxrCuZBNpwZAzKaFsQZAZA6M7v328ggNDWpoCoj",
}

var AddToWaitlistResponseExample1 = AddToWaitlistResponse{
	AccessToken:          "EAAblUoKuqcIBAL9RCiu07m1XJ1EpPFMZA6bySc0WMY9lJpxcNmMPwWVirkVZBWl1GUQ9WP2b6SmtnE3MgHQ1EjogMFGXoZBo37jONoY0aIiNjOoKQfMORtohJ1LofZCfgXK4WUfjg6DFQgYk0z0ZAP8OhhXT2l33FcTHONbxcfUe6Htc8sQUM",
	AccessTokenExpiresAt: 5184000,
}

var ValidateUsernameParamsExample1 = ValidateUsernameParams{
	Username: "chad",
}

var ValidateUsernameResponseExample1 = ValidateUsernameResponse{}
