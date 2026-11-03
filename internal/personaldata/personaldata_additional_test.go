package personaldata

import (
	"io"
	"os"
	"testing"
)

func TestPrintRoundsWithoutChangingPersonal(t *testing.T) {
	p := Personal{Name: "Анна", Weight: 65.456, Height: 1.756}
	before := p
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()
	stdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdout }()

	p.Print()
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	const want = "Имя: Анна\nВес: 65.46\nРост: 1.76\n"
	if string(output) != want {
		t.Fatalf("Print() = %q, want %q", output, want)
	}
	if p != before {
		t.Fatalf("Print() changed Personal: got %+v, want %+v", p, before)
	}
}
