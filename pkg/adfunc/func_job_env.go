package adfunc

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/mritd/goadmission/pkg/conf"
	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"strings"
)

func init() {
	register(AdmissionFunc{
		Type: AdmissionTypeMutating,
		Path: "/envs",
		Func: AddEnvs,
	})
}

func AddEnvs(request *admissionv1.AdmissionRequest) (*admissionv1.AdmissionResponse, error) {
	//获取属性Kind为Node
	namespace := strings.Split(conf.Namespaces, ",")
	switch request.Kind.Kind {
	case "Job":
		var job v1.Job
		err := jsoniter.Unmarshal(request.Object.Raw, &job)
		if err != nil {
			errMsg := fmt.Sprintf("[route.Mutating] /AddReplicas: failed to unmarshal object: %v", err)
			klog.Error(errMsg)
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Code:    http.StatusBadRequest,
					Message: errMsg,
				},
			}, nil
		}
		//判断是否忽略相关namespace
		for _, v := range namespace {
			//klog.Info(v)
			if job.Namespace == v {
				klog.Info("the namespace is in ignore namespace ignoring" + v)
				return &admissionv1.AdmissionResponse{
					Allowed: true,
					Result: &metav1.Status{
						Code:    http.StatusOK,
						Message: "success",
					},
				}, nil
			}
		}

		return ReplaceJobEnvs(job)

	default:
		errMsg := fmt.Sprintf("[route.Mutating] /AddReplicas: received wrong kind request: %s, Only support Kind: Deployment", request.Kind.Kind)
		logger.Error(errMsg)
		return &admissionv1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Code:    http.StatusForbidden,
				Message: errMsg,
			},
		}, nil
	}

}

func ReplaceJobEnvs(job v1.Job) (*admissionv1.AdmissionResponse, error) {
	value := os.Getenv("env")
	env := GetJobEnvs(job)
	tmp := corev1.EnvVar{
		Name:  "ENV",
		Value: value,
	}
	env = append(env, tmp)
	patches := []Patch{
		{
			Option: PatchOptionAdd,
			Path:   "/spec/template/spec/containers/0/env",
			Value:  env,
		},
	}
	patch, err := jsoniter.Marshal(patches)
	if err != nil {
		errMsg := fmt.Sprintf("[route.Mutating] /AddReplicas: failed to unmarshal object: %v", err)
		klog.Error(errMsg)
	}
	return &admissionv1.AdmissionResponse{
		Allowed:   true,
		Patch:     patch,
		PatchType: JSONPatch(),
		Result: &metav1.Status{
			Code:    http.StatusOK,
			Message: "success",
		},
	}, nil

}

func GetJobEnvs(job v1.Job) []corev1.EnvVar {
	for k, v := range job.Spec.Template.Spec.Containers {
		if k == 0 {
			klog.Info(v.Env)
			return v.Env
		}
	}
	return nil
}
