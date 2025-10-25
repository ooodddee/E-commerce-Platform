# AWS Bedrock 迁移指南

## 概述

本项目已从火山舟大模型服务平台迁移至 AWS Bedrock。以下是主要变化和配置说明。

## 主要变化

### 1. 代码变更
- 移除了 `github.com/cloudwego/eino-ext/components/model/ark` 依赖
- 添加了 AWS SDK v2 依赖：
  - `github.com/aws/aws-sdk-go-v2`
  - `github.com/aws/aws-sdk-go-v2/config`
  - `github.com/aws/aws-sdk-go-v2/service/bedrockruntime`
- 实现了新的 `BedrockChatModel` 替换原有的 `ArkChatModel`

### 2. 环境变量变更
原有火山舟环境变量：
```bash
ARK_CHAT_MODEL=ep-20250201132604-lfwhm
ARK_API_KEY=9727f906-23a1-4f66-902f-d6e1d5d45950
```

新的AWS Bedrock环境变量：
```bash
AWS_BEDROCK_MODEL=anthropic.claude-3-sonnet-20240229-v1:0
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-aws-access-key-id
AWS_SECRET_ACCESS_KEY=your-aws-secret-access-key
```

### 3. 支持的模型
AWS Bedrock 支持多种大模型：
- **Anthropic Claude 3**: `anthropic.claude-3-sonnet-20240229-v1:0`
- **Anthropic Claude 2**: `anthropic.claude-v2`
- **Meta Llama 2**: `meta.llama2-70b-chat-v1`
- **Amazon Titan**: `amazon.titan-text-express-v1`

## 配置说明

### 本地开发环境

1. 复制环境变量文件：
```bash
cp app/llm/.env.example app/llm/.env
```

2. 编辑 `.env` 文件，填入你的AWS凭证：
```bash
AWS_BEDROCK_MODEL=anthropic.claude-3-sonnet-20240229-v1:0
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-actual-access-key-id
AWS_SECRET_ACCESS_KEY=your-actual-secret-access-key
```

### Kubernetes部署

1. 创建AWS凭证Secret：
```bash
# 替换为你的实际凭证
kubectl create secret generic aws-credentials \
  --from-literal=access-key-id=your-actual-access-key-id \
  --from-literal=secret-access-key=your-actual-secret-access-key
```

2. 或者通过文件创建：
```bash
echo -n "your-actual-access-key-id" | base64
echo -n "your-actual-secret-access-key" | base64
```

然后编辑 `deploy/k8s/app.yaml` 中的Secret配置，替换base64编码的值。

### AWS IAM权限

确保你的AWS凭证具有以下权限：
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "bedrock:InvokeModel",
                "bedrock:InvokeModelWithResponseStream"
            ],
            "Resource": "*"
        }
    ]
}
```

## 测试验证

### 1. 构建项目
```bash
cd app/llm
go mod tidy
go build -v -o server
```

### 2. 运行测试
```bash
go test ./biz/service/...
```

### 3. 启动服务
```bash
./server
```

## 故障排除

### 常见问题

1. **认证失败**
   - 检查AWS凭证是否正确
   - 验证IAM权限配置
   - 确认AWS区域设置

2. **模型不可用**
   - 检查模型名称是否正确
   - 验证所选区域是否支持该模型
   - 确认是否已申请模型访问权限

3. **网络连接问题**
   - 检查防火墙设置
   - 验证AWS服务端点连通性

### 日志调试

启用详细日志：
```bash
export AWS_SDK_LOAD_CONFIG=1
export AWS_LOG_LEVEL=debug
```

## 性能优化

1. **区域选择**: 选择距离最近的AWS区域以降低延迟
2. **模型选择**: 根据需求平衡性能和成本
3. **连接池**: AWS SDK 自动管理连接池

## 成本优化

1. 监控API调用量和Token使用量
2. 根据实际需求选择合适的模型
3. 实施请求缓存机制
4. 使用AWS成本管理工具监控费用

## 迁移验证清单

- [ ] 代码编译无错误
- [ ] 单元测试通过
- [ ] LLM服务正常响应
- [ ] 流式响应功能正常
- [ ] 环境变量正确配置
- [ ] K8s部署成功
- [ ] AWS凭证权限验证
- [ ] 端到端测试通过

## 回滚方案

如需回滚到火山舟服务：

1. 恢复 `go.mod` 中的ark依赖
2. 还原 `model.go` 和相关文件
3. 更新环境变量配置
4. 重新部署服务

## 联系支持

如遇到问题，请联系：
- 技术支持团队
- AWS技术支持（如有AWS企业支持计划）