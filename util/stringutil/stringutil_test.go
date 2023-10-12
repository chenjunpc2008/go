package stringutil

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_CutByMaxLen(t *testing.T) {
    //
    src := "1234567890ABC"
    dest := CutByMaxLen(src, 5)
    expected := "12345"

    assert.Equal(t, expected, dest)

    //
    src = "123456"
    dest = CutByMaxLen(src, 6)
    expected = "123456"

    assert.Equal(t, expected, dest)

    //
    src = "12345"
    dest = CutByMaxLen(src, 10)
    expected = "12345"

    assert.Equal(t, expected, dest)
}

func Test_CutLastByMaxLen(t *testing.T) {
    //
    src := "1234567890ABC"
    dest := CutLastByMaxLen(src, 5)
    expected := "90ABC"

    assert.Equal(t, expected, dest)

    //
    src = "123456"
    dest = CutLastByMaxLen(src, 6)
    expected = "123456"

    assert.Equal(t, expected, dest)

    //
    src = "123456"
    dest = CutLastByMaxLen(src, 7)
    expected = "123456"

    assert.Equal(t, expected, dest)

    //
    src = "12345"
    dest = CutLastByMaxLen(src, 10)
    expected = "12345"

    assert.Equal(t, expected, dest)
}
