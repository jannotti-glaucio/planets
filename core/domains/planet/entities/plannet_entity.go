package entities

import (
	"time"

	"github.com/google/uuid"
)

//Planet ...
type Planet struct {
	UUID      string     `json:"uuid" bson:"_id"`
	Name      string     `json:"name" bson:"name"`
	Terrain   string     `json:"terrain" bson:"terrain"`
	Climate   string     `json:"climate" bson:"climate"`
	Films     int        `json:"films" bson:"films"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

//Planets ...
type Planets []Planet

//PlanetPublic ...
type PlanetPublic struct {
	UUID    string `json:"uuid" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Terrain string `json:"terrain" bson:"terrain"`
	Climate string `json:"climate" bson:"climate"`
	Films   int    `json:"films" bson:"films"`
}

//New return a object
func (planet *Planet) New() *Planet {
	planet.UUID = uuid.New().String()
	planet.CreatedAt = time.Now()
	return planet
}

//PublicPlanet return a Planet without non important attributes
func (planet *Planet) PublicPlanet() interface{} {
	return PlanetPublic{
		UUID:    planet.UUID,
		Name:    planet.Name,
		Terrain: planet.Terrain,
		Climate: planet.Climate,
		Films:   planet.Films,
	}
}

//PublicPlanets return a list of planets without important attributes
func (planets *Planets) PublicPlanets() []interface{} {
	result := make([]interface{}, len(*planets))
	for index, user := range *planets {
		result[index] = user.PublicPlanet()
	}
	return result
}
