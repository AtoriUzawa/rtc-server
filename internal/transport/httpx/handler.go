package httpx

import "github.com/gin-gonic/gin"

// Handle is a generic request handler that binds the request to a typed struct,
// executes the business logic, and sends the resulting error (or nil) as a JSON response.
func Handle[T any](
	c *gin.Context,
	bind func(*gin.Context, *T) error,
	logic func(*gin.Context, *T) error,
) {
	var req T
	if err := bind(c, &req); err != nil {
		Err(c, err)
		return
	}

	Err(c, logic(c, &req))
}
