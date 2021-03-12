package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

const version = "1.0"

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("calculator available at /calculator")
	})

	app.Get("/calculator", func(c *fiber.Ctx) error {
		h := c.Query("h", "")
		w := c.Query("w", "")
		format := c.Query("format", "string")

		if h == "" || w == "" {
			return c.SendString("missing weight(w) or height(h)")
		}

		height, hErr := strconv.Atoi(h)
		if hErr != nil {
			return c.SendString("height(h) must be integer value in CM")
		}

		weight, wErr := strconv.Atoi(w)
		if wErr != nil {
			return c.SendString("weight(w) must be integer value in KG")
		}

		bmi := calculator(height, weight)

		if format == "json" {
			return c.JSON(fiber.Map{
				"bmi":     fmt.Sprintf("%.2f", bmi),
				"status":  getBMIStatus(bmi),
				"version": version,
			})
		} else {
			return c.SendString(
				fmt.Sprintf("BMI: %.2f, status: %s, version: %s", bmi, getBMIStatus(bmi), version),
			)
		}
	})

	log.Println(app.Listen(":3000"))
}

//  BMI = kg/m2
func calculator(height, weight int) float64 {
	heightInMetre := float64(height) / 100.0
	return float64(weight) / (heightInMetre * heightInMetre)
}

// BMI status
// reference https://images.theconversation.com/files/349366/original/file-20200724-25-osy3a3.PNG?ixlib=rb-1.1.0&q=30&auto=format&w=600&h=304&fit=crop&dpr=2
func getBMIStatus(bmi float64) string {
	if bmi < 18.5 {
		return "Underweight"
	} else if bmi > 18.5 && bmi < 24.9 {
		return "Normal Weight"
	} else if bmi > 25.0 && bmi < 29.9 {
		return "Overweight"
	} else if bmi > 30.0 && bmi < 34.9 {
		return "Obesity class I"
	} else if bmi > 35.0 && bmi < 39.9 {
		return "Obesity class II"
	} else {
		return "Obesity class III"
	}
}
