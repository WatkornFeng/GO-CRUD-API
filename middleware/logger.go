package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger(verbose bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Continue to next middleware/handler
		err := c.Next()

		duration := time.Since(start)

		// Basic log
		fmt.Printf("[%s] %s %s %d\n",
			time.Now().Format(time.RFC3339),
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
		)

		// Verbose log
		if verbose {
			fmt.Printf(" - Duration: %s\n", duration)
			fmt.Printf("Request Headers: %v\n", c.GetReqHeaders())
			fmt.Printf("Request Body: %s\n", c.Body())
			fmt.Println()
		}

		return err
	}
}
