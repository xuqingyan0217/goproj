package locationKnowledge

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
	"github.com/xuri/excelize/v2"
)

type DocumentTransformerImpl struct {
	config *DocumentTransformerConfig
}

type DocumentTransformerConfig struct {
	SheetName string // 指定要处理的工作表名称，为空则处理所有工作表
	HasHeader bool   // 是否包含表头
}

// newDocumentTransformer component initialization function of node 'ExcelTransformer' in graph 'locationKnowledge'
func newDocumentTransformer(ctx context.Context) (tfr document.Transformer, err error) {
	// 配置HasHeader为true，表示第一行为表头
	config := &DocumentTransformerConfig{
		HasHeader: true,
	}
	tfr = &DocumentTransformerImpl{config: config}
	return tfr, nil
}

func (impl *DocumentTransformerImpl) Transform(ctx context.Context, src []*schema.Document, opts ...document.TransformerOption) ([]*schema.Document, error) {
	var ret []*schema.Document

	for _, doc := range src {
		docs, err := impl.processExcelDocument(doc)
		if err != nil {
			return nil, err
		}
		ret = append(ret, docs...)
	}

	return ret, nil
}

func (impl *DocumentTransformerImpl) processExcelDocument(doc *schema.Document) ([]*schema.Document, error) {
	// 解析Excel文件内容
	xlFile, err := excelize.OpenReader(strings.NewReader(doc.Content))
	if err != nil {
		return nil, err
	}
	defer xlFile.Close()

	// 获取所有工作表
	sheets := xlFile.GetSheetList()
	if len(sheets) == 0 {
		return nil, nil
	}

	// 确定要处理的工作表
	sheetName := sheets[0]
	if impl.config.SheetName != "" {
		sheetName = impl.config.SheetName
	}

	// 获取所有行
	rows, err := xlFile.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	var ret []*schema.Document

	// 处理表头
	startIdx := 0
	var headers []string
	if impl.config.HasHeader && len(rows) > 0 {
		headers = rows[0]
		startIdx = 1
	}

	// 处理数据行
	for i := startIdx; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		// 将行数据转换为字符串
		contentParts := make([]string, len(row))
		for j, cell := range row {
			contentParts[j] = strings.TrimSpace(cell)
		}
		content := strings.Join(contentParts, "\t")

		// 创建新的Document
		nDoc := &schema.Document{
			ID:       doc.ID,
			Content:  content,
			MetaData: deepCopyMap(doc.MetaData),
		}

		// 如果有表头，将数据添加到元数据中
		if impl.config.HasHeader {
			if nDoc.MetaData == nil {
				nDoc.MetaData = make(map[string]any)
			}
			for j, header := range headers {
				if j < len(row) {
					nDoc.MetaData[header] = row[j]
				}
			}
		}

		ret = append(ret, nDoc)
	}

	return ret, nil
}

func deepCopyMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	ret := make(map[string]any)
	for k, v := range m {
		ret[k] = v
	}
	return ret
}
