package test

import (
	"net/http"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestCreateUnknownConfigFails(t *testing.T) {

	requestBody := strings.NewReader(`{"max_limit":1000}`)
	_, updateHTTP := performRequest(http.MethodPost, "/configs/unknown_config/payments", requestBody, true)

	assert.Equal(t, http.StatusBadRequest, updateHTTP.StatusCode)
}

func TestCreatePaymentConfigSuccess(t *testing.T) {

	requestBody := strings.NewReader(`{"max_limit":1000,"enabled":true}`)
	response, createHTTP := performRequest(http.MethodPost, "/configs/payment_config/payments", requestBody, true)

	assert.Equal(t, http.StatusCreated, createHTTP.StatusCode)

	assert.Equal(t, "payment_config", response["schema"])
	assert.Equal(t, "payments", response["name"])

	versionFloat := response["version"].(float64)
	assert.Equal(t, 1, int(versionFloat))

	properties := response["data"].(map[string]interface{})

	// Enabled
	enabledVal, ok := properties["enabled"].(bool)
	if assert.True(t, ok, "enabled should be a bool") {
		assert.Equal(t, true, enabledVal)
	}

	// Max limit
	maxLimitFloat, ok := properties["max_limit"].(float64)
	if assert.True(t, ok, "max_limit should be a number") {
		assert.Equal(t, 1000, int(maxLimitFloat))
	}
}

func TestCreatePaymentConfigFails(t *testing.T) {

	requestBody := strings.NewReader(`{"max_limit":1000}`)
	_, updateHTTP := performRequest(http.MethodPost, "/configs/payment_config/payments", requestBody, true)

	assert.Equal(t, http.StatusBadRequest, updateHTTP.StatusCode)
}

func TestUpdatePaymentConfigSuccess(t *testing.T) {

	// 1. Create initial config
	requestBody := strings.NewReader(`{"max_limit":1000,"enabled":true}`)

	createResp, createHTTP := performRequest(http.MethodPost, "/configs/payment_config/payments", requestBody, true)
	assert.Equal(t, http.StatusCreated, createHTTP.StatusCode)

	assert.Equal(t, "payment_config", createResp["schema"])
	assert.Equal(t, "payments", createResp["name"])
	assert.Equal(t, 1, int(createResp["version"].(float64)))

	// 2. Update config
	secondRequestBody := strings.NewReader(`{"max_limit":5000,"enabled":false}`)

	updateResp, updateHTTP := performRequest(http.MethodPut, "/configs/payment_config/payments", secondRequestBody, false)
	assert.Equal(t, http.StatusOK, updateHTTP.StatusCode)

	assert.Equal(t, "payment_config", updateResp["schema"])
	assert.Equal(t, "payments", updateResp["name"])
	assert.Equal(t, 2, int(updateResp["version"].(float64)))

	props := updateResp["data"].(map[string]interface{})
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
	assert.Equal(t, http.StatusCreated, createHttp.StatusCode)

	fetchResp, fetchHttp := performRequest(http.MethodGet, "/configs/payment_config/payments", nil, false)
	assert.Equal(t, http.StatusOK, fetchHttp.StatusCode)

	assert.Equal(t, "payment_config", fetchResp["schema"])
	assert.Equal(t, "payments", fetchResp["name"])
	assert.Equal(t, 1, int(fetchResp["version"].(float64)))

	props := fetchResp["data"].(map[string]interface{})
	assert.Equal(t, true, props["enabled"].(bool))
	assert.Equal(t, 500, int(props["max_limit"].(float64)))
}

func TestFetchSpecificVersionPaymentConfigSuccess(t *testing.T) {

	createRequestBody := strings.NewReader(`{"max_limit":500,"enabled":true}`)
	_, createHttp := performRequest(http.MethodPost, "/configs/payment_config/payments", createRequestBody, true)
	assert.Equal(t, http.StatusCreated, createHttp.StatusCode)

	fetchRequestBody := strings.NewReader(`{"version":1}`)
	fetchResp, fetchHttp := performRequest(http.MethodGet, "/configs/payment_config/payments", fetchRequestBody, false)
	assert.Equal(t, http.StatusOK, fetchHttp.StatusCode)

	assert.Equal(t, "payment_config", fetchResp["schema"])
	assert.Equal(t, "payments", fetchResp["name"])
	assert.Equal(t, 1, int(fetchResp["version"].(float64)))

	props := fetchResp["data"].(map[string]interface{})
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
	assert.Equal(t, http.StatusCreated, createHttp.StatusCode)

	secondCreateRequestBody := strings.NewReader(`{"max_limit":1000,"enabled":true}`)
	_, secondCreateHttp := performRequest(http.MethodPut, "/configs/payment_config/payments", secondCreateRequestBody, false)
	assert.Equal(t, http.StatusOK, secondCreateHttp.StatusCode)

	fetchResp, fetchHttp := performRequest(http.MethodGet, "/configs/payment_config/payments/versions", nil, false)
	assert.Equal(t, http.StatusOK, fetchHttp.StatusCode)

	assert.Equal(t, "payment_config", fetchResp["schema"])
	assert.Equal(t, "payments", fetchResp["name"])

	versions := fetchResp["configVersions"].([]interface{})
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
	assert.Equal(t, http.StatusCreated, createHTTP.StatusCode)

	assert.Equal(t, "payment_config", createResp["schema"])
	assert.Equal(t, "payments", createResp["name"])
	assert.Equal(t, 1, int(createResp["version"].(float64)))

	// 2. Update config
	secondRequestBody := strings.NewReader(`{"max_limit":5000,"enabled":false}`)

	updateResp, updateHTTP := performRequest(http.MethodPut, "/configs/payment_config/payments", secondRequestBody, false)
	assert.Equal(t, http.StatusOK, updateHTTP.StatusCode)

	assert.Equal(t, "payment_config", updateResp["schema"])
	assert.Equal(t, "payments", updateResp["name"])
	assert.Equal(t, 2, int(updateResp["version"].(float64)))

	// 3. Rollback config
	thirdRequestBody := strings.NewReader(`{"version":1}`)

	rollbackResp, rollbackHttp := performRequest(http.MethodPost, "/configs/payment_config/payments/rollback", thirdRequestBody, false)
	assert.Equal(t, http.StatusOK, rollbackHttp.StatusCode)

	assert.Equal(t, "payment_config", rollbackResp["schema"])
	assert.Equal(t, "payments", rollbackResp["name"])
	assert.Equal(t, 3, int(rollbackResp["version"].(float64)))

	props := rollbackResp["data"].(map[string]interface{})
	assert.True(t, props["enabled"].(bool))
	assert.Equal(t, 1000, int(props["max_limit"].(float64)))
}
