package controllers

import (
	"app/models"
	"app/repositories"
	"app/services"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func initTestRouter(dbHandler *sql.DB) *gin.Engine {
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	usersRepository := repositories.NewUsersRepository(dbHandler)
	runnersService := services.NewRunnersService(runnersRepository, nil)
	usersServices := services.NewUsersService(usersRepository)
	runnersController := NewRunnersController(runnersService, usersServices)
	router := gin.Default()
	router.GET("/runner", runnersController.GetRunnersBatch)
	return router
}

func TestGetRunnersResponse(t *testing.T) {
	dbHandler, mock, _ := sqlmock.New()
	defer dbHandler.Close()
	columns := []string{"id", "first_name", "last_name", "age",
		"is_active", "country", "personal_best",
		"season_best"}
	mock.ExpectQuery("SELECT *").WillReturnRows(
		sqlmock.NewRows(columns).
			AddRow("1", "John", "Smith", 30, true,
				"United States", "02:00:41",
				"02:13:13").
			AddRow("2", "Marijana", "Komatinovic", 30,
				true, "Serbia", "01:18:28",
				"01:18:28"))
	columnsUsers := []string{"user_role"}
	mock.ExpectQuery("SELECT user_role").WillReturnRows(
		sqlmock.NewRows(columnsUsers).AddRow("runner"))
	router := initTestRouter(dbHandler)
	request, _ := http.NewRequest("GET", "/runner", nil)
	recorder := httptest.NewRecorder()
	request.Header.Set("token", "token")
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK,
		recorder.Result().StatusCode)
	var runners []*models.Runner
	json.Unmarshal(recorder.Body.Bytes(), &runners)
	assert.NotEmpty(t, runners)
	assert.Equal(t, 2, len(runners))
}
