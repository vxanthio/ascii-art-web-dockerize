package sanitize

import "testing"

func TestHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "script tag",
			input: "<script>alert('XSS')</script>",
			want:  "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;",
		},
		{
			name:  "img tag with onerror",
			input: "<img src=x onerror=alert(1)>",
			want:  "&lt;img src=x onerror=alert(1)&gt;",
		},
		{
			name:  "ampersand",
			input: "Hello & goodbye",
			want:  "Hello &amp; goodbye",
		},
		{
			name:  "less than and greater than",
			input: "5 < 10 > 3",
			want:  "5 &lt; 10 &gt; 3",
		},
		{
			name:  "double quotes",
			input: `He said "hello"`,
			want:  "He said &#34;hello&#34;",
		},
		{
			name:  "single quotes",
			input: "It's working",
			want:  "It&#39;s working",
		},
		{
			name:  "mixed dangerous chars",
			input: `<div class="test" onclick='alert("XSS")'>`,
			want:  "&lt;div class=&#34;test&#34; onclick=&#39;alert(&#34;XSS&#34;)&#39;&gt;",
		},
		{
			name:  "safe text unchanged",
			input: "Hello World 123",
			want:  "Hello World 123",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "javascript protocol",
			input: `<a href="javascript:alert('XSS')">Click</a>`,
			want:  "&lt;a href=&#34;javascript:alert(&#39;XSS&#39;)&#34;&gt;Click&lt;/a&gt;",
		},
		{
			name:  "svg with script",
			input: "<svg onload=alert(1)>",
			want:  "&lt;svg onload=alert(1)&gt;",
		},
		{
			name:  "iframe injection",
			input: "<iframe src=evil.com></iframe>",
			want:  "&lt;iframe src=evil.com&gt;&lt;/iframe&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HTML(tt.input)
			if got != tt.want {
				t.Errorf("HTML() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTMLPreservesNewlines(t *testing.T) {
	input := "Line 1\nLine 2\nLine 3"
	want := "Line 1\nLine 2\nLine 3"
	
	got := HTML(input)
	if got != want {
		t.Errorf("HTML() should preserve newlines, got %q, want %q", got, want)
	}
}

func TestHTMLMultipleEscapes(t *testing.T) {
	input := "<script>alert('test')</script>"
	
	// First escape
	escaped1 := HTML(input)
	expected1 := "&lt;script&gt;alert(&#39;test&#39;)&lt;/script&gt;"
	if escaped1 != expected1 {
		t.Errorf("First escape failed: got %q, want %q", escaped1, expected1)
	}
	
	// Second escape (should escape the & symbols)
	escaped2 := HTML(escaped1)
	if escaped2 == escaped1 {
		t.Error("Second escape should change the string")
	}
}
