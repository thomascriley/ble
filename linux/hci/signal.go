package hci

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/thomascriley/ble/linux/hci/cmd"
	"github.com/thomascriley/ble/linux/hci/evt"
	"github.com/thomascriley/ble/linux/l2cap"
)

// Signal ...
type Signal interface {
	Code() int
	Marshal() []byte
	Unmarshal([]byte) error
}

type sigCmd []byte

func (s sigCmd) code() int    { return int(s[0]) }
func (s sigCmd) id() uint8    { return s[1] }
func (s sigCmd) len() int     { return int(binary.LittleEndian.Uint16(s[2:4])) }
func (s sigCmd) data() []byte { return s[4 : 4+s.len()] }

// Signal ...
func (c *Conn) Signal(req Signal, rsp Signal, timeout time.Duration) error {

	logger.Info("Signaling (request: %X, response: %X)\n", req.Code(), rsp.Code())

	// The value of this timer is implementation-dependent but the minimum
	// initial value is 1 second and the maximum initial value is 60 seconds.
	// One RTX timer shall exist for each outstanding signaling request,
	// including each Echo Request. The timer disappears on the final expiration,
	// when the response is received, or the physical link is lost. The maximum
	// elapsed time between the initial start of this timer and the initiation
	// of channel termination (if no response is received) is 60 seconds.
	// [Vol 3, Part A, 6.2.1 ]
	if timeout > 60*time.Second {
		timeout = time.Duration(60 * time.Second)
	}
	if timeout < 1*time.Second {
		timeout = time.Duration(1 * time.Second)
	}

	data := req.Marshal()
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, uint16(4+len(data)))
	binary.Write(buf, binary.LittleEndian, uint16(c.sigCID))
	binary.Write(buf, binary.LittleEndian, uint8(req.Code()))
	binary.Write(buf, binary.LittleEndian, uint8(c.sigID))
	binary.Write(buf, binary.LittleEndian, uint16(len(data)))
	binary.Write(buf, binary.LittleEndian, data)

	// add a buffer of 1 in case the response occurs before we have a chance
	// to wait on sigSent
	c.sigSent = make(chan []byte, 1)
	defer close(c.sigSent)

	if _, err := c.writePDU(buf.Bytes()); err != nil {
		return err
	}

	var s sigCmd
	select {
	case s = <-c.sigSent:
	case <-time.After(timeout):
		return errors.New("signaling request timed out")
	}

	if s.id() != c.sigID {
		return errors.New("mismatched signaling id")
	}
	c.sigID++
	if rsp == nil {
		return nil
	}
	if s.code() != rsp.Code() {
		return errors.New("mismatched signaling response")
	}
	return rsp.Unmarshal(s.data())
}

func (c *Conn) sendResponse(code uint8, id uint8, r Signal) (int, error) {
	data := r.Marshal()
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, uint16(4+len(data)))
	binary.Write(buf, binary.LittleEndian, uint16(c.sigCID))
	binary.Write(buf, binary.LittleEndian, uint8(code))
	binary.Write(buf, binary.LittleEndian, uint8(id))
	binary.Write(buf, binary.LittleEndian, uint16(len(data)))
	if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
		return 0, err
	}
	logger.Debug("sig", "send", fmt.Sprintf("[%X]", buf.Bytes()))
	return c.writePDU(buf.Bytes())
}

