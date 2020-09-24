package commons

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/Knetic/govaluate"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/urfave/cli"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

// SubjectOperator is a common type to describe different operations performed on a different subjects.
// It can be used to perform list, delete, pure and other actions of queues, exchanges, etc.
type SubjectOperator struct {
	ExecuteAction SubjectActionFn
	GetName       SubjectNameFn
	Type          string
	PrintHeader   HeaderPrinterFn
}

// SubjectActionFn function is a wrapper for a function that need to be applied on a given subject
type SubjectActionFn func(client *rabbithole.Client, subject *interface{}, w *tabwriter.Writer)

// SubjectNameFn function that returns name of the subject in question
type SubjectNameFn func(subject *interface{}) string

// HeaderPrinterFn is responsible for printing header upon executing action on multiple subjects
type HeaderPrinterFn func(w *tabwriter.Writer)

// ExecuteOperation is a main function that executes different actions on different subjects.
// In general it provides following functionalities:
// - handles progress bar display (useful for time consuming operations)
// - prints header for the operation that is about to be executed
// - iterates over provided subjects and uses provided filter to determine if an action should be applied to a subject
// - executes action on a subject if it matches the filter or skips the action in case common dry-run flag was provided
// - flushes the output and cleans progress bar after completing operation on all subjects
func ExecuteOperation(ctx *cli.Context, client *rabbithole.Client, subjects *[]interface{}, subjectOperator SubjectOperator) {
	filter := ctx.String("filter")
	dryRun := ctx.Bool("dry-run")

	if len(*subjects) <= 0 {
		fmt.Println("No subjects found to act upon.")
		return
	}

	matchExpression, err := govaluate.NewEvaluableExpressionWithFunctions(filter, GetCustomFilterFunctions())
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
				Fprintf(w, "Skipping %s operation: %s in dry-run mode\n", subjectOperator.Type, subjectOperator.GetName(&subject))
				bar.Increment()
				bar.DecoratorEwmaUpdate(time.Since(start))
				continue
			}
			matchingSubjectsCount++
			subjectOperator.ExecuteAction(client, &subject, w)

			Fprintln(w)
		}
		bar.Increment()
		// since ewma decorator is used, we need to pass time.Since(start)
		bar.DecoratorEwmaUpdate(time.Since(start))
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
		mpb.BarRemoveOnComplete(),
	)
	return p, bar
}
