package endpoints

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Brodobrey-inc/TestService/database"
	"github.com/Brodobrey-inc/TestService/database/structs"
	"github.com/Brodobrey-inc/TestService/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UpdatePersonWebhook(c *gin.Context) {
	personUUID, err := uuid.Parse(c.Param("person_uuid"))
	if err != nil {
		logging.LogError(err, "Failed to parse uuid from query url", "url", c.Request.RequestURI)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var person structs.Person
	if err = database.DB.Get(&person, `SELECT 
	uuid, name, surname, patronymic, age, gender, nation
	FROM person
	WHERE uuid=$1`, personUUID); err != nil {
		logging.LogDebug("Not found person in database")
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	jsonByteArray, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logging.LogError(err, "Failed to read request body")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var newPersonData structs.Person
	if err = json.Unmarshal(jsonByteArray, &newPersonData); err != nil {
		logging.LogError(err, "Failed to read data from request body")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if newPersonData.Name != "" {
		person.Name = newPersonData.Name
	}

	if newPersonData.Surname != "" {
		person.Surname = newPersonData.Surname
	}

	if newPersonData.Patronymic != "" {
		person.Patronymic = newPersonData.Patronymic
	}

	if newPersonData.Age != 0 {
		person.Age = newPersonData.Age
	}

	if newPersonData.Gender != "" {
		person.Gender = newPersonData.Gender
	}

	if newPersonData.Nation != "" {
		person.Nation = newPersonData.Nation
	}

	if _, err = database.DB.NamedExec(`UPDATE person
	SET name=:name, surname=:surname, patronymic=:patronymic,
 	age=:age, gender=:gender, nation=:nation
	WHERE uuid=:uuid`, person); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, person)
}
