package fanuc

import (
	"fmt"
	"testing"
)

func TestParseFrames(t *testing.T) {
	src := `Tool Frame
   0.0     0.0     0.0     0.0     0.0     0.0 Eoat1
   1.0     0.0     0.0     0.0     0.0     0.0 Eoat2
   2.0     0.0     0.0     0.0     0.0     0.0 Eoat3
   3.0     0.0     0.0     0.0     0.0     0.0 Eoat4
   4.0     0.0     0.0     0.0     0.0     0.0 Eoat5
   5.0     0.0     0.0     0.0     0.0     0.0 Eoat6
   6.0     0.0     0.0     0.0     0.0     0.0 Eoat7
   7.0     0.0     0.0     0.0     0.0     0.0 Eoat8
   8.0     0.0     0.0     0.0     0.0     0.0 Eoat9
   9.0     0.0     0.0     0.0     0.0     0.0 Eoat10

Jog Frame
   0.0     1.0     0.0     0.0     0.0     0.0 Jog1
   0.0     2.0     0.0     0.0     0.0     0.0 Jog2
   0.0     3.0     0.0     0.0     0.0     0.0 Jog3
   0.0     4.0     0.0     0.0     0.0     0.0 Jog4
   0.0     5.0     0.0     0.0     0.0     0.0 Jog5

User Frame
   0.0     0.0     6.0     0.0     0.0     0.0 UFrame1
   0.0     0.0     7.0     0.0     0.0     0.0 UFrame2
   0.0     0.0     8.0     0.0     0.0     0.0 UFrame3
   0.0     0.0     9.0     0.0     0.0     0.0 UFrame4
   0.0     0.0    10.0     0.0     0.0     0.0 UFrame5
   0.0     0.0    11.0     0.0     0.0     0.0 UFrame6
   0.0     0.0    12.0     0.0     0.0     0.0 UFrame7
   0.0     0.0    13.0     0.0     0.0     0.0 UFrame8
   0.0     0.0    14.0     0.0     0.0     0.0 UFrame9`

	frames, err := parseFrames(src)
	if err != nil {
		t.Fatal(err)
	}

	if len(frames) != 24 {
		t.Fatalf("Got %d frames. Want 24", len(frames))
	}

	// tool frames
	for i := 0; i < 10; i++ {
		tf := frames[i]
		if tf.Type != ToolFrame {
			t.Errorf("Bad type. Got %s, want %s", tf.Type, ToolFrame)
		}
		if tf.Id != i+1 {
			t.Errorf("Bad id. Got %d, want %d", tf.Id, i+1)
		}
		x := fmt.Sprintf("%d.0", i)
		if tf.X != x {
			t.Errorf("Bad x. Got %s, want %s", tf.X, x)
		}
		comment := fmt.Sprintf("Eoat%d", i+1)
		if tf.Comment != comment {
			t.Errorf("Bad comment. Got %q, want %q", tf.Comment, comment)
		}
	}

	// jog frames
	for i := 10; i < 15; i++ {
		f := frames[i]
		if f.Type != JogFrame {
			t.Errorf("Bad type. Got %s, want %s", f.Type, JogFrame)
		}
		if f.Id != i-9 {
			t.Errorf("Bad id. Got %d, want %d", f.Id, i-9)
		}
		y := fmt.Sprintf("%d.0", i-9)
		if f.Y != y {
			t.Errorf("Bad y. Got %s, want %s", f.Y, y)
		}
		comment := fmt.Sprintf("Jog%d", i-9)
		if f.Comment != comment {
			t.Errorf("Bad comment. Got %q, want %q", f.Comment, comment)
		}
	}

	// user frames
	for i := 15; i < 24; i++ {
		f := frames[i]
		if f.Type != UserFrame {
			t.Errorf("Bad type. Got %s, want %s", f.Type, UserFrame)
		}
		if f.Id != i-14 {
			t.Errorf("Bad id. Got %d, want %d", f.Id, i-14)
		}
		z := fmt.Sprintf("%d.0", i-9)
		if f.Z != z {
			t.Errorf("Bad z. Got %s, want %s", f.Z, z)
		}
		comment := fmt.Sprintf("UFrame%d", i-14)
		if f.Comment != comment {
			t.Errorf("Bad comment. Got %q, want %q", f.Comment, comment)
		}
	}
}

func TestParseFramesError(t *testing.T) {
	src := `Tool Frame
   0.0     0.0     0.0     0.0     0.0     0.0 Eoat1
   1.0     0.0     0.0     0.0     0.0     0.0 Eoat2
   2.0     0.0     0.0     0.0     0.0     0.0 Eoat3
   3.0     0.0     0.0     0.0     0.0     0.0 Eoat4
   4.0     0.0     0.0     0.0     0.0     0.0 Eoat5
   5.0     0.0     0.0     0.0     0.0     0.0 Eoat6
   6.0     0.0     0.0     0.0     0.0     0.0 Eoat7
   7.0     0.0     0.0     0.0     0.0     0.0 Eoat8
   8.0     0.0     0.0     0.0     0.0     0.0 Eoat9
   9.0     0.0     0.0     0.0     0.0     0.0 Eoat10`
	_, err := parseFrames(src)
	if err == nil {
		t.Fatal("want an error")
	}
	if err.Error() != "Invalid frame.dg" {
		t.Fatalf("invalid error message. Got %q", err.Error())
	}

}
