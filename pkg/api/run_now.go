package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/flanksource/canary-checker/api/context"
	"github.com/flanksource/canary-checker/checks"
	"github.com/flanksource/canary-checker/pkg"
	"github.com/flanksource/canary-checker/pkg/db"
	"github.com/flanksource/canary-checker/pkg/topology"
	"github.com/flanksource/commons/logger"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CheckErrorMessage struct {
	Description string `json:"description"`
	Error       string `json:"error"`
}

// RunCanaryResponse represents the response body for a run now request
type RunCanaryResponse struct {
	Total   int                 `json:"total"`
	Failed  int                 `json:"failed"`
	Success int                 `json:"success"`
	Errors  []CheckErrorMessage `json:"errors,omitempty"`
}

func (t *RunCanaryResponse) FromCheckResults(result []*pkg.CheckResult) {
	t.Total = len(result)
	for _, r := range result {
		if r.Pass {
			t.Success++
			continue
		}

		t.Failed++
		if r.Error != "" {
			t.Errors = append(t.Errors, CheckErrorMessage{
				Description: r.GetDescription(),
				Error:       r.Error,
			})
		}
	}
}

func RunCanaryHandler(c echo.Context) error {
	id := c.Param("id")

	canaryModel, err := db.GetCanary(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorResonse(c, fmt.Errorf("canary with id=%s was not found", id), http.StatusNotFound)
		}

		return errorResonse(c, err, http.StatusInternalServerError)
	}

	canary, err := canaryModel.ToV1()
	if err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}

	kommonsClient, err := pkg.NewKommonsClient()
	if err != nil {
		logger.Warnf("failed to get kommons client, checks that read kubernetes configs will fail: %v", err)
	}

	ctx := context.New(kommonsClient, *canary)
	result := checks.RunChecks(ctx)

	var response RunCanaryResponse
	response.FromCheckResults(result)
	return c.JSON(http.StatusOK, response)
}

func RunTopologyHandler(c echo.Context) error {
	id := c.Param("id")

	topologyRunDepth := 10
	_depth := c.QueryParam("depth")
	if _depth != "" {
		num, err := strconv.Atoi(_depth)
		if err != nil {
			return errorResonse(c, err, http.StatusBadRequest)
		}

		if num < 0 {
			return errorResonse(c, fmt.Errorf("depth must be greater than 0"), http.StatusBadRequest)
		}

		topologyRunDepth = num
	}

	systemTemplate, err := db.GetSystemTemplate(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorResonse(c, fmt.Errorf("topology with id=%s was not found", id), http.StatusNotFound)
		}

		return errorResonse(c, err, http.StatusInternalServerError)
	}

	kommonsClient, err := pkg.NewKommonsClient()
	if err != nil {
		logger.Warnf("failed to get kommons client, checks that read kubernetes configs will fail: %v", err)
	}

	opts := topology.TopologyRunOptions{
		Client:    kommonsClient,
		Depth:     topologyRunDepth,
		Namespace: systemTemplate.Namespace,
	}
	if err := topology.SyncComponents(opts, *systemTemplate); err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}