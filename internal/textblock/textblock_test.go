package textblock

import "testing"

func TestUpsertIdempotent(t *testing.T) {
	t.Parallel()
	begin := "# >>> a >>>"
	end := "# <<< a <<<"
	block := begin + "\nX\n" + end

	got, changed, err := Upsert("hello\n", begin, end, block)
	if err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}
	if !changed {
		t.Fatal("expected changed=true")
	}

	got2, changed2, err := Upsert(got, begin, end, block)
	if err != nil {
		t.Fatalf("2nd Upsert failed: %v", err)
	}
	if changed2 {
		t.Fatal("expected idempotent upsert")
	}
	if got2 != got {
		t.Fatal("expected identical output on 2nd upsert")
	}
}

func TestRemoveMalformed(t *testing.T) {
	t.Parallel()
	begin := "# >>> a >>>"
	end := "# <<< a <<<"

	if _, _, err := Remove("x\n"+end+"\n", begin, end); err == nil {
		t.Fatal("expected malformed markers error")
	}
}

func TestFindSingleRejectsDuplicateBegin(t *testing.T) {
	t.Parallel()
	begin := "# >>> a >>>"
	end := "# <<< a <<<"
	content := begin + "\n1\n" + begin + "\n2\n" + end
	if _, _, err := FindSingle(content, begin, end); err == nil {
		t.Fatal("expected malformed markers for duplicate begin")
	}
}
