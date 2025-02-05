package refdir

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/devnev/refdir/analysis/refdir/color"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "refdir",
	Doc:      "Report potential reference-to-decleration ordering issues",
	Run:      run,
	Flags:    flag.FlagSet{},
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var (
	verbose  bool
	colorize bool
)

type RefKind string

const (
	Func     RefKind = "func"
	Type     RefKind = "type"
	RecvType RefKind = "recvtype"
	Var      RefKind = "var"
	Const    RefKind = "const"
)

type Direction string

const (
	Down   Direction = "down"
	Up     Direction = "up"
	Ignore Direction = "ignore"
)

var RefOrder = map[RefKind]Direction{
	Func:     Down,
	Type:     Down,
	RecvType: Up,
	Var:      Down,
	Const:    Down,
}

func init() {
	Analyzer.Flags.BoolVar(&verbose, "verbose", false, `print all details`)
	Analyzer.Flags.BoolVar(&colorize, "color", true, `colorize terminal`)
	addDirectionFlag := func(kind RefKind, desc string) {
		Analyzer.Flags.Func(string(kind)+"-dir", fmt.Sprintf("%s (default %s)", desc, RefOrder[kind]), func(s string) error {
			switch dir := Direction(s); dir {
			case Down, Up, Ignore:
				RefOrder[kind] = dir
				return nil
			default:
				return fmt.Errorf("must be %s, %s, or %s", Up, Down, Ignore)
			}
		})
	}
	addDirectionFlag(Func, "direction of references to functions and methods")
	addDirectionFlag(Type, "direction of type references, excluding references to the receiver type")
	addDirectionFlag(RecvType, "direction of references to the receiver type")
	addDirectionFlag(Var, "direction of references to var declarations")
	addDirectionFlag(Const, "direction of references to const declarations")
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	var printer Printer = SimplePrinter{Pass: pass}
	if colorize {
		printer = ColorPrinter{
			Pass:       pass,
			ColorError: color.Red,
			ColorInfo:  color.Gray,
			ColorOk:    color.Green,
		}
	}
	printer = VerbosePrinter{Verbose: verbose, Printer: printer}
	printer = &SortedPrinter{Pass: pass, Printer: printer}
	defer printer.Flush()

	check := func(ref *ast.Ident, def token.Pos, kind RefKind) {
		if !def.IsValid() {
			// So far only seen on calls to Error method of error interface
			printer.Info(ref.Pos(), fmt.Sprintf("got invalid definition position for %q", ref.Name))
			return
		}

		if RefOrder[kind] == Ignore {
			printer.Info(ref.Pos(), fmt.Sprintf("%s reference %s ignored by options", kind, ref.Name))
		}

		if pass.Fset.File(ref.Pos()).Name() != pass.Fset.File(def).Name() {
			printer.Info(ref.Pos(), fmt.Sprintf(`%s reference %s is to definition in separate file (%s)`, kind, ref.Name, pass.Fset.Position(def)))
			return
		}

		refLine, defLine := pass.Fset.Position(ref.Pos()).Line, pass.Fset.Position(def).Line
		if refLine == defLine {
			printer.Ok(ref.Pos(), fmt.Sprintf(`%s reference %s is on same line as definition (%s)`, kind, ref.Name, pass.Fset.Position(def)))
			return
		}

		refBeforeDef := refLine < defLine
		order := "before"
		if !refBeforeDef {
			order = "after"
		}
		var message string
		if verbose {
			message = fmt.Sprintf(`%s reference %s is %s definition (%s)`, kind, ref.Name, order, pass.Fset.Position(def))
		} else {
			message = fmt.Sprintf(`%s reference %s is %s definition`, kind, ref.Name, order)
		}

		if orderOk := refBeforeDef == (RefOrder[kind] == Down); orderOk {
			printer.Ok(ref.Pos(), message)
		} else {
			printer.Error(ref.Pos(), message)
		}
	}

	// State for keeping track of the receiver type.
	// No need for a stack as method declarations can only be at file scope.
	var (
		funcDecl       *ast.FuncDecl
		recvType       *types.TypeName
		beforeFuncType bool
	)

	inspect.Nodes(nil, func(n ast.Node, push bool) (proceed bool) {
		if !push {
			if funcDecl == n {
				funcDecl = nil
			}
			return true
		}

		switch node := n.(type) {
		case *ast.File:
			if ast.IsGenerated(node) {
				printer.Info(node.Pos(), "skipping generated file")
				return false
			}

		case *ast.FuncDecl:
			if funcDecl == nil {
				funcDecl = node
				beforeFuncType = true
			}

		case *ast.FuncType:
			beforeFuncType = false

		case *ast.SelectorExpr:
			sel := pass.TypesInfo.Selections[node]
			if sel == nil {
				// Based on TypesInfo.Selection docs this should only be the
				// case for "qualified identifiers", which I think means
				// references to out-of-package identifiers, which we don't care
				// about anyway. Logging just in case.
				printer.Info(node.Pos(), fmt.Sprintf("skipping selector %s with missing Selections", node.Sel.String()))
				break
			}

			obj := sel.Obj()
			switch sel.Kind() {
			case types.MethodVal:
				check(node.Sel, obj.Pos(), Func)
			case types.FieldVal:
			case types.MethodExpr:
				check(node.Sel, obj.Pos(), Func)
			default:
				// No other enum values are defined, logging just in case.
				printer.Info(node.Pos(), fmt.Sprintf("unknown selection kind %v", sel.Kind()))
			}

		case *ast.Ident:
			switch def := pass.TypesInfo.Uses[node].(type) {
			case *types.Var:
				def = def.Origin()
				if def.IsField() {
					printer.Info(node.Pos(), fmt.Sprintf("skipping var ident %s for field %s", node.Name, pass.Fset.Position(def.Pos())))
				} else if def.Parent() != def.Pkg().Scope() {
					printer.Info(node.Pos(), fmt.Sprintf("skipping var ident %s with inner parent scope %s", node.Name, pass.Fset.Position(def.Parent().Pos())))
				} else {
					check(node, def.Pos(), Var)
				}
			case *types.Const:
				if def.Parent() != def.Pkg().Scope() {
					printer.Info(node.Pos(), fmt.Sprintf("skipping var ident %s with inner parent scope %s", node.Name, pass.Fset.Position(def.Parent().Pos())))
				} else {
					check(node, def.Pos(), Const)
				}
			case *types.Func:
				def = def.Origin()
				if def.Parent() != nil && def.Parent() != def.Pkg().Scope() {
					printer.Info(node.Pos(), fmt.Sprintf("skipping func ident %s with inner parent scope %s", node.Name, pass.Fset.Position(def.Parent().Pos())))
				} else {
					check(node, def.Pos(), Func)
				}
			case *types.TypeName:
				if funcDecl != nil && beforeFuncType {
					// We're in a file-level func decl before getting to the
					// function type, so this must be an identifier in the type
					// of the receiver.
					recvType = def
					printer.Info(node.Pos(), fmt.Sprintf("skipping ident %s in recv list", node.Name))
					break
				}
				if funcDecl != nil && recvType == def {
					// Reference to the receiver type within a method type or body
					check(node, def.Pos(), RecvType)
					break
				}
				check(node, def.Pos(), Type)
			default:
				printer.Info(node.Pos(), fmt.Sprintf("unexpected ident def type %T for %q", pass.TypesInfo.Uses[node], node.Name))
			}
		}

		return true
	})

	return nil, nil
}
