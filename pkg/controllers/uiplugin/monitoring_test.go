package uiplugin

import (
	"testing"

	"gotest.tools/v3/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	uiv1alpha1 "github.com/rhobs/observability-operator/pkg/apis/uiplugin/v1alpha1"
)

var namespace = "openshift-operators"
var name = "monitoring"
var image = "quay.io/monitoring-foo-test:123"

var pluginConfigAll = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			Alertmanager: uiv1alpha1.AlertmanagerReference{
				Url: "https://alertmanager.open-cluster-management-observability.svc:9095",
			},
			ThanosQuerier: uiv1alpha1.ThanosQuerierReference{
				Url: "https://rbac-query-proxy.open-cluster-management-observability.svc:8443",
			},
			Perses: uiv1alpha1.PersesReference{
				Name:      "perses-api-http",
				Namespace: "perses-operator",
			},
		},
	},
}

var pluginConfigPerses = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			Perses: uiv1alpha1.PersesReference{
				Name:      "perses-api-http",
				Namespace: "perses-operator",
			},
		},
	},
}

var pluginConfigACM = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			Alertmanager: uiv1alpha1.AlertmanagerReference{
				Url: "https://alertmanager.open-cluster-management-observability.svc:9095",
			},
			ThanosQuerier: uiv1alpha1.ThanosQuerierReference{
				Url: "https://rbac-query-proxy.open-cluster-management-observability.svc:8443",
			},
		},
	},
}

var pluginConfigThanos = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			ThanosQuerier: uiv1alpha1.ThanosQuerierReference{
				Url: "https://rbac-query-proxy.open-cluster-management-observability.svc:8443",
			},
		},
	},
}

var pluginConfigAlertmanager = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			Alertmanager: uiv1alpha1.AlertmanagerReference{
				Url: "https://alertmanager.open-cluster-management-observability.svc:9095",
			},
		},
	},
}

var pluginMalformed = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type:       "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{},
	},
}

func containsProxy(pluginInfo *UIPluginInfo) (bool, bool, bool) {
	alertmanagerFound, thanosFound, persesFound := false, false, false

	for _, proxy := range pluginInfo.Proxies {
		if proxy.Alias == "alertmanager-proxy" {
			alertmanagerFound = true
		}
		if proxy.Alias == "thanos-proxy" {
			thanosFound = true
		}
		if proxy.Alias == "perses" {
			persesFound = true
		}
	}
	return alertmanagerFound, thanosFound, persesFound
}

func getPluginInfo(plugin *uiv1alpha1.UIPlugin, features []string) (*UIPluginInfo, error) {
	return createMonitoringPluginInfo(plugin, namespace, name, image, features)
}

func TestCreateMonitoringPluginInfo(t *testing.T) {
	t.Run("Test createMonitoringPluginInfo with all monitoring configurations", func(t *testing.T) {
		var features = []string{"perses-dashboards", "acm-alerting"}
		pluginInfo, error := getPluginInfo(pluginConfigAll, features)
		alertmanagerFound, thanosFound, persesFound := containsProxy(pluginInfo)

		assert.Assert(t, alertmanagerFound == true)
		assert.Assert(t, thanosFound == true)
		assert.Assert(t, persesFound == true)
		assert.Assert(t, error == nil)
	})

	t.Run("Test createMonitoringPluginInfo with AMC configuration only", func(t *testing.T) {
		var features = []string{"acm-alerting"}
		pluginInfo, error := getPluginInfo(pluginConfigACM, features)
		alertmanagerFound, thanosFound, persesFound := containsProxy(pluginInfo)

		assert.Assert(t, alertmanagerFound == true)
		assert.Assert(t, thanosFound == true)
		assert.Assert(t, persesFound == false)
		assert.Assert(t, error == nil)

	})

	t.Run("Test createMonitoringPluginInfo with Perses configuration only", func(t *testing.T) {
		var features = []string{"perses-dashboards"}
		pluginInfo, error := getPluginInfo(pluginConfigPerses, features)
		alertmanagerFound, thanosFound, persesFound := containsProxy(pluginInfo)

		assert.Assert(t, error == nil)
		assert.Assert(t, alertmanagerFound == false)
		assert.Assert(t, thanosFound == false)
		assert.Assert(t, persesFound == true)
	})

	t.Run("Test createMonitoringPluginInfo with missing URLs from thanos and alertmanager", func(t *testing.T) {
		var features = []string{"acm-alerting", "perses-dashboards"}

		// this should throw and error because thanosQuerier.URL and alertManager.URL are not set
		pluginInfo, error := getPluginInfo(pluginConfigPerses, features)
		alertmanagerFound, thanosFound, persesFound := containsProxy(pluginInfo)

		assert.Assert(t, error == nil)
		assert.Assert(t, alertmanagerFound == false)
		assert.Assert(t, thanosFound == false)
		assert.Assert(t, persesFound == true)
	})

	t.Run("Test createMonitoringPluginInfo with missing URL from thanos", func(t *testing.T) {
		var features = []string{"acm-alerting", "perses-dashboards"}

		// this should throw and error because thanosQuerier.URL is not set
		pluginInfo, error := getPluginInfo(pluginConfigAlertmanager, features)
		assert.Assert(t, pluginInfo == nil)
		assert.Assert(t, error != nil)
		assert.Equal(t, error.Error(), "thanosQuerier location can not be empty for plugin type monitoring")
	})

	t.Run("Test createMonitoringPluginInfo with missing URL from alertmanager ", func(t *testing.T) {
		var features = []string{"acm-alerting", "perses-dashboards"}

		// this should throw and error because alertManager.URL is not set
		pluginInfo, error := getPluginInfo(pluginConfigThanos, features)
		assert.Assert(t, pluginInfo == nil)
		assert.Assert(t, error != nil)
		assert.Equal(t, error.Error(), "alertmanager location can not be empty for plugin type monitoring")
	})

	t.Run("Test createMonitoringPluginInfo with malform UIPlugin custom resource", func(t *testing.T) {
		var features = []string{"acm-alerting", "perses-dashboards"}

		// this should throw and error because UIPlugin doesn't include alertmanager, thanos, or perses
		pluginInfo, error := getPluginInfo(pluginMalformed, features)
		assert.Assert(t, pluginInfo == nil)
		assert.Assert(t, error != nil)
		assert.Equal(t, error.Error(), "alertmanager location can not be empty for plugin type monitoring")
	})
}
