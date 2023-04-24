package models

import (
	"fmt"
	"github.com/televi-go/televi/models/render"
	"github.com/televi-go/televi/models/render/results"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/bot"
	"github.com/televi-go/televi/telegram/dto"
	"github.com/televi-go/televi/util"
	"sync"
	"time"
)

type ResultLine struct {
	Line []*render.CompletedResult
	// stored controller-wide
	ActiveKbResult *results.ReplyKeyboardResult
}

type analysisNode struct {
	CompletedResult *render.CompletedResult
	MarkedDelete    bool
	MarkedInsert    bool
}

func (resultLine *ResultLine) getAnalysisNodes() (result []*analysisNode) {
	for _, completedResult := range resultLine.Line {
		result = append(result, &analysisNode{
			CompletedResult: completedResult,
			MarkedDelete:    false,
		})
	}
	return
}

type analysisSections struct {
	new      []render.IResult
	orphaned []*analysisNode
	common   []commonNode
}

type commonNode struct {
	NewResult render.IResult
	*analysisNode
}

func buildSections(nodes []*analysisNode, line []render.IResult) analysisSections {
	commonLength := util.Min(len(nodes), len(line))
	sections := analysisSections{
		common: make([]commonNode, commonLength),
	}

	for i := 0; i < commonLength; i++ {
		sections.common[i] = commonNode{
			NewResult:    line[i],
			analysisNode: nodes[i],
		}
	}

	if len(nodes) > commonLength {
		sections.orphaned = nodes[commonLength:]
	}

	if len(line) > commonLength {
		sections.new = line[commonLength:]
	}

	return sections
}

type CompareResult struct {
	Consecutive []BoundRequest
	Parallel    []BoundRequest
}

type BoundRequest struct {
	Request telegram.Request
	Slot    *analysisNode
}

func bind(requests []telegram.Request, slot *analysisNode) []BoundRequest {
	var result []BoundRequest
	for _, request := range requests {
		result = append(result, BoundRequest{
			Request: request,
			Slot:    slot,
		})
	}
	return result
}

func compare(destination telegram.Destination, replaceMode bool, node commonNode) (bool, CompareResult) {
	cr := CompareResult{}
	if !replaceMode {
		canBeChanged, changes := node.CompletedResult.CompareTo(node.NewResult, destination)
		if canBeChanged {
			cr.Parallel = bind(changes, node.analysisNode)
			return false, cr
		}
	}
	node.analysisNode.MarkedDelete = true
	cr.Parallel = bind(node.CompletedResult.Cleanup(destination), node.analysisNode)
	cr.Consecutive = bind([]telegram.Request{
		node.NewResult.InitAction(destination),
	}, &analysisNode{
		CompletedResult: &render.CompletedResult{
			MessageIds:   nil,
			LatestResult: node.NewResult,
			MountedAt:    time.Now(),
		},
		MarkedDelete: false,
		MarkedInsert: true,
	})
	return true, cr
}

func (resultLine *ResultLine) CompareAndProduce(
	destination telegram.Destination,
	newLine []render.IResult,
	globalReplaceMode bool,
) (result CompareResult) {
	nodes := resultLine.getAnalysisNodes()
	sections := buildSections(nodes, newLine)
	replaceMode := globalReplaceMode
	for _, commonNode := range sections.common {
		var resultSection CompareResult
		replaceMode, resultSection = compare(destination, replaceMode, commonNode)
		commonNode.analysisNode.CompletedResult.LatestResult = commonNode.NewResult
		result.Parallel = append(result.Parallel, resultSection.Parallel...)
		result.Consecutive = append(result.Consecutive, resultSection.Consecutive...)
	}

	for _, node := range sections.orphaned {
		node.MarkedDelete = true
		result.Parallel = append(result.Parallel, bind(node.CompletedResult.Cleanup(destination), node)...)
	}

	for _, node := range sections.new {
		result.Consecutive = append(result.Consecutive, bind([]telegram.Request{node.InitAction(destination)}, &analysisNode{
			CompletedResult: &render.CompletedResult{
				MessageIds:   nil,
				LatestResult: node,
				MountedAt:    time.Now(),
			},
			MarkedDelete: false,
			MarkedInsert: true,
		})...)
	}

	return
}

func (resultLine *ResultLine) Run(result CompareResult, api *bot.Api) error {

	var toDelete []*render.CompletedResult
	wg := sync.WaitGroup{}
	wg.Add(len(result.Parallel))
	for _, request := range result.Parallel {

		if request.Slot.MarkedDelete {
			toDelete = append(toDelete, request.Slot.CompletedResult)
		}

		go func(boundRequest BoundRequest) {
			defer wg.Done()
			response, err := api.Request(boundRequest.Request)
			if err != nil {
				fmt.Println("error with", boundRequest.Request, err)
				return
			}
			messageList, err := telegram.ParseAs[dto.MessageList](response)
			if err != nil {
				return
			}
			boundRequest.Slot.CompletedResult.MessageIds = messageList.CollectIds()
		}(request)
	}

	for _, request := range result.Consecutive {
		response, err := api.Request(request.Request)
		if err != nil {
			fmt.Println("error with", request.Request, err)
			continue
		}
		messageList, err := telegram.ParseAs[dto.MessageList](response)
		request.Slot.CompletedResult.MessageIds = messageList.CollectIds()
		if request.Slot.MarkedInsert {
			resultLine.Line = append(resultLine.Line, request.Slot.CompletedResult)
		}
	}

	var newResults []*render.CompletedResult

	for _, completedResult := range resultLine.Line {
		isToDelete := false
		for _, toDeleteEntry := range toDelete {
			if completedResult == toDeleteEntry {
				isToDelete = true
			}
		}
		if !isToDelete {
			newResults = append(newResults, completedResult)
		}
	}

	resultLine.Line = newResults
	wg.Wait()
	return nil
}
