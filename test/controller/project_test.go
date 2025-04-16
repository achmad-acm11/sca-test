package controller

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
	"sca-integrator/app/dbo/cli"
	"sca-integrator/app/dbo/repository"
	"sca-integrator/app/dbo/repository/project"
	request2 "sca-integrator/app/dto/request"
	"sca-integrator/app/middleware"
	"sca-integrator/app/service"
	"sca-integrator/test"
	"sca-integrator/test/dummy"
	"testing"
	"time"
)

func setupInitialProject(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Router init
	route := gin.Default()

	route.Use(middleware.ErrorHandler())
	router := route.Group("api/v1")

	projectRepository := repository.NewProjectRepository()
	exclusionRepository := project.NewExclusionRepository()
	optionRepository := project.NewFilterOptionRepository()
	resultRepository := repository.NewResultRepository()
	trivyCli := cli.NewTrivyCli()

	projectService := service.NewProjectService(projectRepository,
		optionRepository,
		exclusionRepository,
		resultRepository,
		trivyCli,
		validator.New(),
		db)

	projectController := controller.NewProjectController(projectService)

	router.GET("projects", projectController.GetAllHandler)
	router.GET("project/:id", projectController.GetDetailByIdHandler)
	router.POST("project", projectController.CreateProjectHandler)
	router.POST("project/scan", projectController.ScanProjectHandler)
	router.PUT("project/:id", projectController.UpdateProjectHandler)
	router.DELETE("project/:id", projectController.DeleteProjectHandler)

	return route
}

func TestGetAllProject(t *testing.T) {
	sqlDB, db, mock := test.DbMock(t)
	defer sqlDB.Close()

	router := setupInitialProject(db)

	p1 := dummy.ProjectDummyList[0]
	pRow1 := dummy.MappingProjectStore(p1, 1)
	t.Run("Success Test", func(t *testing.T) {
		projects := sqlmock.NewRows(dummy.ProjectCols).AddRow(pRow1...)

		mock.ExpectQuery(regexp.QuoteMeta(test.SelectAllProjectSQL)).WillReturnRows(projects)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/projects", nil)
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

func TestGetDetailByIdProject(t *testing.T) {
	sqlDB, db, mock := test.DbMock(t)
	defer sqlDB.Close()

	router := setupInitialProject(db)

	p1 := dummy.ProjectDummyList[0]
	pRow1 := dummy.MappingProjectStore(p1, 1)
	t.Run("Success Test", func(t *testing.T) {
		project := sqlmock.NewRows(dummy.ProjectCols).AddRow(pRow1...)

		mock.ExpectQuery(test.SelectOneByIdProjectSQL).WithArgs(1).WillReturnRows(project)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/project/1", nil)
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
		project := sqlmock.NewRows(dummy.ProjectCols)

		mock.ExpectQuery(test.SelectOneByIdProjectSQL).WithArgs(1).WillReturnRows(project)

		request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/project/1", nil)
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

func TestCreateProject(t *testing.T) {
	sqlDB, db, mock := test.DbMock(t)
	defer sqlDB.Close()

	router := setupInitialProject(db)

	p1 := dummy.ProjectDummyList[0]
	args := test.AnyThingArgs(9)

	t.Run("Success Test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(test.InsertProjectSQL).
			WithArgs(args...).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at"}).
				AddRow(1, time.Now(), time.Now(), time.Now()))
		mock.ExpectCommit()

		projectRequest := request2.CreateProjectRequest{
			Name:        p1.Name,
			Description: p1.Description,
			Repo_type:   p1.RepoType,
			Url:         p1.Url,
			Branch_name: p1.BranchName,
			Visibility:  p1.Visibility,
		}

		reqJson, _ := json.Marshal(projectRequest)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/project", bytes.NewBuffer(reqJson))
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
		projectRequest := request2.CreateProjectRequest{}

		reqJson, _ := json.Marshal(projectRequest)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/project", bytes.NewBuffer(reqJson))
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

func TestUpdateProject(t *testing.T) {
	sqlDB, db, mock := test.DbMock(t)
	defer sqlDB.Close()

	router := setupInitialProject(db)

	p1 := dummy.ProjectDummyList[0]
	pRow1 := dummy.MappingProjectStore(p1, 1)
	args := test.AnyThingArgs(10)

	t.Run("Success Test", func(t *testing.T) {
		project := sqlmock.NewRows(dummy.ProjectCols).AddRow(pRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(test.SelectOneByIdProjectSQL).WithArgs(1).WillReturnRows(project)
		mock.
			ExpectExec(test.UpdateProjectSQL).
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		projectRequest := request2.UpdateProjectRequest{
			Name:        p1.Name,
			Description: p1.Description,
			Repo_type:   p1.RepoType,
			Url:         p1.Url,
			Branch_name: p1.BranchName,
			Visibility:  p1.Visibility,
		}

		reqJson, _ := json.Marshal(projectRequest)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/v1/project/1", bytes.NewBuffer(reqJson))
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
		project := sqlmock.NewRows(dummy.ProjectCols)

		mock.ExpectBegin()
		mock.ExpectQuery(test.SelectOneByIdProjectSQL).WithArgs(1).WillReturnRows(project)
		mock.ExpectRollback()

		projectRequest := request2.UpdateProjectRequest{
			Name:        p1.Name,
			Description: p1.Description,
			Repo_type:   p1.RepoType,
			Url:         p1.Url,
			Branch_name: p1.BranchName,
			Visibility:  p1.Visibility,
		}

		reqJson, _ := json.Marshal(projectRequest)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/v1/project/1", bytes.NewBuffer(reqJson))
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
		projectRequest := request2.UpdateProjectRequest{}

		reqJson, _ := json.Marshal(projectRequest)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/v1/project/1", bytes.NewBuffer(reqJson))
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

	router := setupInitialProject(db)

	p1 := dummy.ProjectDummyList[0]
	pRow1 := dummy.MappingProjectStore(p1, 1)
	args := test.AnyThingArgs(2)

	t.Run("Success Test", func(t *testing.T) {
		project := sqlmock.NewRows(dummy.ProjectCols).AddRow(pRow1...)

		mock.ExpectQuery(test.SelectOneByIdProjectSQL).WithArgs(1).WillReturnRows(project)
		mock.ExpectBegin()
		mock.
			ExpectExec(test.UpdateProjectSQL).
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/v1/project/1", nil)
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
		project := sqlmock.NewRows(dummy.ProjectCols)

		mock.ExpectQuery(test.SelectOneByIdProjectSQL).WithArgs(1).WillReturnRows(project)

		request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/v1/project/1", nil)
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
