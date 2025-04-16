package project

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sca-integrator/app/controller"
	"sca-integrator/app/dbo/repository"
	"sca-integrator/app/dbo/repository/project"
	request2 "sca-integrator/app/dto/request"
	"sca-integrator/app/middleware"
	"sca-integrator/app/service"
	"sca-integrator/test"
	"sca-integrator/test/dummy"
	project2 "sca-integrator/test/dummy/project"
	"testing"
	"time"
)

func setupInitialFilterOption(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Router init
	route := gin.Default()

	route.Use(middleware.ErrorHandler())
	router := route.Group("api/v1")

	optionRepository := project.NewFilterOptionRepository()
	projectRepository := repository.NewProjectRepository()

	optionService := service.NewProjectFilterOptionService(optionRepository, projectRepository, validator.New(), db)

	optionController := controller.NewProjectFilterOptionController(optionService)

	router.GET("filter-options", optionController.GetAllHandler)
	router.GET("filter-option/:id", optionController.GetDetailByIdHandler)
	router.POST("filter-option", optionController.CreateOptionHandler)
	router.PUT("filter-option/:id", optionController.UpdateOptionHandler)
	router.DELETE("filter-option/:id", optionController.DeleteOptionHandler)

	return route
}

func TestGetAllFilterOption(t *testing.T) {
	sqlDB, db, mock := test.DbMock(t)
	defer sqlDB.Close()

	router := setupInitialFilterOption(db)

	f1 := project2.FilterOptionDummyList[0]
	fRow1 := project2.MappingFilterOptionStore(f1, 1)
	t.Run("Success Test", func(t *testing.T) {
		options := sqlmock.NewRows(project2.FilterOptionCols).AddRow(fRow1...)

		mock.ExpectQuery(regexp.QuoteMeta(test.SelectAllFilterOptionSQL)).WillReturnRows(options)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/filter-options", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		bodyByte, err := ioutil.ReadAll(response.Body)

		t.Logf("%+v", string(bodyByte))

		assert.Equal(t, 200, response.StatusCode)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestGetDetailByIdFilterOption(t *testing.T) {
	sqlDB, db, mock := test.DbMock(t)
	defer sqlDB.Close()

	router := setupInitialFilterOption(db)

	f1 := project2.FilterOptionDummyList[0]
	fRow1 := project2.MappingFilterOptionStore(f1, 1)
	t.Run("Success Test", func(t *testing.T) {
		option := sqlmock.NewRows(project2.FilterOptionCols).AddRow(fRow1...)

		mock.ExpectQuery(test.SelectOneByIdFilterOptionSQL).WithArgs(1).WillReturnRows(option)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/filter-option/1", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		bodyByte, err := ioutil.ReadAll(response.Body)

		t.Logf("%+v", string(bodyByte))

		assert.Equal(t, 200, response.StatusCode)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		option := sqlmock.NewRows(project2.FilterOptionCols)

		mock.ExpectQuery(test.SelectOneByIdFilterOptionSQL).WithArgs(1).WillReturnRows(option)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/filter-option/1", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		bodyByte, err := ioutil.ReadAll(response.Body)

		t.Logf("%+v", string(bodyByte))

		assert.Equal(t, 404, response.StatusCode)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCreateFilterOption(t *testing.T) {
	sqlDB, db, mock := test.DbMock(t)
	defer sqlDB.Close()

	router := setupInitialFilterOption(db)

	f1 := project2.FilterOptionDummyList[0]
	args := test.AnyThingArgs(4)

	p1 := dummy.ProjectDummyList[0]
	pRow1 := dummy.MappingProjectStore(p1, 1)
	t.Run("Success Test", func(t *testing.T) {
		project := sqlmock.NewRows(dummy.ProjectCols).AddRow(pRow1...)

		mock.ExpectQuery(test.SelectOneByIdProjectSQL).WithArgs(1).WillReturnRows(project)

		mock.ExpectBegin()
		mock.ExpectQuery(test.InsertFilterOptionSQL).
			WithArgs(args...).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at"}).
				AddRow(1, time.Now(), time.Now(), time.Now()))
		mock.ExpectCommit()

		optionRequest := request2.CreateOptionRequest{
			Project_id:  1,
			Filter_type: f1.FilterType,
			Value:       f1.Value,
		}

		reqJson, _ := json.Marshal(optionRequest)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/filter-option", bytes.NewBuffer(reqJson))
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		bodyByte, err := ioutil.ReadAll(response.Body)

		t.Logf("%+v", string(bodyByte))

		assert.Equal(t, 200, response.StatusCode)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Bad Request", func(t *testing.T) {
		optionRequest := request2.CreateOptionRequest{}

		reqJson, _ := json.Marshal(optionRequest)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/filter-option", bytes.NewBuffer(reqJson))
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		bodyByte, err := ioutil.ReadAll(response.Body)

		t.Logf("%+v", string(bodyByte))

		assert.Equal(t, 400, response.StatusCode)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Project Not Found", func(t *testing.T) {
		project := sqlmock.NewRows(dummy.ProjectCols)

		mock.ExpectQuery(test.SelectOneByIdProjectSQL).WithArgs(1).WillReturnRows(project)

		optionRequest := request2.CreateOptionRequest{
			Project_id:  1,
			Filter_type: f1.FilterType,
			Value:       f1.Value,
		}

		reqJson, _ := json.Marshal(optionRequest)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/filter-option", bytes.NewBuffer(reqJson))
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		bodyByte, err := ioutil.ReadAll(response.Body)

		t.Logf("%+v", string(bodyByte))

		assert.Equal(t, 404, response.StatusCode)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateFilterOption(t *testing.T) {
	sqlDB, db, mock := test.DbMock(t)
	defer sqlDB.Close()

	router := setupInitialFilterOption(db)

	f1 := project2.FilterOptionDummyList[0]
	fRow1 := project2.MappingFilterOptionStore(f1, 1)
	args := test.AnyThingArgs(7)

	p1 := dummy.ProjectDummyList[0]
	pRow1 := dummy.MappingProjectStore(p1, 1)

	t.Run("Success Test", func(t *testing.T) {
		project := sqlmock.NewRows(dummy.ProjectCols).AddRow(pRow1...)
		option := sqlmock.NewRows(project2.FilterOptionCols).AddRow(fRow1...)

		mock.ExpectQuery(test.SelectOneByIdProjectSQL).WithArgs(1).WillReturnRows(project)
		mock.ExpectQuery(test.SelectOneByIdFilterOptionSQL).WithArgs(1).WillReturnRows(option)

		mock.ExpectBegin()
		mock.
			ExpectExec(test.UpdateFilterOptionSQL).
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		optionRequest := request2.UpdateOptionRequest{
			Project_id:  1,
			Filter_type: f1.FilterType,
			Value:       f1.Value,
		}

		reqJson, _ := json.Marshal(optionRequest)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/v1/filter-option/1", bytes.NewBuffer(reqJson))
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		bodyByte, err := ioutil.ReadAll(response.Body)

		t.Logf("%+v", string(bodyByte))

		assert.Equal(t, 200, response.StatusCode)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Project Not Found", func(t *testing.T) {
		project := sqlmock.NewRows(dummy.ProjectCols)

		mock.ExpectQuery(test.SelectOneByIdProjectSQL).WithArgs(1).WillReturnRows(project)

		optionRequest := request2.UpdateOptionRequest{
			Project_id:  1,
			Filter_type: f1.FilterType,
			Value:       f1.Value,
		}

		reqJson, _ := json.Marshal(optionRequest)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/v1/filter-option/1", bytes.NewBuffer(reqJson))
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		bodyByte, err := ioutil.ReadAll(response.Body)

		t.Logf("%+v", string(bodyByte))

		assert.Equal(t, 404, response.StatusCode)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Filter Option Not Found", func(t *testing.T) {
		project := sqlmock.NewRows(dummy.ProjectCols).AddRow(pRow1...)
		option := sqlmock.NewRows(project2.FilterOptionCols)

		mock.ExpectQuery(test.SelectOneByIdProjectSQL).WithArgs(1).WillReturnRows(project)
		mock.ExpectQuery(test.SelectOneByIdFilterOptionSQL).WithArgs(1).WillReturnRows(option)

		optionRequest := request2.UpdateOptionRequest{
			Project_id:  1,
			Filter_type: f1.FilterType,
			Value:       f1.Value,
		}

		reqJson, _ := json.Marshal(optionRequest)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/v1/filter-option/1", bytes.NewBuffer(reqJson))
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		bodyByte, err := ioutil.ReadAll(response.Body)

		t.Logf("%+v", string(bodyByte))

		assert.Equal(t, 404, response.StatusCode)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Bad Request", func(t *testing.T) {
		optionRequest := request2.UpdateOptionRequest{}

		reqJson, _ := json.Marshal(optionRequest)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/v1/filter-option/1", bytes.NewBuffer(reqJson))
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		bodyByte, err := ioutil.ReadAll(response.Body)

		t.Logf("%+v", string(bodyByte))

		assert.Equal(t, 400, response.StatusCode)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteProject(t *testing.T) {
	sqlDB, db, mock := test.DbMock(t)
	defer sqlDB.Close()

	router := setupInitialFilterOption(db)

	f1 := project2.FilterOptionDummyList[0]
	fRow1 := project2.MappingFilterOptionStore(f1, 1)
	args := test.AnyThingArgs(2)

	t.Run("Success Test", func(t *testing.T) {
		option := sqlmock.NewRows(project2.FilterOptionCols).AddRow(fRow1...)

		mock.ExpectQuery(test.SelectOneByIdFilterOptionSQL).WithArgs(1).WillReturnRows(option)
		mock.ExpectBegin()
		mock.
			ExpectExec(test.UpdateFilterOptionSQL).
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/v1/filter-option/1", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		bodyByte, err := ioutil.ReadAll(response.Body)

		t.Logf("%+v", string(bodyByte))

		assert.Equal(t, 200, response.StatusCode)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		option := sqlmock.NewRows(project2.FilterOptionCols)

		mock.ExpectQuery(test.SelectOneByIdFilterOptionSQL).WithArgs(1).WillReturnRows(option)

		request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/v1/filter-option/1", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		bodyByte, err := ioutil.ReadAll(response.Body)

		t.Logf("%+v", string(bodyByte))

		assert.Equal(t, 404, response.StatusCode)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
