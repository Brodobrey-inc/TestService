package endpoints

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Brodobrey-inc/TestService/database"
	"github.com/Brodobrey-inc/TestService/database/structs"
	"github.com/Brodobrey-inc/TestService/logging"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

var (
	defaultPersonsPerList = 2
)

func ListPersonsWebhook(c *gin.Context) {
	var (
		page  = 1
		limit = defaultPersonsPerList
	)

	if pageNumber, err := strconv.Atoi(c.Query("page")); err != nil {
		logging.LogDebug("Ignore page number, set to first page")
	} else {
		logging.LogDebug("Got page number from url")
		page = pageNumber
	}

	if dataLimit, err := strconv.Atoi(c.Query("limit")); err != nil {
		logging.LogDebug("Ignore limit for data, use default value")
	} else {
		logging.LogDebug("Got client data limit")
		limit = dataLimit
	}

	resultCondition := "true"

	if ageFilter := c.Query("age"); ageFilter != "" {
		newConditions, err := parseAgeFilter(ageFilter)
		if err != nil {
			logging.LogError(err, "Failed to parse filter", "column", "age", "filter", ageFilter)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		resultCondition = strings.Join(append(newConditions, resultCondition), " and ")
	}

	if genderFilter := c.Query("gender"); genderFilter != "" {
		newConditions, err := parseGenderFilter(genderFilter)
		if err != nil {
			logging.LogError(err, "Failed to parse filter", "column", "gender", "filter", genderFilter)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		resultCondition = strings.Join(append(newConditions, resultCondition), " and ")
	}

	stringFields := []string{"nation", "name", "surname", "patronymic"}
	for _, field := range stringFields {
		if filter := c.Query(field); filter != "" {
			newConditions, err := parseStringFieldFilter(filter, field)
			if err != nil {
				logging.LogError(err, "Failed to parse filter", "column", field, "filter", filter)
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			resultCondition = strings.Join(append(newConditions, resultCondition), " and ")
		}
	}

	query := fmt.Sprintf(`SELECT 
	uuid, name, surname, patronymic, age, gender, nation
	FROM person
	WHERE %s
	ORDER BY id ASC
	LIMIT $1 OFFSET $2`, resultCondition)
	var personList []structs.Person
	if err := database.DB.Select(
		&personList,
		query,
		limit, (page-1)*limit,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Status(http.StatusNoContent)
		}
		logging.LogError(err, "Failed to execute query")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, personList)
}

func parseAgeFilter(filter string) ([]string, error) {
	if !strings.Contains(filter, ",") {
		requestedAge, err := strconv.Atoi(filter)
		if err != nil {
			return nil, err
		}
		return []string{fmt.Sprintf("age=%d", requestedAge)}, nil
	}

	params := strings.Split(filter, ",")
	if len(params) != 2 {
		return nil, fmt.Errorf("expected 2 arguments for age filter, but got %d", len(params))
	}

	minAge, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, err
	}

	maxAge, err := strconv.Atoi(params[1])
	if err != nil {
		return nil, err
	}

	var parsedConditions []string
	if minAge > 0 {
		parsedConditions = append(parsedConditions, fmt.Sprintf("age>%d", minAge))
	}
	if maxAge > 0 {
		parsedConditions = append(parsedConditions, fmt.Sprintf("age<%d", maxAge))
	}

	return parsedConditions, nil
}

func parseGenderFilter(filter string) ([]string, error) {
	options := []string{string(structs.Male), string(structs.Female)}
	if !slices.Contains(options, filter) {
		return nil, errors.New("option for filter does not exist")
	}

	return []string{fmt.Sprintf("gender='%s'", filter)}, nil
}

func parseStringFieldFilter(filter string, field string) ([]string, error) {
	nationsList := strings.Split(filter, ",")
	parsedConditions := []string{""}

	for _, nation := range nationsList {
		if strings.HasPrefix(nation, "!") {
			parsedConditions = append(parsedConditions, fmt.Sprintf("%s<>'%s'", field, nation))
		} else if parsedConditions[0] != "" {
			parsedConditions[0] = strings.Join([]string{fmt.Sprintf("%s='%s'", field, nation), parsedConditions[0]}, " or ")
		} else {
			parsedConditions[0] = fmt.Sprintf("%s='%s'", field, nation)
		}
	}

	if parsedConditions[0] == "" {
		parsedConditions = parsedConditions[1:]
	}

	return parsedConditions, nil
}
