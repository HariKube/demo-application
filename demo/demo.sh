#! /bin/bash -e

: ${HARIKUBE_URL:=https://harikube.info}

export KIND_CLUSTER=kind
REGISTRY_PASSWORD=$(head -1 demo/credential)
export IMG=controller:dev

exe() {
    local display_cmd="$@"
    display_cmd="${display_cmd//\$\{REGISTRY_PASSWORD\}/****}"
    display_cmd="${display_cmd//$REGISTRY_PASSWORD/****}"
    echo -e "\n\033[1;36m$ $display_cmd\033[0m"
    [[ -n $DEBUG ]] && read -p ""
    eval "$@"
}

exe echo ' \
kubebuilder init --domain example.com --repo example.com/my-project \
kubebuilder create api --group example --version v1 --kind Report --resource=true --controller=true \
kubebuilder create api --group example --version v1 --kind Email --resource=true --controller=true \
kubebuilder create webhook --group example --version v1 --kind Report --programmatic-validation --defaulting \
'
exe kind create cluster
exe kubectl wait --for=condition=Ready node/kind-control-plane --timeout=2m

KINEIP=$(kubectl get no kind-control-plane -o jsonpath='{.status.addresses[0].address}')

exe kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.16.3/cert-manager.yaml
exe kubectl apply -f https://github.com/prometheus-operator/prometheus-operator/releases/download/v0.77.1/stripped-down-crds.yaml
exe kubectl wait -n cert-manager --for=jsonpath='{.status.readyReplicas}'=1 deployment/cert-manager-webhook --timeout=2m

exe kubectl create namespace harikube
exe kubectl create secret generic -n harikube harikube-license --from-file=demo/license
exe kubectl create secret docker-registry harikube-registry-secret \
--docker-server=registry.harikube.info \
--docker-username=harikube \
--docker-password='${REGISTRY_PASSWORD}' \
--namespace=harikube
exe kubectl apply -f ${HARIKUBE_URL}/manifests/harikube-operator-release-v1.0.1.yaml
exe kubectl apply -f ${HARIKUBE_URL}/manifests/harikube-middleware-vcluster-workload-release-v1.0.3.yaml
exe kubectl wait -n harikube --for=jsonpath='{.status.readyReplicas}'=1 deployment/operator-controller-manager --timeout=2m
exe kubectl wait -n harikube --for=jsonpath='{.status.readyReplicas}'=1 statefulset/harikube --timeout=5m

exe "echo '
apiVersion: harikube.info/v1
kind: TopologyConfig
metadata:
  name: topologyconfig-report
  namespace: harikube
spec:
  targetSecret: harikube/topology-config
  backends:
  - name: report
    endpoint: sqlite:///db/report.db?_journal=WAL&cache=shared
    customresource:
      group: example.example.com/v1
      kind: reports
' | kubectl apply -f -
"
sleep 2
exe "kubectl logs -n harikube -l app=harikube | grep 'Backends registered' | tail -1"

exe vcluster connect harikube

exe kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.16.3/cert-manager.yaml
exe kubectl wait -n cert-manager --for=jsonpath='{.status.readyReplicas}'=1 deployment/cert-manager-webhook --timeout=2m

exe make manifests generate install

exe "echo '
apiVersion: example.example.com/v1
kind: Report
metadata:
  name: report-sample
  labels:
    example.example.com/priority: \"3\"
    example.example.com/deadline: \"1755587875\"
spec:
  priority: 3
  details: Sample report details
  deadline: 2025-08-19T09:17:55Z
' | kubectl apply -f -
"

exe 'kubectl get reports \
--selector "example.example.com/deadline>1755587874" \
--selector "example.example.com/priority=3" \
--field-selector "spec.priority=3" \
--field-selector "spec.repoState=Pending"
'

exe kubectl apply -f config/admission/report-admission-config.yaml

exe "echo '
apiVersion: example.example.com/v1
kind: Report
metadata:
  name: report-sample-2
  labels:
    example.example.com/priority: \"4\"
    example.example.com/deadline: \"1755587875\"
spec:
  priority: 4
  details: Sample report details
  deadline: 2025-08-19T09:17:55Z
' | kubectl apply -f -
"

exe 'kubectl get reports report-sample-2 -o yaml'

exe make docker-build
exe kind load docker-image $IMG

exe make manifests generate deploy

exe kubectl wait -n demo-application-system --for=jsonpath='{.status.readyReplicas}'=1 deployment/demo-application-controller-manager --timeout=2m

exe "echo '
apiVersion: example.example.com/v1
kind: Report
metadata:
  name: report-sample-3
spec:
  priority: 2
  details: Sample report details
  deadline: 2025-08-19T09:17:55Z
' | kubectl apply -f -
"

exe 'kubectl get reports --selector "example.example.com/priority=2" -o yaml'

exe "kubectl patch reports report-sample-3 --type='json' -p='[{\"op\":\"replace\", \"path\":\"/spec/repoState\", \"value\":\"Finished\"}]'"

exe kubectl get emails report-sample-3-updated -o yaml

exe echo 'Do you want more?'

exe "echo '
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: serverless-kube-watch-trigger-emails-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: serverless-kube-watch-trigger-emails-role
subjects:
- kind: ServiceAccount
  name: serverless-kube-watch-trigger-controller-manager
  namespace: serverless-kube-watch-trigger-system
' | kubectl apply -f -
"
exe "echo '
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: serverless-kube-watch-trigger-emails-role
rules:
- apiGroups:
  - example.example.com
  resources:
  - emails
  verbs:
  - get
  - list
  - watch
' | kubectl apply -f -
"
exe kubectl apply -f https://github.com/HariKube/serverless-kube-watch-trigger/releases/download/beta-v1.0.0-7/bundle.yaml

exe kubectl wait -n serverless-kube-watch-trigger-system --for=jsonpath='{.status.readyReplicas}'=1 deployment/serverless-kube-watch-trigger-controller-manager --timeout=2m

exe "echo '
apiVersion: triggers.harikube.info/v1
kind: HTTPTrigger
metadata:
  name: demo-application-email
  namespace: demo-application-system
spec:
  resource:
    apiVersion: example.example.com/v1
    kind: Email
  eventTypes:
    - ADDED
  url:
    service:
      name: gateway
      namespace: openfaas
      portName: http
      scheme: http
      uri:
        static: /function/email
  method: POST
  body:
    contentType: application/json
    template: |
      {{ toJson . }}
  delivery:
    timeout: 10s
    retries: 3
' | kubectl apply -f -
"

exe helm repo add openfaas https://openfaas.github.io/faas-netes/
exe helm repo update
exe kubectl apply -f https://raw.githubusercontent.com/openfaas/faas-netes/master/namespaces.yml
exe helm upgrade openfaas --install openfaas/openfaas \
--namespace openfaas \
--set basic_auth=true \
--set functionNamespace=demo-application-system \
--set serviceType=NodePort \
--set gateway.nodePort=32767
exe kubectl wait -n openfaas --for=jsonpath='{.status.readyReplicas}'=1 deployment/gateway --timeout=2m

OPENFAASPWD=$(kubectl get secret -n openfaas basic-auth -o jsonpath='{.data.basic-auth-password}'| base64 -d)

pushd function
exe faas-cli template store pull python3-http
exe faas-cli build
exe faas-cli push
exe faas-cli login --password ${OPENFAASPWD} --gateway http://${KINEIP}:32767
exe faas-cli deploy --gateway http://${KINEIP}:32767
popd

exe 'kubectl patch deployment email -n demo-application-system --type=json -p='\''[{"op": "replace", "path": "/spec/template/spec/serviceAccountName", "value": "demo-application-controller-manager"}]'\'''
exe kubectl rollout status deployment/email -n demo-application-system --timeout=1m

exe kubectl delete reports report-sample-3

exe kubectl get emails report-sample-3-deleted -o yaml
