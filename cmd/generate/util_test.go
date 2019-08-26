package generate

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestInsertFileContent(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.WriteString("test1 \n test2 \n \n test3 \n test3 \n test4")

	name := "insert_file_content.txt"
	err := writeFile(name, buf)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(name)

	err = insertFileContent(name, "test2", "test3", "test5 \n")
	if err != nil {
		t.Fatal(err)
	}

	nbuf, err := readFile(name)
	if err != nil {
		t.Fatal(err)
	}

	ss := strings.Split(nbuf.String(), "\n")
	if len(ss) != 7 {
		t.Errorf("不符合预期：%v", ss)
		return
	}

	expect := []string{"test1", "test2", "", "test3", "test3", "test5", "test4"}
	for i := 0; i < len(ss); i++ {
		if expect[i] != strings.TrimSpace(ss[i]) {
			t.Errorf("不符合预期：%v", ss)
			return
		}
	}

}
