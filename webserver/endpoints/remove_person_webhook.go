package endpoints

import (
	"net/http"

	"github.com/Brodobrey-inc/TestService/database"
	"github.com/Brodobrey-inc/TestService/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RemovePersonWebhook(c *gin.Context) {
	personUUID, err := uuid.Parse(c.Param("person_uuid"))
	if err != nil {
		logging.LogError(err, "Failed to parse uuid from query url", "url", c.Request.RequestURI)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if _, err = database.DB.Exec(`DELETE
	FROM person WHERE uuid=$1`, personUUID); err != nil {
		logging.LogError(err, "Failed to delete person from database")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
}
