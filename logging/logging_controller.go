package logging

import (
	"encoding/json"
	"net/http"
)

type changeLevelRequest struct {
	LvlFromRequest string `json:"lvl"`
	PackageName    string `json:"packageName"`
}

func ChangeLogLevel(w http.ResponseWriter, r *http.Request) {
	controllerLogger := GetLogger("logging")
	defer r.Body.Close()
	controllerLogger.Info("Start Change Log Level API call")
	decoder := json.NewDecoder(r.Body)

	var data changeLevelRequest
	if err := decoder.Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err := setLogLevel(data.LvlFromRequest, data.PackageName)
	if err != nil {
		controllerLogger.Warn("Log Level didn't changed")
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		controllerLogger.Info("Successfully change logLevel to %s", data.LvlFromRequest)
		respondWithJson(w, http.StatusOK, "Successfully change logLevel")
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
