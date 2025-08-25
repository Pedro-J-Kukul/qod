// Filename: cmd/api/handlers.go
// Description: HTTP request handlers for the API

package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

type civ6Quote struct {
	Technology string
	Quote      string
	Source     string
}

func getQoute() string {
	civ6Quotes := []civ6Quote{
		{
			Technology: "Pottery",
			Quote:      "No man ever wetted clay and then left it, as if there would be bricks by chance and fortune.",
			Source:     "Plutarch",
		},
		{
			Technology: "Animal Husbandry",
			Quote:      "If there are no dogs in Heaven, then when I die I want to go where they went.",
			Source:     "Will Rogers",
		},
		{
			Technology: "Mining",
			Quote:      "When you find yourself in a hole, quit digging.",
			Source:     "Will Rogers",
		},
		{
			Technology: "Sailing",
			Quote:      "Vessels large may venture more, but little boats should keep near shore.",
			Source:     "Benjamin Franklin",
		},
		{
			Technology: "Astrology",
			Quote:      "I don’t believe in astrology; I’m a Sagittarius and we’re skeptical.",
			Source:     "Arthur C. Clarke",
		},
		{
			Technology: "Irrigation",
			Quote:      "Thousands have lived without love, not one without water.",
			Source:     "W. H. Auden",
		},
		{
			Technology: "Archery",
			Quote:      "I shot an arrow into the air. It fell to earth, I knew not where.",
			Source:     "Henry Wadsworth Longfellow",
		},
		{
			Technology: "Writing",
			Quote:      "Writing means sharing. It’s part of the human condition to want to share things – thoughts, ideas, opinions.",
			Source:     "Paulo Coelho",
		},
		{
			Technology: "Masonry",
			Quote:      "Each of us is carving a stone, erecting a column, or cutting a piece of stained glass in the construction of something much bigger than ourselves.",
			Source:     "Adrienne Clarkson",
		},
		{
			Technology: "Bronze Working",
			Quote:      "Bronze is the mirror of the form, wine of the mind.",
			Source:     "Aeschylus",
		},
		{
			Technology: "The Wheel",
			Quote:      "Sometimes the wheel turns slowly, but it turns.",
			Source:     "Lorne Michaels",
		},
		{
			Technology: "Celestial Navigation",
			Quote:      "And all I ask is a tall ship and a star to steer her by.",
			Source:     "John Masefield",
		},
		{
			Technology: "Currency",
			Quote:      "Wealth consists not in having great possessions, but in having few wants.",
			Source:     "Epictetus",
		},
		{
			Technology: "Horseback Riding",
			Quote:      "No hour of life is wasted that is spent in the saddle.",
			Source:     "Winston S. Churchill",
		},
		{
			Technology: "Iron Working",
			Quote:      "The Lord made us all out of iron. Then he turns up the heat to forge some of us into steel.",
			Source:     "Marie Osmond",
		},
		{
			Technology: "Shipbuilding",
			Quote:      "I cannot imagine any condition which would cause a ship to founder … Modern shipbuilding has gone beyond that.",
			Source:     "Capt. E.J. Smith, RMS Titanic",
		},
		{
			Technology: "Mathematics",
			Quote:      "Without mathematics, there’s nothing you can do. Everything around you is mathematics. Everything around you is numbers.",
			Source:     "Shakuntala Devi",
		},
		{
			Technology: "Construction",
			Quote:      "Create with the heart; build with the mind.",
			Source:     "Criss Jami",
		},
		{
			Technology: "Engineering",
			Quote:      "One man’s ‘magic’ is another man’s engineering.",
			Source:     "Robert Heinlein",
		},
		{
			Technology: "Military Tactics",
			Quote:      "Tactics mean doing what you can with what you have.",
			Source:     "Saul Alinsky",
		},
	}

	// Get a random quote from the slice
	randomIndex := rand.Intn(len(civ6Quotes))
	randomQuote := civ6Quotes[randomIndex]

	quote := fmt.Sprintf(`{ "Technology": %q, "Quote": %q, "Source": %q }`, randomQuote.Technology, randomQuote.Quote, randomQuote.Source)

	return quote
}

func (app *application) quoteHandler(w http.ResponseWriter, r *http.Request) {
	quote := getQoute()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(quote))
}

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	js := `{"status": "available", 
	"environment": %q, 
	"version": %q}`
	js = fmt.Sprintf(js, app.config.env, app.config.version)

	// Content-Type is text/plain by default

	w.Header().Set("Content-Type", "application/json")
	// Write the JSON as the HTTP response body.
	w.Write([]byte(js))

}
