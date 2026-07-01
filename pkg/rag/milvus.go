package rag

import (
	"context"
	"fmt"
	"strconv"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	MilvusIDField        = "id"
	MilvusTitleField     = "title"
	MilvusCityField      = "city"
	MilvusTagsField      = "tags"
	MilvusContentField   = "content"
	MilvusEmbeddingField = "embedding"
)

// BuildMilvusSchema 构造旅行资料集合结构。
func BuildMilvusSchema(collectionName string, dim int) *entity.Schema {
	return entity.NewSchema().
		WithName(collectionName).
		WithDescription("旅行助手资料库").
		WithAutoID(false).
		WithField(entity.NewField().
			WithName(MilvusIDField).
			WithDataType(entity.FieldTypeVarChar).
			WithIsPrimaryKey(true).
			WithIsAutoID(false).
			WithMaxLength(128)).
		WithField(entity.NewField().
			WithName(MilvusTitleField).
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(512)).
		WithField(entity.NewField().
			WithName(MilvusCityField).
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(128)).
		WithField(entity.NewField().
			WithName(MilvusTagsField).
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(512)).
		WithField(entity.NewField().
			WithName(MilvusContentField).
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(4096)).
		WithField(entity.NewField().
			WithName(MilvusEmbeddingField).
			WithDataType(entity.FieldTypeFloatVector).
			WithDim(int64(dim)))
}

// BuildMilvusColumns 把资料记录转为 Milvus 列数据。
func BuildMilvusColumns(records []SeedRecord, dim int) []entity.Column {
	ids := make([]string, 0, len(records))
	titles := make([]string, 0, len(records))
	cities := make([]string, 0, len(records))
	tags := make([]string, 0, len(records))
	contents := make([]string, 0, len(records))
	embeddings := make([][]float32, 0, len(records))

	for _, record := range records {
		ids = append(ids, record.ID)
		titles = append(titles, record.Title)
		cities = append(cities, record.City)
		tags = append(tags, record.Tags)
		contents = append(contents, record.Content)
		embeddings = append(embeddings, record.Embedding)
	}

	return []entity.Column{
		entity.NewColumnVarChar(MilvusIDField, ids),
		entity.NewColumnVarChar(MilvusTitleField, titles),
		entity.NewColumnVarChar(MilvusCityField, cities),
		entity.NewColumnVarChar(MilvusTagsField, tags),
		entity.NewColumnVarChar(MilvusContentField, contents),
		entity.NewColumnFloatVector(MilvusEmbeddingField, dim, embeddings),
	}
}

// MilvusWriter 负责把旅行资料写入 Milvus。
type MilvusWriter struct {
	client client.Client
}

// NewMilvusWriter 创建写入器。
func NewMilvusWriter(client client.Client) *MilvusWriter {
	return &MilvusWriter{client: client}
}

// EnsureCollection 确保集合存在；reset 为 true 时会先删除旧集合。
func (w *MilvusWriter) EnsureCollection(ctx context.Context, collectionName string, dim int, reset bool) error {
	if collectionName == "" {
		return fmt.Errorf("集合名称不能为空")
	}
	if dim <= 0 {
		return fmt.Errorf("向量维度必须大于 0，当前值：%s", strconv.Itoa(dim))
	}

	hasCollection, err := w.client.HasCollection(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("检查集合失败：%w", err)
	}
	if hasCollection && reset {
		if err := w.client.DropCollection(ctx, collectionName); err != nil {
			return fmt.Errorf("删除旧集合失败：%w", err)
		}
		hasCollection = false
	}
	if hasCollection {
		return nil
	}

	if err := w.client.CreateCollection(ctx, BuildMilvusSchema(collectionName, dim), entity.DefaultShardNumber); err != nil {
		return fmt.Errorf("创建集合失败：%w", err)
	}
	index, err := entity.NewIndexIvfFlat(entity.L2, 128)
	if err != nil {
		return fmt.Errorf("创建索引参数失败：%w", err)
	}
	if err := w.client.CreateIndex(ctx, collectionName, MilvusEmbeddingField, index, false); err != nil {
		return fmt.Errorf("创建向量索引失败：%w", err)
	}
	return nil
}

// UpsertRecords 写入或更新资料记录。
func (w *MilvusWriter) UpsertRecords(ctx context.Context, collectionName string, records []SeedRecord, dim int) error {
	if len(records) == 0 {
		return nil
	}
	if _, err := w.client.Upsert(ctx, collectionName, "", BuildMilvusColumns(records, dim)...); err != nil {
		return fmt.Errorf("写入资料失败：%w", err)
	}
	if err := w.client.Flush(ctx, collectionName, false); err != nil {
		return fmt.Errorf("刷新集合失败：%w", err)
	}
	return nil
}
