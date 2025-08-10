package test

import (
	"net/http"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUnknownConfigFails(t *testing.T) {

	requestBody := strings.NewReader(`{"max_limit":1000}`)
	_, updateHTTP := performRequest(http.MethodPost, "/configs/unknown_config/payments", requestBody, true)

	assert.Equal(t, http.StatusBadRequest, updateHTTP.StatusCode)
}

func TestCreatePaymentConfigSuccess(t *testing.T) {

	requestBody := strings.NewReader(`{"max_limit":1000,"enabled":true}`)
	response, createHTTP := performRequest(http.MethodPost, "/configs/payment_config/payments", requestBody, true)

	assert.Equal(t, http.StatusOK, createHTTP.StatusCode)

	createdData := response["data"].(map[string]interface{})

	assert.Equal(t, "payment_config", createdData["schema"])
	assert.Equal(t, "payments", createdData["name"])

	versionFloat := createdData["version"].(float64)
	assert.Equal(t, 1, int(versionFloat))

	properties := createdData["data"].(map[string]interface{})
	assert.Equal(t, true, properties["enabled"])
	maxLimitFloat := properties["max_limit"].(float64)
	assert.Equal(t, 1000, int(maxLimitFloat))
}

func TestCreatePaymentConfigFails(t *testing.T) {

	requestBody := strings.NewReader(`{"max_limit":1000}`)
	_, updateHTTP := performRequest(http.MethodPut, "/configs/payment_config/payments", requestBody, true)

	assert.Equal(t, http.StatusBadRequest, updateHTTP.StatusCode)
}

func TestUpdatePaymentConfigSuccess(t *testing.T) {

	// 1. Create initial config
	requestBody := strings.NewReader(`{"max_limit":1000,"enabled":true}`)

	createResp, createHTTP := performRequest(http.MethodPost, "/configs/payment_config/payments", requestBody, true)
	assert.Equal(t, http.StatusOK, createHTTP.StatusCode)

	createdData := createResp["data"].(map[string]interface{})
	assert.Equal(t, "payment_config", createdData["schema"])
	assert.Equal(t, "payments", createdData["name"])
	assert.Equal(t, 1, int(createdData["version"].(float64)))

	// 2. Update config
	secondRequestBody := strings.NewReader(`{"max_limit":5000,"enabled":false}`)

	updateResp, updateHTTP := performRequest(http.MethodPut, "/configs/payment_config/payments", secondRequestBody, false)
	assert.Equal(t, http.StatusOK, updateHTTP.StatusCode)

	updatedData := updateResp["data"].(map[string]interface{})
	assert.Equal(t, "payment_config", updatedData["schema"])
	assert.Equal(t, "payments", updatedData["name"])
	assert.Equal(t, 2, int(updatedData["version"].(float64)))

	props := updatedData["data"].(map[string]interface{})
	assert.Equal(t, false, props["enabled"].(bool))
	assert.Equal(t, 5000, int(props["max_limit"].(float64)))
}

func TestUpdateNonExistPaymentConfigFail(t *testing.T) {

	secondRequestBody := strings.NewReader(`{"max_limit":5000,"enabled":false}`)

	_, updateHTTP := performRequest(http.MethodPut, "/configs/payment_config/payments", secondRequestBody, true)
	assert.Equal(t, http.StatusBadRequest, updateHTTP.StatusCode)

}

func TestFetchLatestPaymentConfigSuccess(t *testing.T) {

	createRequestBody := strings.NewReader(`{"max_limit":500,"enabled":true}`)
	_, createHttp := performRequest(http.MethodPost, "/configs/payment_config/payments", createRequestBody, true)
	assert.Equal(t, http.StatusOK, createHttp.StatusCode)

	fetchResp, fetchHttp := performRequest(http.MethodGet, "/configs/payment_config/payments", nil, false)
	assert.Equal(t, http.StatusOK, fetchHttp.StatusCode)

	fetchData := fetchResp["data"].(map[string]interface{})
	assert.Equal(t, "payment_config", fetchData["schema"])
	assert.Equal(t, "payments", fetchData["name"])
	assert.Equal(t, 1, int(fetchData["version"].(float64)))

	props := fetchData["data"].(map[string]interface{})
	assert.Equal(t, true, props["enabled"].(bool))
	assert.Equal(t, 500, int(props["max_limit"].(float64)))
}

func TestFetchNonexistConfigFails(t *testing.T) {

	fetchRequestBody := strings.NewReader(`{"version":1}`)
	_, fetchHttp := performRequest(http.MethodGet, "/configs/payment_config/payments", fetchRequestBody, true)
	assert.Equal(t, http.StatusBadRequest, fetchHttp.StatusCode)
}

func TestFetchAllVersionsPaymentConfigSuccess(t *testing.T) {

	createRequestBody := strings.NewReader(`{"max_limit":500,"enabled":true}`)
	_, createHttp := performRequest(http.MethodPost, "/configs/payment_config/payments", createRequestBody, true)
	assert.Equal(t, http.StatusOK, createHttp.StatusCode)

	secondCreateRequestBody := strings.NewReader(`{"max_limit":1000,"enabled":true}`)
	_, secondCreateHttp := performRequest(http.MethodPut, "/configs/payment_config/payments", secondCreateRequestBody, false)
	assert.Equal(t, http.StatusOK, secondCreateHttp.StatusCode)

	fetchResp, fetchHttp := performRequest(http.MethodGet, "/configs/payment_config/payments/versions", nil, false)
	assert.Equal(t, http.StatusOK, fetchHttp.StatusCode)

	fetchData := fetchResp["data"].(map[string]interface{})
	assert.Equal(t, "payment_config", fetchData["schema"])
	assert.Equal(t, "payments", fetchData["name"])

	versions := fetchData["ConfigVersions"].([]interface{})
	assert.Equal(t, 2, len(versions))

	firstVersion := versions[0].(map[string]interface{})
	assert.Equal(t, 1, int(firstVersion["version"].(float64)))
	assert.Equal(t, float64(500), firstVersion["data"].(map[string]interface{})["max_limit"])

	secondVersion := versions[1].(map[string]interface{})
	assert.Equal(t, 2, int(secondVersion["version"].(float64)))
	assert.Equal(t, float64(1000), secondVersion["data"].(map[string]interface{})["max_limit"])
}

func TestRollbackPaymentFunctionSuccess(t *testing.T) {

	// 1. Create initial config
	requestBody := strings.NewReader(`{"max_limit":1000,"enabled":true}`)

	createResp, createHTTP := performRequest(http.MethodPost, "/configs/payment_config/payments", requestBody, true)
	assert.Equal(t, http.StatusOK, createHTTP.StatusCode)

	createdData := createResp["data"].(map[string]interface{})
	assert.Equal(t, "payment_config", createdData["schema"])
	assert.Equal(t, "payments", createdData["name"])
	assert.Equal(t, 1, int(createdData["version"].(float64)))

	// 2. Update config
	secondRequestBody := strings.NewReader(`{"max_limit":5000,"enabled":false}`)

	updateResp, updateHTTP := performRequest(http.MethodPut, "/configs/payment_config/payments", secondRequestBody, false)
	assert.Equal(t, http.StatusOK, updateHTTP.StatusCode)

	updatedData := updateResp["data"].(map[string]interface{})
	assert.Equal(t, "payment_config", updatedData["schema"])
	assert.Equal(t, "payments", updatedData["name"])
	assert.Equal(t, 2, int(updatedData["version"].(float64)))

	// 3. Rollback config
	thirdRequestBody := strings.NewReader(`{"version":1}`)

	rollbackResp, rollbackHttp := performRequest(http.MethodPost, "/configs/payment_config/payments/rollback", thirdRequestBody, false)
	assert.Equal(t, http.StatusOK, rollbackHttp.StatusCode)

	rollbackData := rollbackResp["data"].(map[string]interface{})
	assert.Equal(t, "payment_config", rollbackData["schema"])
	assert.Equal(t, "payments", rollbackData["name"])
	assert.Equal(t, 3, int(rollbackData["version"].(float64)))

	props := rollbackData["data"].(map[string]interface{})
	assert.Equal(t, true, props["enabled"].(bool))
	assert.Equal(t, 1000, int(props["max_limit"].(float64)))
}
