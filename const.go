package ble

// DefaultMTU defines the default MTU of ATT protocol including 3 bytes of ATT header.
const DefaultMTU = 23

// MaxMTU is maximum of ATT_MTU, which is 512 bytes of value length, plus 3 bytes of ATT header.
// The maximum length of an attribute value shall be 512 octets [Vol 3, Part F, 3.2.9]
const MaxMTU = 512 + 3

// MaxACLMTU is max MTU for ACL-U logical links,
//  The default MTU was selected based on the payload
// carried by two baseband DH5 packets (2*341=682 octets) minus the baseband ACL
// headers (2*2=4 octets) and a 6-octet L2CAP header. Note that the L2CAP header
// length is 4 octets (see Section 3.3.1) but for historical reasons an L2CAP header
// length of 6 bytes is used. [Vol 3, Part A, 5.1]
const MaxACLMTU = 672

// DefaultACLMTU
const DefaultACLMTU = 48

// UUIDs ...
var (
	GAPUUID         = UUID16(0x1800) // Generic Access
	GATTUUID        = UUID16(0x1801) // Generic Attribute
	CurrentTimeUUID = UUID16(0x1805) // Current Time Service
	DeviceInfoUUID  = UUID16(0x180A) // Device Information
	BatteryUUID     = UUID16(0x180F) // Battery Service
	HIDUUID         = UUID16(0x1812) // Human Interface Device

	PrimaryServiceUUID   = UUID16(0x2800)
	SecondaryServiceUUID = UUID16(0x2801)
	IncludeUUID          = UUID16(0x2802)
	CharacteristicUUID   = UUID16(0x2803)

	ClientCharacteristicConfigUUID = UUID16(0x2902)
	ServerCharacteristicConfigUUID = UUID16(0x2903)

	DeviceNameUUID        = UUID16(0x2A00)
	AppearanceUUID        = UUID16(0x2A01)
	PeripheralPrivacyUUID = UUID16(0x2A02)
	ReconnectionAddrUUID  = UUID16(0x2A03)
	PeferredParamsUUID    = UUID16(0x2A04)
	ServiceChangedUUID    = UUID16(0x2A05)
)
