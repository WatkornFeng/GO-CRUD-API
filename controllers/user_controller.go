package controllers

import (
	"errors"
	"project_restfulApi_go/models"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ==>gorm.Model’s ID field is uint, which is typically
// either 32-bit or 64-bit depending ,on the platform (on 64-bit systems it’s 64-bit)
// ==>The database side stores it as BIGINT (64-bit),
// but your Go struct field is uint
// ==> You should parse the incoming id as uint64
// (*to cover all possible BIGINT values in the DB), then convert to uint for Go usage:
// Let GORM handle the DB translation between Go uint and SQL BIGINT
// When GORM talks to the database, it maps Go’s uint to the database’s BIGINT
func GetUserHandler(c *fiber.Ctx, db *gorm.DB) error {
	idParam := c.Params("id")
	// uint is work only on positve numbers otherwise ParseUint will show error
	// 64bit value also suit BIGINT in DB
	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param format",
		})
	}
	// uint is either 32 /64 depend on system
	// if pass 64bit value into 32bit system data might loss ,but modern system usually 64bit
	id := uint(idUint64)
	user, err := models.GetUser(db, id)
	if err != nil {
		//gorm.ErrRecordNotFound  is only returned by GORM methods that perform 'queries'
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}
	return c.JSON(user)
}
func GetUsersHandler(c *fiber.Ctx, db *gorm.DB) error {
	users, err := models.GetAllUsers(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	return c.JSON(users)
}

var validate = validator.New()

type CreateUserInput struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

func CreateUserHandler(c *fiber.Ctx, db *gorm.DB) error {
	userInput := new(CreateUserInput)

	if err := c.BodyParser(userInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	if err := validate.Struct(userInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}
	// 🔁 Convert to models.User
	user := &models.User{
		Name:  userInput.Name,
		Email: userInput.Email,
	}

	if err := models.CreateUser(db, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create user successful",
		"data":    user,
	})
}

type UpdateUserInput struct {
	Name string `json:"name" validate:"required,min=3"`
}

func UpdateUserHandler(c *fiber.Ctx, db *gorm.DB) error {
	idParam := c.Params("id")
	// validate Param
	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param format",
		})
	}
	id := uint(idUint64)

	userInput := new(UpdateUserInput)
	if err := c.BodyParser(userInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	if err := validate.Struct(userInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	updatedUser, err := models.UpdateUserName(db, id, userInput.Name)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Update user successful",
		"data":    updatedUser,
	})
}

func DeleteUserHandler(c *fiber.Ctx, db *gorm.DB) error {
	idParam := c.Params("id")

	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param format",
		})
	}

	id := uint(idUint64)

	err = models.DeleteUser(db, id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
