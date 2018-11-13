{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 24 -}}
{{- end -}}

{{/*
Expand the service name.
*/}}
{{- define "servicename" -}}
{{- printf "%s" .Values.service.name | trunc 24 -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate application name at 24 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride | trunc 24 -}}
{{- printf "%s-%s" $name .Release.Namespace -}}
{{- end -}}

{{/*
Create a fully qualified app name.
We truncate application name at 24 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "app.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride | trunc 24 -}}
{{- printf "%s-%s" $name "v0" -}}
{{- end -}}
