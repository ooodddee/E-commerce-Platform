package embedding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// BedrockEmbedder AWS Bedrock Embedding实现
type BedrockEmbedder struct {
	client    *bedrockruntime.Client
	modelName string
}

// TitanEmbeddingRequest Amazon Titan Embedding请求格式
type TitanEmbeddingRequest struct {
	InputText string `json:"inputText"`
}

// TitanEmbeddingResponse Amazon Titan Embedding响应格式
type TitanEmbeddingResponse struct {
	Embedding     []float64 `json:"embedding"`
	InputTextTokenCount int   `json:"inputTextTokenCount"`
}

// EmbedText 实现embedding.Embedder接口的EmbedText方法
func (b *BedrockEmbedder) EmbedText(ctx context.Context, text string) ([]float64, error) {
	// 根据模型类型构建不同的请求
	var requestBody []byte
	var err error
	
	if b.isAmazonTitanModel() {
		request := TitanEmbeddingRequest{
			InputText: text,
		}
		requestBody, err = json.Marshal(request)
	} else {
		// 对于其他模型，可能需要不同的请求格式
		return nil, fmt.Errorf("unsupported embedding model: %s", b.modelName)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 调用Bedrock API
	response, err := b.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     &b.modelName,
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        requestBody,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to invoke embedding model: %w", err)
	}

	// 解析响应
	if b.isAmazonTitanModel() {
		var titanResp TitanEmbeddingResponse
		if err := json.Unmarshal(response.Body, &titanResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal titan response: %w", err)
		}
		return titanResp.Embedding, nil
	}

	return nil, fmt.Errorf("unsupported embedding model response format: %s", b.modelName)
}

// EmbedStrings 批量嵌入文本 - 实现embedding.Embedder接口
func (b *BedrockEmbedder) EmbedStrings(ctx context.Context, texts []string) ([][]float64, error) {
	embeddings := make([][]float64, len(texts))
	
	// 简单的串行处理，生产环境可以考虑并发处理
	for i, text := range texts {
		embedding, err := b.EmbedText(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to embed text at index %d: %w", i, err)
		}
		embeddings[i] = embedding
	}
	
	return embeddings, nil
}

// isAmazonTitanModel 检查是否是Amazon Titan模型
func (b *BedrockEmbedder) isAmazonTitanModel() bool {
	return b.modelName == "amazon.titan-embed-text-v1" || 
		   b.modelName == "amazon.titan-embed-text-v2:0"
}

// GetDimensions 获取嵌入向量的维度
func (b *BedrockEmbedder) GetDimensions() int {
	switch b.modelName {
	case "amazon.titan-embed-text-v1":
		return 1536
	case "amazon.titan-embed-text-v2:0":
		return 1024 // 或其他维度，需要根据实际情况调整
	default:
		return 1536 // 默认维度
	}
}