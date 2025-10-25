# TikTok E-commerce Platform

## Project Overview

This is a comprehensive microservices-based e-commerce platform built with modern cloud-native technologies, featuring AI-powered services for enhanced user experience.

## AI-Powered Features

**🤖 Advanced LLM Integration with AWS Bedrock**

### Key Features

1. **Intelligent Customer Support**
   - AI-powered chat service using Claude 3 Sonnet
   - Real-time conversation handling with streaming responses
   - Context-aware product recommendations

2. **Smart Product Search**
   - Vector-based product embedding using Amazon Titan
   - Semantic search capabilities for better product discovery
   - Personalized recommendations based on user behavior

3. **AWS Bedrock Configuration**
   ```bash
   AWS_BEDROCK_MODEL=anthropic.claude-3-sonnet-20240229-v1:0
   AWS_BEDROCK_EMBEDDING_MODEL=amazon.titan-embed-text-v1
   AWS_REGION=us-east-1
   AWS_ACCESS_KEY_ID=your-access-key-id
   AWS_SECRET_ACCESS_KEY=your-secret-access-key
   ```

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

### Service Architecture

```
├── app/
│   ├── cart/           # Shopping cart service
│   ├── checkout/       # Checkout service  
│   ├── email/          # Email notification service
│   ├── gateway/        # API Gateway with load balancing
│   ├── llm/            # AI-powered chat service (AWS Bedrock)
│   ├── order/          # Order management service
│   ├── payment/        # Payment processing service
│   ├── product/        # Product catalog with AI search
│   └── user/           # User management service
├── common/             # Shared libraries and utilities
├── deploy/             # Docker and Kubernetes configurations
├── idl/                # Protocol buffer definitions
└── rpc_gen/            # Generated RPC code
```

### Technology Stack

- **Language**: Go 1.23
- **RPC Framework**: CloudWeGo Kitex
- **API Gateway**: CloudWeGo Hertz
- **AI Services**: AWS Bedrock (Claude 3 Sonnet, Amazon Titan)
- **Databases**: MySQL, MongoDB, Redis
- **Search Engine**: Elasticsearch
- **Message Queue**: Built-in async processing
- **Deployment**: Docker, Kubernetes
- **Monitoring**: Prometheus, Grafana
- **Distributed Tracing**: Jaeger

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

### Prerequisites

1. **Go Environment**
   - Go 1.23+ required
   - Docker and Docker Compose for local development
   - Kubernetes cluster for production deployment

2. **AWS Bedrock Setup**
   - AWS account with Bedrock access
   - IAM permissions: `bedrock:InvokeModel`, `bedrock:InvokeModelWithResponseStream`
   - Model access for Claude 3 Sonnet and Amazon Titan (request in AWS Console)

3. **Infrastructure Requirements**
   - MySQL 8.0+ for transactional data
   - MongoDB for document storage
   - Redis for caching and sessions
   - Elasticsearch for product search

## Features

### Core E-commerce Functionality
- **User Management**: Registration, authentication, profile management
- **Product Catalog**: Advanced search with AI-powered recommendations
- **Shopping Cart**: Real-time cart management and persistence
- **Order Processing**: Complete order lifecycle management
- **Payment Integration**: Secure payment processing
- **Email Notifications**: Automated email communications

### AI-Powered Capabilities
- **Intelligent Chat Support**: Claude 3 powered customer service
- **Smart Product Search**: Vector-based semantic search using Amazon Titan
- **Personalized Recommendations**: ML-driven product suggestions
- **Real-time Streaming**: Live chat with streaming AI responses

## Troubleshooting

### Common Issues

1. **Service Discovery Issues**
   - Verify service registry connectivity
   - Check network policies in Kubernetes
   - Ensure proper service mesh configuration

2. **Database Connection Problems**
   - Validate connection strings and credentials
   - Check database service health
   - Verify network connectivity

3. **AI Service Issues**
   - Confirm AWS Bedrock credentials and permissions
   - Check model availability in your region
   - Verify API rate limits and quotas

### Monitoring and Debugging

- **Service Logs**: `kubectl logs -f deployment/{service-name}`
- **Service Health**: `kubectl get pods -l app={service-name}`
- **Metrics Dashboard**: Access Grafana at `http://localhost:3000`
- **Distributed Tracing**: View traces in Jaeger UI

## Contributing

1. Fork the project
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

