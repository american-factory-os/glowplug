package sparkplug

import (
	"fmt"
	"strings"
)

var (
	ErrInvalidTopic = fmt.Errorf("invalid sparkplug topic")
	ErrTopicEmpty   = fmt.Errorf("topic is empty")
)

// Topic represents values in a sparkplug topic.
// ex: spBv1.0/groupId/DDATA/edgeNodeId/deviceId
type Topic struct {
	Command     Command `json:"command"`
	GroupId     string  `json:"group_id"`
	EdgeNodeId  string  `json:"edge_node_id"`
	DeviceId    string  `json:"device_id"`
	HasDevice   bool    `json:"has_device"`
	ScadaNodeId string  `json:"scada_node_id,omitempty"`
}

// ToTopic parses a MQTT topic into a Sparkplug topic data structure
func ToTopic(topic string) (*Topic, error) {
	var groupId, edgeNodeId, deviceId string

	if len(topic) == 0 {
		return nil, ErrTopicEmpty
	}

	if !IsValidSparkplugBTopic(topic) {
		return nil, ErrInvalidTopic
	}

	if strings.Contains(topic, string(STATE)) {
		fields := strings.Split(topic, "/")
		return &Topic{
			Command:     Command(fields[1]),
			ScadaNodeId: fields[2],
			HasDevice:   false,
		}, nil
	}

	fields := strings.Split(topic, "/")
	if len(fields) < 4 || len(fields) > 5 {
		return nil, ErrInvalidTopic
	}

	hasDevice := len(fields) == 5

	groupId = fields[0]
	edgeNodeId = fields[3]
	if len(groupId) == 0 || len(edgeNodeId) == 0 {
		return nil, ErrInvalidTopic
	}

	if hasDevice {
		deviceId = fields[4]
		if len(deviceId) == 0 {
			return nil, ErrInvalidTopic
		}

		// ex: spBv1.0/groupId/DDATA/edgeNodeId/deviceId
		return &Topic{
			GroupId:    fields[1],
			Command:    Command(fields[2]),
			EdgeNodeId: fields[3],
			DeviceId:   fields[4],
			HasDevice:  true,
		}, nil
	}

	// ex: spBv1.0/groupId/DDATA/edgeNodeId
	return &Topic{
		GroupId:    fields[1],
		Command:    Command(fields[2]),
		EdgeNodeId: fields[3],
		HasDevice:  false,
	}, nil

}

// EdgeNodeBirthTopic returns the namespace for a Sparkplug B node birth
// namespace/group_id/NBIRTH/edge_node_id
func EdgeNodeBirthTopic(groupId, edgeNodeId string) string {
	return fmt.Sprintf("%s/%s/%s/%s", SPB_NS, groupId, NBIRTH, edgeNodeId)
}

// EdgeNodeDataTopic returns the namespace for Sparkplug B node data
// namespace/group_id/NDATA/edge_node_id
func EdgeNodeDataTopic(groupId, edgeNodeId string) string {
	return fmt.Sprintf("%s/%s/%s/%s", SPB_NS, groupId, NDATA, edgeNodeId)
}

// EdgeNodeDeathTopic returns the namespace for a Sparkplug B node death
// namespace/group_id/NDEATH/edge_node_id
func EdgeNodeDeathTopic(groupId, edgeNodeId string) string {
	return fmt.Sprintf("%s/%s/%s/%s", SPB_NS, groupId, NDEATH, edgeNodeId)
}

// DeviceBirthTopic returns the namespace for a Sparkplug B device birth
// namespace/group_id/DBIRTH/edge_node_id/device_id
func DeviceBirthTopic(groupId, edgeNodeId, deviceId string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", SPB_NS, groupId, DBIRTH, edgeNodeId, deviceId)
}

// DeviceDataTopic returns the namespace for Sparkplug B device data
// namespace/group_id/DDATA/edge_node_id/device_id
func DeviceDataTopic(groupId, edgeNodeId, deviceId string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", SPB_NS, groupId, DDATA, edgeNodeId, deviceId)
}

// DeviceDeathTopic returns the namespace for a Sparkplug B device death
// namespace/group_id/DDEATH/edge_node_id/device_id
func DeviceDeathTopic(groupId, edgeNodeId, deviceId string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", SPB_NS, groupId, DDEATH, edgeNodeId, deviceId)
}

// StateCommandTopic returns the namespace for a Sparkplug B state command
// The format of the scada_host_id can be valid String with the exception of reserved characters.
// STATE/scada_host_id
func StateCommandTopic(scadaNodeId string) string {
	return fmt.Sprintf("%s/%s", STATE, scadaNodeId)
}
