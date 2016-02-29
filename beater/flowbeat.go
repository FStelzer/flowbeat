package beater

import (
	"bytes"
	"net"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/fstelzer/sflow"
)

type Flowbeat struct {
	FbConfig ConfigSettings
	events   publisher.Client

	listen string
	conn   *net.UDPConn

	done chan struct{}
}

func New() *Flowbeat {
	return &Flowbeat{}
}

func (fb *Flowbeat) Config(b *beat.Beat) error {

	err := cfgfile.Read(&fb.FbConfig, "")
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}

	if fb.FbConfig.Input.Listen != nil {
		fb.listen = *fb.FbConfig.Input.Listen
	} else {
		fb.listen = ":6343"
	}

	logp.Debug("flowbeat", "Init flowbeat")
	logp.Debug("flowbeat", "Listening on %s\n", fb.listen)

	return nil
}

func (fb *Flowbeat) Setup(b *beat.Beat) error {
	fb.events = b.Events
	fb.done = make(chan struct{})

	addr, err := net.ResolveUDPAddr("udp", fb.listen)
	if err != nil {
		return err
	}
	fb.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	return nil
}

func (fb *Flowbeat) Run(b *beat.Beat) error {
	var err error
	packetbuffer := make([]byte, 65535)
	reader := bytes.NewReader(packetbuffer)
	decoder := sflow.NewDecoder(reader)

	for {
		select {
		case <-fb.done:
			return nil
		default:
		}

		reader.Seek(0, 0) //Reset the reader on our buffer

		// Listen for sflow datagrams
		size, addr, err := fb.conn.ReadFromUDP(packetbuffer)
		logp.Debug("flowbeat", "Received UDP Packet with Size: %d", size)
		if err != nil {
			return err
		}

		dgram, err := decoder.Decode()
		if err != nil {
			logp.Warn("Error decoding sflow packet: %s", err)
			continue
		}

		for _, sample := range dgram.Samples {
			event := common.MapStr{
				"@timestamp":     common.Time(time.Now()),
				"datagramSource": addr.IP,
				"agent":          dgram.IpAddress,
				"subAgentId":     dgram.SubAgentId,
				"sequenceNumber": dgram.SequenceNumber,
				"uptime":         dgram.Uptime,
			}

			switch sample.SampleType() {
			case sflow.TypeFlowSample:
				event["type"] = "flow"
				sample := sample.(*sflow.FlowSample)
				event["sequenceNum"] = sample.SequenceNum
				event["samplingRate"] = sample.SamplingRate
				event["samplePool"] = sample.SamplePool
				event["drops"] = sample.Drops
				event["input"] = sample.Input
				event["output"] = sample.Output
			case sflow.TypeCounterSample:
				event["type"] = "counter"
				sample := sample.(*sflow.CounterSample)
				event["sequenceNum"] = sample.SequenceNum
			case sflow.TypeExpandedFlowSample:
				event["type"] = "extended_flow"
			case sflow.TypeExpandedCounterSample:
				event["type"] = "extended_counter"
			default:
				event["type"] = "unknown"
			}

			for _, record := range sample.GetRecords() {
				event[record.RecordName()] = record
			}

			fb.events.PublishEvent(event)
		}
	}

	return err
}

func (fb *Flowbeat) Cleanup(b *beat.Beat) error {
	if fb.conn != nil {
		fb.conn.Close()
	}
	return nil
}

func (fb *Flowbeat) Stop() {
	close(fb.done)
}
