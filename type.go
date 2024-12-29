package mcp

import "github.com/sivchari/go-mcp/apis"

type TextContent apis.TextContent

func NewTextContent(text string) TextContent {
	return TextContent{
		Type: "text",
		Text: text,
	}
}

func (t *TextContent) WithAnnotations(annotation *apis.TextContentAnnotations) {
	t.Annotations = annotation
}

type ImageContent apis.ImageContent

func NewImageContent(encodedData, mimeType string) ImageContent {
	return ImageContent{
		Type:     "image",
		Data:     encodedData,
		MimeType: mimeType,
	}
}

func (i *ImageContent) WithAnnotations(annotation *apis.ImageContentAnnotations) {
	i.Annotations = annotation
}

type Resource apis.EmbeddedResource

func NewResource[T interface {
	apis.TextResourceContents | apis.BlobResourceContents
}](resource T) Resource {
	return Resource{
		Type:     "resource",
		Resource: resource,
	}
}

func (e *Resource) WithAnnotations(annotation *apis.EmbeddedResourceAnnotations) {
	e.Annotations = annotation
}

func Ptr[T any](v T) *T {
	return &v
}

type ItemType string

const (
	ItemTypeString  ItemType = "string"
	ItemTypeNumber  ItemType = "number"
	ItemTypeBoolean ItemType = "boolean"
)

func (i ItemType) String() string {
	return string(i)
}

type ToolInput struct {
	Properties map[string]map[string]any `json:"properties"`
	Required   []string                  `json:"required"`
}

func NewToolInput() *ToolInput {
	return &ToolInput{
		Properties: make(map[string]map[string]any),
	}
}

func (t *ToolInput) WithString(name string) *ToolInput {
	t.Properties[name] = map[string]any{"type": ItemTypeString.String()}
	return t
}

func (t *ToolInput) WithNumber(name string) *ToolInput {
	t.Properties[name] = map[string]any{"type": ItemTypeNumber.String()}
	return t
}

func (t *ToolInput) WithBoolean(name string) *ToolInput {
	t.Properties[name] = map[string]any{"type": ItemTypeBoolean.String()}
	return t
}

func (t *ToolInput) WithArray(name string, item ItemType) *ToolInput {
	t.Properties[name] = map[string]any{"type": "array", "items": map[string]string{"type": item.String()}}
	return t
}

func (t *ToolInput) WithEnum(name string, enums []string) *ToolInput {
	t.Properties[name] = map[string]any{"type": "array", "items": map[string][]string{"enum": enums}}
	return t
}

func (t *ToolInput) WithRequired(names ...string) *ToolInput {
	t.Required = names
	return t
}

func (t *ToolInput) Build() *apis.ToolInputSchema {
	return &apis.ToolInputSchema{
		Type:       "object",
		Properties: t.Properties,
		Required:   t.Required,
	}
}
