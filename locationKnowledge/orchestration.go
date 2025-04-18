package locationKnowledge

import (
	"context"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
)

func BuildLocationKnowledge(ctx context.Context) (r compose.Runnable[document.Source, []string], err error) {
	const (
		ExcelLoader      = "ExcelLoader"
		ExcelTransformer = "ExcelTransformer"
		Indexer1         = "Indexer1"
	)
	g := compose.NewGraph[document.Source, []string]()
	excelLoaderKeyOfLoader, err := newLoader(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLoaderNode(ExcelLoader, excelLoaderKeyOfLoader)
	excelTransformerKeyOfDocumentTransformer, err := newDocumentTransformer(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddDocumentTransformerNode(ExcelTransformer, excelTransformerKeyOfDocumentTransformer)
	indexer1KeyOfIndexer, err := newIndexer(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddIndexerNode(Indexer1, indexer1KeyOfIndexer)
	_ = g.AddEdge(compose.START, ExcelLoader)
	_ = g.AddEdge(Indexer1, compose.END)
	_ = g.AddEdge(ExcelLoader, ExcelTransformer)
	_ = g.AddEdge(ExcelTransformer, Indexer1)
	r, err = g.Compile(ctx, compose.WithGraphName("locationKnowledge"), compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
