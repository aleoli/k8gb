/*
Copyright 2021 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/
package metrics

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/AbsaOSS/k8gb/controllers/depresolver"
	"github.com/stretchr/testify/assert"
)

const (
	namespace = "ns"
	gslbName  = "test-gslb"
)

var (
	defaultGslb   = new(k8gbv1beta1.Gslb)
	defaultConfig = depresolver.Config{K8gbNamespace: namespace}
)

func TestMetricsSingletonIsNotNil(t *testing.T) {
	// arrange
	// act
	m := Metrics()
	// assert
	assert.NotNil(t, m)
	assert.Equal(t, DefaultMetricsNamespace, m.config.K8gbNamespace)
}

func TestMetricsSingletonInitTwice(t *testing.T) {
	// arrange
	c1 := &depresolver.Config{K8gbNamespace: "c1"}
	c2 := &depresolver.Config{K8gbNamespace: "c2"}
	// act
	Init(c1)
	Init(c2)
	m := Metrics()
	// assert
	assert.Equal(t, c1.K8gbNamespace, m.config.K8gbNamespace)
}

func TestPrometheusRegistry(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	fieldCnt := reflect.TypeOf(metrics.metrics).NumField()
	// act
	registry := m.registry()
	// assert
	assert.Equal(t, len(registry), fieldCnt, "not all metrics are initialised, check init() function")
	for n, i := range registry {
		assert.True(t, strings.HasPrefix(n, namespace+"_"+gslbSubsystem+"_"))
		assert.NotNil(t, i, n, "is declared but not initialized. Check init() function")
	}
}

func TestMetricsRegister(t *testing.T) {
	// arrange
	m := Metrics()
	// act
	err := m.Register()
	m.Unregister()
	// assert
	assert.NoError(t, err)
}

func TestReconciliationTotal(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	name := fmt.Sprintf("%s_%s_reconciliation_total", namespace, gslbSubsystem)
	cnt1 := testutil.ToFloat64(m.Get(name).AsCounter())
	// act
	m.ReconciliationIncrement()
	// assert
	cnt2 := testutil.ToFloat64(m.Get(name).AsCounter())
	assert.Equal(t, cnt1+1.0, cnt2)
}

func TestHealthyRecords(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	name := fmt.Sprintf("%s_%s_healthy_records", namespace, gslbSubsystem)
	data := map[string][]string{
		"roundrobin.cloud.example.com":      {"10.0.0.1", "10.0.0.2", "10.0.0.3"},
		"roundrobin-test.cloud.example.com": {"10.0.0.4", "10.0.0.5", "10.0.0.6"},
	}
	// act
	cnt1 := testutil.ToFloat64(m.Get(name).AsGaugeVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	m.UpdateHealthyRecordsMetric(defaultGslb, data)
	cnt2 := testutil.ToFloat64(m.Get(name).AsGaugeVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	// assert
	assert.Equal(t, 0.0, cnt1)
	assert.Equal(t, 6.0, cnt2)
}

func TestEmptyHealthyRecords(t *testing.T) {
	// arrange
	var data map[string][]string
	m := newPrometheusMetrics(defaultConfig)
	name := fmt.Sprintf("%s_%s_healthy_records", namespace, gslbSubsystem)
	// act
	cnt1 := testutil.ToFloat64(m.Get(name).AsGaugeVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	m.UpdateHealthyRecordsMetric(defaultGslb, data)
	cnt2 := testutil.ToFloat64(m.Get(name).AsGaugeVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	// assert
	assert.Equal(t, 0.0, cnt1)
	assert.Equal(t, 0.0, cnt2)
}

func TestZoneUpdate(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	name := fmt.Sprintf("%s_%s_zone_update_total", namespace, gslbSubsystem)
	cnt1 := testutil.ToFloat64(m.Get(name).AsCounter())
	// act
	m.ZoneUpdateIncrement()
	// assert
	cnt2 := testutil.ToFloat64(m.Get(name).AsCounter())
	assert.Equal(t, cnt1+1.0, cnt2)
}

func TestUpgradeIngressHost(t *testing.T) {
	// arrange
	name := fmt.Sprintf("%s_%s_ingress_hosts_per_status", namespace, gslbSubsystem)
	m := newPrometheusMetrics(defaultConfig)
	var serviceHealth = map[string]string{
		"roundrobin.cloud.example.com": HealthyStatus,
		"failover.cloud.example.com":   HealthyStatus,
		"unhealthy.cloud.example.com":  UnhealthyStatus,
		"notfound.cloud.example.com":   NotFoundStatus,
	}
	// act
	cntHealthy1 := testutil.ToFloat64(m.Get(name).AsGaugeVec().With(
		prometheus.Labels{"namespace": namespace, "name": gslbName, "status": HealthyStatus}))
	cntUnhealthy1 := testutil.ToFloat64(
		m.Get(name).AsGaugeVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": UnhealthyStatus}))
	cntNotFound1 := testutil.ToFloat64(
		m.Get(name).AsGaugeVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": NotFoundStatus}))
	m.UpdateIngressHostsPerStatusMetric(defaultGslb, serviceHealth)
	cntHealthy2 := testutil.ToFloat64(
		m.Get(name).AsGaugeVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": HealthyStatus}))
	ctnUnhealthy2 := testutil.ToFloat64(
		m.Get(name).AsGaugeVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": UnhealthyStatus}))
	cntNotFound2 := testutil.ToFloat64(
		m.Get(name).AsGaugeVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": NotFoundStatus}))
	// assert
	assert.Equal(t, .0, cntHealthy1)
	assert.Equal(t, .0, cntUnhealthy1)
	assert.Equal(t, .0, cntNotFound1)
	assert.Equal(t, 2., cntHealthy2)
	assert.Equal(t, 1., ctnUnhealthy2)
	assert.Equal(t, 1., cntNotFound2)
}

func TestMain(m *testing.M) {
	defaultGslb.Name = gslbName
	defaultGslb.Namespace = namespace
	m.Run()
}