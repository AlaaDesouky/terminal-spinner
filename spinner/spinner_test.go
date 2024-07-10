package spinner_test

import (
	"bytes"
	"io"
	"terminal-spinner/spinner"
	"testing"
	"time"
)

func TestSpinnerStart(t *testing.T) {
	testCases := []struct{
		name string
		duration time.Duration
		expects string
	}{
		{
			name: "should write correct values after 2 frames",
			duration: time.Microsecond * 25,
			expects: "-\b\\\b",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			s := spinner.New(spinner.Config{
				Writer: buf,
				Rate: time.Microsecond * 20,
			})

			s.Start()
			time.Sleep(tc.duration)
			s.Stop()
			
			data, err := io.ReadAll(buf)
			if err != nil {
				t.Errorf(err.Error())
			}

			if string(data) != tc.expects {
				t.Errorf("expected %+v, got %+v", tc.expects, string(data))
			}
		})
	}
}

func TestSpinnerWorksAsync(t *testing.T) {
	buf := &bytes.Buffer{}

	s := spinner.New(spinner.Config{
		Writer: buf,
		Rate: time.Microsecond * 5,
	})

	done := make(chan struct{})

	go func(){
		s.Start()
		time.Sleep(time.Millisecond * 10)
		s.Stop()
		close(done)
	}()

	select {
	case <- time.After(time.Millisecond * 200):
		t.FailNow()
	case <- done:
	}

}

func TestWaitingAfterStop(t *testing.T) {
	testCases := []struct{
		name string
		expects string
	}{
		{
			name: "should stop after waiting",
			expects: "-\b\\\b",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			s := spinner.New(spinner.Config{
				Writer: buf,
				Rate: time.Microsecond * 20,
			})

			s.Start()
			time.Sleep(time.Microsecond * 35)
			s.Stop()
			time.Sleep(time.Microsecond * 10)
			
			data, err := io.ReadAll(buf)
			if err != nil {
				t.Errorf(err.Error())
			}

			if string(data) != tc.expects {
				t.Errorf("expected %+v, got %+v", tc.expects, string(data))
			}
		})
	}
}

func TestStop(t *testing.T){
	t.Run("calling stop on non started spinner should do nothing", func(t *testing.T) {
		buf := &bytes.Buffer{}

		s := spinner.New(spinner.Config{
			Writer: buf,
			Rate: time.Microsecond * 20,
		})
		
		s.Stop()
	})
}

func TestStart(t *testing.T){
	t.Run("calling start on a started spinner should do nothing", func(t *testing.T) {
		buf := &bytes.Buffer{}

		s := spinner.New(spinner.Config{
			Writer: buf,
			Rate: time.Microsecond * 5,
		})

		s.Start()
		s.Start()

		time.Sleep(time.Millisecond * 6)
		s.Stop()
			
		data, err := io.ReadAll(buf)
		if err != nil {
			t.Errorf(err.Error())
		}

		expects := "-\b\\\b"
		if string(data) != expects {
			t.Errorf("expected %+v, got %+v", expects, string(data))
		}
	})
}

func TestRestart(t *testing.T){
	t.Run("calling start on a stopped spinner should restart", func(t *testing.T) {
		buf := &bytes.Buffer{}

		s := spinner.New(spinner.Config{
			Writer: buf,
			Rate: time.Microsecond * 5,
		})

		s.Start()
		s.Stop()
		s.Start()

		time.Sleep(time.Millisecond * 6)
		s.Stop()
			
		data, err := io.ReadAll(buf)
		if err != nil {
			t.Errorf(err.Error())
		}

		expects := "-\b-\b\\\b"
		if string(data) != expects {
			t.Errorf("expected %+v, got %+v", expects, string(data))
		}
	})
}