func (c *Conn) handleSignal(p pdu) error {

	logger.Debug("sig", "recv", fmt.Sprintf("[%X]", p))

	s := sigCmd(p.payload())

	// When multiple commands are included in an L2CAP packet and the packet
	// exceeds the signaling MTU (MTUsig) of the receiver, a single Command Reject
	// packet shall be sent in response. The identifier shall match the first Request
	// command in the L2CAP packet. If only Responses are recognized, the packet
	// shall be silently discarded. [Vol3, Part A, 4.1]
	if p.dlen() > c.sigRxMTU {
		c.sendResponse(l2cap.SignalCommandReject, s.id(),
			&l2cap.CommandReject{
				Reason:       l2cap.ReasonSignalingMTUExceeded, // Signaling MTU exceeded.
				ActualSigMTU: uint16(c.sigRxMTU)})
		return nil
	}

	var err error
	for len(s) > 0 {
		// Check if it's a supported request.
		switch s.code() {
		case l2cap.SignalCommandReject:
			err = c.handleCommandReject(s)
		case l2cap.SignalConnectionRequest:
			err = c.handleConnectionRequest(s)
		case l2cap.SignalConfigurationRequest:
			err = c.handleConfigurationRequest(s)
		case l2cap.SignalDisconnectRequest:
			err = c.handleDisconnectRequest(s)
		case l2cap.SignalEchoRequest:
			err = c.handleEchoRequest(s)
		case l2cap.SignalInformationRequest:
			err = c.handleInformationRequest(s)
		case l2cap.SignalCreateChannelRequest:
			err = c.handleCreateChannelRequest(s)
		case l2cap.SignalMoveChannelRequest:
			err = c.handleMoveChannelRequest(s)
		case l2cap.SignalConnectionParameterUpdateRequest:
			err = c.handleConnectionParameterUpdateRequest(s)
		case l2cap.SignalLECreditBasedConnectionRequest:
			err = c.LECreditBasedConnectionRequest(s)
		case l2cap.SignalLEFlowControlCredit:
			err = c.LEFlowControlCredit(s)
		case l2cap.SignalConnectionResponse:
			err = c.handleConnectionResponse(s)
		case l2cap.SignalConfigurationResponse:
			err = c.handleConfigurationResponse(s)
		default:
			// Check if it's a response to a sent command.
			select {
			case c.sigSent <- s:
			default:
				_, err = c.sendResponse(
					l2cap.SignalCommandReject,
					s.id(),
					&l2cap.CommandReject{
						Reason: 0x0000, // Command not understood.
					})
			}
		}
		if err != nil {
			return err
		}
		s = s[4+s.len():] // advance to next the packet.
	}
	return nil
}

// handleCommandReject ...
func (c *Conn) handleCommandReject(s sigCmd) error {
	var rej l2cap.CommandReject
	if err := rej.Unmarshal(s.data()); err != nil {
		return err
	}

	switch rej.Reason {
	case l2cap.ReasonInvalidCID:
		return fmt.Errorf("Command rejected: Invalid CID ( resp:[source: %04X, Destination: %04X], conn: [source:%02X, destination:%02X]) c", rej.SourceCID, rej.DestinationCID, c.SourceID, c.DestinationID)
	case l2cap.ReasonCommandNotUnderstood:
		return errors.New("Command rejected: Command not understood")
	case l2cap.ReasonSignalingMTUExceeded:
		return fmt.Errorf("Command rejected: Signaling MTU Exceeded (Actual: %d)", rej.ActualSigMTU)
	default:
		return errors.New("Command rejected for unknown reason")
	}
}

// handleConnectionRequest ...
func (c *Conn) handleConnectionRequest(s sigCmd) error {
	var req l2cap.ConnectionRequest
	if err := req.Unmarshal(s.data()); err != nil {
		return err
	}

	// TODO: Add authentication, PSM check, etc
	_, err := c.sendResponse(l2cap.SignalConnectionResponse, s.id(),
		&l2cap.ConnectionResponse{
			DestinationCID: c.SourceID,
			SourceCID:      c.DestinationID,
			Status:         l2cap.ConnectionStatusNoInfo,
			Result:         l2cap.ConnectionResultSuccessful})
	return err
}

