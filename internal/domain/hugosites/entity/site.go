package entity

import (
	"context"
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/deps"
	fsVO "github.com/dddplayer/hugoverse/internal/domain/fs/valueobject"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	langEntity "github.com/dddplayer/hugoverse/internal/domain/language/entity"
	"github.com/dddplayer/hugoverse/internal/domain/template"
	"github.com/dddplayer/hugoverse/pkg/bufferpool"
	"github.com/dddplayer/hugoverse/pkg/lazy"
	"github.com/dddplayer/hugoverse/pkg/paths"
	"io"
	"path"
	"sort"
	"strings"
	"sync"
)

type Site struct {
	deps.Deps

	// The owning container. When multiple languages, there will be multiple
	// sites .
	H *HugoSites

	Language *langEntity.Language

	// Output formats defined in site config per Page Kind, or some defaults
	// if not set.
	// Output formats defined in Page front matter will override these.
	OutputFormats map[string]valueobject.Formats

	// The output formats that we need to render this site in. This slice
	// will be fixed once set.
	// This will be the union of Site.Pages' outputFormats.
	// This slice will be sorted.
	RenderFormats valueobject.Formats

	// All the output formats and media types available for this site.
	// These values will be merged from the Hugo defaults, the site config and,
	// finally, the language settings.
	OutputFormatsConfig valueobject.Formats
	MediaTypesConfig    valueobject.Types

	Publisher hugosites.Publisher

	*PageCollections
}

func (s *Site) process() error {
	if err := s.readAndProcessContent(); err != nil {
		err = fmt.Errorf("readAndProcessContent: %w", err)

		return err
	}

	return nil
}

func (s *Site) readAndProcessContent() error {
	proc := newPagesProcessor(s.H)
	c := newPagesCollector(proc, s.Deps.ContentFs())
	if err := c.Collect(); err != nil {
		return err
	}

	return nil
}

func (s *Site) newPage(n *contentNode, kind string, sections ...string) *pageState {
	p, err := newPageFromMeta(
		n,
		&pageMeta{
			s:        s,
			kind:     kind,
			sections: sections,
		})
	if err != nil {
		panic(err)
	}

	return p
}

func (s *Site) sectionsFromFile(fi hugosites.File) []string {
	dirname := fi.Dir()

	dirname = strings.Trim(dirname, paths.FilePathSeparator)
	if dirname == "" {
		return nil
	}
	parts := strings.Split(dirname, paths.FilePathSeparator)

	if fii, ok := fi.(*fileInfo); ok {
		if len(parts) > 0 && fii.FileInfo().Meta().Classifier == fsVO.ContentClassLeaf {
			// my-section/mybundle/index.md => my-section
			return parts[:len(parts)-1]
		}
	}

	return parts
}

func (s *Site) kindFromFileInfoOrSections(fi *fileInfo, sections []string) string {
	if fi.TranslationBaseName() == "_index" {
		if fi.Dir() == "" {
			return hugosites.KindHome
		}
		return s.kindFromSections(sections)
	}

	return hugosites.KindPage
}

func (s *Site) kindFromSections(sections []string) string {
	if len(sections) == 0 {
		return hugosites.KindHome
	}

	return s.kindFromSectionPath(path.Join(sections...))
}

func (s *Site) kindFromSectionPath(sectionPath string) string {
	return hugosites.KindSection
}

// This is all the kinds we can expect to find in .Site.Pages.
var allKindsInPages = []string{
	valueobject.KindPage,
	valueobject.KindHome,
	valueobject.KindSection}

func (s *Site) initRenderFormats() {
	formatSet := make(map[string]bool)
	formats := valueobject.Formats{}

	// media type - format
	// site output format - render format
	// Add the per kind configured output formats
	for _, kind := range allKindsInPages {
		if siteFormats, found := s.OutputFormats[kind]; found {
			for _, f := range siteFormats {
				if !formatSet[f.Name] {
					formats = append(formats, f)
					formatSet[f.Name] = true
				}
			}
		}
	}

	sort.Sort(formats)

	// HTML
	s.RenderFormats = formats
}

func (s *Site) preparePagesForRender() error {
	var err error
	s.PageMap.withEveryBundlePage(func(p *pageState) bool {
		if err = p.initOutputFormat(); err != nil {
			return true
		}
		return false
	})
	return nil
}

func (s *Site) initInit(init *lazy.Init) bool {
	_, err := init.Do()
	if err != nil {
		fmt.Printf("fatal error %v", err)
	}
	return err == nil
}

func (s *Site) render() (err error) {
	if err = s.renderPages(); err != nil {
		return
	}
	return
}

// renderPages renders pages each corresponding to a markdown file.
func (s *Site) renderPages() error {
	numWorkers := 3

	results := make(chan error)
	pages := make(chan *pageState, numWorkers) // buffered for performance
	errs := make(chan error)

	go s.errorCollator(results, errs)

	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go pageRenderer(s, pages, results, wg)
	}

	var count int
	s.PageMap.PageTrees.Walk(func(ss string, n *contentNode) bool {
		select {
		default:
			count++
			fmt.Println("777 count: ", count, ss, n, n.p)
			pages <- n.p
		}

		return false
	})

	close(pages)

	wg.Wait()

	close(results)

	err := <-errs
	if err != nil {
		return fmt.Errorf("failed to render pages: %w", err)
	}
	return nil
}

func (s *Site) errorCollator(results <-chan error, errs chan<- error) {
	var errors []error
	for e := range results {
		errors = append(errors, e)
	}

	if len(errors) > 0 {
		errs <- fmt.Errorf("failed to render pages: %v", errors)
	}

	close(errs)
}

func pageRenderer(s *Site, pages <-chan *pageState, results chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range pages {
		fmt.Printf(">>>> page: %#+v\n", p)

		templ, found, err := p.resolveTemplate()
		if err != nil {
			fmt.Println("failed to resolve template")
			continue
		}

		if !found { // layout: "", kind: section, name: HTML
			fmt.Printf("layout: %s, kind: %s, name: %s", p.Layout(), p.Kind(), p.f.Name)
			continue
		}

		targetPath := p.TargetPaths().TargetFilename

		if err := s.renderAndWritePage(targetPath, p, templ); err != nil {
			fmt.Println(" render err")
			fmt.Printf("%#v", err)
			results <- err
		}
	}
}

func (s *Site) renderAndWritePage(targetPath string, p *pageState, templ template.Template) error {
	renderBuffer := bufferpool.GetBuffer()
	defer bufferpool.PutBuffer(renderBuffer)

	of := p.outputFormat()

	if err := s.renderForTemplate(p.Kind(), of.Name, p, renderBuffer, templ); err != nil {
		return err
	}

	if renderBuffer.Len() == 0 {
		return nil
	}

	pd := hugosites.Descriptor{
		Src:          renderBuffer,
		TargetPath:   targetPath,
		OutputFormat: p.outputFormat(),
	}

	return s.Publisher.Publish(pd)
}

func (s *Site) renderForTemplate(name, outputFormat string, d any, w io.Writer, templ template.Template) (err error) {
	if templ == nil {
		fmt.Printf("missing layout name: %s, output format: %s", name, outputFormat)
		return nil
	}

	if err = s.Tmpl().ExecuteWithContext(context.Background(), templ, w, d); err != nil {
		return fmt.Errorf("render of %q failed: %w", name, err)
	}
	return
}
