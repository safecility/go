package lib

// Device minimally identifies an IoT resource
// a devices message typically identifies itself with a UID
// the Meta information can then be accessed and added to the struct as well as Groups, an open array of structural helpers
type Device struct {
	DeviceUID   string `firestore:"uid,omitempty" json:"uid"`
	*DeviceMeta `datastore:",omitempty" firestore:",omitempty" json:",omitempty"`
}

// DeviceMeta helps further identify a Device.
// Deployed devices *must* have an UID that identifies the physical unit.
//
// # The metadata can be added to passed messages so a cache is not needed to query the structure needed for microservice
//
// DeviceName is a human-readable identifier for the device which may or may not be unique
// DeviceTag identifies the function of the device in its environment and should remain constant if the device is
// switched out due to failure or replacement but the function of the device remains constant
//
// so for example:
//
//	the cooler in room 4B is still the cooler in 4B even if it is upgraded from PhilipsKA4500 to PhilipsKA4501
//
// CompanyUID and LocationUID remain as basic organizational elements
type DeviceMeta struct {
	DeviceName string     `datastore:",omitempty" firestore:"name,omitempty" json:"deviceName,omitempty"`
	DeviceTag  string     `datastore:",omitempty" firestore:"tag,omitempty" json:"deviceTag,omitempty"`
	DeviceType DeviceType `datastore:",omitempty" firestore:"type,omitempty" json:"type,omitempty"`

	CompanyUID  string `datastore:",omitempty" firestore:"companyUID,omitempty" json:"companyUID,omitempty"`
	LocationUID string `datastore:",omitempty" firestore:"locationUID,omitempty" json:"locationUID,omitempty"`

	Version    *Firmware  `datastore:",omitempty" firestore:"firmware,omitempty" json:"firmware,omitempty"`
	Processors *Processor `datastore:",omitempty" firestore:"processors,omitempty" json:"processors,omitempty"`
}

// Processor allows implementations to save a number of identifiers for different processing options
type Processor map[string]interface{}

type Firmware struct {
	FirmwareName    string `datastore:",omitempty" firestore:"name,omitempty" json:"name,omitempty"`
	FirmwareVersion string `datastore:",omitempty" firestore:"version,omitempty" json:"version,omitempty"`
}

// DeviceSkeleton is used by microservices before they have full admin details for a device - this allows a device's data
// to be stored and traced even if it is not currently identified in admin.
type DeviceSkeleton struct {
	DeviceUID      string
	ServiceUID     string
	DeviceCategory Category
}

type DeviceType string
type DeviceGroup string

type Category struct {
	DeviceType
	DeviceGroup
}

const (
	Power    DeviceType  = "power"
	Lighting DeviceType  = "lighting"
	Meter    DeviceGroup = "meter"
	DaliEL   DeviceGroup = "daliEL"
)
