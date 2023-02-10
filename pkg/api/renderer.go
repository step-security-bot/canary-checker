package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/flanksource/canary-checker/pkg"
	"github.com/flanksource/canary-checker/pkg/topology"
	babel "github.com/jvatic/goja-babel"
	"github.com/labstack/echo/v4"
)

var jsComponentTpl *template.Template

type component struct {
	Name string
	JS   string
}

// GetCustomRenderer returns an application/javascript HTTP response
// with custom components and a registry.
// This registry needs to be used to select custom components
// for rendering of properties and cards.
func GetCustomRenderer(ctx echo.Context) error {
	// 1. Read the template of the topology
	params := topology.NewTopologyParams(ctx.QueryParams())
	results, err := topology.QueryRenderComponents(ctx.Request().Context(), params.ID)
	if err != nil {
		return errorResonse(ctx, err, http.StatusBadRequest)
	}

	// 2. Create a registry of all the components
	var components = make(map[string]component)
	for _, r := range results {
		if err := compileComponents(components, r.Components, false); err != nil {
			return errorResonse(ctx, err, http.StatusInternalServerError)
		}

		if err := compileComponents(components, r.Properties, true); err != nil {
			return errorResonse(ctx, err, http.StatusInternalServerError)
		}
	}

	registryResp, err := renderComponents(components)
	if err != nil {
		return errorResonse(ctx, err, http.StatusInternalServerError)
	}

	return ctx.Stream(http.StatusOK, "application/javascript", registryResp)
}

func compileComponents(output map[string]component, components []pkg.RenderComponent, isProp bool) error {
	if err := babel.Init(len(components)); err != nil {
		return fmt.Errorf("failed to init babel; %w", err)
	}

	for _, c := range components {
		res, err := babel.TransformString(c.JSX, map[string]interface{}{
			"plugins": []string{
				"transform-react-jsx",
				"transform-block-scoping",
			},
		})
		if err != nil {
			return fmt.Errorf("error transforming jsx; %w", err)
		}

		output[componentKey(isProp, c)] = component{
			Name: c.Name,
			JS:   res,
		}
	}

	return nil
}

func componentKey(isProp bool, c pkg.RenderComponent) string {
	prefix := "component"
	if isProp {
		prefix = "property"
	}

	if c.Type != "" {
		return fmt.Sprintf("%s_%s_%s", prefix, c.Type, c.Name)
	}

	return fmt.Sprintf("%s_%s", prefix, c.Name)
}

func renderComponents(components map[string]component) (io.Reader, error) {
	var buf bytes.Buffer
	if err := jsComponentTpl.Execute(&buf, components); err != nil {
		return nil, fmt.Errorf("error generating components; %w", err)
	}

	return &buf, nil
}

const jsComponentRegistryTpl = `
{{range $k, $v := .}}
const {{$k}} = {{$v.JS}}
{{end}}

const componentRegistry = {
	{{range $k, $v := .}}"{{$k}}": {{$k}},
	{{end}}
};
`

func init() {
	tpl, err := template.New("registry").Parse(jsComponentRegistryTpl)
	if err != nil {
		panic(fmt.Errorf("error parsing template jsComponentRegistryTpl. %w", err))
	}

	jsComponentTpl = tpl
}