// handleConfigurationRequest ...
func (c *Conn) handleConfigurationRequest(s sigCmd) error {
	rsp := &l2cap.ConfigurationResponse{
		SourceCID: c.DestinationID,
		Flags:     0x0000,
		Result:    l2cap.ConfigurationResultSuccessful,
	}

	var req l2cap.ConfigurationRequest
	if err := req.Unmarshal(s.data()); err != nil {
		rsp.Result = l2cap.ConfigurationResultFailureRejected
		c.sendResponse(l2cap.SignalConfigurationResponse, s.id(), rsp)
		return err
	}

	for _, option := range req.ConfigurationOptions {
		switch option.Type() {
		case l2cap.MTUOptionType:
			c.txMTU = int(option.(*l2cap.MTUOption).MTU)
		default:
			if option.Hint() == 0x00 {
				rsp.Result = l2cap.ConfigurationResultFailureUnknown
			}
		}
		rsp.ConfigurationOptions = append(rsp.ConfigurationOptions, option)
	}
	if _, err := c.sendResponse(l2cap.SignalConfigurationResponse, s.id(), rsp); err != nil {
		return err
	}

	select {
	case <-c.cfgRequest:
	default:
		close(c.cfgRequest)
	}
	return nil
}

// handleEchoRequest ...
func (c *Conn) handleEchoRequest(s sigCmd) error {
	// TODO: Allow user to supply own data response
	_, err := c.sendResponse(l2cap.SignalEchoResponse, s.id(),
		&l2cap.EchoResponse{Data: 0x00})
	return err
}

// handleInformationRequest ...
func (c *Conn) handleInformationRequest(s sigCmd) error {
	var req l2cap.InformationRequest
	if err := req.Unmarshal(s.data()); err != nil {
		return err
	}

	rsp := &l2cap.InformationResponse{
		InfoType: req.InfoType,
		Result:   l2cap.InfoResponseResultSuccess}

	switch req.InfoType {
	case l2cap.InfoTypeConnectionlessMTU:
		rsp.ConnectionlessMTU = uint16(c.txMTU)
	case l2cap.InfoTypeExtendedFeatures:
		rsp.ExtendedFeatureMask = c.extendedFeatures
	case l2cap.InfoTypeFixedChannels:
		rsp.FixedChannels = c.fixedChannels
	default:
		rsp.Result = l2cap.InfoResponseResultNotSupported
	}
	_, err := c.sendResponse(l2cap.SignalInformationResponse, s.id(), rsp)
	return err
}

// handleCreateChannelRequest ...
func (c *Conn) handleCreateChannelRequest(s sigCmd) error {
	var req l2cap.CreateChannelRequest
	if err := req.Unmarshal(s.data()); err != nil {
		return err
	}

	// TODO: Add authentication, PSM check, creating the channel, etc
	_, err := c.sendResponse(l2cap.SignalCreateChannelResponse, s.id(),
		&l2cap.CreateChannelResponse{
			DestinationCID: c.SourceID,
			SourceCID:      c.DestinationID,
			Status:         l2cap.CreateChannelStatusNoInfo,
			Result:         l2cap.CreateChannelResultSuccessful})
	return err
}

// handleMoveChannelRequest ...
func (c *Conn) handleMoveChannelRequest(s sigCmd) error {
	var req l2cap.MoveChannelRequest
	if err := req.Unmarshal(s.data()); err != nil {
		return err
	}

	// TODO: check for enhanced retransmission mode or streaming mode to allow
	// changing of the cids
	_, err := c.sendResponse(l2cap.SignalMoveChannelResponse, s.id(),
		&l2cap.MoveChannelResponse{
			InitiatorCID: req.InitiatorCID,
			Result:       l2cap.MoveChannelResultNotAllowed})
	return err
}

// DisconnectRequest implements Disconnect Request (0x06) [Vol 3, Part A, 4.6].
func (c *Conn) handleDisconnectRequest(s sigCmd) error {
	var req l2cap.DisconnectRequest
	if err := req.Unmarshal(s.data()); err != nil {
		return err
	}

	// Send Command Reject when the DCID is unrecognized.
	if req.DestinationCID != cidLEAtt {
		_, err := c.sendResponse(l2cap.SignalCommandReject, s.id(),
			&l2cap.CommandReject{
				Reason:         l2cap.ReasonInvalidCID,
				SourceCID:      req.SourceCID,
				DestinationCID: req.DestinationCID})
		return err
	}

	// Silently discard the request if SCID failed to find the same match.
	if req.SourceCID != cidLEAtt {
		return nil
	}

	_, err := c.sendResponse(l2cap.SignalDisconnectResponse, s.id(),
		&l2cap.DisconnectResponse{
			DestinationCID: req.DestinationCID,
			SourceCID:      req.SourceCID})
	return err
}

