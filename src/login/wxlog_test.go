package login

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestWxAccessToken(t *testing.T) {
	// t.Logf("\n<<content begin>>\n%s\n<<content end>>\n","name")
	router := gin.Default()
	router.GET("/api/:code", func(c *gin.Context) {
		code := c.Param("code")
		accToken, openID, refreshToken, err := fetchWxAccesTockenOpenID(code)
		m := map[string]interface{}{}
		if err != nil {
			m["error"] = err.Error()
		} else {
			m["accToken"] = accToken
			m["openID"] = openID
			m["refreshToken"] = refreshToken
		}
		c.JSON(http.StatusOK, m)
	})
	router.Run(":8080")
}
func TestWxFreshToken(t *testing.T) {
	router := gin.Default()
	router.GET("/api/:code", func(c *gin.Context) {
		code := c.Param("code")
		accToken, openID, err := refreshToken(code)
		m := map[string]interface{}{}
		if err != nil {
			m["error"] = err.Error()
		} else {
			m["accToken"] = accToken
			m["openID"] = openID
		}
		c.JSON(http.StatusOK, m)
	})
	router.Run(":8080")
}
func TestWxUserInfo(t *testing.T) {
	router := gin.Default()
	router.GET("/api/:accTocken/:openID", func(c *gin.Context) {
		accToken := c.Param("accTocken")
		openID := c.Param("openID")
		user, err := fetchUserInfo(accToken,openID)
		m := map[string]interface{}{}
		if err != nil {
			m["error"] = err.Error()
			c.JSON(http.StatusOK, m)
		} else {
			c.JSON(http.StatusOK,user)
		}
	})
	router.Run(":8080")
}
