package adfunc

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"net/http"
	"strings"
)

func init() {
	register(AdmissionFunc{
		Type: AdmissionTypeValidating,
		Path: "/check-deploy",
		Func: checkDeploy,
	})
}

func checkDeploy(request *admissionv1.AdmissionRequest) (*admissionv1.AdmissionResponse, error) {
	switch request.Kind.Kind {
	case "Deployment":
		var deploy appsv1.Deployment
		err := jsoniter.Unmarshal(request.Object.Raw, &deploy)
		if err != nil {
			errMsg := fmt.Sprintf("[route.Validating] /check-deploy: failed to unmarshal object: %v", err)
			klog.Error(errMsg)
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Code:    http.StatusBadRequest,
					Message: errMsg,
				},
			}, nil
		}
		//namespace
		if len(deploy.Namespace) == 0 {
			klog.Error(deploy.Name + "Namespace is null")
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Code:    http.StatusBadRequest,
					Message: deploy.Name + "Namespace is null",
				},
			}, nil
		}
		if deploy.Namespace == "kube-system" || deploy.Namespace == "oversold" || deploy.Namespace == "kubeapps" || deploy.Namespace == "kubeapps2" {
			return &admissionv1.AdmissionResponse{
				Allowed: true,
				Result: &metav1.Status{
					Code:    http.StatusOK,
					Message: "success",
				},
			}, nil
		}

		//Tolerations
		if len(deploy.Spec.Template.Spec.Tolerations) == 0 {
			klog.Error(deploy.Name + "Tolerations is null ")
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Code:    http.StatusBadRequest,
					Message: deploy.Name + "Tolerations is null ",
				},
			}, nil
		}
		//Probe and Resources
		for _, v := range deploy.Spec.Template.Spec.Containers {
			//ReadinessProbe or LivenessProbe

			if v.Name != "fluent-bit-container" {
				if v.ReadinessProbe == nil || v.LivenessProbe == nil {
					klog.Warningf(deploy.Name + " ReadinessProbe or LivenessProbe is null ")
					return &admissionv1.AdmissionResponse{
						Allowed: false,
						Result: &metav1.Status{
							Code:    http.StatusBadRequest,
							Message: deploy.Name + v.Name + " ReadinessProbe or LivenessProbe is null ",
						},
					}, nil
				}
			}
			//ignore web
			if cheklego(deploy.Name) != "lego" {
				if len(v.Resources.Limits) == 0 || len(v.Resources.Requests) == 0 {
					klog.Error(deploy.Name + "Resource is nil")
					return &admissionv1.AdmissionResponse{
						Allowed: false,
						Result: &metav1.Status{
							Code:    http.StatusBadRequest,
							Message: deploy.Name + "Resources is null ",
						},
					}, nil
				}
			}

		}

		return &admissionv1.AdmissionResponse{
			Allowed: true,
			Result: &metav1.Status{
				Code:    http.StatusOK,
				Message: "success",
			},
		}, nil
	default:
		errMsg := fmt.Sprintf("[route.Validating] /check-deploy: received wrong kind request: %s, Only support Kind: Deployment", request.Kind.Kind)
		klog.Error(errMsg)
		//logger.Error(errMsg)
		return &admissionv1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Code:    http.StatusForbidden,
				Message: errMsg,
			},
		}, nil
	}
}

func cheklego(name string) string {
	name1 := strings.Split(name, "-")
	return name1[len(name1)-1]
}
