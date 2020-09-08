package appearance

import (
	"time"

	"github.com/kelvins/sunrisesunset"
)

func (m *Manager) getSunriseSunset(t time.Time, latitude, longitude float64) (time.Time, time.Time, error) {
	timeUtc := t.UTC()
	_, offsetSec := t.Zone()
	utcOffset := float64(offsetSec / 3600.0)
	sunrise, sunset, err := sunrisesunset.GetSunriseSunset(latitude, longitude, utcOffset, timeUtc)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	sunriseT := time.Date(t.Year(), t.Month(), t.Day(),
		sunrise.Hour(), sunrise.Minute(), sunrise.Second(),
		0, m.loc)
	sunsetT := time.Date(t.Year(), t.Month(), t.Day(),
		sunset.Hour(), sunset.Minute(), sunset.Second(),
		0, m.loc)
	return sunriseT, sunsetT, nil
}

func isDaytime(t, sunriseT, sunsetT time.Time) bool {
	return sunriseT.Before(t) && t.Before(sunsetT)
}

func getThemeAutoName(isDaytime bool) string {
	if isDaytime {
		return "deepin"
	}
	return "deepin-dark"
}

func (m *Manager) getThemeAutoChangeTime(t time.Time, latitude, longitude float64) (time.Time, error) {
	sunrise, sunset, err := m.getSunriseSunset(t, latitude, longitude)
	if err != nil {
		return time.Time{}, err
	}
	logger.Debugf("t: %v, sunrise: %v, sunset: %v", t, sunrise, sunset)
	if t.Before(sunrise) || t.Equal(sunrise) {
		// t <= sunrise
		return sunrise, nil
	}

	if t.Before(sunset) || t.Equal(sunset) {
		// t <= sunset
		return sunset, nil
	}

	nextDay := t.AddDate(0, 0, 1)
	nextDaySunrise, _, err := m.getSunriseSunset(nextDay, latitude, longitude)
	logger.Debug("next day sunrise:", nextDaySunrise)
	if err != nil {
		return time.Time{}, err
	}

	return nextDaySunrise, nil
}
