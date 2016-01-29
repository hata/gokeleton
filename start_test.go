package main

import (
    "testing"
)

func Test_stringsToMap_Nothing(t *testing.T) {
    m := stringsToMap("", ",")
    if len(m) != 0 {
        t.Error("Verify no map key-value found.")
    }
}

func Test_stringsToMap_Simple(t *testing.T) {
    m := stringsToMap("foo=bar", ",")
    if m["foo"] != "bar" {
        t.Error("Verify there is a key value pair.")
    }
}

func Test_stringsToMap_KeyOnly(t *testing.T) {
    m := stringsToMap("foo", ",")
    if m["foo"] != "" {
        t.Error("Verify no value is set")
    }
    m = stringsToMap("foo=", ",")
    if m["foo"] != "" {
        t.Error("Verify no value is set")
    }
}

func Test_stringsToMap_MultiKeywords(t *testing.T) {
    m := stringsToMap("foo=bar,bar=hoge", ",")

    if len(m) != 2 {
        t.Error("Verify there are two entries.")
    }

    if m["foo"] != "bar" {
        t.Error("Verify there is a key value pair.")
    }

    if m["bar"] != "hoge" {
        t.Error("Verify there is a key value pair.")
    }
}
