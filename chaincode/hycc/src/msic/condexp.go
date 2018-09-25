package msic

import (
	"container/list"
	"github.com/jenner/chaincode/hycc/src/utils"
)

//=============================================================================
//表达式元素类型
type ElemType int

const (
	ET_UNK ElemType = 0 //未知
	ET_VAR ElemType = 1 //变量
	ET_OPT ElemType = 2 //操作符
	ET_BRK ElemType = 3 //括号
)

//=============================================================================
//条件表达式
type CondExp struct {
	Post   int
	StrExp string
	InExp  []string
	SufExp []string
}

//-----------------------------------------------------------------------------
//新建条件表达式
func NewCondExp(exp string) *CondExp {
	return &CondExp{StrExp: exp}
}

//转换为中缀表达式
func (ce *CondExp) toInExp() error {
	etype, elem := ET_UNK, ""
	for i := 0; i < len(ce.StrExp); i++ {
		ch := ce.StrExp[i]
		if ch == '(' || ch == ')' {
			if etype != ET_UNK {
				ce.InExp = append(ce.InExp, elem)
				elem = ""
			}
			etype = ET_BRK
			elem = string(ch)
		} else if ch == '|' || ch == '&' {
			if etype != ET_UNK {
				ce.InExp = append(ce.InExp, elem)
				elem = ""
			}
			if ch != ce.StrExp[i+1] {
				return utils.Errorf("Wrong expression unkonwn elem: %s", string(ch))
			} else {
				etype = ET_OPT
				elem = string(ch) + string(ch)
				i += 1
			}
		} else if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			if etype != ET_UNK && etype != ET_VAR {
				ce.InExp = append(ce.InExp, elem)
				elem = ""
			}
			etype = ET_VAR
			elem += string(ch)
		} else if ch == ' ' {
			continue
		} else {
			return utils.Errorf("Wrong expression unkonwn elem: %s", string(ch))
		}
	}
	if elem != "" {
		ce.InExp = append(ce.InExp, elem)
	}
	return nil
}

//转换为后缀表达式
func (ce *CondExp) toSufExp() error {
	stack := list.New()
	for _, elem := range ce.InExp {
		if elem == "(" || elem == "||" || elem == "&&" {
			stack.PushBack(elem)
		} else if elem == ")" {
			e := stack.Back()
			stack.Remove(e)
			for e.Value != "(" {
				ce.SufExp = append(ce.SufExp, e.Value.(string))
				e = stack.Back()
				stack.Remove(e)
			}
		} else {
			ce.SufExp = append(ce.SufExp, elem)
		}
	}
	for stack.Len() > 0 {
		e := stack.Back()
		stack.Remove(e)
		ce.SufExp = append(ce.SufExp, e.Value.(string))
	}
	return nil
}

//计算表达式结果
func (ce *CondExp) Calc(args []string) (rst bool, err error) {
	//参数转换
	if err = ce.toInExp(); err != nil {
		return false, utils.Error(err)
	}
	if err = ce.toSufExp(); err != nil {
		return false, utils.Error(err)
	}
	//替换参数
	exp := make([]string, len(ce.SufExp))
	copy(exp, ce.SufExp)
	for i := 0; i < len(exp); i++ {
		if exp[i] == "(" || exp[i] == ")" ||
			exp[i] == "||" || exp[i] == "&&" {
			continue
		}
		for _, arg := range args {
			if arg == exp[i] {
				exp[i] = "1"
				break
			}
		}
		if exp[i] != "1" {
			exp[i] = "0"
		}
	}
	//计算结果
	stack := list.New()
	for _, e := range exp {
		if e == "||" || e == "&&" {
			if stack.Len() < 2 {
				return false, utils.Errorf("Wrong expression")
			}
			e1 := stack.Back()
			stack.Remove(e1)
			e2 := stack.Back()
			stack.Remove(e2)
			if (e1.Value != "1" && e1.Value != "0") || (e2.Value != "1" && e2.Value != "0") {
				return false, utils.Errorf("Wrong expression")
			}
			if e == "||" {
				if e1.Value == "1" || e2.Value == "1" {
					stack.PushBack("1")
				} else {
					stack.PushBack("0")
				}
			} else if e == "&&" {
				if e1.Value == "1" && e2.Value == "1" {
					stack.PushBack("1")
				} else {
					stack.PushBack("0")
				}
			}
		} else {
			stack.PushBack(e)
		}
	}
	if stack.Len() != 1 {
		return false, utils.Errorf("Wrong expression")
	}
	return stack.Back().Value == "1", nil
}

