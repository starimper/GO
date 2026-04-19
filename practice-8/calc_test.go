package main

import "testing"


func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5
	if got != want {
		t.Errorf("Add(2, 3) = %d; want %d", got, want)
	}
}

func TestAddTableDriven(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    int
	}{
		{"both positive", 2, 3, 5},
		{"positive + zero", 5, 0, 5},
		{"negative + positive", -1, 4, 3},
		{"both negative", -2, -3, -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}


func TestSubtractTableDriven(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    int
	}{
		{"both positive", 10, 3, 7},
		{"positive minus zero", 5, 0, 5},
		{"negative minus positive", -4, 3, -7},
		{"both negative", -2, -5, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Subtract(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}


func TestDivideTableDriven(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    int
		wantErr bool
		errMsg  string
	}{
		{"both positive", 10, 2, 5, false, ""},
		{"positive divided by 1", 7, 1, 7, false, ""},
		{"negative dividend", -9, 3, -3, false, ""},
		{"both negative", -8, -4, 2, false, ""},
		{"zero numerator", 0, 5, 0, false, ""},
		{"division by zero", 10, 0, 0, true, "division by zero"},
		{"negative divided by zero", -4, 0, 0, true, "division by zero"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Divide(%d, %d): expected error, got nil", tt.a, tt.b)
				}
				if err.Error() != tt.errMsg {
					t.Errorf("error = %q; want %q", err.Error(), tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("Divide(%d, %d): unexpected error: %v", tt.a, tt.b, err)
			}
			if got != tt.want {
				t.Errorf("Divide(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
