package users

import "github.com/zoommix/fasthttp_template/utils"

// HelloHandler ...
func HelloHandler(c *utils.Client) {
	c.SendJSON(map[string]string{"Message": "Hi there!"})
}
