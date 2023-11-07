/*
Copyright 2020 The SuperEdge Authors.

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

package cache

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"gotest.tools/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

func TestWatchDecode(t *testing.T) {
	str := `{"type":"ADDED","object":{"kind":"EndpointSlice","apiVersion":"discovery.k8s.io/v1","metadata":{"name":"externalnode-proxy-l28pc","generateName":"externalnode-proxy-","namespace":"kube-system","uid":"19c23c9e-58c6-4676-ab46-ce85c96ce2d9","resourceVersion":"5493506153","generation":2,"creationTimestamp":"2023-11-06T08:58:24Z","labels":{"endpointslice.kubernetes.io/managed-by":"endpointslice-controller.k8s.io","kubernetes.io/service-name":"externalnode-proxy"},"annotations":{"endpoints.kubernetes.io/last-change-trigger-time":"2023-11-06T08:58:42Z"},"ownerReferences":[{"apiVersion":"v1","kind":"Service","name":"externalnode-proxy","uid":"28540fa8-22cc-4f07-8393-81bacc1ee305","controller":true,"blockOwnerDeletion":true}],"managedFields":[{"manager":"kube-controller-manager","operation":"Update","apiVersion":"discovery.k8s.io/v1","time":"2023-11-06T08:58:42Z","fieldsType":"FieldsV1","fieldsV1":{"f:addressType":{},"f:endpoints":{},"f:metadata":{"f:annotations":{".":{},"f:endpoints.kubernetes.io/last-change-trigger-time":{}},"f:generateName":{},"f:labels":{".":{},"f:endpointslice.kubernetes.io/managed-by":{},"f:kubernetes.io/service-name":{}},"f:ownerReferences":{".":{},"k:{\"uid\":\"28540fa8-22cc-4f07-8393-81bacc1ee305\"}":{}}},"f:ports":{}}}]},"addressType":"IPv4","endpoints":[{"addresses":["172.24.0.7"],"conditions":{"ready":true,"serving":true,"terminating":false},"targetRef":{"kind":"Pod","namespace":"kube-system","name":"externalnode-proxy-56c4d977d8-zxknh","uid":"4d8cffb5-5c4c-4604-9395-b4efe905ddf0"},"nodeName":"10.0.0.29","zone":"540001"}],"ports":[{"name":"tcp-80","protocol":"TCP","port":80},{"name":"tcp-443","protocol":"TCP","port":443}]}}`
	//str := `{"Type":"ADDED","Object":{"kind":"Service","apiVersion":"v1","metadata":{"name":"edge-kube-dns","namespace":"edge-system","uid":"9f0975b1-6a99-4804-a7d8-6721ee112328","resourceVersion":"5493396552","creationTimestamp":"2023-11-06T08:52:44Z","labels":{"app.kubernetes.io/managed-by":"Helm","k8s-app":"edge-kube-dns","kubernetes.io/cluster-service":"true","kubernetes.io/name":"CoreDNS"},"annotations":{"meta.helm.sh/release-name":"externaledge","meta.helm.sh/release-namespace":"kube-system","prometheus.io/port":"9153","prometheus.io/scrape":"true"},"managedFields":[{"manager":"tke-application-controller","operation":"Update","apiVersion":"v1","time":"2023-11-06T08:52:44Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:meta.helm.sh/release-name":{},"f:meta.helm.sh/release-namespace":{},"f:prometheus.io/port":{},"f:prometheus.io/scrape":{}},"f:labels":{".":{},"f:app.kubernetes.io/managed-by":{},"f:k8s-app":{},"f:kubernetes.io/cluster-service":{},"f:kubernetes.io/name":{}}},"f:spec":{"f:internalTrafficPolicy":{},"f:ports":{".":{},"k:{\"port\":53,\"protocol\":\"TCP\"}":{".":{},"f:name":{},"f:port":{},"f:protocol":{},"f:targetPort":{}},"k:{\"port\":53,\"protocol\":\"UDP\"}":{".":{},"f:name":{},"f:port":{},"f:protocol":{},"f:targetPort":{}},"k:{\"port\":9153,\"protocol\":\"TCP\"}":{".":{},"f:name":{},"f:port":{},"f:protocol":{},"f:targetPort":{}}},"f:sessionAffinity":{},"f:type":{}}}}]},"spec":{"ports":[{"name":"dns","protocol":"UDP","port":53,"targetPort":60053},{"name":"dns-tcp","protocol":"TCP","port":53,"targetPort":60053},{"name":"metrics","protocol":"TCP","port":9153,"targetPort":9153}],"clusterIP":"172.24.253.228","clusterIPs":["172.24.253.228"],"type":"ClusterIP","sessionAffinity":"None","ipFamilies":["IPv4"],"ipFamilyPolicy":"SingleStack","internalTrafficPolicy":"Cluster"},"status":{"loadBalancer":{}}}}`
	buff := bytes.NewBufferString(str)

	decoder := getWatchDecoder(ioutil.NopCloser(buff))
	eventType, obj, err := decoder.Decode()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", eventType)
	svc := &v1.Service{}
	runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, svc)
	t.Logf("new watch event, type=%s, svc=%v", eventType, svc)
	assert.Equal(t, eventType, watch.Added)
	event := &watch.Event{
		Type:   eventType,
		Object: svc,
	}
	data, err := json.Marshal(event)
	if err != nil {
		t.Logf(" error=%s", err.Error())
	}
	t.Logf("data : %v", string(data))

	accessor := meta.NewAccessor()
	kind, err := accessor.Kind(obj)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, kind, "Service")
}
