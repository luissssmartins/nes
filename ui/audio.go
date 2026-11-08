package ui

import "github.com/gordonklaus/portaudio"

type Audio struct {
	stream         *portaudio.Stream
	sampleRate     float64
	outputChannels int
	channel        chan float32
	last           float32
}

func NewAudio() *Audio {
	return &Audio{}
}

func (a *Audio) Start() error {
	host, err := portaudio.DefaultHostApi()
	if err != nil {
		return err
	}
	parameters := portaudio.LowLatencyParameters(nil, host.DefaultOutputDevice)
	a.channel = make(chan float32, int(parameters.SampleRate/10))
	stream, err := portaudio.OpenStream(parameters, a.Callback)
	if err != nil {
		return err
	}
	if err := stream.Start(); err != nil {
		return err
	}
	a.stream = stream
	a.sampleRate = parameters.SampleRate
	a.outputChannels = parameters.Output.Channels
	return nil
}

func (a *Audio) Stop() error {
	return a.stream.Close()
}

func (a *Audio) Callback(out []float32) {
	output := a.last
	for i := range out {
		if i%a.outputChannels == 0 {
			select {
			case sample := <-a.channel:
				output = sample
			default:
			}
		}
		out[i] = output
	}
	a.last = output
}
