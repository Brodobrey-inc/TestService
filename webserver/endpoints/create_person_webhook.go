package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/Brodobrey-inc/TestService/database"
	"github.com/Brodobrey-inc/TestService/database/structs"
	"github.com/Brodobrey-inc/TestService/logging"
	"github.com/gin-gonic/gin"
)

type PersonRequestBody struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type AgeResponseBody struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   uint   `json:"age"`
}

type GenderResponseBody struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float32 `json:"probability"`
}

type CountryProbability struct {
	CountryId   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type NationResponseBody struct {
	Count  int                  `json:"count"`
	Name   string               `json:"name"`
	County []CountryProbability `json:"country"`
}

func CreatePersonWebhook(c *gin.Context) {
	logging.LogDebug("Got request for CreatePersonWebhook", "requestHost", c.Request.Host)

	jsonDataBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logging.LogError(err, "Failed to read request body")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var person PersonRequestBody
	err = json.Unmarshal(jsonDataBytes, &person)
	if err != nil {
		logging.LogError(err, "Failed to read data from request body")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if person.Name == "" || person.Surname == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	personWithExtra := structs.Person{
		Name:       person.Name,
		Surname:    person.Surname,
		Patronymic: person.Patronymic,
	}
	logging.LogDebug("Try to get person from database")

	query, args, err := database.DB.BindNamed(`SELECT 
		uuid, name, surname, patronymic, age, gender, nation
		FROM person
		WHERE name=:name and surname=:surname and patronymic=:patronymic`,
		personWithExtra,
	)
	if err != nil {
		logging.LogError(err, "Can not prepare sql query")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = database.DB.Get(&personWithExtra, query, args...); err == nil {
		logging.LogDebug("Found person in database")
		c.JSON(http.StatusOK, personWithExtra)
		return
	}

	personWithExtra.Age, err = getEstimatedAgeForName(person.Name)
	if err != nil {
		logging.LogError(err, "Failed to get age for person", "name", person.Name, "surname", person.Surname)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	gender, err := getEstimatedGenderForName(person.Name)
	if err != nil {
		logging.LogError(err, "Failed to get gender for person", "name", person.Name, "surname", person.Surname)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	personWithExtra.Gender = structs.GenderType(gender)

	personWithExtra.Nation, err = getEstimatedNationForName(person.Name)
	if err != nil {
		logging.LogError(err, "Failed to get nation for person", "name", person.Name, "surname", person.Surname)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	query, args, err = database.DB.BindNamed(`INSERT INTO person
		(name, surname, patronymic, age, gender, nation) VALUES
		(:name, :surname, :patronymic, :age, :gender, :nation)
		RETURNING uuid`, personWithExtra,
	)
	if err != nil {
		logging.LogError(err, "Can not prepare sql query")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = database.DB.Get(
		&personWithExtra.UUID, query, args...,
	); err != nil {
		logging.LogError(err, "Failed update database info")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, personWithExtra)
}

// TODO: Maybe generalize getting extra info?
func getEstimatedAgeForName(name string) (uint, error) {
	logging.LogDebug("Request estimated age for person", "name", name)
	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		logging.LogDebug("Failed to get estimated age from side service", "name", name, "err", err)
		return 0, err
	}
	jsonDataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.LogDebug("Failed to get response body", "name", name, "err", err)
		return 0, err
	}

	var estimatedAge AgeResponseBody
	err = json.Unmarshal(jsonDataBytes, &estimatedAge)
	if err != nil {
		logging.LogDebug("Failed to get data from response body", "name", name, "err", err)
		return 0, err
	}
	logging.LogDebug("Got age for requested person", "name", name, "age", estimatedAge)

	return estimatedAge.Age, nil
}

func getEstimatedGenderForName(name string) (string, error) {
	logging.LogDebug("Request estimated gender for person", "name", name)
	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		logging.LogDebug("Failed to get estimated gender from side service", "name", name, "err", err)
		return "", err
	}
	jsonDataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.LogDebug("Failed to get response body", "name", name, "err", err)
		return "", err
	}

	var estimatedGender GenderResponseBody
	err = json.Unmarshal(jsonDataBytes, &estimatedGender)
	if err != nil {
		logging.LogDebug("Failed to get data from response body", "name", name, "err", err)
		return "", err
	}
	logging.LogDebug("Got gender for requested person", "name", name, "gender", estimatedGender)

	return estimatedGender.Gender, nil
}

func getEstimatedNationForName(name string) (string, error) {
	logging.LogDebug("Request estimated nation for person", "name", name)
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		logging.LogDebug("Failed to get estimated nation from side service", "name", name, "err", err)
		return "", err
	}
	jsonDataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.LogDebug("Failed to get response body", "name", name, "err", err)
		return "", err
	}

	var estimatedNation NationResponseBody
	err = json.Unmarshal(jsonDataBytes, &estimatedNation)
	if err != nil {
		logging.LogDebug("Failed to get data from response body", "name", name, "err", err)
		return "", err
	}
	logging.LogDebug("Got nation for requested person", "name", name, "age", estimatedNation)

	sort.Slice(estimatedNation.County, func(i, j int) bool {
		return estimatedNation.County[i].Probability > estimatedNation.County[j].Probability
	})

	return estimatedNation.County[0].CountryId, nil
}
