# TikTok E-commerce Platform

## Project Overview

This is a microservices-based e-commerce platform with multiple core services.

## Important Update: LLM Service Migration

**🔄 Successfully migrated from Volcengine ARK Platform to AWS Bedrock**

### Key Changes

1. **LLM Service Provider Migration**
   - From: Volcengine ARK Platform
   - To: AWS Bedrock

2. **Environment Variables Update**
   ```bash
   # Old Configuration (Removed)
   ARK_CHAT_MODEL=
   ARK_API_KEY=
   
   # New Configuration
   AWS_BEDROCK_MODEL=anthropic.claude-3-sonnet-20240229-v1:0
   AWS_REGION=us-east-1
   AWS_ACCESS_KEY_ID=your-access-key-id
   AWS_SECRET_ACCESS_KEY=your-secret-access-key
   ```

3. **Technology Stack Updates**
   - Added: AWS SDK v2
   - Removed: Volcengine SDK
   - New: Custom Bedrock Chat Model Implementation

### Quick Start

1. **Configure Environment Variables**
   ```bash
   cp app/llm/.env.example app/llm/.env
   # Edit .env file and add your AWS credentials
   ```

2. **Build and Run**
   ```bash
   cd app/llm
   go mod tidy
   go build -o server
   ./server
   ```

3. **Kubernetes Deployment**
   ```bash
   # Create AWS credentials secret
   kubectl create secret generic aws-credentials \
     --from-literal=access-key-id=your-access-key-id \
     --from-literal=secret-access-key=your-secret-access-key
   
   # Deploy application
   kubectl apply -f deploy/k8s/
   ```

### Detailed Documentation

For complete migration guide and configuration instructions, see: [AWS_BEDROCK_MIGRATION_GUIDE.md](./AWS_BEDROCK_MIGRATION_GUIDE.md)

### Service Architecture

```
├── app/
│   ├── cart/           # Shopping cart service
│   ├── checkout/       # Checkout service  
│   ├── email/          # Email service
│   ├── gateway/        # API Gateway
│   ├── llm/            # LLM service (migrated to AWS Bedrock)
│   ├── order/          # Order service
│   ├── payment/        # Payment service
│   ├── product/        # Product service
│   └── user/           # User service
├── common/             # Common libraries
├── deploy/             # Deployment configurations
├── idl/                # Interface definitions
└── rpc_gen/            # RPC generated code
```

### Technology Stack

- **Language**: Go
- **RPC Framework**: Kitex
- **Gateway**: Hertz
- **Databases**: MySQL, MongoDB, Redis
- **LLM**: AWS Bedrock (Claude 3)
- **Deployment**: Docker, Kubernetes
- **Monitoring**: Prometheus, Grafana
- **Tracing**: Jaeger

## Development Guide

### Local Development

1. Start infrastructure services:
   ```bash
   docker-compose up -d mysql redis mongodb
   ```

2. Configure environment variables (`.env` files in each service directory)

3. Start individual services:
   ```bash
   cd app/{service-name}
   go run .
   ```

### Deployment

#### Docker Build
```bash
# Build all services
make docker-build

# Build specific service
docker build -f deploy/Dockerfile.svc --build-arg SVC=llm -t llm:latest .
```

#### Kubernetes Deployment
```bash
# Deploy to Kind cluster
make deploy-kind

# Direct deployment
kubectl apply -f deploy/k8s/
```

## Environment Configuration

### AWS Bedrock Requirements

1. **AWS Account and Permissions**
   - AWS account required
   - IAM permissions: `bedrock:InvokeModel`, `bedrock:InvokeModelWithResponseStream`
   - Model access permissions (must be requested in AWS Console)

2. **Supported AWS Regions**
   - us-east-1 (N. Virginia)
   - us-west-2 (Oregon)  
   - eu-west-1 (Ireland)

3. **Supported Models**
   - Anthropic Claude 3 Sonnet
   - Anthropic Claude 2
   - Amazon Titan
   - Meta Llama 2

## Troubleshooting

### LLM Service Issues

1. **AWS Authentication Failed**
   - Check AWS credentials configuration
   - Verify IAM permissions
   - Confirm AWS region settings

2. **Model Access Denied**
   - Request model access in AWS Console
   - Verify model name and region support

3. **Network Connection Issues**
   - Check VPC and security group configuration
   - Verify internet connectivity

### General Issues

- View service logs: `kubectl logs -f deployment/{service-name}`
- Check service status: `kubectl get pods`
- Verify configuration: Check ConfigMap and Secret

## Contributing

1. Fork the project
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

[Add license information]