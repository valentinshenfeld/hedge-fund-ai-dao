# Makefile for managing AI Agents
# Location: agents/

AGENTS = agent-analyst agent-risk agent-trader manager mcp-server-evm mcp-server-x
DOCKER_IMAGE_PREFIX = hedge-fund
DOCKER_TAG = latest

.PHONY: all test build docker-build deploy clean $(AGENTS)

# Default target: test and build all agents
all: test build

# --- Testing ---

test: $(addprefix test-, $(AGENTS))
	@echo "✅ All tests completed."

# --- Local Build ---

build: $(addprefix build-, $(AGENTS))
	@echo "✅ All binaries built."

# --- Docker Build ---

docker-build: $(addprefix docker-build-, $(AGENTS))
	@echo "✅ All Docker images built."

# --- Deployment ---

deploy: docker-build
	@echo "🚀 Agents are ready for deployment."
	@echo "Note: Current deployment script builds Docker images."
	@echo "To deploy to a cluster, use: make deploy-to-k8s (if implemented)"

# --- Cleanup ---

clean:
	@echo "🧹 Cleaning up binaries..."
	@for agent in $(AGENTS); do \
		rm -rf $$agent/bin; \
	done

# --- Individual Agent Rules ---

define AGENT_RULES
test-$(1):
	@echo "🧪 Testing $(1)..."
	@cd $(1) && go test ./...

build-$(1):
	@echo "🏗️ Building $(1)..."
	@mkdir -p $(1)/bin
	@cd $(1) && go build -o bin/$(1) main.go

docker-build-$(1):
	@echo "🐳 Docker building $(1)..."
	@docker build -t $(DOCKER_IMAGE_PREFIX)-$(1):$(DOCKER_TAG) -f $(1)/Dockerfile ..
endef

# Generate rules for each agent
$(foreach agent,$(AGENTS),$(eval $(call AGENT_RULES,$(agent))))

# Help target
help:
	@echo "Available targets:"
	@echo "  make test          Run tests for all agents"
	@echo "  make build         Build local binaries for all agents"
	@echo "  make docker-build  Build Docker images for all agents"
	@echo "  make deploy        Prepare agents for deployment"
	@echo "  make clean         Remove built binaries"
	@echo "  make test-<name>   Run tests for a specific agent (e.g., test-manager)"
	@echo "  make build-<name>  Build a specific agent binary"
	@echo "  make docker-build-<name> Build a specific agent Docker image"
