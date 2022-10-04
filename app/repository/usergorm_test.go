package repository

import (
	"database/sql"
	"errors"
	"exampleclean.com/refactor/app/domain"
	rest_structs "exampleclean.com/refactor/app/rest-structs"
	"exampleclean.com/refactor/app/utils"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
)

var columns = []string{"id", "email", "password", "firstname", "lastname"}

type SuiteGorm struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository userDatabaseGorm
	user       *domain.Users
}

func (s *SuiteGorm) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("mysql", db)
	require.NoError(s.T(), err)

	s.DB.LogMode(true)

	s.repository = userDatabaseGorm{DB: s.DB}
}

func (s *SuiteGorm) TestRepository_deleteUserFailed() {
	db, mockTemp, _ := sqlmock.New()
	database, _ := gorm.Open("mysql", db)
	database.LogMode(true)
	repo := userDatabaseGorm{DB: database}

	user := domain.Users{
		Id:        uint(8),
		Email:     "mameddram@gmail.com",
		Password:  "1234",
		Firstname: "mamed",
		Lastname:  "dram",
	}

	mockTemp.ExpectBegin()
	mockTemp.ExpectExec(regexp.QuoteMeta("DELETE FROM `users`")).
		WithArgs().
		WillReturnError(errors.New("cannot delete user"))
	mockTemp.ExpectCommit()
	err := repo.Delete(user)
	s.Error(err)
}

func (s *SuiteGorm) TestRepository_deleteUserSuccess() {
	db, mockTemp, _ := sqlmock.New()
	database, _ := gorm.Open("mysql", db)
	database.LogMode(true)
	repo := userDatabaseGorm{DB: database}

	user := domain.Users{
		Id:        uint(8),
		Email:     "mameddram@gmail.com",
		Password:  "1234",
		Firstname: "mamed",
		Lastname:  "dram",
	}

	mockTemp.ExpectBegin()
	mockTemp.ExpectExec(regexp.QuoteMeta("DELETE FROM `users`")).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockTemp.ExpectCommit()
	err := repo.Delete(user)
	s.NoError(err)
}

func (s *SuiteGorm) TestRepository_SaveUserSuccess() {
	db, mockTemp, _ := sqlmock.New()
	database, _ := gorm.Open("mysql", db)
	database.LogMode(true)
	repo := userDatabaseGorm{DB: database}

	user := rest_structs.RequestSignup{
		Email:     "mameddram@gmail.com",
		Password:  "1234",
		Firstname: "mamed",
		Lastname:  "dram",
	}

	mockTemp.ExpectBegin()
	mockTemp.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`email`,`password`,`firstname`,`lastname`) VALUES (?,?,?,?)")).
		WithArgs(user.Email, utils.HashThisSHA1(user.Password), user.Firstname, user.Lastname).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockTemp.ExpectCommit()
	err := repo.Save(user)

	s.NoError(err)
}

func (s *SuiteGorm) TestRepository_SaveUserFailed() {
	user := rest_structs.RequestSignup{
		Email:     "mameddram@gmail.com",
		Password:  "1234",
		Firstname: "mamed",
		Lastname:  "dram",
	}

	s.mock.ExpectBegin() // start transaction
	s.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `users` (`email`,`password`,`firstname`,`lastname`) VALUES (?,?,?,?)")).
		WithArgs(user.Email, utils.HashThisSHA1(user.Password), user.Firstname, user.Lastname).
		WillReturnError(errors.New("cannot create new user"))

	err := s.repository.Save(user)
	s.mock.ExpectCommit() // commit transaction

	s.Error(err)
	s.NotNil(err)
}

