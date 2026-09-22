package crypto 

import ( 
	"encoding/hex" 
	"testing" 
) 

func TestKeccak256(t *testing.T) { 
	tests := []struct { 
		name string 
		input string 
		expected string 
}{ 
		{ 
				name: "Empty String Hash", 
				input: "", 
				expected: "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470", 
		},
		{ 		name: "Hello String Hash", 
				input: "hello", 
				expected: "1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8", 
		}, 
} 

for _, tt := range tests { 
		t.Run(tt.name, func(t *testing.T) { 
			hash := Keccak256([]byte(tt.input)) 
			gotHex := hex.EncodeToString(hash) 
			if gotHex != tt.expected { 
					t.Errorf("Keccak256(%q) = %s; want %s", tt.input, gotHex, tt.expected) 
			} 
		}) 
	} 
}