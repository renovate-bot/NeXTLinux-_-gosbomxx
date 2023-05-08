package packages

import (
	"context"
	"fmt"

	"github.com/wagoodman/go-partybus"

	"github.com/nextlinux/stereoscope"
	"github.com/nextlinux/gosbom/cmd/gosbom/cli/eventloop"
	"github.com/nextlinux/gosbom/cmd/gosbom/cli/options"
	"github.com/nextlinux/gosbom/internal"
	"github.com/nextlinux/gosbom/internal/bus"
	"github.com/nextlinux/gosbom/internal/config"
	"github.com/nextlinux/gosbom/internal/log"
	"github.com/nextlinux/gosbom/internal/ui"
	"github.com/nextlinux/gosbom/internal/version"
	"github.com/nextlinux/gosbom/gosbom"
	"github.com/nextlinux/gosbom/gosbom/artifact"
	"github.com/nextlinux/gosbom/gosbom/event"
	"github.com/nextlinux/gosbom/gosbom/formats/template"
	"github.com/nextlinux/gosbom/gosbom/sbom"
	"github.com/nextlinux/gosbom/gosbom/source"
)

func Run(_ context.Context, app *config.Application, args []string) error {
	err := ValidateOutputOptions(app)
	if err != nil {
		return err
	}

	writer, err := options.MakeWriter(app.Outputs, app.File, app.OutputTemplatePath)
	if err != nil {
		return err
	}

	defer func() {
		if err := writer.Close(); err != nil {
			log.Warnf("unable to write to report destination: %w", err)
		}
	}()

	// could be an image or a directory, with or without a scheme
	userInput := args[0]
	si, err := source.ParseInputWithName(userInput, app.Platform, app.Name, app.DefaultImagePullSource)
	if err != nil {
		return fmt.Errorf("could not generate source input for packages command: %w", err)
	}

	eventBus := partybus.NewBus()
	stereoscope.SetBus(eventBus)
	gosbom.SetBus(eventBus)
	subscription := eventBus.Subscribe()

	return eventloop.EventLoop(
		execWorker(app, *si, writer),
		eventloop.SetupSignals(),
		subscription,
		stereoscope.Cleanup,
		ui.Select(options.IsVerbose(app), app.Quiet)...,
	)
}

func execWorker(app *config.Application, si source.Input, writer sbom.Writer) <-chan error {
	errs := make(chan error)
	go func() {
		defer close(errs)

		src, cleanup, err := source.New(si, app.Registry.ToOptions(), app.Exclusions)
		if cleanup != nil {
			defer cleanup()
		}
		if err != nil {
			errs <- fmt.Errorf("failed to construct source from user input %q: %w", si.UserInput, err)
			return
		}

		s, err := GenerateSBOM(src, errs, app)
		if err != nil {
			errs <- err
			return
		}

		if s == nil {
			errs <- fmt.Errorf("no SBOM produced for %q", si.UserInput)
		}

		bus.Publish(partybus.Event{
			Type:  event.Exit,
			Value: func() error { return writer.Write(*s) },
		})
	}()
	return errs
}

func GenerateSBOM(src *source.Source, errs chan error, app *config.Application) (*sbom.SBOM, error) {
	tasks, err := eventloop.Tasks(app)
	if err != nil {
		return nil, err
	}

	s := sbom.SBOM{
		Source: src.Metadata,
		Descriptor: sbom.Descriptor{
			Name:          internal.ApplicationName,
			Version:       version.FromBuild().Version,
			Configuration: app,
		},
	}

	buildRelationships(&s, src, tasks, errs)

	return &s, nil
}

func buildRelationships(s *sbom.SBOM, src *source.Source, tasks []eventloop.Task, errs chan error) {
	var relationships []<-chan artifact.Relationship
	for _, task := range tasks {
		c := make(chan artifact.Relationship)
		relationships = append(relationships, c)
		go eventloop.RunTask(task, &s.Artifacts, src, c, errs)
	}

	s.Relationships = append(s.Relationships, MergeRelationships(relationships...)...)
}

func MergeRelationships(cs ...<-chan artifact.Relationship) (relationships []artifact.Relationship) {
	for _, c := range cs {
		for n := range c {
			relationships = append(relationships, n)
		}
	}

	return relationships
}

func ValidateOutputOptions(app *config.Application) error {
	var usesTemplateOutput bool
	for _, o := range app.Outputs {
		if o == template.ID.String() {
			usesTemplateOutput = true
			break
		}
	}

	if usesTemplateOutput && app.OutputTemplatePath == "" {
		return fmt.Errorf(`must specify path to template file when using "template" output format`)
	}

	return nil
}
