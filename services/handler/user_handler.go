package handler

import (
	"fmt"
	"net/http"
	db "server/internal/database"
	"server/models"
	"server/services/util/psql"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func CreateUser(c *gin.Context) {
	var user models.NewUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}

	conn, err := db.New()
	defer db.Close(conn)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": "server issue : database connection failed !!"})
		return
	}

	_, err = psql.GetUserByEmail(user.Email, conn)

	if err == nil {

		c.JSON(http.StatusConflict, gin.H{"errorMessage": "User  with this email already exists"})

		return
	} else if err.Error() == pgx.ErrNoRows.Error() { // user does not exist

		err = psql.InsertUserIntoDb(&user, conn)

		if err != nil {

			fmt.Println("error:InsertUserIntoDb():", err)
			c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": "Server error"})
			return
		}

	} else if err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok {

			switch pgErr.Code {

			case "42P01":
				err := psql.CreateUserTable(conn)
				if err != nil {
					fmt.Println("error:CreateUserTable():", err)
					c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": "Server error"})
					return
				}
				err = psql.InsertUserIntoDb(&user, conn)

				if err != nil {

					fmt.Println("error:InsertUserIntoDb():", err)
					c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": "Server"})
					return
				}

			case "23505": // Unique violation
				fmt.Println("Error: Unique constraint violation.")
				c.JSON(http.StatusConflict, gin.H{"errorMessage": "Unique constraint violation"})
				return

			}

		} else {
			fmt.Println("errr:GetUserByEmail():", err)
			c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": "server  error"})
			return
		}
	}
	fmt.Printf("User created successfully: first_name: %s, last_name: %s, email: %s\n", user.FirstName, user.LastName, user.Email)
	c.JSON(http.StatusCreated, gin.H{"sucessMessage": "User  created successfully", "first_name": user.FirstName, "last_name": user.LastName, "email": user.Email})
}

// fetch user handler
func GetUser(c *gin.Context) {
	conn, err := db.New()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": "Server error"})
		return
	}
	defer db.Close(conn)
	idStr := c.Param("id")
	if len(idStr) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "Invalid ID"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errorMessage": "Un authorized user : Invalid id type"})
		return
	}

	
	user, err := psql.GetUserById(uint(id), conn)
	if err != nil {
		if err.Error() == fmt.Sprintf("no user found with id: %d", id) {
			c.JSON(http.StatusNotFound, gin.H{"errorMessage": "User  not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": "Internal server error"})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}
