package l2cap

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	// [Vol 3, Part A, 4.3]
	ReasonCommandNotUnderstood uint16 = 0x0000
	ReasonSignalingMTUExceeded uint16 = 0x0001
	ReasonInvalidCID           uint16 = 0x0002

	// [Vol 3, Part A, 4.3]
	ConnectionResultSuccessful      uint16 = 0x0000
	ConnectionResultPending         uint16 = 0x0001
	ConnectionResultPSMNotSupported uint16 = 0x0002
	ConnectionResultSecurityBlock   uint16 = 0x0003
	ConnectionResultNoResources     uint16 = 0x0004
	ConnectionStatusNoInfo          uint16 = 0x0000
	ConnectionStatusAuthentication  uint16 = 0x0001
	ConnectionStatusAuthorization   uint16 = 0x0002

	// [Vol 3, Part A, 4.5]
	ConfigurationResultSuccessful              uint16 = 0x0000
	ConfigurationResultFailureUnacceptable     uint16 = 0x0001
	ConfigurationResultFailureRejected         uint16 = 0x0002
	ConfigurationResultFailureUnknown          uint16 = 0x0003
	ConfigurationResultPending                 uint16 = 0x0004
	ConfigurationResultFailureFlowSpecRejected uint16 = 0x0005

	// [Vol 3, Part A, 4.10]
	InfoTypeConnectionlessMTU uint16 = 0x0001
	InfoTypeExtendedFeatures  uint16 = 0x0002
	InfoTypeFixedChannels     uint16 = 0x0003

	// [Vol 3, Part A, 4.11]
	InfoResponseResultSuccess      uint16 = 0x0000
	InfoResponseResultNotSupported uint16 = 0x0001

	// [Vol 3, Part A, 4.12]
	ExtendedFeatureFlowControlModeSupported            uint32 = 0x00000001
	ExtendedFeatureRetransmissionModeSupported         uint32 = 0x00000002
	ExtendedFeatureBidirectionalQoSSupported           uint32 = 0x00000004
	ExtendedFeatureEnhancedRetransmissionModeSupported uint32 = 0x00000008
	ExtendedFeatureStreamingMode                       uint32 = 0x00000010
	ExtendedFeatureFCSOption                           uint32 = 0x00000020
	ExtendedFeatureExtendedFlowSpecification           uint32 = 0x00000040
	ExtendedFeatureFixedChannels                       uint32 = 0x00000080
	ExtendedFeatureExtendedWindowSize                  uint32 = 0x00000100
	ExtendedFeatureUnicastDataReception                uint32 = 0x00000200

	// [Vol 3, Part A, 4.15]
	CreateChannelResultSuccessful             uint16 = 0x0000
	CreateChannelResultPending                uint16 = 0x0001
	CreateChannelResultPSMNotSupported        uint16 = 0x0002
	CreateChannelResultSecurityBlock          uint16 = 0x0003
	CreateChannelResultNoResources            uint16 = 0x0004
	CreateChannelResultControllerNotSupported uint16 = 0x0005
	CreateChannelStatusNoInfo                 uint16 = 0x0000
	CreateChannelStatusAuthentication         uint16 = 0x0001
	CreateChannelStatusAuthorization          uint16 = 0x0002

	// [Vol 3, Part A, 4.17]
	MoveChannelResultSuccess                uint16 = 0x0000
	MoveChannelResultPending                uint16 = 0x0001
	MoveChannelResultControllerNotSupported uint16 = 0x0002
	MoveChannelResultControllerSame         uint16 = 0x0003
	MoveChannelResultConfigNotSupported     uint16 = 0x0004
	MoveChannelResultCollision              uint16 = 0x0005
	MoveChannelResultNotAllowed             uint16 = 0x0006

	// [ Vol 3, Part A, 3.2]
	DefaultConnectionlessMTU uint16 = 0x0030
)

// Marshal Serializes the struct into binary data int LittleEndian order
func (s *CommandReject) Marshal() []byte {
	var b []byte
	switch s.Reason {
	case ReasonCommandNotUnderstood:
		b = make([]byte, 2)
	case ReasonSignalingMTUExceeded:
		b = make([]byte, 4)
		binary.LittleEndian.PutUint16(b[2:], s.ActualSigMTU)
	case ReasonInvalidCID:
		b = make([]byte, 6)
		binary.LittleEndian.PutUint16(b[2:], s.SourceCID)
		binary.LittleEndian.PutUint16(b[4:], s.DestinationCID)
	}
	binary.LittleEndian.PutUint16(b[0:], s.Reason)
	return b
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *CommandReject) Unmarshal(b []byte) error {
	switch s.Reason = binary.LittleEndian.Uint16(b[0:]); s.Reason {
	case ReasonCommandNotUnderstood:
	case ReasonSignalingMTUExceeded:
		s.ActualSigMTU = binary.LittleEndian.Uint16(b[2:])
	case ReasonInvalidCID:
		s.SourceCID = binary.LittleEndian.Uint16(b[2:])
		s.DestinationCID = binary.LittleEndian.Uint16(b[4:])
	}
	return nil
}