// handleConnectionResponse ...
func (c *Conn) handleConnectionResponse(s sigCmd) error {
	var rsp l2cap.ConnectionResponse
	if err := rsp.Unmarshal(s.data()); err != nil {
		return err
	}

	// wait for a non pending result
	if rsp.Result == l2cap.ConnectionResultPending {
		switch rsp.Status {
		case l2cap.ConnectionStatusAuthentication:
		case l2cap.ConnectionStatusAuthorization:
		case l2cap.ConnectionStatusNoInfo:
		}
		return nil
	}

	c.DestinationID = rsp.DestinationCID

	select {
	case c.sigSent <- s:
	default:
	}
	return nil
}

// handleConfigurationResponse ...
func (c *Conn) handleConfigurationResponse(s sigCmd) error {
	rsp := &l2cap.ConfigurationResponse{}
	if err := rsp.Unmarshal(s.data()); err != nil {
		logger.Error("Configuration Response Error: %s\n", err)
		c.Close()
		return err
	}

	// wait for a non pending result
	if rsp.Result == l2cap.ConfigurationResultPending {
		return nil
	}

	select {
	case c.sigSent <- s:
	default:
		logger.Error("Configuration Response error: signal channel buffer full\n")
	}
	return nil
}

// ConnectionParameterUpdateRequest implements Connection Parameter Update Request (0x12) [Vol 3, Part A, 4.20].
func (c *Conn) handleConnectionParameterUpdateRequest(s sigCmd) error {
	// This command shall only be sent from the LE slave device to the LE master
	// device and only if one or more of the LE slave Controller, the LE master
	// Controller, the LE slave Host and the LE master Host do not support the
	// Connection Parameters Request Link Layer Control Procedure ([Vol. 6] Part B,
	// Section 5.1.7). If an LE slave Host receives a Connection Parameter Update
	// Request packet it shall respond with a Command Reject packet with reason
	// 0x0000 (Command not understood).
	if c.param.Role() != roleMaster {
		_, err := c.sendResponse(
			l2cap.SignalCommandReject,
			s.id(),
			&l2cap.CommandReject{
				Reason: 0x0000, // Command not understood.
			})
		return err
	}
	var req l2cap.ConnectionParameterUpdateRequest
	if err := req.Unmarshal(s.data()); err != nil {
		return err
	}

	// LE Connection Update (0x08|0x0013) [Vol 2, Part E, 7.8.18]
	return c.hci.Send(&cmd.LEConnectionUpdate{
		ConnectionHandle:   c.param.ConnectionHandle(),
		ConnIntervalMin:    req.IntervalMin,
		ConnIntervalMax:    req.IntervalMax,
		ConnLatency:        req.SlaveLatency,
		SupervisionTimeout: req.TimeoutMultiplier,
		MinimumCELength:    0, // Informational, and spec doesn't specify the use.
		MaximumCELength:    0, // Informational, and spec doesn't specify the use.
	}, nil)
}

func (c *Conn) handleLEConnectionUpdateComplete(e evt.LEConnectionUpdateComplete) error {
	// Currently, we (as a slave host) accept all the parameters and forward
	// it to the controller. The controller might update all, partial or even
	// none (ignore) of the parameters. The slave(remote) host will be indicated
	// by its controller if the update actually happens.
	// TODO: allow users to implement what parameters to accept.
	_, err := c.sendResponse(
		l2cap.SignalConnectionParameterUpdateResponse,
		c.sigID,
		&l2cap.ConnectionParameterUpdateResponse{
			Result: 0, // Accept.
		})
	return err
}

// LECreditBasedConnectionRequest ...
func (c *Conn) LECreditBasedConnectionRequest(s sigCmd) error {
	// TODO:
	return nil
}

