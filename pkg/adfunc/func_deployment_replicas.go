package adfunc

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"strconv"
)

func init() {
	register(AdmissionFunc{
		Type: AdmissionTypeMutating,
		Path: "/depreplicas",
		Func: AddReplicas,
	})
}

func AddReplicas(request *admissionv1.AdmissionRequest) (*admissionv1.AdmissionResponse, error) {
	//获取属性Kind为Node
	switch request.Kind.Kind {
	case "Deployment":
		var deploy appsv1.Deployment
		env := os.Getenv("env")
		//klog.Info(env)
		err := jsoniter.Unmarshal(request.Object.Raw, &deploy)
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
		//判断是否deployment 副本是否需要扩容
		key := env + ".cloudglab.cn/replicas"
		//klog.Info(key)
		//klog.Info("key 值为 " + deploy.Annotations[key])
		if deploy.Annotations[key] != "" {
			repicas, _ := strconv.Atoi(deploy.Annotations[key])
			klog.Info("the tag is " + key)
			klog.Info("the deployment " + deploy.Name + "want have " + deploy.Annotations[key] + " replicas")
			patches := []Patch{
				{
					Option: PatchOptionReplace,
					Path:   "/spec/replicas",
					Value:  int32(repicas),
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
		} else if int(*deploy.Spec.Replicas) > 1 {
			klog.Info("the deployment is default replicas")
			return &admissionv1.AdmissionResponse{
				Allowed: true,
				Result: &metav1.Status{
					Code:    http.StatusOK,
					Message: "success",
				},
			}, nil
		}
		patches := []Patch{
			{
				Option: PatchOptionReplace,
				Path:   "/spec/replicas",
				Value:  int32(2),
			},
		}
		patch, err := jsoniter.Marshal(patches)
		if err != nil {
			errMsg := fmt.Sprintf("[route.Mutating]/AddReplicas: failed to unmarshal object: %v", err)
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
