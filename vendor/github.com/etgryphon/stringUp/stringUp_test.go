package stringUp

import (
    "testing"
)

func TestCamelCased(t *testing.T) {
  const cameled,upCameled = "thisIsIt", "ThisIsItBob"
  if x := CamelCase(cameled); x != cameled {
    t.Errorf("CamelCase(%v) = %v, want %v", cameled, x, cameled)
  }
  if x := CamelCase(upCameled); x != upCameled {
    t.Errorf("CamelCase(%v) = %v, want %v", upCameled, x, upCameled)
  }
}

func TestCamelCaseSpaced(t *testing.T) {
  const src,upSrc = "this is it", "This Is It Bob"
  if x := CamelCase(src); x != "thisIsIt" {
    t.Errorf("CamelCase(%v) = %v, want %v", src, x, "thisIsIt")
  }
  if x := CamelCase(upSrc); x != "ThisIsItBob" {
    t.Errorf("CamelCase(%v) = %v, want %v", upSrc, x, "ThisIsItBob")
  }
}

func TestCamelCaseUnderscored(t *testing.T) {
  const src,upSrc = "this_is_it", "This_Is_It_Bob"
  if x := CamelCase(src); x != "thisIsIt" {
    t.Errorf("CamelCase(%v) = %v, want %v", src, x, "thisIsIt")
  }
  if x := CamelCase(upSrc); x != "ThisIsItBob" {
    t.Errorf("CamelCase(%v) = %v, want %v", upSrc, x, "ThisIsItBob")
  }
}

func TestCamelCaseDashed(t *testing.T) {
  const src,upSrc = "this-is-it", "This-Is-It-Bob"
  if x := CamelCase(src); x != "thisIsIt" {
    t.Errorf("CamelCase(%v) = %v, want %v", src, x, "thisIsIt")
  }
  if x := CamelCase(upSrc); x != "ThisIsItBob" {
    t.Errorf("CamelCase(%v) = %v, want %v", upSrc, x, "ThisIsItBob")
  }
}

func TestCamelCaseMixed(t *testing.T) {
  const src,upSrc = "-this is_it", "This Is_It-Bob"
  if x := CamelCase(src); x != "thisIsIt" {
    t.Errorf("CamelCase(%v) = %v, want %v", src, x, "thisIsIt")
  }
  if x := CamelCase(upSrc); x != "ThisIsItBob" {
    t.Errorf("CamelCase(%v) = %v, want %v", upSrc, x, "ThisIsItBob")
  }
}