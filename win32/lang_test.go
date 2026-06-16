package win32

import "testing"

func TestLang(t *testing.T) {
	t.Run("MAKELANGID", func(t *testing.T) {
		id := MAKELANGID(
			LANG_ENGLISH,
			SUBLANG_ENGLISH_US,
		)

		if id != 0x0409 {
			t.Fatalf("got %#x want %#x", id, LANGID(0x0409))
		}
	})

	t.Run("PRIMARYLANGID", func(t *testing.T) {
		id := LANGID(0x0804)

		if got := PRIMARYLANGID(id); got != LANG_CHINESE {
			t.Fatalf("got %#x want %#x", got, LANG_CHINESE)
		}
	})

	t.Run("SUBLANGID", func(t *testing.T) {
		id := LANGID(0x0804)

		if got := SUBLANGID(id); got != SUBLANG_CHINESE_SIMPLIFIED {
			t.Fatalf("got %#x want %#x", got, SUBLANG_CHINESE_SIMPLIFIED)
		}
	})

	t.Run("RoundTripLANGID", func(t *testing.T) {
		cases := []struct {
			primary LANGID_PRIMARY
			sub     LANGID_SUB
		}{
			{LANG_ENGLISH, SUBLANG_ENGLISH_US},
			{LANG_ENGLISH, SUBLANG_ENGLISH_UK},
			{LANG_CHINESE, SUBLANG_CHINESE_SIMPLIFIED},
			{LANG_CHINESE, SUBLANG_CHINESE_TRADITIONAL},
			{LANG_JAPANESE, SUBLANG_JAPANESE_JAPAN},
		}

		for _, c := range cases {
			id := MAKELANGID(c.primary, c.sub)

			if got := PRIMARYLANGID(id); got != c.primary {
				t.Fatalf("PRIMARYLANGID(%#x)=%#x want %#x",
					id, got, c.primary)
			}

			if got := SUBLANGID(id); got != c.sub {
				t.Fatalf("SUBLANGID(%#x)=%#x want %#x",
					id, got, c.sub)
			}
		}
	})

	t.Run("MAKELCID", func(t *testing.T) {
		lang := LANGID(0x0409)

		lcid := MAKELCID(lang, SORT_DEFAULT)

		if lcid != 0x00000409 {
			t.Fatalf("got %#x want %#x", lcid, LCID(0x00000409))
		}
	})

	t.Run("MAKESORTLCID", func(t *testing.T) {
		lcid := MAKESORTLCID(
			LANGID(0x0409),
			SORT_GERMAN_PHONE_BOOK,
			2,
		)

		if LANGIDFROMLCID(lcid) != 0x0409 {
			t.Fatal("LANGIDFROMLCID failed")
		}

		if SORTIDFROMLCID(lcid) != SORT_GERMAN_PHONE_BOOK {
			t.Fatal("SORTIDFROMLCID failed")
		}

		if SORTVERSIONFROMLCID(lcid) != 2 {
			t.Fatal("SORTVERSIONFROMLCID failed")
		}
	})

	t.Run("RoundTripLCID", func(t *testing.T) {
		for _, sortID := range []SORTID{
			SORT_DEFAULT,
			SORT_GERMAN_PHONE_BOOK,
			SORT_CHINESE_PRC,
			SORT_CHINESE_BOPOMOFO,
		} {
			lcid := MAKESORTLCID(
				MAKELANGID(
					LANG_ENGLISH,
					SUBLANG_ENGLISH_US,
				),
				sortID,
				7,
			)

			if got := LANGIDFROMLCID(lcid); got != 0x0409 {
				t.Fatalf("LANGIDFROMLCID got %#x", got)
			}

			if got := SORTIDFROMLCID(lcid); got != sortID {
				t.Fatalf("SORTIDFROMLCID got %#x want %#x",
					got, sortID)
			}

			if got := SORTVERSIONFROMLCID(lcid); got != 7 {
				t.Fatalf("SORTVERSIONFROMLCID got %#x want 7",
					got)
			}
		}
	})
}

func TestLangIDKnownValues(t *testing.T) {
	if got := MAKELANGID(
		LANG_ENGLISH,
		SUBLANG_ENGLISH_US,
	); got != 0x0409 {
		t.Fatal()
	}

	if got := MAKELANGID(
		LANG_CHINESE,
		SUBLANG_CHINESE_SIMPLIFIED,
	); got != 0x0804 {
		t.Fatal()
	}

	if got := MAKELANGID(
		LANG_CHINESE,
		SUBLANG_CHINESE_TRADITIONAL,
	); got != 0x0404 {
		t.Fatal()
	}
}
