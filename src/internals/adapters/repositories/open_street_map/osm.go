package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	osmmodels "github.com/gcarrenho/guidemysteps/src/internals/adapters/repositories/open_street_map/models"
	"github.com/gcarrenho/guidemysteps/src/internals/core/models"
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

// Número máximo de reintentos
//const maxRetries = 3

// Tiempo de espera entre reintentos (usamos exponencial backoff)
//const retryDelay = 2 * time.Second

/*
func (r *routingRepo) GetRouting(routesRequest models.RoutesRequest) (*osmmodels.OsmResponse, error) {
	var osmResponse osmmodels.OsmResponse
	url := r.baseUrl + "/" + routesRequest.DrivingMode + "" + "/route/v1/driving/" +
		fmt.Sprintf("%f", routesRequest.Start.Latitud) + "," + fmt.Sprintf("%f", routesRequest.Start.Longitud) + ";" +
		fmt.Sprintf("%f", routesRequest.Destination.Latitud) + "," + fmt.Sprintf("%f", routesRequest.Destination.Longitud) +
		"?overview=false&geometries=geojson&steps=true"

	fmt.Println(url)

	var resp *http.Response
	var err error

	// Intentamos hacer la solicitud hasta un máximo de `maxRetries` veces
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Hacemos la solicitud HTTP
		resp, err = r.cli.Get(url)
		if err == nil {
			defer resp.Body.Close()
			break // Salimos del bucle si no hay error
		}

		// Si estamos en el último intento, devolvemos el error
		if attempt == maxRetries {
			fmt.Println("Error después de múltiples intentos:", err)
			return nil, err
		}

		// Mostramos el error y esperamos antes de volver a intentar
		fmt.Printf("Intento %d fallido, reintentando... Error: %v\n", attempt, err)
		time.Sleep(retryDelay * time.Duration(attempt)) // Exponencial backoff
	}

	// Leemos el cuerpo de la respuesta
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Deserializamos la respuesta en la estructura `osmResponse`
	err = json.Unmarshal(respBytes, &osmResponse)
	if err != nil {
		fmt.Println("Error al deserializar la respuesta:", err)
		return nil, err
	}

	return &osmResponse, nil
}*/
/*
func (r *routingRepo) GetRouting(routesRequest models.RoutesRequest) (*osmmodels.OsmResponse, error) {
	var osmResponse osmmodels.OsmResponse
	//https://routing.openstreetmap.de/routed-foot/route/v1/driving/41.930021,2.263037;41.930492,2.254357?overview=false&geometries=geojson&steps=true
	//"https://routing.openstreetmap.de/routed-foot/route/v1/driving/2.20801,41.40289;2.170084,41.3865465?overview=false&geometries=geojson&steps=true"
	url := r.baseUrl + "/" + routesRequest.DrivingMode + "" + "/route/v1/driving/" + fmt.Sprintf("%f", routesRequest.Start.Latitud) + "," + fmt.Sprintf("%f", routesRequest.Start.Longitud) + ";" + fmt.Sprintf("%f", routesRequest.Destination.Latitud) + "," + fmt.Sprintf("%f", routesRequest.Destination.Longitud) + "?overview=false&geometries=geojson&steps=true"
	fmt.Println(url)
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
}*/

func (r *routingRepo) GetRouting(routesRequest models.RoutesRequest) (*osmmodels.OsmResponse, error) {
	const maxRetries = 3             // Definir un número máximo de reintentos
	var retryDelay = time.Second * 3 // Tiempo de espera entre reintentos

	var osmResponse osmmodels.OsmResponse

	url := r.baseUrl + "/" + routesRequest.DrivingMode + "/route/v1/driving/" +
		fmt.Sprintf("%f", routesRequest.Start.Longitud) + "," +
		fmt.Sprintf("%f", routesRequest.Start.Latitud) + ";" +
		fmt.Sprintf("%f", routesRequest.Destination.Longitud) + "," +
		fmt.Sprintf("%f", routesRequest.Destination.Latitud) +
		"?overview=false&geometries=geojson&steps=true"

	var resp *http.Response
	var err error
	fmt.Println(url)
	// Implementar un bucle de reintentos en caso de fallo
	for i := 0; i < maxRetries; i++ {
		resp, err = r.cli.Get(url)

		if err != nil {
			if errors.Is(err, io.EOF) {
				// Mostramos el error y esperamos antes de volver a intentar
				fmt.Printf("Intento %d fallido, reintentando... Error: %v\n", i, err)
				// Si se trata de un error de EOF, espera y reintenta
				time.Sleep(retryDelay)
				continue
			}
			return nil, err // Otros errores se devuelven directamente
		}
		if resp.StatusCode == 500 {
			fmt.Printf("Intento %d fallido, reintentando... Error: %v\n", i, err)

			time.Sleep(retryDelay)
			continue
		}
		fmt.Println("exito ", resp)
		break // Si la solicitud es exitosa, sale del bucle
	}

	if resp == nil {
		return nil, fmt.Errorf("respuesta nula del servidor después de %d intentos", maxRetries)
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBytes, &osmResponse)
	if err != nil {
		return nil, fmt.Errorf("error deserializando la respuesta: %w", err)
	}

	return &osmResponse, nil
}

/*var drivingMode = map[string]string{
	"foot":             "routed-foot",
	"public-transport": "",
	"car":              "routed-car",
}*/
