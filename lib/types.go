package lib

// Device minimally identifies an IoT resource
// a devices message typically identifies itself with a UID
// the Meta information can then be accessed and added to the struct
type Device struct {
	DeviceUID   string
	*DeviceMeta `datastore:",omitempty"`
	*Group      `datastore:",omitempty"`
}

// DeviceMeta helps further identify a Device.
// Deployed devices typically have an UID that identifies the physical unit.
// DeviceName is a human readable identifier for the device which may or may not be unique
// DeviceTag identifies the function of the device in its environment and should remain constant is the device is
// switched out
// due to failure or replacement but the function of the device remains constant: so for example the cooler in room 4B is still
// the cooler in 4B even if it is upgraded from PhilipsKA4500 to PhilipsKA4501
type DeviceMeta struct {
	DeviceName string
	DeviceTag  string
}

// Group provides Devices with a basic organizational system that can be easily cached and added as metadata early in a
// processing pipeline
// our types use this slightly cumbersome naming convention for attributes to avoid erasure of values in datastore
// and other mechanisms that flatten structures
// a Group always has a GroupUID and belongs to a Company - if the group is a top level group then the ParentUID is
// the same as the CompanyUID or empty - group structures do not need to represent all or any of their children but may
// choose to do so (children may instead name their ParentUID and associate in this way instead)
type Group struct {
	GroupUID   string `datastore:",omitempty"`
	CompanyUID string `datastore:",omitempty"`
	ParentUID  string `datastore:",omitempty"`
	Children   []Group
}
