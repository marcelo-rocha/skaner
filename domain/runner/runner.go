package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/marcelo-rocha/skaner/domain/checker"
	"github.com/marcelo-rocha/skaner/domain/checker/exposure"
	"github.com/marcelo-rocha/skaner/domain/checker/sql"
	"github.com/marcelo-rocha/skaner/domain/checker/xss"
	"github.com/marcelo-rocha/skaner/domain/sourcecode"

	"go.uber.org/zap"
)

type CheckingKind int

const (
	XSSChecking CheckingKind = iota
	SQLChecking
	ExposureChecking
)

const CheckersQty = 3

type Job struct {
	Kind        CheckingKind
	FileContent checker.SourceCode
	FileName    string
}

func worker(id int, checkers []checker.Checker, jobs <-chan Job, results chan<- []checker.Vulnerability,
	wg *sync.WaitGroup, logger *zap.Logger) {
	defer wg.Done()
	for j := range jobs {
		r, err := checkers[j.Kind].Check(j.FileContent)
		if err != nil {
			logger.Error("checking failed", zap.Error(err), zap.Int("worked_id", id))
			continue
		}
		results <- r
	}
}

type Options struct {
	SensitiveText        []string
	DisableXSSCheck      bool
	DisableExposureCheck bool
	DisableSQLCheck      bool
	JsonOutput           bool
	WorkersQty           int
}

func prepareCheckers(opts Options, logger *zap.Logger) []checker.Checker {
	checkers := make([]checker.Checker, CheckersQty)
	nop := &checker.NoOperationChecker{}
	for i := range checkers {
		checkers[i] = nop
	}
	if !opts.DisableExposureCheck {
		c, err := exposure.New(opts.SensitiveText)
		if err != nil {
			logger.Fatal("can't execute exposure checker", zap.Error(err))
		}
		checkers[ExposureChecking] = c
	}

	if !opts.DisableSQLCheck {
		checkers[SQLChecking] = sql.New()
	}

	if !opts.DisableXSSCheck {
		checkers[XSSChecking] = xss.New()
	}
	return checkers
}

func AddJobIfEnabled(enabled bool, queue []Job, kind CheckingKind, checkers []checker.Checker,
	src checker.SourceCode, name string) []Job {
	if enabled && checkers[kind].SupportedFileExtension(path.Ext(name)) {
		queue = append(queue, Job{
			Kind:        kind,
			FileContent: src,
			FileName:    name,
		})
	}
	return queue
}

func printVulnerabilites(output io.Writer, list []checker.Vulnerability, jsonFormat bool) error {
	if jsonFormat {
		JsonOutput, err := json.MarshalIndent(list, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintf(output, "%s\n", JsonOutput)
	} else {
		for _, item := range list {
			fmt.Fprintf(output, "[%s] in file \"%s\" on line %v\n", item.Kind, path.Base(item.FilePath), item.Line)
		}
	}
	return nil
}

func Run(ctx context.Context, output io.Writer, fileNames []string, opts Options, logger *zap.Logger) {

	checkers := prepareCheckers(opts, logger)
	var wg sync.WaitGroup
	var result []checker.Vulnerability
	mustAbort := false
	jobChannel := make(chan Job)
	resultChannel := make(chan []checker.Vulnerability)

	for n := 0; n < opts.WorkersQty; n++ {
		wg.Add(1)
		go worker(n, checkers, jobChannel, resultChannel, &wg, logger)
	}

	activeJobs := 0
_filesLoop:
	for _, name := range fileNames {
		content, err := os.ReadFile(name)
		if err != nil {
			logger.Error("failed to read file", zap.Error(err), zap.String("fileName", name))
			continue
		}
		var src checker.SourceCode = sourcecode.NewSourceCode(content, name)

		var queue = []Job{}
		queue = AddJobIfEnabled(!opts.DisableExposureCheck, queue, ExposureChecking, checkers, src, name)
		queue = AddJobIfEnabled(!opts.DisableSQLCheck, queue, SQLChecking, checkers, src, name)
		queue = AddJobIfEnabled(!opts.DisableXSSCheck, queue, XSSChecking, checkers, src, name)

		i := 0
		for i < len(queue) {
			select {
			case <-ctx.Done():
				logger.Debug("context cancelled")
				mustAbort = true
				break _filesLoop
			case jobChannel <- queue[i]:
				activeJobs++
				i++
			case r := <-resultChannel:
				result = append(result, r...)
				activeJobs--
			}
		}
	}

	close(jobChannel)

	// All jobs was sent, now only waiting for results
	for activeJobs > 0 && !mustAbort {
		select {
		case <-ctx.Done():
			{
				logger.Debug("context cancelled")
				mustAbort = true
			}
		case r := <-resultChannel:
			{
				result = append(result, r...)
				activeJobs--
			}
		}
	}

	wg.Wait()
	close(resultChannel)

	if !mustAbort && len(result) > 0 {
		if err := printVulnerabilites(output, result, opts.JsonOutput); err != nil {
			logger.Error("failed to output result", zap.Error(err))
		}
	}
}
