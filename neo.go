package neows

type NearEarthObject struct {
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
	ID                 string  `json:"id"`
	NeoReferenceID     string  `json:"neo_reference_id"`
	Name               string  `json:"name"`
	NameLimited        string  `json:"name_limited"`
	Designation        string  `json:"designation"`
	NasaJplURL         string  `json:"nasa_jpl_url"`
	AbsoluteMagnitudeH float64 `json:"absolute_magnitude_h"`
	EstimatedDiameter  struct {
		Kilometers struct {
			EstimatedDiameterMin float64 `json:"estimated_diameter_min"`
			EstimatedDiameterMax float64 `json:"estimated_diameter_max"`
		} `json:"kilometers"`
		Meters struct {
			EstimatedDiameterMin float64 `json:"estimated_diameter_min"`
			EstimatedDiameterMax float64 `json:"estimated_diameter_max"`
		} `json:"meters"`
		Miles struct {
			EstimatedDiameterMin float64 `json:"estimated_diameter_min"`
			EstimatedDiameterMax float64 `json:"estimated_diameter_max"`
		} `json:"miles"`
		Feet struct {
			EstimatedDiameterMin float64 `json:"estimated_diameter_min"`
			EstimatedDiameterMax float64 `json:"estimated_diameter_max"`
		} `json:"feet"`
	} `json:"estimated_diameter"`
	IsPotentiallyHazardousAsteroid bool `json:"is_potentially_hazardous_asteroid"`
	CloseApproachData              []struct {
		CloseApproachDate      string `json:"close_approach_date"`
		CloseApproachDateFull  string `json:"close_approach_date_full"`
		EpochDateCloseApproach int64  `json:"epoch_date_close_approach"`
		RelativeVelocity       struct {
			KilometersPerSecond string `json:"kilometers_per_second"`
			KilometersPerHour   string `json:"kilometers_per_hour"`
			MilesPerHour        string `json:"miles_per_hour"`
		} `json:"relative_velocity"`
		MissDistance struct {
			Astronomical string `json:"astronomical"`
			Lunar        string `json:"lunar"`
			Kilometers   string `json:"kilometers"`
			Miles        string `json:"miles"`
		} `json:"miss_distance"`
		OrbitingBody string `json:"orbiting_body"`
	} `json:"close_approach_data"`
	OrbitalData struct {
		OrbitID                   string `json:"orbit_id"`
		OrbitDeterminationDate    string `json:"orbit_determination_date"`
		FirstObservationDate      string `json:"first_observation_date"`
		LastObservationDate       string `json:"last_observation_date"`
		DataArcInDays             int    `json:"data_arc_in_days"`
		ObservationsUsed          int    `json:"observations_used"`
		OrbitUncertainty          string `json:"orbit_uncertainty"`
		MinimumOrbitIntersection  string `json:"minimum_orbit_intersection"`
		JupiterTisserandInvariant string `json:"jupiter_tisserand_invariant"`
		Eccentricity              string `json:"eccentricity"`
		SemiMajorAxis             string `json:"semi_major_axis"`
		AscendingNodeLongitude    string `json:"ascending_node_longitude"`
		OrbitalPeriod             string `json:"orbital_period"`
		Inclination               string `json:"inclination"`
		AxialTilt                 string `json:"axial_tilt"`
		PerihelionDistance        string `json:"perihelion_distance"`
		PerihelionArgument        string `json:"perihelion_argument"`
		AphelionDistance          string `json:"aphelion_distance"`
		PerihelionTime            string `json:"perihelion_time"`
		MeanAnomaly               string `json:"mean_anomaly"`
		MeanMotion                string `json:"mean_motion"`
		Equinox                   string `json:"equinox"`
		OrbitClass                struct {
			Type        string `json:"orbit_class_type"`
			Description string `json:"orbit_class_description"`
			Range       string `json:"orbit_class_range"`
		} `json:"orbit_class"`
	} `json:"orbital_data"`
	IsSentryObject bool `json:"is_sentry_object"`
}
