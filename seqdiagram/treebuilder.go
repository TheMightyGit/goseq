// Adds the parse tree to the model

package seqdiagram

import (
    "fmt"

    "bitbucket.org/lmika/goseq/seqdiagram/parse"
)

var arrowStemMap = map[parse.ArrowStemType]ArrowStem {
    parse.SOLID_ARROW_STEM: SolidArrowStem,
    parse.DASHED_ARROW_STEM: DashedArrowStem,
    parse.THICK_ARROW_STEM: ThickArrowStem,
}

var arrowHeadMap = map[parse.ArrowHeadType]ArrowHead {
    parse.SOLID_ARROW_HEAD: SolidArrowHead,
    parse.OPEN_ARROW_HEAD: OpenArrowHead,
    parse.BARBED_ARROW_HEAD: BarbArrowHead,
    parse.LOWER_BARBED_ARROW_HEAD: LowerBarbArrowHead,
}

var noteAlignmentMap = map[parse.NoteAlignment]NoteAlignment {
    parse.LEFT_NOTE_ALIGNMENT: LeftNoteAlignment,
    parse.RIGHT_NOTE_ALIGNMENT: RightNoteAlignment,
    parse.OVER_NOTE_ALIGNMENT: OverNoteAlignment,
}

var dividerTypeMap = map[parse.GapType]DividerType {
    parse.SPACER_GAP: DTSpacer,
    parse.EMPTY_GAP: DTGap,
    parse.FRAME_GAP: DTFrame,
    parse.LINE_GAP: DTLine,
}

type treeBuilder struct {
    nodeList        *parse.NodeList
    filename        string
}

func (tb *treeBuilder) buildTree(d *Diagram) error {
    for nodeList := tb.nodeList; nodeList != nil; nodeList = nodeList.Tail {
        seqItem, err := tb.toSequenceItem(nodeList.Head, d)
        if err != nil {
            return err
        } else if seqItem != nil {
            d.AddSequenceItem(seqItem)
        }
    }

    return nil
}

func (tb *treeBuilder) nodesToSlice(nodeList *parse.NodeList, d *Diagram) ([]SequenceItem, error) {
    seq := make([]SequenceItem, 0)

    for ; nodeList != nil; nodeList = nodeList.Tail {
        seqItem, err := tb.toSequenceItem(nodeList.Head, d)
        if err != nil {
            return nil, err
        } else if seqItem != nil {
            seq = append(seq, seqItem)
        }
    }

    return seq, nil
}

func (tb *treeBuilder) makeError(msg string) error {
    return fmt.Errorf("%s:%s", tb.filename, msg)
}

func (tb *treeBuilder) toSequenceItem(node parse.Node, d *Diagram) (SequenceItem, error) {
    switch n := node.(type) {
    case *parse.ProcessInstructionNode:
        d.ProcessingInstructions = append(d.ProcessingInstructions, &ProcessingInstruction{
            Prefix: n.Prefix,
            Value: n.Value,
        })
        return nil, nil
    case *parse.TitleNode:
        d.Title = n.Title
        return nil, nil
    case *parse.ActorNode:
        d.GetOrAddActorWithOptions(n.Ident, n.ActorName())
        return nil, nil
    case *parse.ActionNode:
        return tb.addAction(n, d)
    case *parse.NoteNode:
        return tb.addNote(n, d)
    case *parse.GapNode:
        return tb.addGap(n, d)
    case *parse.BlockNode:
        return tb.addBlock(n, d)
    default:
        return nil, tb.makeError("Unrecognised declaration")
    }
}

func (tb *treeBuilder) addAction(an *parse.ActionNode, d *Diagram) (SequenceItem, error) {
    arrow := Arrow{arrowStemMap[an.Arrow.Stem], arrowHeadMap[an.Arrow.Head]}
    action := &Action{d.GetOrAddActor(an.From), d.GetOrAddActor(an.To), arrow, an.Descr}
    return action, nil
}

func (tb *treeBuilder) addNote(nn *parse.NoteNode, d *Diagram) (SequenceItem, error) {
    note := &Note{d.GetOrAddActor(nn.Actor), noteAlignmentMap[nn.Position], nn.Descr}
    return note, nil
}

func (tb *treeBuilder) addGap(gn *parse.GapNode, d *Diagram) (SequenceItem, error) {
    divider := &Divider{gn.Descr, dividerTypeMap[gn.Type]}
    return divider, nil
}

func (tb *treeBuilder) addBlock(bn *parse.BlockNode, d *Diagram) (SequenceItem, error) {
    slice, err := tb.nodesToSlice(bn.SubNodes, d)
    if err != nil {
        return nil, err
    }

    return &Block{bn.Message, slice}, nil
}