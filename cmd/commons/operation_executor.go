package commons

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/urfave/cli"

	"github.com/Knetic/govaluate"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

type SubjectOperator struct {
	ExecuteAction SubjectActionFn
	GetName       SubjectNameFn
	Type          string
	PrintHeader   HeaderPrinterFn
}

type SubjectActionFn func(client *rabbithole.Client, subject *interface{}, w *tabwriter.Writer)
type SubjectNameFn func(subject *interface{}) string

type HeaderPrinterFn func(w *tabwriter.Writer)

func ExecuteOperation(ctx *cli.Context, client *rabbithole.Client, subjects *[]interface{}, subjectOperator SubjectOperator) {
	filter := ctx.String("filter")
	dryRun := ctx.Bool("dry-run")

	matchExpression, err := govaluate.NewEvaluableExpression(filter)
	AbortIfError(err)
	p, bar := initializeProgressBar(subjects)

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	if subjectOperator.PrintHeader != nil {
		subjectOperator.PrintHeader(w)
	}
	matchingSubjectsCount := 0
	for _, subject := range *subjects {
		start := time.Now()
		parameters := make(map[string]interface{}, 1)
		parameters[subjectOperator.Type] = subject
		subjectMatches, err := matchExpression.Evaluate(parameters)
		AbortIfError(err, fmt.Sprintf("Error evaluating following filter expression[%s] err: %s", filter, err))

		if subjectMatches.(bool) {
			if dryRun {
				//fmt.Fprintf(w, "Skipping %s operation: %s/%s in dry-run mode\n", subjectName, subject.Vhost, subject.Name)
				fmt.Fprintf(w, "Skipping %s operation: %s in dry-run mode\n", subjectOperator.Type, subjectOperator.GetName(&subject))
				bar.Increment(time.Since(start))
				continue
			}
			matchingSubjectsCount++
			subjectOperator.ExecuteAction(client, &subject, w)

			fmt.Fprintln(w)
		}
		// since ewma decorator is used, we need to pass time.Since(start)
		bar.Increment(time.Since(start))
	}

	// wait for our bar to complete and flush
	p.Wait()

	w.Flush()
	fmt.Printf("Operation executed on %d matching subjects\n", matchingSubjectsCount)
}

func initializeProgressBar(subjects *[]interface{}) (*mpb.Progress, *mpb.Bar) {
	// initialize progress container, with custom width
	p := mpb.New(mpb.WithWidth(64))
	total := len(*subjects)
	name := "Executing operation:"
	// adding a single bar, which will inherit container's width
	bar := p.AddBar(int64(total),
		mpb.PrependDecorators(
			// display our name with one space on the right
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
			// replace ETA decorator with "done" message, OnComplete event
			decor.OnComplete(
				// ETA decorator with ewma age of 60, and width reservation of 4
				decor.EwmaETA(decor.ET_STYLE_GO, 60, decor.WC{W: 4}), "done",
			),
		),
		mpb.AppendDecorators(decor.Percentage()),
		mpb.BarClearOnComplete(),
		mpb.BarRemoveOnComplete(),
	)
	return p, bar
}
