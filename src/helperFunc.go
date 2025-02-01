package server

import (
	"strconv"

	"github.com/gofrs/uuid"
)

func convertime(s int64) string {
	var convert int
	if s >= 31557600 {
		for s > 0 {
			s = s / 31557600
			convert++
		}
		duree := strconv.Itoa(convert)
		if duree == "1" {
			return duree + " Year"

		} else {
			return duree + " Years"
		}

	} else if s >= 2629744 {
		for s > 0 {
			s = s / 2629744
			convert++
		}
		duree := strconv.Itoa(convert)
		if duree == "1" {
			return duree + " Month"

		} else {
			return duree + " Months"
		}

	} else if s >= 86400 {
		for s > 0 {
			s = s / 86400
			convert++
		}
		duree := strconv.Itoa(convert)
		if duree == "1" {
			return duree + " Day"

		} else {
			return duree + " Days"
		}

	} else if s >= 3600 {
		for s > 0 {
			s = s / 3600
			convert++
		}
		duree := strconv.Itoa(convert)
		if duree == "1" {
			return duree + " Hour"

		} else {
			return duree + " Hours"
		}

	} else if s >= 60 {
		for s > 0 {
			s = s / 60
			convert++
		}
		duree := strconv.Itoa(convert)
		if duree == "1" {
			return duree + " Minute"

		} else {
			return duree + " Minutes"
		}

	} else if s < 60 {
		if s == 0 {
			return "Now"
		}
		return strconv.Itoa(int(s)) + " s"

	}
	return ""
}

func GenerateUUID() (string, error) {
	// Génère un nouvel UUID
	newUUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return newUUID.String(), nil
}