func (s *SuiteGorm) TestRepository_UpdatePasswordSuccess() {

	db, mockTemp, _ := sqlmock.New()
	database, err := gorm.Open("mysql", db)
	database.LogMode(true)
	repo := userDatabaseGorm{DB: database}
	//require.NoError(s.T(), err)

	user := domain.Users{
		Id:        11,
		Email:     "sekolahmu1@gmail.com",
		Password:  "12345",
		Firstname: "sekolah1",
		Lastname:  "mu1",
	}

	mockTemp.ExpectBegin()
	mockTemp.ExpectExec(regexp.QuoteMeta(
		"UPDATE `users`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockTemp.ExpectCommit()
	mockTemp.ExpectationsWereMet()

	res, err := repo.UpdatePassword(user)
	s.NotNil(res)
	s.Nil(err)
	s.Equal(res, int64(1))
}

func (s *SuiteGorm) TestRepository_UpdatePasswordFailed() {
	user := domain.Users{
		Id:        11,
		Email:     "sekolahmu1@gmail.com",
		Password:  "12345",
		Firstname: "sekolah1",
		Lastname:  "mu1",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE `users`")).
		WillReturnError(errors.New("Cannot update user"))
	s.mock.ExpectCommit()
	s.mock.ExpectClose()

	res, err := s.repository.UpdatePassword(user)
	s.NotNil(err)
	s.Equal(res, int64(0))
}

func (s *SuiteGorm) TestRepository_FindByEmailFailed() {

	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE (Email = ?)")).
		WillReturnError(errors.New("user not found"))

	res, err := s.repository.FindByEmail("mamed@gmail.com")
	s.NotNil(err)
	s.Nil(res)
}

func (s *SuiteGorm) TestRepository_FindByEmailSuccess() {
	user := domain.Users{
		Id:        11,
		Email:     "sekolahmu1@gmail.com",
		Password:  "12345",
		Firstname: "sekolah1",
		Lastname:  "mu1",
	}

	rows := sqlmock.NewRows(columns).
		AddRow(user.Id, user.Email, user.Password, user.Firstname, user.Lastname)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE (Email = ?)")).
		WithArgs(user.Email).
		WillReturnRows(rows)

	res, err := s.repository.FindByEmail(user.Email)
	s.NotNil(res)
	s.Nil(err)
}

func (s *SuiteGorm) TestRepository_FindAllFailed() {

	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users`")).
		WillReturnError(errors.New("no user is found"))

	res, err := s.repository.FindAll()
	s.NotNil(err)
	s.Nil(res)
}

func (s *SuiteGorm) TestRepository_FindAllSuccess() {

	users := []domain.Users{
		{
			Id:        11,
			Email:     "sekolahmu1@gmail.com",
			Password:  "12345",
			Firstname: "sekolah1",
			Lastname:  "mu1",
		},
		{
			Id:        12,
			Email:     "sekolahmu2@gmail.com",
			Password:  "12345",
			Firstname: "sekolah2",
			Lastname:  "mu2",
		},
	}

	rows := sqlmock.NewRows(columns).
		AddRow(users[0].Id, users[0].Email, users[0].Password, users[0].Firstname, users[0].Lastname).
		AddRow(users[1].Id, users[1].Email, users[1].Password, users[1].Lastname, users[1].Lastname)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users`")).
		WillReturnRows(rows)

	res, err := s.repository.FindAll()
	s.NotNil(res)
	s.Nil(err)
	s.Equal(cap(res), cap(users))

	for index, temp := range res {
		s.Equal(temp.Id, users[index].Id)
		s.Equal(temp.Email, users[index].Email)
		s.Equal(temp.Password, users[index].Password)

	}
}

func (s *SuiteGorm) TestRepository_FindByIdFailed() {
	id := uint(8)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE (Id = ?)")).
		WithArgs(8).
		WillReturnError(errors.New("user not found"))

	res := s.repository.FindByID(id)
	s.Nil(res)
}

func (s *SuiteGorm) TestRepository_FindByIdSuccess() {
	var (
		id = uint(8)
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE (Id = ?)")).
		WithArgs(8).
		WillReturnRows(s.mock.NewRows(columns).AddRow(id, "sekolahmu@gmail.com", "1234", "sekolah", "mu"))

	res := s.repository.FindByID(id)

	s.NotNil(res)
	s.Equal(res.Id, id)
}

func TestSuiteRepositoryGorm(t *testing.T) {
	suite.Run(t, new(SuiteGorm))
}
