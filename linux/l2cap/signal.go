package l2cap

const (
	// [Vol 3, Part A, 4.3]
	ConnectionResultSuccessful      uint16 = 0x0000
	ConnectionResultPending         uint16 = 0x0001
	ConnectionResultPSMNotSupported uint16 = 0x0002
	ConnectionResultSecurityBlock   uint16 = 0x0003
	ConnectionResultNoResources     uint16 = 0x0004
	ConnectionStatusNoInfo          uint16 = 0x0000
	ConnectionStatusAuthentication  uint16 = 0x0000
	ConnectionStatusAuthorization   uint16 = 0x0000

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
)
