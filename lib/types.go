package lib

// Group our types use this slightly cumbersome naming convention for attributes to avoid erasure of values in datastore
// and other mechanisms that flatten structures
type Group struct {
	GroupID       int64
	GroupParentID int64
	GroupChildren []Group
}

type Device struct {
	DeviceUID  string
	DeviceName string
	CompanyID  int64
	*Group
}
