package kubernetes

import (
	"time"

	agentpayload "github.com/DataDog/agent-payload/gogen"
	agentkubernetes "github.com/DataDog/datadog-agent/pkg/metadata/kubernetes"
	agentkubelet "github.com/DataDog/datadog-agent/pkg/util/kubernetes/kubelet"
	log "github.com/cihub/seelog"

	"github.com/DataDog/datadog-process-agent/util/cache"
)

const (
	cacheKey          = "kubernetes_meta"
	kubernetesMetaTTL = 3 * time.Minute
)

var lastKubeErr string

// GetKubernetesServices returns a mapping of container ID to list of service names
func GetKubernetesServices() (containerServices map[string][]string) {
	containerServices = make(map[string][]string)

	kubeMeta := getKubernetesMeta()
	if kubeMeta == nil {
		return
	}

	ku, err := agentkubelet.GetKubeUtil()
	if err != nil {
		return
	}
	localPods, err := ku.GetLocalPodList()
	if err != nil {
		log.Errorf("Unable to get local pods from kubelet: %s", err)
		return
	}

	for _, p := range localPods {
		services := findServicesForPod(p, kubeMeta)
		for _, c := range p.Status.Containers {
			containerServices[c.ID] = services
		}
	}

	return
}

func findServicesForPod(pod *agentkubelet.Pod, kubeMeta *agentpayload.KubeMetadataPayload) []string {
	names := make([]string, 0)
	for _, s := range kubeMeta.Services {
		if s.Namespace != pod.Metadata.Namespace {
			continue
		}
		match := true
		for k, search := range s.Selector {
			if v, ok := pod.Metadata.Labels[k]; !ok || v != search {
				match = false
				break
			}
		}
		if match {
			names = append(names, s.Name)
		}
	}
	return names
}

func getKubernetesMeta() (kubeMeta *agentpayload.KubeMetadataPayload) {
	if payload, ok := cache.Get(cacheKey); ok {
		kubeMeta = payload.(*agentpayload.KubeMetadataPayload)
	} else {
		if p, err := agentkubernetes.GetPayload(); err == nil {
			kubeMeta = p.(*agentpayload.KubeMetadataPayload)
			cache.SetWithTTL(cacheKey, kubeMeta, kubernetesMetaTTL)
		} else if err.Error() != lastKubeErr {
			// Swallowing this error for now with an error as it shouldn't block collection.
			log.Errorf("Unable to get kubernetes metadata: %s", err)
			// Only log the same error once to prevent noisy logs.
			lastKubeErr = err.Error()
		}
	}
	return
}
