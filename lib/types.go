package lib

// Device minimally identifies an IoT resource
// a devices message typically identifies itself with a UID
// the Meta information can then be accessed and added to the struct as well as Groups, an open array of structural helpers
type Device struct {
	DeviceUID   string
	*DeviceMeta `datastore:",omitempty" firestore:",omitempty"`
	Groups      []Group `datastore:",omitempty" firestore:",omitempty"`
}

// DeviceMeta helps further identify a Device.
// Deployed devices *must* have an UID that identifies the physical unit.
// The metadata can be added to passed messages so a cache is not needed to query the structure needed for microservice
//
// DeviceName is a human-readable identifier for the device which may or may not be unique
// DeviceTag identifies the function of the device in its environment and should remain constant if the device is
// switched out due to failure or replacement but the function of the device remains constant
//
// so for example:
//
//	the cooler in room 4B is still the cooler in 4B even if it is upgraded from PhilipsKA4500 to PhilipsKA4501
//
// CompanyUID and LocationUID are optional structure data that can be used by microservice to store and process data
// without needing to access the device cache

type DeviceMeta struct {
	DeviceName string         `datastore:",omitempty" firestore:",omitempty"`
	DeviceTag  string         `datastore:",omitempty" firestore:",omitempty"`
	DeviceType DeviceType     `datastore:",omitempty" firestore:",omitempty"`
	Listing    *Listing       `datastore:",omitempty" firestore:",omitempty"`
	Version    *DeviceVersion `datastore:",omitempty" firestore:",omitempty"`
}

type Listing struct {
	CompanyUID  string `datastore:",omitempty" firestore:",omitempty"`
	LocationUID string `datastore:",omitempty" firestore:",omitempty"`
}

type DeviceVersion struct {
	FirmwareName    string `datastore:",omitempty" firestore:",omitempty"`
	FirmwareVersion string `datastore:",omitempty" firestore:",omitempty"`
}

// Group provides a basic organizational system that can be added as metadata in a processing pipeline
//
// a Group always has a GroupUID and a GroupType, groups can be flat or, if a groupParentUID is present, hierarchical
type Group struct {
	GroupUID       string
	GroupType      string
	GroupParentUID string `datastore:",omitempty" firestore:",omitempty"`
}

// DeviceSkeleton is used by microservices before they have full admin details for a device - this allows a device's data
// to be stored and traced even if it is not currently identified in admin.
type DeviceSkeleton struct {
	DeviceUID      string
	DeviceCategory Category
	ServiceUID     string
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