// LEFlowControlCredit ...
func (c *Conn) LEFlowControlCredit(s sigCmd) error {
	// TODO:
	return nil
}

// InformationRequest [Vol 3, Part A, 4.10]
func (c *Conn) InformationRequest(infoType uint16, timeout time.Duration) error {
	req := &l2cap.InformationRequest{}
	rsp := &l2cap.InformationResponse{}
	req.InfoType = infoType
	if err := c.Signal(req, rsp, timeout); err != nil {
		return err
	}

	switch infoType {
	case l2cap.InfoTypeConnectionlessMTU:
		c.SetTxMTU(int(rsp.ConnectionlessMTU))
	case l2cap.InfoTypeExtendedFeatures:
		c.extendedFeatures = rsp.ExtendedFeatureMask
	case l2cap.InfoTypeFixedChannels:
		c.fixedChannels = rsp.FixedChannels
	default:
		return errors.New("Invalid infoType")
	}
	return nil
}

// ConnectionRequest [Vol 3, Part A, 4.2]
func (c *Conn) ConnectionRequest(psm uint16, timeout time.Duration) error {
	rsp := &l2cap.ConnectionResponse{}
	req := &l2cap.ConnectionRequest{
		PSM:       psm,
		SourceCID: c.SourceID,
	}
	if err := c.Signal(req, rsp, timeout); err != nil {
		return err
	}
	switch rsp.Result {
	case l2cap.ConnectionResultSuccessful:
	case l2cap.ConnectionResultPending:
		// should never get here since pending results are already handled
	case l2cap.ConnectionResultPSMNotSupported:
		return errors.New("Connection refused - PSM is not supported")
	case l2cap.ConnectionResultNoResources:
		return errors.New("Connection refused - No resources available")
	case l2cap.ConnectionResultSecurityBlock:
		return errors.New("Connection refused - Security block")
	}
	c.DestinationID = rsp.DestinationCID
	return nil
}

// ConfigurationRequest [Vol 3, Part A, 4.4]
func (c *Conn) ConfigurationRequest(options []l2cap.Option, timeout time.Duration) error {
	i := 0

	rsp := &l2cap.ConfigurationResponse{
		Flags: 0x0001,
	}
	req := &l2cap.ConfigurationRequest{
		DestinationCID: c.DestinationID,
	}

	// the options need to be split into chunks and sent
	for i < len(options) || rsp.Flags == 0x0001 {

		// fill the request with options
		length := 0
		for ; i < len(options); i++ {
			b, err := options[i].MarshalBinary()
			if err != nil {
				return err
			}
			if length+len(b) > c.sigTxMTU-8 {
				break
			}

			req.ConfigurationOptions = append(req.ConfigurationOptions, options[i])
			length = length + len(b)
		}

		// if extended flow specification is enabled, continuation bit is 0
		// otherwise if the options can not fit into one request the continuation
		// bit is 1
		if c.extendedFeatures&(l2cap.ExtendedFeatureExtendedFlowSpecification+1) == l2cap.ExtendedFeatureExtendedFlowSpecification+1 {
			req.Flags = 0x0000
		} else if i < len(options) {
			req.Flags = 0x0001
		} else {
			req.Flags = 0x0000
		}

		if err := c.Signal(req, rsp, timeout); err != nil {
			return err
		}

		switch rsp.Result {
		case l2cap.ConfigurationResultSuccessful:

		case l2cap.ConfigurationResultFailureUnacceptable:
			return errors.New("Failure - unacceptable parameters")
		case l2cap.ConfigurationResultFailureRejected:
			return errors.New("Failure - rejected (no reason provided)")
		case l2cap.ConfigurationResultFailureUnknown:
			return errors.New("Failure - unknown options")
		case l2cap.ConfigurationResultPending:

			// should never get here, pending results are prehandled
		case l2cap.ConfigurationResultFailureFlowSpecRejected:
			return errors.New("Failure - flow spec rejected")
		}
	}
	return nil
}
