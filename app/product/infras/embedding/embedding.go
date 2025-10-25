package embedding

import (
	"context"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/cloudwego/eino/components/embedding"
)

var (
	EB   embedding.Embedder
	once sync.Once
)

// BedrockEmbeddingConfig AWS Bedrock Embedding配置
type BedrockEmbeddingConfig struct {
	Model       string
	Region      string
	AccessKeyID string
	SecretKey   string
}

func defaultBedrockEmbeddingConfig(_ context.Context) (*BedrockEmbeddingConfig, error) {
	config := &BedrockEmbeddingConfig{
		Model:       os.Getenv("AWS_BEDROCK_EMBEDDING_MODEL"),
		Region:      os.Getenv("AWS_REGION"),
		AccessKeyID: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:   os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
	
	// 设置默认值
	if config.Model == "" {
		config.Model = "amazon.titan-embed-text-v1"
	}
	if config.Region == "" {
		config.Region = "us-east-1"
	}
	
	return config, nil
}

func GetBedrockEmbedding(ctx context.Context, config *BedrockEmbeddingConfig) (eb embedding.Embedder, err error) {
	once.Do(func() {
		if config == nil {
			config, err = defaultBedrockEmbeddingConfig(ctx)
			if err != nil {
				return
			}
		}
		
		// 创建AWS配置
		var cfg aws.Config
		if config.AccessKeyID != "" && config.SecretKey != "" {
			// 使用显式提供的凭证
			cfg = aws.Config{
				Region: config.Region,
				Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
					return aws.Credentials{
						AccessKeyID:     config.AccessKeyID,
						SecretAccessKey: config.SecretKey,
					}, nil
				}),
			}
		} else {
			// 使用默认凭证链
			cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(config.Region))
			if err != nil {
				return
			}
		}
		
		// 创建Bedrock运行时客户端
		client := bedrockruntime.NewFromConfig(cfg)
		
		// 创建自定义的Bedrock Embedding实现
		EB = &BedrockEmbedder{
			client:    client,
			modelName: config.Model,
		}
	})
	return EB, err
}
