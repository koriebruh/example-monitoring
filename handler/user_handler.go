package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"wyw/entity"
	"wyw/metric"
)

type UserHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	GetUser(c *gin.Context)
}

type UserHandlerImpl struct {
	*gorm.DB
	*metric.AppMetricsExporter
}

func NewUserHandler(DB *gorm.DB, appMetricsExporter *metric.AppMetricsExporter) *UserHandlerImpl {
	return &UserHandlerImpl{DB: DB, AppMetricsExporter: appMetricsExporter}
}

// Login handles user login requests.
// @Summary      User Login
// @Description  Validates user credentials and logs the user in if the credentials are correct.
// @Param        credentials body entity.User true "User credentials (username and password)"
// @Produce      application/json
// @Tags         user
// @Success      200 {object} entity.MsgResponse "Success message indicating user login"
// @Failure      400 {object} entity.ErrorResponse "Error message indicating invalid credentials"
// @Failure      422 {object} entity.ErrorResponse "Error message indicating invalid JSON format"
// @Router       /login [post]
func (u UserHandlerImpl) Login(c *gin.Context) {
	var request entity.User
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "invalid format json ",
		})
		return
	}

	var user entity.User
	if err := u.DB.WithContext(c.Request.Context()).
		Where("username =? AND password =?", request.Username, request.Password).
		First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Username and Password",
		})
		return
	}

	//RECORD METRICS
	u.AppMetricsExporter.RecordBusinessEvent("login", user.Username)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s login successfully", user.Username)})
}

// Register handles user registration requests.
// @Summary      User Registration
// @Description  Registers a new user by saving the provided user credentials to the database.
// @Param        credentials body entity.User true "User credentials (username, password, etc.)"
// @Produce      application/json
// @Tags         user
// @Success      200 {object} entity.MsgResponse "Success message indicating successful registration"
// @Failure      400 {object} entity.ErrorResponse "Error message indicating invalid username or password"
// @Failure      422 {object} entity.ErrorResponse "Error message indicating invalid JSON format"
// @Router       /register [post]
func (u UserHandlerImpl) Register(c *gin.Context) {
	var request entity.User
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "invalid format json ",
		})
		return
	}

	if err := u.DB.WithContext(c.Request.Context()).Create(&request).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Username and Password",
		})
		return
	}

	//RECORD METRICS
	u.AppMetricsExporter.RecordBusinessEvent("register", request.Username)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s register successfully", request.Username)})
}

// GetUser handles fetching all users from the database.
// @Summary      Get Users
// @Description  Retrieves all users from the database.
// @Produce      application/json
// @Tags         user
// @Success      200 {object} []entity.User{} "List of users in the database"
// @Failure      500 {object} entity.ErrorResponse "Error message indicating internal server error"
// @Failure      422 {object} entity.ErrorResponse "Error message indicating no users found"
// @Router       /users [get]
func (u UserHandlerImpl) GetUser(c *gin.Context) {

	var results []entity.User
	if err := u.DB.WithContext(c.Request.Context()).Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error try again later",
		})
		return
	}

	if len(results) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "data not insert yet ",
		})
		return
	}

	u.AppMetricsExporter.RecordBusinessEvent("get_users", "system")
	c.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}

//func loginHandler(c *gin.Context, exporter *AppMetricsExporter) {
//	userID := c.DefaultQuery("user_id", "anonymous")
//
//	// Rekam event login
//	exporter.RecordBusinessEvent("login", userID)
//
//	c.JSON(http.StatusOK, gin.H{
//		"status":  "login successfully",
//		"user_id": userID,
//	})
//}
//
//func registerHandler(c *gin.Context, exporter *AppMetricsExporter) {
//	userID := c.DefaultQuery("user_id", "new_user")
//
//	// Rekam event register
//	exporter.RecordBusinessEvent("register", userID)
//
//	c.JSON(http.StatusOK, gin.H{
//		"status": "register successfully",
//	})
//}