// Marshal Serializes the struct into binary data int LittleEndian order
func (s *ConfigurationRequest) Marshal() []byte {
	length := uint8(6)
	for _, option := range s.ConfigurationOptions {
		length = length + option.Len() + 2
	}
	b := make([]byte, 4, length)

	binary.LittleEndian.PutUint16(b, s.DestinationCID)
	binary.LittleEndian.PutUint16(b[2:], s.Flags)

	for _, option := range s.ConfigurationOptions {
		bo, err := option.MarshalBinary()
		if err != nil {
			logger.Error("Could not marshal option: %s", err)
			continue
		}
		b = append(b, bo...)
	}
	return b
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *ConfigurationRequest) Unmarshal(b []byte) error {
	s.DestinationCID = binary.LittleEndian.Uint16(b[0:])
	s.Flags = binary.LittleEndian.Uint16(b[2:])

	for i := 4; i < len(b); {
		var option Option
		switch optionType := optionTypeFromTypeHint(b[i]); optionType {
		case MTUOptionType:
			option = &MTUOption{}
		case FlushTimeoutOptionType:
			option = &FlushTimeoutOption{}
		case QoSOptionType:
			option = &QoSOption{}
		case RetransmissionAndFlowControlOptionType:
			option = &RetransmissionAndFlowControlOption{}
		case FrameCheckSequenceOptionType:
			option = &FrameCheckSequenceOption{}
		case ExtendedFlowSpecificationOptionType:
			option = &ExtendedFlowSpecificationOption{}
		case ExtendedWindowSizeOptionType:
			option = &ExtendedWindowSizeOption{}
		default:
			logger.Error("Option error: unknown option type %X", optionType)
			i = i + int(b[i+1]) + 2
			continue
		}

		if err := option.UnmarshalBinary(b[i:]); err != nil {
			logger.Error("Could not unmarshal option: %s", err)
			return err
		}
		s.ConfigurationOptions = append(s.ConfigurationOptions, option)
		i = i + int(option.Len()) + 2
	}
	return nil
}

// Marshal Serializes the struct into binary data int LittleEndian order
func (s *ConfigurationResponse) Marshal() []byte {
	length := uint8(6)
	for _, option := range s.ConfigurationOptions {
		length = length + option.Len() + 2
	}
	b := make([]byte, 6, length)

	binary.LittleEndian.PutUint16(b, s.SourceCID)
	binary.LittleEndian.PutUint16(b[2:], s.Flags)
	binary.LittleEndian.PutUint16(b[4:], s.Result)

	for _, option := range s.ConfigurationOptions {
		bo, err := option.MarshalBinary()
		if err != nil {
			logger.Error("Error marshalling option: %s", err)
			continue
		}
		b = append(b, bo...)
	}
	return b
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *ConfigurationResponse) Unmarshal(b []byte) error {
	s.SourceCID = binary.LittleEndian.Uint16(b[0:])
	s.Flags = binary.LittleEndian.Uint16(b[2:])
	s.Result = binary.LittleEndian.Uint16(b[4:])

	for i := 6; i < len(b); {
		var option Option
		switch uint8(b[i] & 0x7F) {
		case MTUOptionType:
			option = &MTUOption{}
		case FlushTimeoutOptionType:
			option = &FlushTimeoutOption{}
		case QoSOptionType:
			option = &QoSOption{}
		case RetransmissionAndFlowControlOptionType:
			option = &RetransmissionAndFlowControlOption{}
		case FrameCheckSequenceOptionType:
			option = &FrameCheckSequenceOption{}
		case ExtendedFlowSpecificationOptionType:
			option = &ExtendedFlowSpecificationOption{}
		case ExtendedWindowSizeOptionType:
			option = &ExtendedWindowSizeOption{}
		default:
			i = i + int(b[i+1]) + 2
			continue
		}

		option.SetHint(uint8(b[i] >> 7 & 0x01))
		if err := option.UnmarshalBinary(b[i+2:]); err != nil {
			return err
		}
		s.ConfigurationOptions = append(s.ConfigurationOptions, option)
		i = i + int(b[i+1]) + 2
	}
	return nil
}

// Marshal Serializes the struct into binary data int LittleEndian order
func (s *InformationResponse) Marshal() []byte {
	var b []byte
	switch s.InfoType {
	case InfoTypeConnectionlessMTU:
		b = make([]byte, 6)
		binary.LittleEndian.PutUint16(b[4:], s.ConnectionlessMTU)
	case InfoTypeExtendedFeatures:
		b = make([]byte, 8)
		binary.LittleEndian.PutUint32(b[4:], s.ExtendedFeatureMask)
	case InfoTypeFixedChannels:
		b = make([]byte, 12)
		binary.LittleEndian.PutUint64(b[4:], s.FixedChannels)
	}
	binary.LittleEndian.PutUint16(b[0:], s.InfoType)
	binary.LittleEndian.PutUint16(b[2:], s.Result)
	return b
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *InformationResponse) Unmarshal(b []byte) error {
	s.InfoType = binary.LittleEndian.Uint16(b[0:])
	s.Result = binary.LittleEndian.Uint16(b[2:])

	switch s.Result {
	case InfoResponseResultNotSupported:
		return errors.New("Not supported")
	case InfoResponseResultSuccess:
		switch s.InfoType {
		case InfoTypeConnectionlessMTU:
			if len(b) > 5 {
				s.ConnectionlessMTU = binary.LittleEndian.Uint16(b[4:])
			}
			if s.ConnectionlessMTU == 0 {
				s.ConnectionlessMTU = DefaultConnectionlessMTU
			}
		case InfoTypeExtendedFeatures:
			if len(b) > 7 {
				s.ExtendedFeatureMask = binary.LittleEndian.Uint32(b[4:])
			}
		case InfoTypeFixedChannels:
			if len(b) > 11 {
				s.FixedChannels = binary.LittleEndian.Uint64(b[4:])
			}
		}
		return nil
	default:
		return fmt.Errorf("Unknown result: %X", s.Result)
	}
}
