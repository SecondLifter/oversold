package adfunc

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/mritd/goadmission/pkg/conf"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func init() {
	register(AdmissionFunc{
		Type: AdmissionTypeMutating,
		Path: "/depreplicas",
		Func: AddReplicas,
	})
}

type ReplicaValue struct {
	Key   string
	Value string
	Name  string
	Env   []corev1.EnvVar
}

func AddReplicas(request *admissionv1.AdmissionRequest) (*admissionv1.AdmissionResponse, error) {
	//获取属性Kind为Node
	namespace := strings.Split(conf.Namespaces, ",")
	env := os.Getenv("env")
	switch request.Kind.Kind {
	case "Deployment":
		var deploy appsv1.Deployment
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
		for _, v := range namespace {
			//klog.Info(v)
			if deploy.Namespace == v {
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
		key := env + ".cloudglab.cn/replicas"
		var ReplicaValue = ReplicaValue{
			Key:   key,
			Value: deploy.Annotations[key],
			Name:  deploy.Name,
			//Env:   env,
		}

		return ReplicaValue.CopyRendering()
	case "StatefulSet":
		var state appsv1.StatefulSet
		err := jsoniter.Unmarshal(request.Object.Raw, &state)
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
		for _, v := range namespace {
			//klog.Info(v)
			if state.Namespace == v {
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
		key := env + ".cloudglab.cn/replicas"
		var ReplicaValue = ReplicaValue{
			Key:   key,
			Value: state.Annotations[key],
			Name:  state.Name,
		}
		return ReplicaValue.CopyRendering()

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

func (A *ReplicaValue) CopyRendering() (*admissionv1.AdmissionResponse, error) {
	//判断是否deployment 副本是否需要扩容
	if A.Value == "" {
		klog.Info("the container is default replicas")
		return &admissionv1.AdmissionResponse{
			Allowed: true,
			Result: &metav1.Status{
				Code:    http.StatusOK,
				Message: "success",
			},
		}, nil
	}

	replicas, _ := strconv.Atoi(A.Value)
	klog.Info("the tag is " + A.Key)
	klog.Info("the " + A.Name + " want have " + A.Value + " replicas")
	var patches []Patch
	// 是否使用非root用户启动
	if conf.RunAsNonRoot {
		patches = []Patch{
			{
				Option: PatchOptionReplace,
				Path:   "/spec/replicas",
				Value:  int32(replicas),
			}, {
				Option: PatchOptionAdd,
				Path:   "/spec/template/spec/securityContext",
				Value: corev1.SecurityContext{
					RunAsNonRoot: &conf.RunAsNonRoot,
					RunAsUser:    &conf.UserID,
					RunAsGroup:   &conf.GroupId,
				},
			},
		}
	} else {
		patches = []Patch{
			{
				Option: PatchOptionReplace,
				Path:   "/spec/replicas",
				Value:  int32(replicas),
			},
		}
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
