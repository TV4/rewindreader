package rewindreader

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRewindReader(t *testing.T) {
	data := []byte("foo bar baz")

	t.Run("ReadFull", func(t *testing.T) {
		rwr := New(bytes.NewReader(data))

		var buf bytes.Buffer

		if _, err := io.Copy(&buf, rwr); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("ReadExact->Rewind->ReadFull", func(t *testing.T) {
		rwr := New(bytes.NewReader(data))

		var buf bytes.Buffer

		if _, err := io.CopyN(&buf, rwr, int64(len(data))); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}

		if err := rwr.Rewind(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		buf.Reset()

		if _, err := io.Copy(&buf, rwr); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("ReadExact->Rewind->ReadExact->Rewind-ReadFull", func(t *testing.T) {
		rwr := New(bytes.NewReader(data))

		for n := 0; n < 2; n++ {
			var buf bytes.Buffer

			if _, err := io.CopyN(&buf, rwr, int64(len(data))); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got, want := buf.Bytes(), data; !bytes.Equal(got, want) {
				t.Errorf("got %q, want %q", got, want)
			}

			if err := rwr.Rewind(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		var buf bytes.Buffer

		if _, err := io.Copy(&buf, rwr); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("ReadExact->Rewind->Rewind->ReadFull", func(t *testing.T) {
		rwr := New(bytes.NewReader(data))

		var buf bytes.Buffer

		if _, err := io.CopyN(&buf, rwr, int64(len(data))); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}

		for n := 0; n < 2; n++ {
			if err := rwr.Rewind(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		buf.Reset()

		if _, err := io.Copy(&buf, rwr); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("ReadPartial->Rewind->ReadFull", func(t *testing.T) {
		rwr := New(bytes.NewReader(data))

		var buf bytes.Buffer

		if _, err := io.CopyN(&buf, rwr, 3); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data[:3]; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}

		if err := rwr.Rewind(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		buf.Reset()

		if _, err := io.Copy(&buf, rwr); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("ReadPartial->Rewind->ReadFull->Rewind", func(t *testing.T) {
		rwr := New(bytes.NewReader(data))

		var buf bytes.Buffer

		if _, err := io.CopyN(&buf, rwr, 3); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data[:3]; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}

		if err := rwr.Rewind(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		buf.Reset()

		if _, err := io.Copy(&buf, rwr); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}

		if err := rwr.Rewind(); err == nil {
			t.Errorf("got nil, want error")
		}
	})

	t.Run("ReadPartial->Rewind->ReadPartial->Rewind->ReadFull", func(t *testing.T) {
		rwr := New(bytes.NewReader(data))

		for n := 0; n < 2; n++ {
			var buf bytes.Buffer

			if _, err := io.CopyN(&buf, rwr, 3); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got, want := buf.Bytes(), data[:3]; !bytes.Equal(got, want) {
				t.Errorf("got %q, want %q", got, want)
			}

			if err := rwr.Rewind(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		var buf bytes.Buffer

		if _, err := io.Copy(&buf, rwr); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("ReadPartial->Rewind->Rewind->ReadFull", func(t *testing.T) {
		rwr := New(bytes.NewReader(data))

		var buf bytes.Buffer

		if _, err := io.CopyN(&buf, rwr, 3); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data[:3]; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}

		for n := 0; n < 2; n++ {
			if err := rwr.Rewind(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		buf.Reset()

		if _, err := io.Copy(&buf, rwr); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Rewind->ReadFull", func(t *testing.T) {
		rwr := New(bytes.NewReader(data))

		if err := rwr.Rewind(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var buf bytes.Buffer

		if _, err := io.Copy(&buf, rwr); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Rewind->Rewind->ReadFull", func(t *testing.T) {
		rwr := New(bytes.NewReader(data))

		for n := 0; n < 2; n++ {
			if err := rwr.Rewind(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		var buf bytes.Buffer

		if _, err := io.Copy(&buf, rwr); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := buf.Bytes(), data; !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func Example() {
	rwr := New(strings.NewReader("foo bar baz"))

	io.CopyN(os.Stdout, rwr, 3)
	fmt.Println()

	if err := rwr.Rewind(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	io.Copy(os.Stdout, rwr)
	fmt.Println()

	// Output:
	// foo
	// foo bar baz
}
