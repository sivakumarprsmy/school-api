package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/periasamy/school-api/internal/handlers"
	"github.com/periasamy/school-api/internal/models"
)

// newTestDB creates a GORM DB backed by sqlmock for unit testing.
func newTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { sqlDB.Close() })

	dialector := postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	require.NoError(t, err)

	return gormDB, mock
}

// newTestRouter wires up a Gin router with student routes for testing.
func newTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handlers.NewStudentHandler(db, zerolog.Nop())
	v1 := r.Group("/api/v1")
	s := v1.Group("/students")
	s.POST("", h.CreateStudent)
	s.GET("", h.GetAllStudents)
	s.GET("/:id", h.GetStudentByID)
	s.DELETE("/:id", h.DeleteStudent)
	return r
}

// studentColumns returns the column list used in SELECT mock rows.
var studentColumns = []string{"id", "name", "email", "age", "grade", "created_at", "updated_at"}

// ---- CreateStudent -------------------------------------------------------

func TestCreateStudent_Success(t *testing.T) {
	db, mock := newTestDB(t)
	router := newTestRouter(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "students"`).
		WithArgs(
			sqlmock.AnyArg(), // name
			sqlmock.AnyArg(), // email
			sqlmock.AnyArg(), // age
			sqlmock.AnyArg(), // grade
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	body, _ := json.Marshal(models.CreateStudentRequest{
		Name: "Alice", Email: "alice@school.com", Age: 15, Grade: "10th",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var got models.Student
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, "Alice", got.Name)
	assert.Equal(t, "alice@school.com", got.Email)
	assert.Equal(t, 15, got.Age)
	assert.Equal(t, "10th", got.Grade)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStudent_MissingRequiredFields(t *testing.T) {
	db, _ := newTestDB(t)
	router := newTestRouter(db)

	// only name provided — email, age, grade missing
	body, _ := json.Marshal(map[string]string{"name": "Alice"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateStudent_InvalidJSON(t *testing.T) {
	db, _ := newTestDB(t)
	router := newTestRouter(db)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateStudent_InvalidEmail(t *testing.T) {
	db, _ := newTestDB(t)
	router := newTestRouter(db)

	body, _ := json.Marshal(models.CreateStudentRequest{
		Name: "Alice", Email: "not-an-email", Age: 15, Grade: "10th",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateStudent_DuplicateEmail(t *testing.T) {
	db, mock := newTestDB(t)
	router := newTestRouter(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "students"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("ERROR: duplicate key value violates unique constraint"))
	mock.ExpectRollback()

	body, _ := json.Marshal(models.CreateStudentRequest{
		Name: "Alice", Email: "alice@school.com", Age: 15, Grade: "10th",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---- GetStudentByID ------------------------------------------------------

func TestGetStudentByID_Success(t *testing.T) {
	db, mock := newTestDB(t)
	router := newTestRouter(db)

	now := time.Now()
	mock.ExpectQuery(`SELECT \* FROM "students"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()). // id, LIMIT 1
		WillReturnRows(
			sqlmock.NewRows(studentColumns).
				AddRow(1, "Alice", "alice@school.com", 15, "10th", now, now),
		)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got models.Student
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, uint(1), got.ID)
	assert.Equal(t, "Alice", got.Name)
	assert.Equal(t, "alice@school.com", got.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStudentByID_NotFound(t *testing.T) {
	db, mock := newTestDB(t)
	router := newTestRouter(db)

	mock.ExpectQuery(`SELECT \* FROM "students"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows(studentColumns)) // empty

	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStudentByID_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	router := newTestRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---- GetAllStudents ------------------------------------------------------

func TestGetAllStudents_ReturnsList(t *testing.T) {
	db, mock := newTestDB(t)
	router := newTestRouter(db)

	now := time.Now()
	mock.ExpectQuery(`SELECT \* FROM "students"`).
		WillReturnRows(
			sqlmock.NewRows(studentColumns).
				AddRow(1, "Alice", "alice@school.com", 15, "10th", now, now).
				AddRow(2, "Bob", "bob@school.com", 16, "11th", now, now),
		)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/students", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []models.Student
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Len(t, got, 2)
	assert.Equal(t, "Alice", got[0].Name)
	assert.Equal(t, "Bob", got[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllStudents_Empty(t *testing.T) {
	db, mock := newTestDB(t)
	router := newTestRouter(db)

	mock.ExpectQuery(`SELECT \* FROM "students"`).
		WillReturnRows(sqlmock.NewRows(studentColumns)) // no rows

	req := httptest.NewRequest(http.MethodGet, "/api/v1/students", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []models.Student
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Empty(t, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllStudents_DBError(t *testing.T) {
	db, mock := newTestDB(t)
	router := newTestRouter(db)

	mock.ExpectQuery(`SELECT \* FROM "students"`).
		WillReturnError(errors.New("connection lost"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/students", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---- DeleteStudent -------------------------------------------------------

func TestDeleteStudent_Success(t *testing.T) {
	db, mock := newTestDB(t)
	router := newTestRouter(db)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "students"`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	mock.ExpectCommit()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/students/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteStudent_NotFound(t *testing.T) {
	db, mock := newTestDB(t)
	router := newTestRouter(db)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "students"`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
	mock.ExpectCommit()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/students/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteStudent_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	router := newTestRouter(db)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/students/xyz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteStudent_DBError(t *testing.T) {
	db, mock := newTestDB(t)
	router := newTestRouter(db)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "students"`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New("connection lost"))
	mock.ExpectRollback()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/students/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
