module github.com/mritd/goadmission

go 1.16

require (
	github.com/coreos/bbolt v1.3.2 // indirect
	github.com/coreos/etcd v3.3.13+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/json-iterator/go v1.1.11
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/prometheus/client_golang v0.9.3 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.2 // indirect
	go.uber.org/atomic v1.8.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.18.1
	gopkg.in/resty.v1 v1.12.0 // indirect
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/klog/v2 v2.8.0
)

replace (
	k8s.io/api => k8s.io/api v0.21.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.2
)

// common replace
//replace (
//	k8s.io/api => k8s.io/api v0.21.2
//	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.2
//	k8s.io/apimachinery => k8s.io/apimachinery v0.21.2
//	k8s.io/apiserver => k8s.io/apiserver v0.21.2
//	k8s.io/cli-runtime => k8s.io/cli-runtime v0.21.2
//	k8s.io/client-go => k8s.io/client-go v0.21.2
//	k8s.io/cloud-provider => k8s.io/cloud-provider v0.21.2
//	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.21.2
//	k8s.io/code-generator => k8s.io/code-generator v0.21.2
//	k8s.io/component-base => k8s.io/component-base v0.21.2
//	k8s.io/cri-api => k8s.io/cri-api v0.21.2
//	k8s.io/csi-api => k8s.io/csi-api v0.21.2
//	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.21.2
//	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.21.2
//	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.21.2
//	k8s.io/kube-proxy => k8s.io/kube-proxy v0.21.2
//	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.21.2
//	k8s.io/kubectl => k8s.io/kubectl v0.21.2
//	k8s.io/kubelet => k8s.io/kubelet v0.21.2
//	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.21.2
//	k8s.io/metrics => k8s.io/metrics v0.21.2
//	k8s.io/node-api => k8s.io/node-api v0.21.2
//	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.21.2
//	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.21.2
//	k8s.io/sample-controller => k8s.io/sample-controller v0.21.2
//)
