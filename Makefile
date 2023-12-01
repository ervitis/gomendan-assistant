.PHONY: tests

help: ## Show this help
	@echo "Help"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[93m %s\n", $$1, $$2}'

run: ## Run the application
	go mod download
	go build -o ./out/webcam -ldflags "-s -w" ./cmd/webcam && \
		(open http://localhost:8880/mendan &) && \
		GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_APPLICATION_CREDENTIALS} ./out/webcam
