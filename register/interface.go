package register

import (
    . "ml/dict"
    "github.com/PuerkitoBio/goquery"
)

type elementDict map[string]*goquery.Selection
type fillFormCallback func(form Dict, elements elementDict)

type IAppleIdRegister interface {
    Close()
    Initialize()
    Signup() int
}
