package handlers

import "testing"

func TestReplaceBadWords(t *testing.T) {
	cases := []struct {
		name  string
		chirp string
		want  string
	}{
		{
			name:  "Test kerfuffle",
			chirp: "This is a kerfuffle opinion I need to share with the world",
			want:  "This is a **** opinion I need to share with the world",
		},
		{
			name:  "Test sharbert",
			chirp: "Get this sharbert thing outta my face",
			want:  "Get this **** thing outta my face",
		},
		{
			name:  "Test fornax",
			chirp: "That's some really interesting fornax you got there",
			want:  "That's some really interesting **** you got there",
		},
		{
			name:  "Test sharbert with exclamation point",
			chirp: "Sharbert! My bad man",
			want:  "Sharbert! My bad man",
		},
		{
			name:  "Test capital fornax",
			chirp: "Fornax is some prime reading material",
			want:  "**** is some prime reading material",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := replaceBadWords(c.chirp)
			if got != c.want {
				t.Errorf("Expected '%v'\ngot '%v'", c.want, got)
			}
		})
	}
}
