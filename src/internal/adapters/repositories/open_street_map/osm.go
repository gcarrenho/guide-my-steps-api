package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	osmmodels "project/guidemysteps/src/internal/adapters/repositories/open_street_map/models"
	"project/guidemysteps/src/internal/core/models"
	"time"
)

type routingRepo struct {
	cli     *http.Client
	baseUrl string
}

func NewRoutingRepo(baseUrl string) *routingRepo {
	return &routingRepo{
		cli: &http.Client{
			Timeout: time.Second * 60,
		},
		baseUrl: baseUrl,
	}
}

func (r *routingRepo) GetRouting(routesRequest models.RoutesRequest) (*osmmodels.OsmResponse, error) {
	var osmResponse osmmodels.OsmResponse
	//"https://routing.openstreetmap.de/routed-foot/route/v1/driving/2.20801,41.40289;2.170084,41.3865465?overview=false&geometries=geojson&steps=true" /*r.baseUrl*/
	url := r.baseUrl + "/" + routesRequest.DrivingMode + "" + "/route/v1/driving/" + fmt.Sprintf("%f", routesRequest.Start.Latitud) + "," + fmt.Sprintf("%f", routesRequest.Start.Longitud) + ";" + fmt.Sprintf("%f", routesRequest.Destination.Latitud) + "," + fmt.Sprintf("%f", routesRequest.Destination.Longitud) + "?overview=false&geometries=geojson&steps=true"

	resp, err := r.cli.Get(url)
	if err != nil {
		fmt.Println("erro ", err)
		return nil, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBytes, &osmResponse)
	if err != nil {
		fmt.Println("error")
		return nil, err
	}

	return &osmResponse, nil
}

/*var drivingMode = map[string]string{
	"foot":             "routed-foot",
	"public-transport": "",
	"car":              "routed-car",
}*/
