/*
 * Canary Checker API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1..1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type Property struct {
	Color string `json:"color,omitempty"`
	Headline bool `json:"headline,omitempty"`
	Icon string `json:"icon,omitempty"`
	Label string `json:"label,omitempty"`
	LastTransition string `json:"lastTransition,omitempty"`
	Links []Link `json:"links,omitempty"`
	Max int32 `json:"max,omitempty"`
	Min int32 `json:"min,omitempty"`
	Name string `json:"name,omitempty"`
	Order int32 `json:"order,omitempty"`
	Status string `json:"status,omitempty"`
	// Either text or value is required, but not both.
	Text string `json:"text,omitempty"`
	Tooltip string `json:"tooltip,omitempty"`
	Type_ string `json:"type,omitempty"`
	// e.g. milliseconds, bytes, millicores, epoch etc.
	Unit string `json:"unit,omitempty"`
	Value int32 `json:"value,omitempty"`
}
