LAST_TAG := $(shell go run $(ROOT_DIR)/pkg/tools/versiongen -file $(VERSION_FILE) -mode=last)
REMOTE_IMAGE ?= $(DOCKER_REPO):$(LAST_TAG)

deploy-vps:
	@echo "Uploading .env to VPS..."
	scp .env $(VPS_USER)@$(VPS_HOST):/root/$(PROJECT_NAME).env

	@echo "Deploying image $(REMOTE_IMAGE) to $(VPS_HOST)..."
	@ssh $(VPS_USER)@$(VPS_HOST) "\
		docker pull $(REMOTE_IMAGE) && \
		(docker stop $(PROJECT_NAME) || true) && \
		(docker rm $(PROJECT_NAME) || true) && \
		docker run -d --name $(PROJECT_NAME) \
        			--env-file /root/$(PROJECT_NAME).env \
        			-p $(REST_PORT):$(REST_PORT) -p $(GRPC_PORT):$(GRPC_PORT) \
        			$(REMOTE_IMAGE)"
	@echo "Deployed $(REMOTE_IMAGE) to $(VPS_USER)@$(VPS_HOST)"