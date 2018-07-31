// +build !ignore_autogenerated

/*
Copyright 2018 The Kmodules Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package api

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"kmodules.xyz/monitoring-agent-api/api.AgentSpec":      schema_kmodulesxyz_monitoring_agent_api_api_AgentSpec(ref),
		"kmodules.xyz/monitoring-agent-api/api.PrometheusSpec": schema_kmodulesxyz_monitoring_agent_api_api_PrometheusSpec(ref),
	}
}

func schema_kmodulesxyz_monitoring_agent_api_api_AgentSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Properties: map[string]spec.Schema{
					"agent": {
						SchemaProps: spec.SchemaProps{
							Description: "Valid values: coreos-prometheus-operator",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"prometheus": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("kmodules.xyz/monitoring-agent-api/api.PrometheusSpec"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"kmodules.xyz/monitoring-agent-api/api.PrometheusSpec"},
	}
}

func schema_kmodulesxyz_monitoring_agent_api_api_PrometheusSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Properties: map[string]spec.Schema{
					"port": {
						SchemaProps: spec.SchemaProps{
							Description: "Port number for the exporter side car.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"namespace": {
						SchemaProps: spec.SchemaProps{
							Description: "Namespace of Prometheus. Service monitors will be created in this namespace.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"labels": {
						SchemaProps: spec.SchemaProps{
							Description: "Labels are key value pairs that is used to select Prometheus instance via ServiceMonitor labels.",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"interval": {
						SchemaProps: spec.SchemaProps{
							Description: "Interval at which metrics should be scraped",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{},
	}
}
