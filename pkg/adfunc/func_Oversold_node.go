package adfunc

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	register(AdmissionFunc{
		Type: AdmissionTypeMutating,
		Path: "/oversold",
		Func: oversold,
	})
}

func oversold(request *admissionv1.AdmissionRequest) (*admissionv1.AdmissionResponse, error) {
	//获取属性Kind为Node
	switch request.Kind.Kind {
	case "Node":
		node := v1.Node{}
		if err := jsoniter.Unmarshal(request.Object.Raw, &node); err != nil {
			errMsg := fmt.Sprintf("[route.Mutating] /oversold: failed to unmarshal object: %v", err)
			klog.Error(errMsg)
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Code:    http.StatusBadRequest,
					Message: errMsg,
				},
			}, nil
		}
		//判断是否该节点需要进行超售
		/*
			if node.Labels["Kubernetes.io/oversold"] =="auto"{
				klog.Info(request.UserInfo.Username +"的label为kubernetes.io/oversold:" +node.Labels["kubernetes.io/oversold"])
				klog.Info(request.UserInfo.Username +"===================该节点开启自动超售========================")

				patches := []Patch{
					{
						Option: PatchOptionReplace,
						Path:   "/status/allocatable/cpu",
						//Value: "32",
						Value:  overcpu(Quantitytostring(node.Status.Allocatable.Cpu()),strconv.Itoa(cpuutil(node.Name))),
					},
					{
						Option: PatchOptionReplace,
						Path:   "/status/allocatable/memory",
						//Value:  "134217728Ki",
						Value: overmem(Quantitytostring(node.Status.Allocatable.Memory()),strconv.Itoa(memutil(node.Name))),
					},
				}
				//klog.Info(overcpu(Quantitytostring(node.Status.Allocatable.Cpu()),node.Labels["kubernetes.io/overcpu"]))
				//klog.Info(overmem(Quantitytostring(node.Status.Allocatable.Memory()),node.Labels["kubernetes.io/overmem"]))
				patch, err := jsoniter.Marshal(patches)
				if err != nil {
					errMsg := fmt.Sprintf("[route.Mutating] /oversold: failed to marshal patch: %v", err)
					logger.Error(errMsg)
					return &admissionv1.AdmissionResponse{
						Allowed: false,
						Result: &metav1.Status{
							Code:    http.StatusInternalServerError,
							Message: errMsg,
						},
					}, nil
				}
				logger.Infof("[route.Mutating] /oversold: patches: %s", string(patch))
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
		*/
		if node.Labels["kubernetes.io/oversold"] != "oversold" {

			return &admissionv1.AdmissionResponse{
				Allowed:   true,
				PatchType: JSONPatch(),
				Result: &metav1.Status{
					Code:    http.StatusOK,
					Message: "节点无需超售",
				},
			}, nil
		}

		klog.Info(request.UserInfo.Username + "的label为kubernetes.io/oversold:" + node.Labels["kubernetes.io/oversold"])
		klog.Info(request.UserInfo.Username + "===================该节点允许超售========================")
		patches := []Patch{
			{
				Option: PatchOptionReplace,
				Path:   "/status/allocatable/cpu",
				//Value: "32",
				Value: overcpu(Quantitytostring(node.Status.Allocatable.Cpu()), node.Labels["kubernetes.io/overcpu"]),
			},
			{
				Option: PatchOptionReplace,
				Path:   "/status/allocatable/memory",
				//Value:  "134217728Ki",
				Value: overmem(Quantitytostring(node.Status.Allocatable.Memory()), node.Labels["kubernetes.io/overmem"]),
			},
		}
		klog.Info(overcpu(Quantitytostring(node.Status.Allocatable.Cpu()), node.Labels["kubernetes.io/overcpu"]))
		klog.Info(overmem(Quantitytostring(node.Status.Allocatable.Memory()), node.Labels["kubernetes.io/overmem"]))
		patch, err := jsoniter.Marshal(patches)
		if err != nil {
			errMsg := fmt.Sprintf("[route.Mutating] /oversold: failed to marshal patch: %v", err)
			logger.Error(errMsg)
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Code:    http.StatusInternalServerError,
					Message: errMsg,
				},
			}, nil
		}
		logger.Infof("[route.Mutating] /oversold: patches: %s", string(patch))
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
		errMsg := fmt.Sprintf("[route.Mutating] /oversold: received wrong kind request: %s, Only support Kind: Deployment", request.Kind.Kind)
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

//*resource.Quantity类型转string
func Quantitytostring(r *resource.Quantity) string {
	return fmt.Sprint(r)
}

//cpu 超售计算
func overcpu(cpu, multiple string) string {
	a, _ := strconv.Atoi(cpu)
	if multiple == "" {
		multiple = "1"
	}
	b, _ := strconv.Atoi(multiple)
	return strconv.Itoa(a * b)
}

// mem 超售计算
func overmem(mem, multiple string) string {
	a, err := strconv.Atoi(strings.Trim(mem, "Ki"))
	if err != nil {
		klog.Error("--------内存超售计算失败----------")
		klog.Error(err)
		return "1"
	}
	if multiple == "" {
		multiple = "1"
	}
	b, _ := strconv.Atoi(multiple)
	c := a * b
	return strconv.Itoa(c) + "Ki"
}

////实际使用率计算
//func cpuutil(name string) int {
//	body :=conntroller.PromPost("sum(kube_pod_container_resource_requests_cpu_cores{node=~\""+ name+"\"})/sum (rate (container_cpu_usage_seconds_total{instance=~\""+name +"\"}[2m]))")
//	cpu:=new(conntroller.AutoGenerated)
//	json.Unmarshal(body,&cpu)
//	for _,v :=range cpu.Data.Result{
//		a,_ :=v.Value[1].(float32)
//		fmt.Println(a)
//		if a <1 {
//			return 1
//		}else if a > 1.5 && a < 2 {
//				return int(a) + 1
//			} else if a >=2{
//				return int(a)
//		}
//	}
//	return 1
//}
//
//func memutil(name string) int {
//	body :=conntroller.PromPost("sum(kube_pod_container_resource_requests_memory_bytes{node=~\""+ name+"\"})/sum (rate (container_memory_working_set_bytes{instance=~\""+name +"\"}[2m]))")
//	cpu:=new(conntroller.AutoGenerated)
//	json.Unmarshal(body,&cpu)
//	for _,v :=range cpu.Data.Result{
//		a,_ :=v.Value[1].(int)
//		fmt.Println(a)
//		if a >=2 {
//			return a
//		}
//	}
//	return 1
//
//}
