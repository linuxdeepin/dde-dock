Sunrise Sunset Package
==========================

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](LICENSE)
[![Build Status](https://travis-ci.org/kelvins/sunrisesunset.svg?branch=master)](https://travis-ci.org/kelvins/sunrisesunset)
[![Coverage Status](https://coveralls.io/repos/github/kelvins/sunrisesunset/badge.svg?branch=master)](https://coveralls.io/github/kelvins/sunrisesunset?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/kelvins/sunrisesunset)](https://goreportcard.com/report/github.com/kelvins/sunrisesunset)

Go package used to calculate the apparent sunrise and sunset times based on latitude, longitude, date and UTC offset.

You can use go get command:

    go get github.com/kelvins/sunrisesunset

![](http://i.imgur.com/hjUZT28.jpg)

General
----

This package was created based on the [Corrected Sunrise, Sunset, Noon Times in Seconds - and Solar Angles][1] Matlab function developed by Richard Droste, that was created based on the [spreadsheets][2] available in the [Earth System Research Laboratory][3] website from the [National Oceanic & Atmospheric Administration (NOAA)][4].

Based on the [Solar Calculation Details][5]:

> The calculations in the NOAA Sunrise/Sunset and Solar Position Calculators are based on equations from Astronomical Algorithms, by Jean Meeus. The sunrise and sunset results are theoretically accurate to within a minute for locations between +/- 72° latitude, and within 10 minutes outside of those latitudes. However, due to variations in atmospheric composition, temperature, pressure and conditions, observed values may vary from calculations.

[Apparent Sunrise/Sunset][6]:

> Due to atmospheric refraction, sunrise occurs shortly before the sun crosses above the horizon. Light from the sun is bent, or refracted, as it enters earth's atmosphere. See [Apparent Sunrise Figure][7]. This effect causes the apparent sunrise to be earlier than the actual sunrise. Similarly, apparent sunset occurs slightly later than actual sunset. The sunrise and sunset times reported in our calculator have been corrected for the approximate effects of atmospheric refraction. However, it should be noted that due to changes in air pressure, relative humidity, and other quantities, we cannot predict the exact effects of atmospheric refraction on sunrise and sunset time. Also note that this possible error increases with higher (closer to the poles) latitudes.

Usage
----

``` go

package main

import (
    "fmt"
    "time"
    "github.com/kelvins/sunrisesunset"
)

func main() {
    // You can use the Parameters structure to set the parameters
    p := sunrisesunset.Parameters{
      Latitude:  -23.545570,
      Longitude: -46.704082,
      UtcOffset: -3.0,
      Date:      time.Date(2017, 3, 23, 0, 0, 0, 0, time.UTC),
    }

    // Calculate the sunrise and sunset times
    sunrise, sunset, err := p.GetSunriseSunset()

    // If no error has occurred, print the results
    if err == nil {
        fmt.Println("Sunrise:", sunrise.Format("15:04:05")) // Sunrise: 06:11:44
        fmt.Println("Sunset:", sunset.Format("15:04:05")) // Sunset: 18:14:27
    } else {
        fmt.Println(err)
    }
}

```

License
----

This project was created under the [MIT license][8]


  [1]: https://www.mathworks.com/matlabcentral/fileexchange/62180-corrected-sunrise--sunset--noon-times-in-seconds-and-solar-angles?requestedDomain=www.mathworks.com
  [2]: https://www.esrl.noaa.gov/gmd/grad/solcalc/calcdetails.html
  [3]: https://www.esrl.noaa.gov/
  [4]: http://www.noaa.gov/
  [5]: https://www.esrl.noaa.gov/gmd/grad/solcalc/calcdetails.html
  [6]: https://www.esrl.noaa.gov/gmd/grad/solcalc/glossary.html#A
  [7]: https://www.esrl.noaa.gov/gmd/grad/solcalc/apparent_sunrise.gif
  [8]: LICENSE
