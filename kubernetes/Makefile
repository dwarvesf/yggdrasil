.PHONY: install apply

install:
	kubectl create -f $(SERVICE)/deployment.yaml
	kubectl create -f $(SERVICE)/service.yaml

apply:
	kubectl apply -f $(SERVICE)/deployment.yaml