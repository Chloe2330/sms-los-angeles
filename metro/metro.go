package metro

// https://github.com/mholt/json-to-go

type MetroPredictions struct {
	Success bool   `json:"success"`
	Route   string `json:"route"`
	Data    struct {
		AgencyKey       string `json:"agencyKey"`
		PredictionsData []struct {
			RouteShortName string `json:"routeShortName"`
			RouteName      string `json:"routeName"`
			RouteID        string `json:"routeId"`
			StopID         string `json:"stopId"`
			StopName       string `json:"stopName"`
			StopCode       int    `json:"stopCode"`
			Destinations   []struct {
				DirectionID string `json:"directionId"`
				Headsign    string `json:"headsign"`
				Predictions []struct {
					Time          int    `json:"time"`
					Sec           int    `json:"sec"`
					Min           int    `json:"min"`
					TripID        string `json:"tripId"`
					VehicleID     string `json:"vehicleId"`
					ScheduleBased bool   `json:"scheduleBased,omitempty"`
				} `json:"predictions"`
			} `json:"destinations"`
			DistanceToStop float64 `json:"distanceToStop"`
		} `json:"predictionsData"`
	} `json:"data"`
}

type MetroPredictionsFormatted struct {
	RouteName        string
	StopName         string
	DestStopName     string
	MinsUntilArrival string
}
