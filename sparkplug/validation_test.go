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
			name:  "valid device DBIRTH topic",
			topic: "spBv1.0/group_id/DBIRTH/edge_node_id/device_id",
			want:  true,
		},
		{
			name:  "valid device DDATA topic",
			topic: "spBv1.0/group_id/DDATA/edge_node_id/device_id",
			want:  true,
		},
		{
			name:  "valid device DDEATH topic",
			topic: "spBv1.0/group_id/DDEATH/edge_node_id/device_id",
			want:  true,
		},
		{
			name:  "valid device NBIRTH topic",
			topic: "spBv1.0/group_id/NBIRTH/edge_node_id",
			want:  true,
		},
		{
			name:  "valid device NDATA topic",
			topic: "spBv1.0/group_id/NDATA/edge_node_id",
			want:  true,
		},
		{
			name:  "valid device NCMD topic",
			topic: "spBv1.0/group_id/NCMD/edge_node_id",
			want:  true,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidSparkplugBTopic(tt.topic); got != tt.want {
				t.Errorf("IsValidSparkplugBTopic() = %v, want %v. Topic: \"%s\"", got, tt.want, tt.topic)
			}
		})
	}
}
