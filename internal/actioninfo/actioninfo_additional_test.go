package actioninfo_test

import (
	"bytes"
	"errors"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/Yandex-Practicum/tracker/internal/actioninfo"
	"github.com/Yandex-Practicum/tracker/internal/daysteps"
	"github.com/Yandex-Practicum/tracker/internal/trainings"
)

var (
	_ actioninfo.DataParser = (*trainings.Training)(nil)
	_ actioninfo.DataParser = (*daysteps.DaySteps)(nil)
)

type parserResult struct {
	data     string
	report   string
	parseErr error
	infoErr  error
}

type parserStub struct {
	t       *testing.T
	results []parserResult
	index   int
	calls   []string
}

func (p *parserStub) Parse(data string) error {
	p.t.Helper()
	p.index++
	if p.index >= len(p.results) {
		p.t.Fatalf("unexpected Parse(%q)", data)
	}
	p.calls = append(p.calls, "Parse:"+data)
	if want := p.results[p.index].data; data != want {
		p.t.Fatalf("Parse argument = %q, want %q", data, want)
	}
	return p.results[p.index].parseErr
}

func (p *parserStub) ActionInfo() (string, error) {
	p.t.Helper()
	if p.index < 0 || p.index >= len(p.results) {
		p.t.Fatal("ActionInfo called before Parse")
	}
	result := p.results[p.index]
	p.calls = append(p.calls, "ActionInfo:"+result.data)
	return result.report, result.infoErr
}

func TestInfoProcessing(t *testing.T) {
	parseErr := errors.New("parse failed")
	infoErr := errors.New("report failed")
	tests := []struct {
		name       string
		results    []parserResult
		wantOutput string
		wantLog    string
		wantCalls  []string
	}{
		{name: "nil dataset"},
		{name: "empty dataset", results: []parserResult{}},
		{
			name:       "successful records preserve report formatting",
			results:    []parserResult{{data: "first", report: "first report\n"}, {data: "second", report: "second report"}},
			wantOutput: "first report\nsecond report",
			wantCalls:  []string{"Parse:first", "ActionInfo:first", "Parse:second", "ActionInfo:second"},
		},
		{
			name:      "parse error skips report",
			results:   []parserResult{{data: "bad", report: "must not print\n", parseErr: parseErr}},
			wantLog:   "parse failed\n",
			wantCalls: []string{"Parse:bad"},
		},
		{
			name:      "report error suppresses returned text",
			results:   []parserResult{{data: "bad", report: "must not print\n", infoErr: infoErr}},
			wantLog:   "report failed\n",
			wantCalls: []string{"Parse:bad", "ActionInfo:bad"},
		},
		{
			name: "mixed records continue after both errors",
			results: []parserResult{
				{data: "first", report: "first report\n"},
				{data: "parse error", report: "must not print\n", parseErr: parseErr},
				{data: "second", report: "second report\n"},
				{data: "report error", report: "must not print\n", infoErr: infoErr},
				{data: "third", report: "third report\n"},
			},
			wantOutput: "first report\nsecond report\nthird report\n",
			wantLog:    "parse failed\nreport failed\n",
			wantCalls: []string{
				"Parse:first", "ActionInfo:first", "Parse:parse error",
				"Parse:second", "ActionInfo:second", "Parse:report error",
				"ActionInfo:report error", "Parse:third", "ActionInfo:third",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dataset []string
			if tt.results != nil {
				dataset = make([]string, len(tt.results))
			}
			for i, result := range tt.results {
				dataset[i] = result.data
			}
			parser := &parserStub{t: t, results: tt.results, index: -1}
			output, logs := captureInfo(t, dataset, parser)
			if output != tt.wantOutput {
				t.Errorf("stdout = %q, want %q", output, tt.wantOutput)
			}
			if logs != tt.wantLog {
				t.Errorf("log = %q, want %q", logs, tt.wantLog)
			}
			if !reflect.DeepEqual(parser.calls, tt.wantCalls) {
				t.Errorf("calls = %v, want %v", parser.calls, tt.wantCalls)
			}
		})
	}
}

func captureInfo(t *testing.T, dataset []string, parser actioninfo.DataParser) (string, string) {
	t.Helper()
	stdout, err := os.CreateTemp(t.TempDir(), "stdout-")
	if err != nil {
		t.Fatal(err)
	}
	oldStdout := os.Stdout
	oldWriter, oldFlags, oldPrefix := log.Writer(), log.Flags(), log.Prefix()
	t.Cleanup(func() {
		os.Stdout = oldStdout
		log.SetOutput(oldWriter)
		log.SetFlags(oldFlags)
		log.SetPrefix(oldPrefix)
		stdout.Close()
	})
	var logs bytes.Buffer
	os.Stdout = stdout
	log.SetOutput(&logs)
	log.SetFlags(0)
	log.SetPrefix("")
	actioninfo.Info(dataset, parser)
	output, err := os.ReadFile(stdout.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(output), logs.String()
}
