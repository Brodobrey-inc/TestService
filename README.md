# Test Service

This service is created for interaction with API's for prediction information about person by his name.

## Start

Before start service you need create and configure `.env` file in repo root. You can configure next options: 

* `DEBUG_LEVEL` - service logging level ("warning", "info", "debug", "trace")
* `SERVICE_PORT` - service port for listening request's
* `POSTGRES_PORT` - postgres database port
* `POSTGRES_HOST` - postgres database host
* `POSTGRES_DB_NAME` - postgres database name
* `POSTGRES_USER` - postgres database user
* `POSTGRES_PASSWORD` - postgres database password

Then you can start service by command

```
go run main.go
```

# Interaction

This service has 4 REST API methods:

* `/create_person` - creates new person in database with predicted age, gender and nation. On success return JSON with information about added person.Expect next JSON in request

```JSON
{
    "name": "Dmitriy",
    "surname": "Ushakov",
    "patronymic": "Vasilevich" // not required
}
```
* `/update_person/:person_uuid` - update fields in database for person with the specified uuid. Expect JSON in format:
```JSON
{
    ...
    "field": "value"
    ...
}
```
* `/update_person/:person_uuid` - remove person from database with the specified uuid
* `/list_persons` - list persons with pagination, without any parameters return first 2 persons from database

## Available parameters for `/list_persons` method

This method has following parameters:

* `page` - number of page
* `limit` - number of elements on page
*  `name`, `surname`, `patronymic`, `nation` - filter applied on corresponding filed, specified in format `filed=value_1,!value_2...` where `!` means to exclude objects with given field value
* `age` - filter on age field specified in format `age=minAge,maxAge` or `age=exactAge`. Id you want ignore `minAge` or `maxAge` set them to -1
* `gender` - filter with to options `male` or `female`