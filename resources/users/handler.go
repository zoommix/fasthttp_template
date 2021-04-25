package users

import (
	"encoding/json"
	"strconv"

	"github.com/valyala/fasthttp"
	"github.com/zoommix/fasthttp_template/utils"
)

type authData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ShowUsersHandler ...
func ShowUsersHandler(c *fasthttp.RequestCtx) {
	// You can fetch current user object like that:
	// currentUser := c.UserValue("user").(*User)
	// currentUser = currentUser.Reload()
	// utils.LogError(fmt.Sprintf("%+v", currentUser))

	id, _ := strconv.Atoi(c.UserValue("id").(string))
	u, err := FindUser(id)

	if err != nil {
		utils.LogError("Unable to find user: ", err)
		utils.RenderNotFoundError(c, err.Error())
		return
	}

	utils.SetStatus(c, fasthttp.StatusOK)
	utils.WriteJSON(c, UserToJSON(u))
	return
}

// CreateUserHandler ...
func CreateUserHandler(c *fasthttp.RequestCtx) {
	user := userParams(c)

	validationErrors := user.Validate()

	if len(validationErrors) > 0 {
		utils.RenderValidationErrors(c, validationErrors)
		return
	}

	user, err := CreateUser(user)

	if err != nil {
		utils.LogError("Unable to create user", err)
		utils.RenderInternalError(c, err.Error())
		return
	}

	utils.SetStatus(c, fasthttp.StatusCreated)
	utils.WriteJSON(c, UserToJSON(user))
	return
}

func userParams(c *fasthttp.RequestCtx) *User {
	params := NewUserJSON()
	json.Unmarshal(c.PostBody(), &params)

	u := &User{
		Email:     params.Email,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Password:  params.Password,
	}

	return u
}

// Authenticate ...
func Authenticate(c *fasthttp.RequestCtx) {
	creds := authParams(c)
	encodedPass := EncodePassword(creds.Password)
	user, err := FindUserByEmail(creds.Email)

	if err != nil || user.PasswordDigest != encodedPass {
		utils.LogError("invalid token or password", err)
		utils.RenderUnauthorized(c, "invalid token or password")
		return
	}

	user.GenerateJWT()

	utils.SetStatus(c, fasthttp.StatusOK)
	utils.WriteJSON(c, UserToJSON(user))
}

func authParams(c *fasthttp.RequestCtx) authData {
	data := &authData{}
	json.Unmarshal(c.PostBody(), &data)
	return *data
}
