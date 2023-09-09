package security_test

import (
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/tools/security"
)

func TestS256Challenge(t *testing.T) {
	scenarios := []struct {
		code     string
		expected string
	}{
		{"", "47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{"123", "pmWkWSBCL51Bfkhn79xPuKBKHz__H6B-mY6G9_eieuM"},
	}

	for _, s := range scenarios {
		t.Run(s.code, func(t *testing.T) {
			result := security.S256Challenge(s.code)

			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
	}
}

func TestMD5(t *testing.T) {
	scenarios := []struct {
		code     string
		expected string
	}{
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"123", "202cb962ac59075b964b07152d234b70"},
	}

	for _, s := range scenarios {
		t.Run(s.code, func(t *testing.T) {
			result := security.MD5(s.code)

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestSHA256(t *testing.T) {
	scenarios := []struct {
		code     string
		expected string
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"123", "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3"},
	}

	for _, s := range scenarios {
		t.Run(s.code, func(t *testing.T) {
			result := security.SHA256(s.code)

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestSHA512(t *testing.T) {
	scenarios := []struct {
		code     string
		expected string
	}{
		{"", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
		{"123", "3c9909afec25354d551dae21590bb26e38d53f2173b8d3dc3eee4c047e7ab1c1eb8b85103e3be7ba613b31bb5c9c36214dc9f14a42fd7a2fdb84856bca5c44c2"},
	}

	for _, s := range scenarios {
		t.Run(s.code, func(t *testing.T) {
			result := security.SHA512(s.code)

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestHS256(t *testing.T) {
	scenarios := []struct {
		text     string
		secret   string
		expected string
	}{
		{" ", "test", "9fb4e4a12d50728683a222b4fc466a69ee977332cfcdd6b9ebb44c7121dbd99f"},
		{" ", "test2", "d792417a504716e22805d940125ec12e68e8cb18fc84674703bd96c59f1e1228"},
		{"hello", "test", "f151ea24bda91a18e89b8bb5793ef324b2a02133cce15a28a719acbd2e58a986"},
		{"hello", "test2", "16436e8dcbf3d7b5b0455573b27e6372699beb5bfe94e6a2a371b14b4ae068f4"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d-%s", i, s.text), func(t *testing.T) {
			result := security.HS256(s.text, s.secret)

			if result != s.expected {
				t.Fatalf("Expected \n%v, \ngot \n%v", s.expected, result)
			}
		})
	}
}

func TestHS512(t *testing.T) {
	scenarios := []struct {
		text     string
		secret   string
		expected string
	}{
		{" ", "test", "eb3bdb0352c95c38880c1f645fc7e1d1332644f938f50de0d73876e42d6f302e599bb526531ba79940e8b314369aaef3675322d8d851f9fc6ea9ed121286d196"},
		{" ", "test2", "8b69e84e9252af78ae8b1c4bed3c9f737f69a3df33064cfbefe76b36d19d1827285e543cdf066cdc8bd556cc0cd0e212d52e9c12a50cd16046181ff127f4cf7f"},
		{"hello", "test", "44f280e11103e295c26cd61dd1cdd8178b531b860466867c13b1c37a26b6389f8af110efbe0bb0717b9d9c87f6fe1c97b3b1690936578890e5669abf279fe7fd"},
		{"hello", "test2", "d7f10b1b66941b20817689b973ca9dfc971090e28cfb8becbddd6824569b323eca6a0cdf2c387aa41e15040007dca5a011dd4e4bb61cfd5011aa7354d866f6ef"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d-%q", i, s.text), func(t *testing.T) {
			result := security.HS512(s.text, s.secret)

			if result != s.expected {
				t.Fatalf("Expected \n%v, \ngot \n%v", s.expected, result)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	scenarios := []struct {
		hash1    string
		hash2    string
		expected bool
	}{
		{"", "", true},
		{"abc", "abd", false},
		{"abc", "abc", true},
	}

	for _, s := range scenarios {
		t.Run(fmt.Sprintf("%qVS%q", s.hash1, s.hash2), func(t *testing.T) {
			result := security.Equal(s.hash1, s.hash2)

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}
