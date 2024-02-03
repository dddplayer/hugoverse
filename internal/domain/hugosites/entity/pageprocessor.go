package entity

import (
	"context"
	"fmt"
	fsVO "github.com/dddplayer/hugoverse/internal/domain/fs/valueobject"
	"golang.org/x/sync/errgroup"
)

func newPagesProcessor(h *HugoSites) *pagesProcessor {
	procs := make(map[string]pagesCollectorProcessorProvider)
	for _, s := range h.Sites {
		procs[s.Language.Lang] = &sitePagesProcessor{
			m:        s.PageMap,
			itemChan: make(chan interface{}, 1),
		}
	}
	return &pagesProcessor{
		procs: procs,
	}
}

type pagesCollectorProcessorProvider interface {
	Process(item any) error
	Start(ctx context.Context) context.Context
	Wait() error
}

type pagesProcessor struct {
	// Per language/Site
	procs map[string]pagesCollectorProcessorProvider
}

func (proc *pagesProcessor) Process(item any) error {
	switch v := item.(type) {
	case fsVO.FileMetaInfo:
		err := proc.getProc().Process(v)
		if err != nil {
			return err
		}
	default:
		panic(fmt.Sprintf("unrecognized item type in Process: %T", item))
	}

	return nil
}

func (proc *pagesProcessor) Start(ctx context.Context) context.Context {
	for _, p := range proc.procs {
		ctx = p.Start(ctx)
	}
	return ctx
}

func (proc *pagesProcessor) Wait() error {
	var err error
	for _, p := range proc.procs {
		if e := p.Wait(); e != nil {
			err = e
		}
	}
	return err
}

type sitePagesProcessor struct {
	m         *PageMap
	ctx       context.Context
	itemChan  chan any
	itemGroup *errgroup.Group
}

func (p *sitePagesProcessor) Process(item any) error {
	select {
	case <-p.ctx.Done():
		return nil
	default:
		p.itemChan <- item
	}
	return nil
}

func (p *sitePagesProcessor) Start(ctx context.Context) context.Context {
	p.itemGroup, ctx = errgroup.WithContext(ctx)
	p.ctx = ctx
	p.itemGroup.Go(func() error {
		for item := range p.itemChan {
			if err := p.doProcess(item); err != nil {
				return err
			}
		}
		return nil
	})
	return ctx
}

func (p *sitePagesProcessor) Wait() error {
	close(p.itemChan)
	return p.itemGroup.Wait()
}

func (p *sitePagesProcessor) doProcess(item any) error {
	fmt.Println("doProcess --- ")
	m := p.m
	switch v := item.(type) {
	case fsVO.FileMetaInfo:
		meta := v.Meta()

		classifier := meta.Classifier
		switch classifier {
		case fsVO.ContentClassContent: // basefs.go createOverlayFs
			if err := m.AddFilesBundle(v); err != nil {
				return err
			}
		case fsVO.ContentClassFile:
			panic("doProcess not support ContentClassFile yet")
		default:
			panic(fmt.Sprintf("invalid classifier: %q", classifier))
		}
	default:
		panic(fmt.Sprintf("unrecognized item type in Process: %T", item))
	}
	return nil
}

var defaultPageProcessor = new(nopPageProcessor)

func (proc *pagesProcessor) getProc() pagesCollectorProcessorProvider {
	if p, found := proc.procs["en"]; found {
		return p
	}
	return defaultPageProcessor
}

type nopPageProcessor int

func (nopPageProcessor) Process(item any) error {
	return nil
}

func (nopPageProcessor) Start(ctx context.Context) context.Context {
	return context.Background()
}

func (nopPageProcessor) Wait() error {
	return nil
}
