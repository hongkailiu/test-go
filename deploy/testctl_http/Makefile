#TODO this is moved from root folder
#need to adjust the file paths below
#
current_oc_context := $(shell oc config current-context)
oc_project := $(shell echo $(current_oc_context) | cut -d "/" -f1)
oc_server := $(shell echo $(current_oc_context) | cut -d "/" -f2)
oc_user := $(shell echo $(current_oc_context) | cut -d "/" -f3)

expected_oc_server := api-starter-us-east-2a-openshift-com:443
expected_oc_user := hongkliu
web_secret_file := /home/hongkliu/repo/me/svt-secret/test_go/web_secret.yaml
grafana_secret_file := /home/hongkliu/repo/me/svt-secret/test_go/grafana_secret.yaml
slack_api_secret_value := $(shell head -n 1 /home/hongkliu/repo/me/svt-secret/test_go/slack_api_secret.txt)

.PHONY : oc-deploy-testctl
oc-deploy-testctl:
	@echo "deploy testctl on openshift starter ... with $(current_oc_context)"
	@echo "oc_project: $(oc_project)"
	@echo "oc_server: $(oc_server)"
	@echo "oc_user: $(oc_user)"
ifeq ($(oc_server),$(expected_oc_server))
	@echo "server match!"
else
	@echo "server do NOT match: exiting ..."
	@echo "expected_oc_server: $(expected_oc_server)"
	false
endif
ifeq ($(oc_user),$(expected_oc_user))
	@echo "user match!"
else
	@echo "user do NOT match: exiting ..."
	@echo "expected_oc_user: $(expected_oc_user)"
	false
endif
	@echo "deploy component http web server ..."
	oc apply -f $(web_secret_file)
	oc apply -f ./deploy/testctl_http/web_deploy.yaml
	oc create configmap -n hongkliu-stage prometheus --from-file=./deploy/testctl_http/prometheus.yml --from-file=alert.rules.yml=./deploy/testctl_http/prometheus_alert.rules.yml --dry-run -o yaml | oc apply -f -
	oc apply -f ./deploy/testctl_http/prometheus_deploy.yaml
	oc apply -f ./deploy/testctl_http/status_deploy.yaml
	oc apply -f $(grafana_secret_file)
	oc create -n hongkliu-stage configmap grafana-config --from-file=./deploy/testctl_http/grafana.ini --dry-run -o yaml | oc apply -f -
	oc create -n hongkliu-stage configmap grafana-datasources --from-file=.datasources.yaml=./deploy/testctl_http/grafana_datasources.yaml --dry-run -o yaml | oc apply -f -
	oc create -n hongkliu-stage configmap grafana-dashboards --from-file=dashboards.yaml=./deploy/testctl_http/grafana_dashboards.yaml --dry-run -o yaml | oc apply -f -
	oc create -n hongkliu-stage configmap grafana-dashboard-test-go --from-file=test-go.json=./deploy/testctl_http/test_go_dashboard.json --dry-run -o yaml | oc apply -f -
	oc apply -f ./deploy/testctl_http/grafana_deploy.yaml
	sed -e "s|{slack_api_secret}|$(slack_api_secret_value)|g" ./deploy/testctl_http/alertmanager.yml > /tmp/alertmanager_decoded.yml
	oc create -n hongkliu-stage configmap alert-manager-config --from-file=alertmanager.yml=/tmp/alertmanager_decoded.yml --from-file=msg.tmpl=./deploy/testctl_http/alert_manager.msg.tmpl --dry-run -o yaml | oc apply -f -
	rm -vf /tmp/alertmanager_decoded.yml
	oc apply -f ./deploy/testctl_http/alert_manager_deploy.yaml
	#https://github.com/kubernetes/kubernetes/issues/13488#issuecomment-481023838
	#kubectl rollout restart #this will be available soon
	@echo "deployed successfully!"
