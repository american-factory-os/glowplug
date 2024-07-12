package sparkplug

import (
	"fmt"
	"testing"
)

func TestToTopic(t *testing.T) {
	tests := []struct {
		name  string
		topic string
		want  *Topic
		err   error
	}{
		{
			name:  "valid STATE topic",
			topic: "spBv1.0/STATE/Ignition",
			want: &Topic{
				Command:     STATE,
				ScadaNodeId: "Ignition",
				HasDevice:   false,
			},
			err: nil,
		},
		{
			name:  "valid device DDATA topic",
			topic: "spBv1.0/group_id/DDATA/edge_node_id/device_id",
			want: &Topic{
				Command:    DDATA,
				GroupId:    "group_id",
				EdgeNodeId: "edge_node_id",
				DeviceId:   "device_id",
				HasDevice:  true,
			},
			err: nil,
		},
		{
			name:  "valid device NDATA topic",
			topic: "spBv1.0/group_id/NDATA/edge_node_id",
			want: &Topic{
				Command:    NDATA,
				GroupId:    "group_id",
				EdgeNodeId: "edge_node_id",
				HasDevice:  false,
			},
			err: nil,
		},
		{
			name:  "invalid topic",
			topic: "spBv1.0/INVALID/edge_node_id/INVALID",
			want:  nil,
			err:   ErrInvalidTopic,
		},
		{
			name:  "invalid topic 2",
			topic: "spBv1.1/group_id/DDATA/edge_node_id/device_id",
			want:  nil,
			err:   ErrInvalidTopic,
		},
		{
			name:  "empty topic",
			topic: "",
			want:  nil,
			err:   ErrTopicEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ToTopic(tt.topic)

			if err != tt.err {
				fmt.Println("GOT", got)
				t.Errorf("ToTopic() error = %v, want %v", err, tt.err)
				return
			}

			if got != nil && tt.want != nil {
				if got.Command != tt.want.Command {
					t.Errorf("ToTopic() got Command = %v, want %v", got.Command, tt.want.Command)
				}
				if got.GroupId != tt.want.GroupId {
					t.Errorf("ToTopic() got GroupId = %v, want %v", got.GroupId, tt.want.GroupId)
				}
				if got.EdgeNodeId != tt.want.EdgeNodeId {
					t.Errorf("ToTopic() got EdgeNodeId = %v, want %v", got.EdgeNodeId, tt.want.EdgeNodeId)
				}
				if got.DeviceId != tt.want.DeviceId {
					t.Errorf("ToTopic() got DeviceId = %v, want %v", got.DeviceId, tt.want.DeviceId)
				}
				if got.HasDevice != tt.want.HasDevice {
					t.Errorf("ToTopic() got HasDevice = %v, want %v", got.HasDevice, tt.want.HasDevice)
				}
			}
		})
	}
}

func TestNodeBirthTopic(t *testing.T) {
	t.Run("Test NodeBirthTopic", func(t *testing.T) {
		got := EdgeNodeBirthTopic("group_id", "edge_node_id")
		want := "spBv1.0/group_id/NBIRTH/edge_node_id"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestNodeDataTopic(t *testing.T) {
	t.Run("Test NodeDataTopic", func(t *testing.T) {
		got := EdgeNodeDataTopic("group_id", "edge_node_id")
		want := "spBv1.0/group_id/NDATA/edge_node_id"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestNodeDeathTopic(t *testing.T) {
	t.Run("Test NodeDeathTopic", func(t *testing.T) {
		got := EdgeNodeDeathTopic("group_id", "edge_node_id")
		want := "spBv1.0/group_id/NDEATH/edge_node_id"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestDeviceBirthTopic(t *testing.T) {
	t.Run("Test DeviceBirthTopic", func(t *testing.T) {
		got := DeviceBirthTopic("group_id", "edge_node_id", "device_id")
		want := "spBv1.0/group_id/DBIRTH/edge_node_id/device_id"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestDeviceDataTopic(t *testing.T) {
	t.Run("Test DeviceDataTopic", func(t *testing.T) {
		got := DeviceDataTopic("group_id", "edge_node_id", "device_id")
		want := "spBv1.0/group_id/DDATA/edge_node_id/device_id"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestDeviceDeathTopic(t *testing.T) {
	t.Run("Test DeviceDeathTopic", func(t *testing.T) {
		got := DeviceDeathTopic("group_id", "edge_node_id", "device_id")
		want := "spBv1.0/group_id/DDEATH/edge_node_id/device_id"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestStateCommandTopic(t *testing.T) {
	t.Run("Test StateCommandTopic", func(t *testing.T) {
		got := StateCommandTopic("scada_host_id")
		want := "STATE/scada_host_id"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestEdgeNodeCommandTopic(t *testing.T) {
	t.Run("Test EdgeNodeCommandTopic", func(t *testing.T) {
		got := EdgeNodeCommandTopic("group_id", "edge_node_id")
		want := "spBv1.0/group_id/NCMD/edge_node_id"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestDeviceCommandTopic(t *testing.T) {
	t.Run("Test DeviceCommandTopic", func(t *testing.T) {
		got := DeviceCommandTopic("group_id", "edge_node_id", "device_id")
		want := "spBv1.0/group_id/DCMD/edge_node_id/device_id"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
