package api

import (
	"net/http"

	"github.com/shirou/gopsutil/disk"
)

func ServerInfoOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
}

func ServerInfoGetHandler(w http.ResponseWriter, r *http.Request) {
	type ServerInfoResponse struct {
		DiskUsed uint64 `json:"disk_used"`
		DiskFree uint64 `json:"disk_free"`
	}

	// #### CHECK IF AUTHENTICATED ####
	_, err := AuthSession(r)
	if err != nil {
		DumpResponse(w, "unauthorized", http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	dstat, err := disk.Usage("/")
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_SERVER_INFO_FAILED, nil)
		return
	}

	resp := &ServerInfoResponse{}
	resp.DiskUsed = dstat.Used
	resp.DiskFree = dstat.Free

	DumpResponse(w, "ok", http.StatusOK, 0, resp)
}
