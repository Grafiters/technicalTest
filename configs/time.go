package configs

import "time"

func TimeToIndo(val time.Time) time.Time {
	const asiaJakarta = 7 * 60 * 60

	return val.Add(time.Second * time.Duration(asiaJakarta))
}
