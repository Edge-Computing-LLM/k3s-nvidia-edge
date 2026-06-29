package edge

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Runner struct {
	opts Options
}

type Step struct {
	Name     string
	Command  string
	Mutating bool
	Host     bool
}

func NewRunner(opts Options) *Runner {
	return &Runner{opts: opts}
}

func (r *Runner) Run(ctx context.Context, step Step) error {
	command := step.Command
	if step.Host && r.opts.Sudo && os.Geteuid() != 0 {
		command = "sudo bash -lc " + shellQuote(command)
	}

	if step.Mutating && !r.opts.Yes {
		fmt.Printf("[dry-run] %s\n  %s\n", step.Name, command)
		return nil
	}

	fmt.Printf("[run] %s\n", step.Name)
	if r.opts.Verbose {
		fmt.Printf("  %s\n", command)
	}

	cmd := exec.CommandContext(ctx, "bash", "-lc", command)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	out := strings.TrimSpace(stdout.String())
	errOut := strings.TrimSpace(stderr.String())
	if r.opts.Verbose && out != "" {
		fmt.Println(out)
	}
	if err != nil {
		if out != "" {
			fmt.Println(out)
		}
		if errOut != "" {
			fmt.Fprintln(os.Stderr, errOut)
		}
		return fmt.Errorf("%s failed: %w", step.Name, err)
	}
	if !r.opts.Verbose && out != "" {
		fmt.Println(out)
	}
	return nil
}

func (r *Runner) Output(ctx context.Context, name, command string) (string, error) {
	cmd := exec.CommandContext(ctx, "bash", "-lc", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("%s failed: %w", name, err)
	}
	return string(out), nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}
