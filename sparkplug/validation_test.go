package sparkplug

import "testing"

func TestIsValidSparkplugBTopic(t *testing.T) {
	tests := []struct {
		name  string
		topic string
		want  bool
	}{
		{
			name:  "valid STATE topic",
			topic: "spBv1.0/STATE/Ignition",
			want:  true,
		},
		{
			name:  "valid DBIRTH topic",
			topic: "spBv1.0/group_id/DBIRTH/edge_node_id/device_id",
			want:  true,
		},
		{
			name:  "valid DDATA topic",
			topic: "spBv1.0/group_id/DDATA/edge_node_id/device_id",
			want:  true,
		},
		{
			name:  "valid DDEATH topic",
			topic: "spBv1.0/group_id/DDEATH/edge_node_id/device_id",
			want:  true,
		},
		{
			name:  "valid NBIRTH topic",
			topic: "spBv1.0/group_id/NBIRTH/edge_node_id",
			want:  true,
		},
		{
			name:  "valid NDATA topic",
			topic: "spBv1.0/group_id/NDATA/edge_node_id",
			want:  true,
		},
		{
			name:  "valid NCMD topic",
			topic: "spBv1.0/group_id/NCMD/edge_node_id",
			want:  true,
		},
		{
			name:  "valid DCMD topic",
			topic: "spBv1.0/group_id/DCMD/edge_node_id/device_id",
			want:  true,
		},
		{
			name:  "invalid DCMD topic",
			topic: "spBv1.0/group_id/DCMD/edge_node_id",
			want:  false,
		},
		{
			name:  "valid device NDEATH topic",
			topic: "spBv1.0/group_id/NDEATH/edge_node_id",
			want:  true,
		},
		{
			name:  "invalid topic",
			topic: "spBv1.0/DDATA/edge_node_id/INVALID",
			want:  false,
		},
		{
			name:  "invalid topic",
			topic: "spBv1.1/group_id/DDATA/edge_node_id/device_id",
			want:  false,
		},
		{
			name:  "invalid topic",
			topic: "spBv1.0/group_id/FAKE/edge_node_id/device_id/extra",
			want:  false,
		},
		{
			name:  "empty topic",
			topic: "",
			want:  false,
		},
		{
			name:  "short topic",
			topic: "spBv1.0",
			want:  false,
		},
		{
			name:  "malformed topic",
			topic: "spBv1.1/group_id+bad/DDATA/edge_node_id/device_id",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidSparkplugBTopic(tt.topic); got != tt.want {
				t.Errorf("IsValidSparkplugBTopic() = %v, want %v. Topic: \"%s\"", got, tt.want, tt.topic)
			}
		})
	}
}
