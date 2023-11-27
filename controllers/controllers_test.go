package controllers

import (
	"net/http"
	"reakgo/utility"
)

type MockHelper struct {
    RenderTemplateCalled    bool
	MockReturnUserDetailsResult                error
	MockGenerateRandomStringStrResult          string
	MockGenerateRandomStringErrorResult        error
	MockSessionGetResult                       interface{}
	MockViewFlashResult                        interface{}
	MockParseDataFromPostRequestToMapResult    map[string]interface{}
	MockParseDataFromPostRequestToMapErrResult error
	MockParseDataFromJsonToMapResult           map[string]interface{}
	MockParseDataFromJsonToMapErrResult        error
	// MockStrictParseDataFromJsonResult          error
	MockStrictParseDataFromPostRequestResult  error
	MockStringInArray                         bool
    MockError   error
}

// type Session struct {
// 	Key   string
// 	Value interface{}
// }

func (m MockHelper) AddFlash(flavour string, message string, w http.ResponseWriter, r *http.Request) {

}

func (m MockHelper) GenerateRandomString(n int) (string, error) {
	return m.MockGenerateRandomStringStrResult, m.MockGenerateRandomStringErrorResult
}

func (m MockHelper) RedirectTo(w http.ResponseWriter, r *http.Request, path string) {

}

func (m MockHelper) SessionGet(r *http.Request, key string) interface{} {
	return m.MockSessionGetResult
}

func (m MockHelper) SessionSet(w http.ResponseWriter, r *http.Request, data utility.Session) {

}

func (m MockHelper) ViewFlash(w http.ResponseWriter, r *http.Request) interface{} {
	return m.MockViewFlashResult
}

func (m MockHelper) RenderTemplate(w http.ResponseWriter, r *http.Request, template string, data interface{}) {
    w.Header().Set("RenderTemplateCalled", "true")
}

func (m MockHelper) ParseDataFromPostRequestToMap(r *http.Request) (map[string]interface{}, error) {
	return m.MockParseDataFromPostRequestToMapResult, m.MockParseDataFromPostRequestToMapErrResult
}

func (m MockHelper) ParseDataFromJsonToMap(r *http.Request) (map[string]interface{}, error) {
	return m.MockParseDataFromJsonToMapResult, m.MockParseDataFromJsonToMapErrResult
}

// func (m MockHelper) StrictParseDataFromJson(r *http.Request, structure interface{}) error {
// 	return m.MockStrictParseDataFromJsonResult
// }

func (m MockHelper) StrictParseDataFromPostRequest(r *http.Request, structure interface{}) error {
	return m.MockStrictParseDataFromPostRequestResult
}

func (m MockHelper) RenderJsonResponse(w http.ResponseWriter, r *http.Request, data interface{}) {

}

func (m MockHelper) RenderTemplateData(w http.ResponseWriter, r *http.Request, template string, data interface{}) {

}

func (m MockHelper) StringInArray(target string, arr []string) bool {
	return m.MockStringInArray
}

func (m MockHelper) StrictParseDataFromJson(r *http.Request, structure interface{}) error {
    return m.MockError
}

