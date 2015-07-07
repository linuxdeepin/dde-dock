package loader

import (
	"pkg.deepin.io/dde-daemon/graph"
	"pkg.deepin.io/lib/log"
)

type DAGBuilder struct {
	modules         map[string]Module
	enablingModules []string
	disableModules  map[string]struct{}
	flag            EnableFlag

	log *log.Logger

	dag *graph.Data

	nodes          map[string]*graph.Node
	handledModules map[string]struct{}
}

func NewDAGBuilder(loader *Loader, enablingModules []string, disableModules []string, flag EnableFlag) *DAGBuilder {
	disableModulesMap := map[string]struct{}{}
	for _, name := range disableModules {
		if _, ok := loader.modules[name]; !ok {
			loader.log.Warningf("disabled module(%s) is no existed", name)
			continue
		}
		disableModulesMap[name] = struct{}{}
	}

	return &DAGBuilder{
		modules:         loader.modules,
		enablingModules: enablingModules,
		disableModules:  disableModulesMap,
		flag:            flag,
		log:             loader.log,
		dag:             graph.New(),
		nodes:           map[string]*graph.Node{},
		handledModules:  map[string]struct{}{},
	}
}

func createNodeIfNeeded(dag *graph.Data, nodes map[string]*graph.Node, name string) *graph.Node {
	node, ok := nodes[name]
	if !ok {
		nodes[name] = graph.NewNode(name)
		node = nodes[name]
		dag.AddNode(node)
	}

	return node
}

func (builder *DAGBuilder) buildDAG(enablingModules []string) error {
	logLevel := builder.log.GetLogLevel()
	for _, name := range enablingModules {
		if _, ok := builder.handledModules[name]; ok {
			continue
		}

		builder.handledModules[name] = struct{}{}

		module, ok := builder.modules[name]
		if !ok {
			if builder.flag.HasFlag(EnableFlagIgnoreMissingModule) {
				if logLevel == log.LevelDebug {
					builder.log.Info("no such a module named", name)
					continue
				}
			} else {
				return &EnableError{ModuleName: name, Code: ErrorMissingModule}
			}
		}

		if _, ok := builder.disableModules[name]; ok {
			if !builder.flag.HasFlag(EnableFlagForceStart) {
				return &EnableError{ModuleName: name, Code: ErrorConflict}
			}

			// TODO: add a flag: skip module whose dependencies is not disabled.
		}

		node := createNodeIfNeeded(builder.dag, builder.nodes, name)
		dependencies := module.GetDependencies()

		for _, dependency := range dependencies {
			if _, ok := builder.modules[dependency]; !ok {
				// TODO: add a flag: skip modules whose dependencies is not meet.
				return &EnableError{ModuleName: name, Code: ErrorNoDependencies, detail: dependency}
			}

			depNode := createNodeIfNeeded(builder.dag, builder.nodes, dependency)
			builder.dag.UpdateEdgeWeight(depNode, node, 0)
		}

		err := builder.buildDAG(dependencies)
		if err != nil {
			return err
		}
	}

	return nil
}

func (builder *DAGBuilder) Execute() (*graph.Data, error) {
	err := builder.buildDAG(builder.enablingModules)
	if err != nil {
		return nil, err
	}

	return builder.dag, nil
}
