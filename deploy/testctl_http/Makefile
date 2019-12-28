#TODO this is moved from root folder
#need to adjust the file paths below
#
current_oc_context := $(shell oc config current-context)
oc_project := $(shell echo $(current_oc_context) | cut -d "/" -f1)
oc_server := $(shell echo $(current_oc_context) | cut -d "/" -f2)
oc_user := $(shell echo $(current_oc_context) | cut -d "/" -f3)

expected_oc_server := api-us-west-1-starter-openshift-online-com:6443
expected_oc_user := hongkliu

secret_path := $(HOME)/repo/me/svt-secret
web_secret_file := $(secret_path)/test_go/web_secret.yaml
grafana_secret_file := $(secret_path)/test_go/grafana_secret.yaml
slack_api_secret_value := $(shell head -n 1 $(secret_path)/test_go/slack_api_secret.txt)


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
	oc apply -f ./web_deploy.yaml
	oc create configmap -n hongkliu-run prometheus --from-file=./prometheus.yml --from-file=alert.rules.yml=./prometheus_alert.rules.yml --dry-run -o yaml | oc apply -f -
	oc apply -f ./prometheus_deploy.yaml
	oc apply -f ./status_rbac.yaml
	oc apply -f ./status_deploy.yaml
	oc apply -f $(grafana_secret_file)
	oc create -n hongkliu-run configmap grafana-config --from-file=./grafana.ini --dry-run -o yaml | oc apply -f -
	oc create -n hongkliu-run configmap grafana-datasources --from-file=.datasources.yaml=./grafana_datasources.yaml --dry-run -o yaml | oc apply -f -
	oc create -n hongkliu-run configmap grafana-dashboards --from-file=dashboards.yaml=./grafana_dashboards.yaml --dry-run -o yaml | oc apply -f -
	oc create -n hongkliu-run configmap grafana-dashboard-test-go --from-file=test-go.json=./test_go_dashboard.json --dry-run -o yaml | oc apply -f -
	oc apply -f ./grafana_deploy.yaml
	@sed -e "s|{slack_api_secret}|$(slack_api_secret_value)|g" ./alertmanager.yml > /tmp/alertmanager_decoded.yml
	oc create -n hongkliu-run configmap alert-manager-config --from-file=alertmanager.yml=/tmp/alertmanager_decoded.yml --from-file=msg.tmpl=./alert_manager.msg.tmpl --dry-run -o yaml | oc apply -f -
	rm -vf /tmp/alertmanager_decoded.yml
	oc apply -f ./alert_manager_deploy.yaml
	##TODO bc with sourceStrategy
	#oc apply -f ./test-go_bc.yaml
	#https://github.com/kubernetes/kubernetes/issues/13488#issuecomment-481023838
	#kubectl rollout restart #this will be available soon
	@echo "deployed successfully!"
