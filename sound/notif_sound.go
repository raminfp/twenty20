package sound

import (
	"errors"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	// This number represents the fixed frequency of the original PC XT's timer chip, which is approximately 1.193 MHz. This number
	// is divided with the desired frequency to obtain a counter value, that is subsequently fed into the timer chip, tied to the PC speaker.
	clockTickRate = 1193180

	// linux/kd.h, start sound generation (0 for off)
	kiocsound = 0x4B2F

	// linux/input-event-codes.h
	evSnd   = 0x12 // Event type
	sndTone = 0x02 // Sound
)

var (
	// Default frequency, in Hz, middle A
	DefaultFreq = 440.0
	// Default duration in milliseconds
	DefaultDuration = 200
)

// inputEvent represents linux/input.h event structure.
type inputEvent struct {
	Time  syscall.Timeval // time in seconds since epoch at which event occurred
	Type  uint16          // event type
	Code  uint16          // event code related to the event type
	Value int32           // event value related to the event type
}

// ioctl system call manipulates the underlying device parameters of special files.
func ioctl(fd, name, data uintptr) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(name), uintptr(data))
	if e != 0 {
		return e
	}

	return nil
}

// Sound beeps the pc speaker (https://en.wikipedia.org/wiki/PC_speaker).
// On Linux it needs permission to access /dev/tty0 or /dev/input/by-path/platform-pcspkr-event-spkr, and pcspkr module must be loaded.
func Sound(freq float64, duration int) error {
	if freq == 0 {
		freq = DefaultFreq
	} else if freq > 20000 {
		freq = 20000
	} else if freq < 0 {
		freq = DefaultFreq
	}

	if duration == 0 {
		duration = DefaultDuration
	}

	period := int(float64(clockTickRate) / freq)

	var evdev bool

	f, err := os.OpenFile("/dev/tty0", os.O_WRONLY, 0644)
	if err != nil {
		e := err
		f, err = os.OpenFile("/dev/input/by-path/platform-pcspkr-event-spkr", os.O_WRONLY, 0644)
		if err != nil {
			e = errors.New("beeep: " + e.Error() + "; " + err.Error())

			// Output the only beep we can
			_, err = os.Stdout.Write([]byte{7})
			if err != nil {
				return errors.New(e.Error() + "; " + err.Error())
			} else {
				return nil
			}
		}
		evdev = true
	}

	if evdev { // Use Linux evdev API
		ev := inputEvent{}
		ev.Type = evSnd
		ev.Code = sndTone
		ev.Value = int32(freq)

		d := *(*[unsafe.Sizeof(ev)]byte)(unsafe.Pointer(&ev))

		// Start beep
		f.Write(d[:])

		time.Sleep(time.Duration(duration) * time.Millisecond)

		ev.Value = 0
		d = *(*[unsafe.Sizeof(ev)]byte)(unsafe.Pointer(&ev))
		f.Write(d[:])
	} else {
		err = ioctl(f.Fd(), kiocsound, uintptr(period))
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(duration) * time.Millisecond)
		err = ioctl(f.Fd(), kiocsound, uintptr(0))
		if err != nil {
			return err
		}
	}

	return f.Close()
}
